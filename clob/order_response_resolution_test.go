package clob

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

func TestResolveTransactionsHashesUsesTradeIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != clobtypes.GET_TRADES {
			t.Fatalf("request path = %s, want %s", r.URL.Path, clobtypes.GET_TRADES)
		}
		if tradeID := r.URL.Query().Get("id"); tradeID != "trade-1" {
			t.Fatalf("trade id = %q, want %q", tradeID, "trade-1")
		}
		_ = json.NewEncoder(w).Encode(clobtypes.Trades{
			Data: []clobtypes.Trade{{
				ID:              "trade-1",
				Status:          "MINED",
				TransactionHash: "0xabc",
			}},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil, nil, nil)
	response := &clobtypes.OrderResponse{TradeIDs: []string{"trade-1"}}
	auth := &sdktypes.AuthOption{
		SingerAddress: "0xsigner",
		ApiKeyCreds: &sdktypes.ApiKeyCreds{
			ApiKey:     "api-key",
			Secret:     "c2VjcmV0",
			Passphrase: "passphrase",
		},
	}

	client.resolveTransactionsHashes(context.Background(), response, auth)
	if want := []string{"0xabc"}; !reflect.DeepEqual(response.TransactionsHashes, want) {
		t.Fatalf("transaction hashes = %v, want %v", response.TransactionsHashes, want)
	}
}

func TestResolveTransactionsHashesPreservesPartialResultsWhenPollingTimesOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("id") {
		case "trade-1":
			_ = json.NewEncoder(w).Encode(clobtypes.Trades{
				Data: []clobtypes.Trade{{
					ID:              "trade-1",
					Status:          "MINED",
					TransactionHash: "0xabc",
				}},
			})
		case "trade-2":
			<-r.Context().Done()
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, nil, nil, nil)
	client.resolveTradesTimeout = 50 * time.Millisecond
	client.resolveTradesPollInterval = time.Millisecond
	response := &clobtypes.OrderResponse{TradeIDs: []string{"trade-1", "trade-2"}}
	startedAt := time.Now()
	client.resolveTransactionsHashes(context.Background(), response, testAuthOption())

	if elapsed := time.Since(startedAt); elapsed > 250*time.Millisecond {
		t.Fatalf("resolution took %s, want it bounded by the polling timeout", elapsed)
	}
	if want := []string{"0xabc"}; !reflect.DeepEqual(response.TransactionsHashes, want) {
		t.Fatalf("transaction hashes = %v, want %v", response.TransactionsHashes, want)
	}
}

func TestResolveTransactionsHashesExcludesFailedTrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(clobtypes.Trades{
			Data: []clobtypes.Trade{{
				ID:     r.URL.Query().Get("id"),
				Status: "FAILED",
			}},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, nil, nil, nil)
	response := &clobtypes.OrderResponse{TradeIDs: []string{"trade-1"}}
	client.resolveTransactionsHashes(context.Background(), response, testAuthOption())

	if len(response.TransactionsHashes) != 0 {
		t.Fatalf("transaction hashes = %v, want none for failed trade", response.TransactionsHashes)
	}
}

func TestResolveTransactionsHashesKeepsExistingHashesWithoutPolling(t *testing.T) {
	client := NewClient("http://127.0.0.1:1", nil, nil, nil)
	response := &clobtypes.OrderResponse{
		TransactionsHashes: []string{"0xexisting"},
		TradeIDs:           []string{"trade-1"},
	}
	client.resolveTransactionsHashes(context.Background(), response, testAuthOption())

	if want := []string{"0xexisting"}; !reflect.DeepEqual(response.TransactionsHashes, want) {
		t.Fatalf("transaction hashes = %v, want %v", response.TransactionsHashes, want)
	}
}

func testAuthOption() *sdktypes.AuthOption {
	return &sdktypes.AuthOption{
		SingerAddress: "0xsigner",
		ApiKeyCreds: &sdktypes.ApiKeyCreds{
			ApiKey:     "api-key",
			Secret:     "c2VjcmV0",
			Passphrase: "passphrase",
		},
	}
}
