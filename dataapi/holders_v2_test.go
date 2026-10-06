package dataapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/require"
)

func TestV2HoldersAllParametersAndFields(t *testing.T) {
	fixture := v2Fixture(t, "holder")
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		require.Equal(t, "/v2/holders", r.URL.Path)
		expected := url.Values{"condition": {v2Condition}, "min_balance": {"0.25"}, "include_pnl": {"true"}}
		if calls == 1 {
			expected.Set("limit", "100")
		} else {
			expected.Set("cursor", "holders+/=")
		}
		require.Equal(t, expected, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			fmt.Fprintf(w, `{"data":[{"token_id":"123","holders":[%s]},{"token_id":"456","holders":[%s]}],"pagination":{"limit":100,"offset":0,"has_more":true,"next_cursor":"holders+/="}}`, fixture, fixture)
		} else {
			fmt.Fprintf(w, `{"data":[{"token_id":"456","holders":[%s]}],"pagination":{"limit":100,"offset":100,"has_more":false,"next_cursor":null}}`, fixture)
		}
	}))
	defer server.Close()
	client := dataapi.NewClient(server.URL, nil)
	q := types.TopHoldersQuery{Condition: []string{v2Condition}, Limit: v2Ptr(100), MinBalance: v2Ptr(0.25), IncludePnL: v2Ptr(true)}
	page, err := client.GetTopHoldersForMarkets(context.Background(), q)
	require.NoError(t, err)
	require.Len(t, page.Data, 2)
	assertV2JSON(t, fixture, page.Data[0].Holders[0])
	require.NotNil(t, page.Data[0].Holders[0].RealizedPnl)
	require.Zero(t, *page.Data[0].Holders[0].RealizedPnl)
	q.Cursor = page.Pagination.NextCursor
	q.Limit = v2Ptr(1001) // Ignored; the signed cursor's window wins.
	page, err = client.GetTopHoldersForMarkets(context.Background(), q)
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	require.Equal(t, "456", page.Data[0].TokenID)
	require.Nil(t, page.Pagination.NextCursor)
}

func TestV2HoldersDefaultsAndNullableEconomics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, url.Values{"condition": {v2Condition}}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"token_id":"123","holders":[{"amount":0,"avg_price":null,"realized_pnl":null}]}],"pagination":{"limit":100,"offset":0,"has_more":false,"next_cursor":null}}`)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetTopHoldersForMarkets(context.Background(), types.TopHoldersQuery{Condition: []string{v2Condition}})
	require.NoError(t, err)
	holder := page.Data[0].Holders[0]
	require.Nil(t, holder.AvgPrice)
	require.Nil(t, holder.RealizedPnl)
	require.Nil(t, holder.EntryCostUSDC)
	require.Nil(t, holder.CurrentPrice)
	require.Nil(t, holder.CurrentValue)
	require.Nil(t, holder.UnrealizedPnl)
	require.Nil(t, holder.TotalPnl)
}

func TestV2HoldersValidation(t *testing.T) {
	many := make([]string, 21)
	for i := range many {
		many[i] = fmt.Sprintf("0x%064x", i)
	}
	for name, q := range map[string]types.TopHoldersQuery{
		"missing condition":       {},
		"invalid condition":       {Condition: []string{"bad"}},
		"too many conditions":     {Condition: many},
		"large limit":             {Condition: many[:2], Limit: v2Ptr(1001)},
		"pnl multiple conditions": {Condition: many[:2], IncludePnL: v2Ptr(true)},
		"pnl large limit":         {Condition: many[:1], IncludePnL: v2Ptr(true), Limit: v2Ptr(101)},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := dataapi.NewClient("http://127.0.0.1:1", nil).GetTopHoldersForMarkets(context.Background(), q)
			require.Error(t, err)
		})
	}
}
