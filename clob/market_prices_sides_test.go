package clob

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
)

func TestMarketPricesPreserveSides(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Path != "/prices" || r.Method != "POST" {
			t.Errorf("request: %s %s", r.Method, r.URL)
		}
		var body []ct.PricesRequest
		if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
			t.Error(e)
		}
		if len(body) == 0 {
			t.Error("empty request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"1":{"BUY":"0.45","SELL":"0.65"},"2":{"BUY":"0.1"}}`))
	}))
	defer server.Close()
	c := NewClient(server.URL, nil, nil, nil)
	req := []ct.PricesRequest{{TokenId: "1", Side: "BUY"}, {TokenId: "1", Side: "SELL"}}
	got, e := c.GetMarketPricesBySide(context.Background(), req)
	if e != nil {
		t.Fatal(e)
	}
	if got["1"]["BUY"] != "0.45" || got["1"]["SELL"] != "0.65" {
		t.Fatal(got)
	}
	before := calls.Load()
	if _, e = c.GetMarketPrices(context.Background(), req); e == nil {
		t.Fatal("ambiguous legacy query accepted")
	}
	if calls.Load() != before {
		t.Fatal("ambiguous request sent")
	}
	for _, side := range []string{"BUY", "SELL"} {
		got, e := c.GetMarketPrices(context.Background(), []ct.PricesRequest{{TokenId: "1", Side: side}, {TokenId: "1", Side: side}, {TokenId: "2", Side: "SELL"}, {TokenId: "missing", Side: "BUY"}})
		want := "0.45"
		if side == "SELL" {
			want = "0.65"
		}
		if e != nil || !reflect.DeepEqual(got, map[string]string{"1": want}) {
			t.Fatalf("%v %v", got, e)
		}
	}
}

func TestMarketPricesErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"bad prices"}`))
	}))
	defer server.Close()
	c := NewClient(server.URL, nil, nil, nil)
	_, e := c.GetMarketPricesBySide(context.Background(), nil)
	var upstream *sdkhttp.UpstreamServiceError
	if !errors.As(e, &upstream) || upstream.StatusCode != 400 {
		t.Fatal(e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := c.GetMarketPricesBySide(ctx, nil); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
