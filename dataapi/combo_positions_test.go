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

const comboPositionFixture = `{"combo_condition_id":"0x03","combo_position_id":"123456789012345678901234567890","proxy_wallet":"0xwallet","outcome_index":0,"outcome_label":"YES","current_size":28.380001,"entry_avg_price_usdc":0.75,"entry_cost_usdc":19.5,"gross_entry_cost_usdc":20.000001,"entry_fees_usdc":0.500001,"realized_payout_usdc":0,"status":"OPEN","redeemable":false,"first_entry_at":"2026-10-19T00:00:00Z","first_entry_at_micros":null,"updated_at":"2026-10-19T00:00:01Z","updated_at_micros":1792368001000000,"resolved_at":null,"legs_total":2,"legs_resolved":0,"legs_pending":2,"legs":[{"leg_index":0,"leg_position_id":"111","leg_condition_id":"0xleg","leg_outcome_index":1,"leg_outcome_label":"NO","leg_status":"OPEN","leg_current_price":0.59,"leg_resolved_at":null,"market":{"market_id":"123","slug":"a","title":"match","outcome":"NO","image_url":"image","icon_url":"icon","category":"sports","subcategory":"tennis","tags":[],"end_date":"2026-10-20T00:00:00Z","event":{"event_id":"1","event_slug":"match","event_title":"Match","event_image":"image"},"group_item_title":"Player","line":null,"outcomes":["YES","NO"],"question":"Win?","sports_market_type":"moneyline"}}]}`

func TestComboPositionsParametersFieldsAndCursor(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/positions/combos", r.URL.Path)
		if calls == 1 {
			require.Equal(t, url.Values{"user": {v2User}, "limit": {"2"}, "condition": {v2Condition}, "status": {"OPEN"}, "sort_by": {"ENTRY_COST"}, "sort_direction": {"DESC"}, "updated_after": {"0"}, "updated_before": {"1792368000"}}, r.URL.Query())
		} else {
			require.Equal(t, url.Values{"user": {v2User}, "cursor": {"next+/="}}, r.URL.Query())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"pagination":{"limit":2,"offset":0,"has_more":true,"next_cursor":"next+/="}}`, comboPositionFixture)
	}))
	defer server.Close()
	c := dataapi.NewClient(server.URL, nil)
	page, err := c.GetComboPositions(context.Background(), types.ComboPositionsQuery{User: " " + v2User + " ", Condition: []string{v2Condition}, Limit: v2Ptr(2), Status: v2Ptr("OPEN"), SortBy: v2Ptr("ENTRY_COST"), SortDirection: v2Ptr(types.SortDESC), UpdatedAfter: v2Ptr(int64(0)), UpdatedBefore: v2Ptr(int64(1792368000))})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assertV2JSON(t, comboPositionFixture, page.Data[0])
	require.Equal(t, "28.380001", page.Data[0].CurrentSize.String())
	require.Equal(t, "20.000001", page.Data[0].GrossEntryCostUSDC.String())
	_, err = c.GetComboPositions(context.Background(), types.ComboPositionsQuery{User: v2User, Cursor: page.Pagination.NextCursor, Limit: v2Ptr(-1)})
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestComboPositionsValidationAndErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"code":"rate_limited","error":"busy"}`)
	}))
	defer server.Close()
	c := dataapi.NewClient(server.URL, nil)
	for _, q := range []types.ComboPositionsQuery{
		{}, {User: v2User, Limit: v2Ptr(1001)}, {User: v2User, Condition: []string{"bad"}}, {User: v2User, Status: v2Ptr("CLOSED")}, {User: v2User, Status: v2Ptr("REDEEMABLE,OPEN")}, {User: v2User, SortBy: v2Ptr("PRICE")}, {User: v2User, SortDirection: v2Ptr(types.SortDirection("bad"))}, {User: v2User, UpdatedAfter: v2Ptr(int64(-1))}, {User: v2User, UpdatedAfter: v2Ptr(int64(2)), UpdatedBefore: v2Ptr(int64(1))},
	} {
		_, err := c.GetComboPositions(context.Background(), q)
		require.Error(t, err)
	}
	_, err := c.GetComboPositions(context.Background(), types.ComboPositionsQuery{User: v2User})
	require.ErrorContains(t, err, "429")
}
