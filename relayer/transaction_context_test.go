package relayer

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/override-coder/go-polymarket-sdk/relayer/types"
)

func TestPollUntilStateCancelsInFlightRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer server.Close()
	defer close(release)
	c := NewClient(server.URL, big.NewInt(137), nil, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := c.PollUntilState(ctx, "id", nil, nil, 1, time.Second)
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("poll did not cancel its HTTP request")
	}
}

func TestGetTransactionCompatibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != types.GET_TRANSACTION || r.URL.Query().Get("id") != "id" {
			t.Errorf("unexpected request: %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"transactionID":"id","state":"STATE_CONFIRMED"}]`))
	}))
	defer server.Close()
	c := NewClient(server.URL, big.NewInt(137), nil, nil)
	for _, fetch := range []func() ([]types.RelayerTransaction, error){
		func() ([]types.RelayerTransaction, error) { return c.GetTransaction("id") },
		func() ([]types.RelayerTransaction, error) { return c.GetTransactionContext(context.Background(), "id") },
	} {
		result, err := fetch()
		if err != nil || len(result) != 1 || string(result[0].State) != "STATE_CONFIRMED" {
			t.Fatalf("result=%v err=%v", result, err)
		}
	}
}
