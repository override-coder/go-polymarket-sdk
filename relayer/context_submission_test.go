package relayer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	rt "github.com/override-coder/go-polymarket-sdk/relayer/types"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
)

func TestGetTransactionsRejectsBuilderSignatureError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); fmt.Fprint(w, "[]") }))
	defer server.Close()
	c := NewClient(server.URL, big.NewInt(137), nil, &sdk.BuilderApiKeyCreds{Key: "key", Secret: "invalid%base64"})
	if _, err := c.GetTransactions(); err == nil {
		t.Fatal("signing error lost")
	}
	if _, err := c.GetTransactionsContext(context.Background()); err == nil {
		t.Fatal("context variant lost signing error")
	}
	if calls.Load() != 0 {
		t.Fatal("invalid credentials reached HTTP")
	}
}

func TestRelayerContextCancellationStages(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		invoke     func(*Client, context.Context, *sdk.AuthOption) error
		deployed   bool
	}{
		{"safe deployed", rt.GET_DEPLOYED, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteContext(ctx, []rt.SafeTransaction{{To: a.SingerAddress, Value: "0", Data: "0x"}}, "", a)
			return e
		}, true},
		{"safe nonce", rt.GET_NONCE, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteContext(ctx, []rt.SafeTransaction{{To: a.SingerAddress, Value: "0", Data: "0x"}}, "", a)
			return e
		}, true},
		{"safe submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteContext(ctx, []rt.SafeTransaction{{To: a.SingerAddress, Value: "0", Data: "0x"}}, "", a)
			return e
		}, true},
		{"safe signed submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteByTxContext(ctx, &rt.TransactionRequest{}, a)
			return e
		}, true},
		{"safe deploy submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.DeployContext(ctx, a)
			return e
		}, false},
		{"deposit nonce", rt.GET_NONCE, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteDepositWalletBatchContext(ctx, []rt.DepositWalletCall{{Target: a.SingerAddress, Value: "0", Data: "0x"}}, "2000000000", nil, a)
			return e
		}, true},
		{"deposit submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteDepositWalletBatchContext(ctx, []rt.DepositWalletCall{{Target: a.SingerAddress, Value: "0", Data: "0x"}}, "2000000000", nil, a)
			return e
		}, true},
		{"deposit signed submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.ExecuteDepositWalletBatchByTxContext(ctx, &rt.DepositWalletBatchRequest{}, a)
			return e
		}, true},
		{"deposit deploy submit", rt.SUBMIT_TRANSACTION, func(c *Client, ctx context.Context, a *sdk.AuthOption) error {
			_, e := c.DeployDepositWalletContext(ctx, a)
			return e
		}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			started := make(chan struct{}, 1)
			release := make(chan struct{})
			var signs, submits atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == rt.SUBMIT_TRANSACTION {
					submits.Add(1)
				}
				if r.URL.Path == tc.path {
					select {
					case started <- struct{}{}:
					default:
					}
					select {
					case <-r.Context().Done():
					case <-release:
					}
					return
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case rt.GET_DEPLOYED:
					fmt.Fprintf(w, `{"deployed":%t}`, tc.deployed)
				case rt.GET_NONCE:
					fmt.Fprint(w, `{"nonce":"1"}`)
				default:
					t.Errorf("unexpected endpoint %s", r.URL.Path)
				}
			}))
			defer server.Close()
			defer close(release)
			sign := func(string, []byte) ([]byte, error) {
				signs.Add(1)
				sig := make([]byte, 65)
				sig[64] = 27
				return sig, nil
			}
			c := NewClient(server.URL, big.NewInt(137), sign, &sdk.BuilderApiKeyCreds{Key: "key", Secret: "c2VjcmV0", Passphrase: "pass"})
			a := &sdk.AuthOption{SingerAddress: "0x0000000000000000000000000000000000000001"}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- tc.invoke(c, ctx, a) }()
			select {
			case <-started:
			case e := <-done:
				t.Fatalf("stopped before target stage: %v", e)
			case <-time.After(2 * time.Second):
				t.Fatal("target stage not reached")
			}
			cancel()
			select {
			case e := <-done:
				if !errors.Is(e, context.Canceled) {
					t.Fatalf("error chain lost cancellation: %v", e)
				}
			case <-time.After(time.Second):
				t.Fatal("cancellation did not stop request")
			}
			if tc.path != rt.SUBMIT_TRANSACTION && (signs.Load() != 0 || submits.Load() != 0) {
				t.Fatal("cancelled preparation reached signing/submission")
			}
			if tc.path == rt.SUBMIT_TRANSACTION && submits.Load() != 1 {
				t.Fatalf("submission retried after cancellation: %d", submits.Load())
			}
		})
	}
}

func TestRelayerCancelledSignerDoesNotSubmit(t *testing.T) {
	var submits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case rt.GET_DEPLOYED:
			fmt.Fprint(w, `{"deployed":true}`)
		case rt.GET_NONCE:
			fmt.Fprint(w, `{"nonce":"1"}`)
		default:
			submits.Add(1)
			fmt.Fprint(w, `{"transactionID":"unexpected"}`)
		}
	}))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sign := func(string, []byte) ([]byte, error) { cancel(); sig := make([]byte, 65); sig[64] = 27; return sig, nil }
	c := NewClient(server.URL, big.NewInt(137), sign, &sdk.BuilderApiKeyCreds{Key: "key", Secret: "c2VjcmV0"})
	a := &sdk.AuthOption{SingerAddress: "0x0000000000000000000000000000000000000001"}
	_, err := c.ExecuteContext(ctx, []rt.SafeTransaction{{To: a.SingerAddress, Value: "0", Data: "0x"}}, "", a)
	if !errors.Is(err, context.Canceled) || submits.Load() != 0 {
		t.Fatalf("err=%v submits=%d", err, submits.Load())
	}
}

func TestRelayerPreCancelledContextVariants(t *testing.T) {
	c := NewClient("http://invalid", big.NewInt(137), nil, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	checks := []func() error{
		func() error { _, e := c.GetNonceContext(ctx, "", rt.TransactionTypeSAFE); return e },
		func() error { _, e := c.GetDeployedContext(ctx, ""); return e },
		func() error { _, e := c.GetTransactionsContext(ctx); return e },
		func() error { _, e := c.DeployContext(ctx, nil); return e },
		func() error { _, e := c.BuildDeployTxContext(ctx, nil); return e },
		func() error { _, e := c.BuildTxContext(ctx, nil, nil, "", nil); return e },
		func() error { _, e := c.ExecuteContext(ctx, nil, "", nil); return e },
		func() error { _, e := c.ExecuteByTxContext(ctx, nil, nil); return e },
		func() error { _, e := c.DeployDepositWalletContext(ctx, nil); return e },
		func() error { _, e := c.BuildDepositWalletBatchContext(ctx, nil, "", nil, nil, nil); return e },
		func() error { _, e := c.ExecuteDepositWalletBatchContext(ctx, nil, "", nil, nil); return e },
		func() error { _, e := c.ExecuteDepositWalletBatchByTxContext(ctx, nil, nil); return e },
	}
	for i, check := range checks {
		if e := check(); !errors.Is(e, context.Canceled) {
			t.Fatalf("variant %d: %v", i, e)
		}
	}
}
