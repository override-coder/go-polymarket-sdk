package combos

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
)

type idleWSSession struct {
	conn   *websocket.Conn
	closed chan struct{}
}

func idleWSServer(t *testing.T) (string, <-chan idleWSSession) {
	t.Helper()
	sessions := make(chan idleWSSession, 100)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		if _, _, err = conn.ReadMessage(); err != nil {
			return
		} // auth
		s := idleWSSession{conn, make(chan struct{})}
		sessions <- s
		defer close(s.closed)
		for {
			if _, _, err = conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	return "ws" + strings.TrimPrefix(server.URL, "http"), sessions
}
func nextIdleSession(t *testing.T, sessions <-chan idleWSSession) idleWSSession {
	t.Helper()
	select {
	case s := <-sessions:
		return s
	case <-time.After(time.Second):
		t.Fatal("server did not authenticate")
	}
	return idleWSSession{}
}
func wsTestAuth() *sdk.AuthOption {
	return &sdk.AuthOption{ApiKeyCreds: &sdk.ApiKeyCreds{ApiKey: "test", Secret: "test", Passphrase: "test"}}
}

func TestWSReadCancellationClosesCapturedConnection(t *testing.T) {
	host, sessions := idleWSServer(t)
	c := NewWSClient(host)
	defer c.Close()
	for i := 0; i < 30; i++ {
		if err := c.Connect(context.Background(), wsTestAuth()); err != nil {
			t.Fatal(err)
		}
		session := nextIdleSession(t, sessions)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result := make(chan error, 1)
		go func() { _, err := c.Read(ctx); result <- err }()
		select {
		case err := <-result:
			if !errors.Is(err, context.Canceled) {
				t.Fatal(err)
			}
		case <-time.After(time.Second):
			t.Fatal("cancelled Read blocked")
		}
		select {
		case <-session.closed:
		case <-time.After(time.Second):
			t.Fatal("cancelled Read left socket open")
		}
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()
		if conn != nil {
			t.Fatal("cancelled connection still installed")
		}
	}
}

func TestWSStaleReaderCannotCloseReplacement(t *testing.T) {
	host, sessions := idleWSServer(t)
	c := NewWSClient(host)
	defer c.Close()
	if err := c.Connect(context.Background(), wsTestAuth()); err != nil {
		t.Fatal(err)
	}
	first := nextIdleSession(t, sessions)
	c.mu.Lock()
	old := c.conn
	c.mu.Unlock()
	if err := c.Connect(context.Background(), wsTestAuth()); err != nil {
		t.Fatal(err)
	}
	second := nextIdleSession(t, sessions)
	select {
	case <-first.closed:
	case <-time.After(time.Second):
		t.Fatal("reconnect leaked old socket")
	}
	// Cleanup of a reader that captured the first connection can run after Connect.
	c.closeConnection(old)
	if err := second.conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"unknown","marker":"new"}`)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	event, err := c.Read(ctx)
	if err != nil || event == nil || !strings.Contains(string(event.Raw), "new") {
		t.Fatalf("replacement event=%v err=%v", event, err)
	}
}
