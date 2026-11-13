package clob

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
)

func TestOrderQueryAssetRequiresMarket(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer server.Close()
	client := NewClient(server.URL, nil, nil, nil)
	for _, market := range []string{"", " ", "\t\n"} {
		result, err := client.GetOrders(context.Background(), clobtypes.GetActiveOrdersRequest{AssetID: "123", Market: market}, testAuthOption())
		if err == nil || !strings.Contains(err.Error(), "asset_id requires market") || result != nil {
			t.Fatalf("market=%q: result=%v err=%v", market, result, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("invalid requests reached server")
	}
}

func TestOrderQueriesWithoutAssociatedTrades(t *testing.T) {
	for _, req := range []clobtypes.GetActiveOrdersRequest{{}, {ID: "order-1"}, {Market: "condition-1"}, {Market: "condition-1", AssetID: "123"}} {
		t.Run(fmt.Sprintf("%+v", req), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != clobtypes.GET_OPEN_ORDERS {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				q := r.URL.Query()
				if q.Get("market") != req.Market || q.Get("asset_id") != req.AssetID || q.Get("id") != req.ID {
					t.Errorf("query=%v", q)
				}
				if req.AssetID == "" && q.Has("asset_id") {
					t.Error("empty asset sent")
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[{"id":"order-1","size_matched":"2","asset_id":"123"}],"next_cursor":"LTE=","limit":100,"count":1}`)
			}))
			defer server.Close()
			result, err := NewClient(server.URL, nil, nil, nil).GetOrders(context.Background(), req, testAuthOption())
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Data) != 1 || result.Data[0].ID != "order-1" || result.Data[0].SizeMatched != "2" || result.Data[0].AssociateTrades != nil || result.NextCursor != "LTE=" {
				t.Fatalf("unexpected result %+v", result)
			}
		})
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != clobtypes.GET_ORDER+"order-1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"order-1","size_matched":"2","asset_id":"123"}`)
	}))
	defer server.Close()
	result, err := NewClient(server.URL, nil, nil, nil).GetOrder(context.Background(), "order-1", clobtypes.GetOrderRequest{}, testAuthOption())
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "order-1" || result.SizeMatched != "2" || result.AssociateTrades != nil {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestTradeQueryFiltersPaginationAndOrderRelations(t *testing.T) {
	market, asset, before, after, cursor := "condition-1", "123", "200", "100", "cursor+/="
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/data/trades" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		// Assert the documented path and its HMAC independently from SDK constants
		// and signing helpers. Query parameters are not part of the signed path.
		mac := hmac.New(sha256.New, []byte("secret"))
		mac.Write([]byte(r.Header.Get("POLY_TIMESTAMP") + "GET/data/trades"))
		if want := base64.URLEncoding.EncodeToString(mac.Sum(nil)); r.Header.Get("POLY_SIGNATURE") != want {
			t.Errorf("trade signature does not cover the documented path")
		}
		q := r.URL.Query()
		for key, want := range map[string]string{"market": market, "asset_id": asset, "before": before, "after": after, "maker_address": "wallet"} {
			if q.Get(key) != want {
				t.Errorf("%s=%q, want %q", key, q.Get(key), want)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if q.Get("next_cursor") == "" {
			fmt.Fprintf(w, `{"data":[{"id":"trade-1","taker_order_id":"taker-1","asset_id":"123","status":"CONFIRMED","transaction_hash":"0xabc","maker_orders":[{"order_id":"maker-1","maker_address":"wallet","asset_id":"123","side":"SELL","matched_amount":"1.25","price":"0.4"}]}],"next_cursor":%q,"limit":1,"count":1}`, cursor)
		} else {
			if q.Get("next_cursor") != cursor {
				t.Errorf("cursor=%q", q.Get("next_cursor"))
			}
			fmt.Fprint(w, `{"data":[{"id":"trade-2","taker_order_id":"taker-2"}],"next_cursor":"LTE=","limit":1,"count":1}`)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, nil, nil, nil)
	req := clobtypes.GetTradesRequest{MakerAddress: "wallet", Market: &market, AssetID: &asset, Before: &before, After: &after}
	page, err := client.GetTrades(context.Background(), req, testAuthOption())
	if err != nil {
		t.Fatal(err)
	}
	if page.NextCursor != cursor || len(page.Data) != 1 {
		t.Fatalf("page=%+v", page)
	}
	trade := page.Data[0]
	wantMaker := map[string]any{"order_id": "maker-1", "maker_address": "wallet", "asset_id": "123", "side": "SELL", "matched_amount": "1.25", "price": "0.4"}
	if trade.ID != "trade-1" || trade.TakerOrderID != "taker-1" || trade.Status != "CONFIRMED" || trade.TransactionHash != "0xabc" || len(trade.MakerOrders) != 1 || !reflect.DeepEqual(trade.MakerOrders[0], wantMaker) {
		t.Fatalf("trade=%+v", trade)
	}
	req.NextCursor = &page.NextCursor
	page, err = client.GetTrades(context.Background(), req, testAuthOption())
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != 1 || page.Data[0].ID != "trade-2" || page.NextCursor != "LTE=" || calls.Load() != 2 {
		t.Fatalf("page=%+v calls=%d", page, calls.Load())
	}
}

func TestQueryErrorsAreNotEmptySuccess(t *testing.T) {
	for _, malformed := range []bool{false, true} {
		t.Run(fmt.Sprint(malformed), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if malformed {
					fmt.Fprint(w, `{"data":`)
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprint(w, `{"error":"invalid params"}`)
				}
			}))
			defer server.Close()
			client := NewClient(server.URL, nil, nil, nil)
			if result, err := client.GetTrades(context.Background(), clobtypes.GetTradesRequest{}, testAuthOption()); err == nil || result != nil {
				t.Fatalf("trades=%v err=%v", result, err)
			}
			if result, err := client.GetOrders(context.Background(), clobtypes.GetActiveOrdersRequest{}, testAuthOption()); err == nil || result != nil {
				t.Fatalf("orders=%v err=%v", result, err)
			}
			if result, err := client.GetOrder(context.Background(), "order-1", clobtypes.GetOrderRequest{}, testAuthOption()); err == nil || result != nil {
				t.Fatalf("order=%v err=%v", result, err)
			}
		})
	}
}
