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

func TestV2ActivityAllParametersAndFields(t *testing.T) {
	fixture := v2Fixture(t, "activity")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v2/activity", r.URL.Path)
		require.Equal(t, url.Values{
			"user": {v2User}, "limit": {"1000"}, "condition": {v2Condition},
			"type": {"DEPOSIT,WITHDRAWAL,TIP"}, "exclude_deposits_withdrawals": {"false"},
			"start": {"1"}, "end": {"0"}, "sort_by": {"TIMESTAMP"}, "sort_direction": {"ASC"}, "side": {"SELL"},
		}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"pagination":{"limit":1000,"offset":0,"has_more":true,"next_cursor":"activity+/="}}`, fixture)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetUserActivity(context.Background(), types.ActivityQuery{
		User: v2User, Limit: v2Ptr(1000), Condition: []string{v2Condition},
		Type: []types.ActivityType{types.ActivityDEPOSIT, types.ActivityWITHDRAWAL, types.ActivityTIP}, ExcludeDepositsWithdrawals: v2Ptr(false),
		Start: v2Ptr(int64(1)), End: v2Ptr(int64(0)), SortBy: v2Ptr(types.ActivitySortTIMESTAMP), SortDirection: v2Ptr(types.SortASC), Side: v2Ptr(types.SideSELL),
	})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assertV2JSON(t, fixture, page.Data[0])
	require.Equal(t, "activity+/=", *page.Pagination.NextCursor)
}

func TestV2ActivityCursorEventsAndTIP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, url.Values{"user": {v2User}, "event_id": {"12,34"}, "cursor": {"activity+/="}, "type": {"TIP"}, "start": {"1"}}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"type":"TIP","side":"IN","size":1.25,"outcome_index":999}],"pagination":{"limit":100,"offset":100,"has_more":false,"next_cursor":null}}`)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetUserActivity(context.Background(), types.ActivityQuery{
		User: v2User, EventID: []int64{12, 34}, Cursor: v2Ptr("activity+/="), Type: []types.ActivityType{types.ActivityTIP}, Start: v2Ptr(int64(1)),
	})
	require.NoError(t, err)
	require.Nil(t, page.Pagination.NextCursor)
	require.Equal(t, "IN", page.Data[0].Side)
	require.Equal(t, 1.25, page.Data[0].Size)
	require.Nil(t, page.Data[0].IsCombo)
}

func TestV2ActivityDefaultsAndValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, url.Values{"user": {v2User}}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"pagination":{"limit":100,"offset":0,"has_more":false,"next_cursor":null}}`)
	}))
	defer server.Close()
	_, err := dataapi.NewClient(server.URL, nil).GetUserActivity(context.Background(), types.ActivityQuery{User: v2User})
	require.NoError(t, err)
	for name, q := range map[string]types.ActivityQuery{
		"missing user":        {},
		"condition and event": {User: v2User, Condition: []string{v2Condition}, EventID: []int64{1}},
		"invalid condition":   {User: v2User, Condition: []string{"bad"}},
		"limit past cap":      {User: v2User, Limit: v2Ptr(1001)},
		"CASH sort":           {User: v2User, SortBy: v2Ptr(types.ActivitySortBy("CASH"))},
		"TOKENS sort":         {User: v2User, SortBy: v2Ptr(types.ActivitySortBy("TOKENS"))},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := dataapi.NewClient("http://127.0.0.1:1", nil).GetUserActivity(context.Background(), q)
			require.Error(t, err)
		})
	}
}
