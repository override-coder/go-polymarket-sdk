package clob

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
)

func TestGetOrdersCursorPages(t *testing.T) {
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		q := r.URL.Query()
		if q.Get("market") != "condition" || q.Get("asset_id") != "123" || q.Get("id") != "order" {
			t.Errorf("filters lost: %v", q)
		}
		w.Header().Set("Content-Type", "application/json")
		switch q.Get("next_cursor") {
		case "":
			if q.Has("next_cursor") {
				t.Error("empty first-page cursor sent")
			}
			fmt.Fprint(w, `{"data":[{"id":"one"}],"next_cursor":"cursor+/="}`)
		case "cursor+/=":
			fmt.Fprint(w, `{"data":[],"next_cursor":"last"}`)
		case "last":
			fmt.Fprint(w, `{"data":[{"id":"two"}],"next_cursor":"LTE="}`)
		default:
			t.Errorf("unexpected cursor: %v", q)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, nil, nil, nil)
	req := ct.GetActiveOrdersRequest{ID: "order", Market: "condition", AssetID: "123"}
	var ids []string
	for {
		page, err := client.GetOrders(context.Background(), req, testAuthOption())
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range page.Data {
			ids = append(ids, o.ID)
		}
		if page.NextCursor == "" || page.NextCursor == "LTE=" {
			break
		}
		req.NextCursor = &page.NextCursor
		if count.Load() > 3 {
			t.Fatal("cursor did not advance")
		}
	}
	if count.Load() != 3 || len(ids) != 2 || ids[0] != "one" || ids[1] != "two" {
		t.Fatalf("requests=%d ids=%v", count.Load(), ids)
	}
}
