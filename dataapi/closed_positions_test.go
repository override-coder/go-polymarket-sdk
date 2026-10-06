package dataapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/require"
)

func TestV2ClosedPositionsAllParametersAndFields(t *testing.T) {
	fixture := strings.Replace(v2Fixture(t, "position"), `"status": "REDEEMABLE"`, `"status": "CLOSED"`, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v2/positions", r.URL.Path)
		require.Equal(t, url.Values{
			"user": {v2User}, "condition": {v2Condition}, "event_id": {"12,34"}, "status": {"CLOSED"},
			"limit": {"1000"}, "title": {"Election"}, "filter_type": {"TOKENS"}, "filter_amount": {"0"},
			"include_archived": {"false"}, "sort_by": {"TIMESTAMP"}, "sort_direction": {"ASC"}, "start": {"1"}, "end": {"1700000000"},
		}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"pagination":{"limit":1000,"offset":0,"has_more":false,"next_cursor":null}}`, fixture)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetClosedPositions(context.Background(), types.ClosedPositionsQuery{
		User: v2User, Condition: []string{v2Condition}, EventID: []int64{12, 34}, Limit: v2Ptr(1000), Title: v2Ptr("Election"),
		FilterType: v2Ptr(types.FilterTOKENS), FilterAmount: v2Ptr(0.0), IncludeArchived: v2Ptr(false),
		SortBy: v2Ptr(types.SortByTIMESTAMP), SortDirection: v2Ptr(types.SortASC), Start: v2Ptr(int64(1)), End: v2Ptr(int64(1700000000)),
	})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assertV2JSON(t, fixture, page.Data[0])
	require.Nil(t, page.Pagination.NextCursor)
}

func TestV2ClosedPositionsDefaultsAndCursor(t *testing.T) {
	for _, cursor := range []string{"", "closed+/="} {
		t.Run(cursor, func(t *testing.T) {
			expected := url.Values{"user": {v2User}, "status": {"CLOSED"}}
			var cursorParam *string
			if cursor != "" {
				expected.Set("cursor", cursor)
				cursorParam = &cursor
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, expected, r.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[],"pagination":{"limit":100,"offset":0,"has_more":false,"next_cursor":null}}`)
			}))
			defer server.Close()
			_, err := dataapi.NewClient(server.URL, nil).GetClosedPositions(context.Background(), types.ClosedPositionsQuery{User: v2User, Cursor: cursorParam})
			require.NoError(t, err)
		})
	}
}

func TestV2ClosedPositionsRejectsConflictingStatus(t *testing.T) {
	_, err := dataapi.NewClient("http://127.0.0.1:1", nil).GetClosedPositions(context.Background(), types.ClosedPositionsQuery{User: v2User, Status: v2Ptr(types.PositionOPEN)})
	require.ErrorContains(t, err, "status=CLOSED")
}
