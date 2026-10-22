package gamma_test

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/gamma"
	"github.com/override-coder/go-polymarket-sdk/gamma/types"
)

func TestGetMarketsByKeysetPositionIDs(t *testing.T) {
	const first = "548324018002783558052327573095137442627908907682778868240145789358304133120"
	const second = "548324018002783558052327573095137442627908907682778868240145789358304133121"
	for _, tc := range []struct {
		name string
		ids  []string
		want []string
	}{
		{name: "omitted"},
		{name: "empty", ids: []string{}},
		{name: "blank", ids: []string{"", " \t "}},
		{name: "single", ids: []string{first}, want: []string{first}},
		{name: "multiple trimmed", ids: []string{" " + first + " ", "", "\t" + second + "\n"}, want: []string{first, second}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/markets/keyset" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				query := r.URL.Query()
				if got := query["position_ids"]; !reflect.DeepEqual(got, tc.want) {
					t.Errorf("position_ids = %v, want %v", got, tc.want)
				}
				if got := query["condition_ids"]; !reflect.DeepEqual(got, []string{"condition-1", "condition-2"}) {
					t.Errorf("condition_ids = %v", got)
				}
				if query.Get("clob_token_ids") != "token-1" || query.Get("after_cursor") != "cursor-1" {
					t.Errorf("existing query parameters changed: %v", query)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"markets":[],"next_cursor":"cursor-2"}`)
			}))
			defer server.Close()
			cursor := "cursor-1"
			client := gamma.NewClient(server.URL, big.NewInt(137))
			response, err := client.GetMarketsByKeyset(context.Background(), &types.GetMarketsKeysetParams{
				PositionIDs:  tc.ids,
				ConditionIDs: []string{"condition-1", "condition-2"},
				CLOBTokenIDs: []string{"token-1"},
				AfterCursor:  &cursor,
			})
			if err != nil {
				t.Fatal(err)
			}
			if response == nil || response.NextCursor == nil || *response.NextCursor != "cursor-2" {
				t.Fatalf("unexpected response: %+v", response)
			}
		})
	}
}
