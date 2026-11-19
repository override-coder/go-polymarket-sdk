package gamma_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/gamma"
	"github.com/override-coder/go-polymarket-sdk/gamma/types"
)

func TestProtocolFilters(t *testing.T) {
	for _, versions := range [][]string{nil, {"v1"}, {"v2"}, {"v1", "v2"}} {
		for _, endpoint := range []string{"/markets", "/markets/keyset", "/events/keyset"} {
			t.Run(fmt.Sprint(endpoint, versions), func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != endpoint || !reflect.DeepEqual(r.URL.Query()["version"], versions) {
						t.Errorf("query: %s", r.URL)
					}
					if endpoint == "/markets/keyset" && (r.URL.Query().Get("resolution_status") != "resolved" || r.URL.Query().Get("after_cursor") != "next" || r.URL.Query().Get("position_ids") != "123") {
						t.Errorf("lost filters: %s", r.URL)
					}
					w.Header().Set("Content-Type", "application/json")
					switch endpoint {
					case "/markets":
						fmt.Fprint(w, `[]`)
					case "/markets/keyset":
						fmt.Fprint(w, `{"markets":[],"next_cursor":"last"}`)
					default:
						fmt.Fprint(w, `{"events":[],"next_cursor":"last"}`)
					}
				}))
				defer server.Close()
				c := gamma.NewClient(server.URL, nil)
				var err error
				status, cursor := "resolved", "next"
				switch endpoint {
				case "/markets":
					_, err = c.GetMarkets(context.Background(), &types.GetMarketsParams{Versions: versions})
				case "/markets/keyset":
					_, err = c.GetMarketsByKeyset(context.Background(), &types.GetMarketsKeysetParams{Versions: versions, ResolutionStatus: &status, AfterCursor: &cursor, PositionIDs: []string{"123"}})
				default:
					_, err = c.GetEventsByKeyset(context.Background(), &types.GetEventsKeysetParams{Versions: versions})
				}
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestProtocolInvalidFilters(t *testing.T) {
	c := gamma.NewClient("http://invalid", nil)
	for _, versions := range [][]string{{"v3"}, {""}, {"v1", "v2", "v1"}} {
		if _, e := c.GetMarkets(context.Background(), &types.GetMarketsParams{Versions: versions}); e == nil {
			t.Fatal("expected version error")
		}
		if _, e := c.GetMarketsByKeyset(context.Background(), &types.GetMarketsKeysetParams{Versions: versions}); e == nil {
			t.Fatal("expected version error")
		}
		if _, e := c.GetEventsByKeyset(context.Background(), &types.GetEventsKeysetParams{Versions: versions}); e == nil {
			t.Fatal("expected version error")
		}
	}
	status := "unknown"
	if _, e := c.GetMarketsByKeyset(context.Background(), &types.GetMarketsKeysetParams{ResolutionStatus: &status}); e == nil {
		t.Fatal("expected status error")
	}
	if _, e := c.GetEventsByKeyset(context.Background(), nil); e == nil {
		t.Fatal("expected nil params error")
	}
}

func TestProtocolResponseFields(t *testing.T) {
	for _, body := range []string{`{}`, `{"version":null,"resolutionStatus":null}`, `{"version":"v1","umaResolutionStatus":"resolved","clobTokenIds":"[\"1\",\"2\"]"}`, `{"version":"v2","resolutionStatus":"resolved","positionIds":["1","2"]}`, `{"version":"future","resolutionStatus":"future-state"}`} {
		var want map[string]json.RawMessage
		json.Unmarshal([]byte(body), &want)
		var market types.Market
		if e := json.Unmarshal([]byte(body), &market); e != nil {
			t.Fatal(e)
		}
		encoded, _ := json.Marshal(market)
		var got map[string]json.RawMessage
		json.Unmarshal(encoded, &got)
		for _, key := range []string{"version", "resolutionStatus", "umaResolutionStatus", "clobTokenIds", "positionIds"} {
			if v, ok := want[key]; ok && string(v) != string(got[key]) {
				t.Errorf("%s lost: %s", key, encoded)
			}
		}
		if _, ok := want["version"]; !ok && market.Version != nil {
			t.Fatal("missing version defaulted")
		}
	}
	const event = `{"version":"v2","markets":[{"version":"v2","resolutionStatus":"active","positionIds":["1"]}]}`
	var response types.SearchResponse
	if e := json.Unmarshal([]byte(`{"events":[`+event+`]}`), &response); e != nil {
		t.Fatal(e)
	}
	if len(response.Events) != 1 || response.Events[0].Version == nil || *response.Events[0].Version != "v2" || *response.Events[0].Markets[0].ResolutionStatus != "active" {
		t.Fatalf("nested fields lost: %+v", response)
	}
}
