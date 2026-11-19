package clob

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMarketCacheIndependentCancellation(t *testing.T) {
	for _, resolveToken := range []bool{false, true} {
		for _, cancelLeader := range []bool{false, true} {
			t.Run(fmt.Sprintf("token=%v/leader=%v", resolveToken, cancelLeader), func(t *testing.T) {
				started := make(chan struct{})
				release := make(chan struct{})
				var blocked atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					isToken := r.URL.Path == "/markets-by-token/1"
					if isToken == resolveToken {
						if blocked.Add(1) == 1 {
							close(started)
						}
						select {
						case <-release:
						case <-r.Context().Done():
							return
						}
					}
					w.Header().Set("Content-Type", "application/json")
					if isToken {
						fmt.Fprint(w, `{"condition_id":"condition"}`)
					} else {
						fmt.Fprint(w, `{"mts":0.01,"t":[{"t":"1","o":"yes"}]}`)
					}
				}))
				defer server.Close()
				var releaseOnce sync.Once
				openGate := func() { releaseOnce.Do(func() { close(release) }) }
				defer openGate()
				c := NewClient(server.URL, nil, nil, nil)
				load := func(ctx context.Context) error {
					condition := "condition"
					var hint *string
					if !resolveToken {
						hint = &condition
					}
					_, err := c.orderMarketInfo(ctx, "1", hint)
					return err
				}
				ctxA, cancelA := context.WithCancel(context.Background())
				defer cancelA()
				ctxBParent, cancelB := context.WithCancel(context.Background())
				ctxB := &observedMarketWaiter{Context: ctxBParent, waiting: make(chan struct{})}
				defer cancelB()
				a := make(chan error, 1)
				b := make(chan error, 1)
				go func() { a <- load(ctxA) }()
				select {
				case <-started:
				case <-time.After(time.Second):
					t.Fatal("fetch did not start")
				}
				go func() { b <- load(ctxB) }()
				select {
				case <-ctxB.waiting:
				case <-time.After(time.Second):
					t.Fatal("B did not join shared lookup")
				}
				cancelled, active := b, a
				if cancelLeader {
					cancelA()
					cancelled, active = a, b
				} else {
					cancelB()
				}
				select {
				case err := <-cancelled:
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("cancelled err=%v", err)
					}
				case <-time.After(time.Second):
					t.Fatal("cancelled waiter blocked on shared request")
				}
				openGate()
				select {
				case err := <-active:
					if err != nil {
						t.Fatal(err)
					}
				case <-time.After(time.Second):
					t.Fatal("active waiter failed to finish")
				}
				if blocked.Load() != 1 {
					t.Fatalf("duplicate shared fetches: %d", blocked.Load())
				}
			})
		}
	}
}

func TestMarketCachePreCancelledDoesNotFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := NewClient("http://invalid", nil, nil, nil)
	if _, err := c.orderMarketInfo(ctx, "1", nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := c.getClobMarketInfoCached(ctx, "condition"); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestSharedMarketLookupHasBoundedIndependentContext(t *testing.T) {
	type key struct{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), key{}, "route"))
	defer cancel()
	var group singleflight.Group
	fetchDone := make(chan struct{})
	_, err := sharedMarketLookup(ctx, &group, "key", func(fetchCtx context.Context) (any, error) {
		defer close(fetchDone)
		cancel()
		if fetchCtx.Err() != nil {
			t.Error("caller cancelled the shared fetch")
		}
		deadline, ok := fetchCtx.Deadline()
		if !ok || time.Until(deadline) > 15*time.Second || time.Until(deadline) <= 0 {
			t.Error("shared fetch has no bounded lifetime")
		}
		if fetchCtx.Value(key{}) != "route" {
			t.Error("transport routing context lost")
		}
		return nil, nil
	})
	select {
	case <-fetchDone:
	case <-time.After(time.Second):
		t.Fatal("shared fetch checks did not finish")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled caller received %v", err)
	}
}

// Done is evaluated by sharedMarketLookup only after DoChan has registered this
// caller. Observing it avoids relying on goroutine scheduling or sleep durations.
type observedMarketWaiter struct {
	context.Context
	waiting chan struct{}
	once    sync.Once
}

func (ctx *observedMarketWaiter) Done() <-chan struct{} {
	ctx.once.Do(func() { close(ctx.waiting) })
	return ctx.Context.Done()
}
