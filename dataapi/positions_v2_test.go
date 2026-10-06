package dataapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/require"
)

const v2User = "0x1111111111111111111111111111111111111111"

var v2Condition = "0x" + strings.Repeat("a", 64)

func v2Ptr[T any](v T) *T { return &v }

// Fixtures enumerate the official V2 OpenAPI response fields (2026-10-06).
// Comparing a decoded-and-encoded row detects silently dropped/renamed fields.
func v2Fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name + "_v2.json")
	require.NoError(t, err)
	return string(b)
}
func assertV2JSON(t *testing.T, expected string, value any) {
	t.Helper()
	actual, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, expected, string(actual))
}

func TestV2PositionsAllParametersAndFields(t *testing.T) {
	fixture := v2Fixture(t, "position")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v2/positions", r.URL.Path)
		require.Equal(t, url.Values{
			"user": {v2User}, "condition": {v2Condition}, "event_id": {"12,34"},
			"limit": {"1000"}, "status": {"REDEEMABLE"}, "title": {"测试%"},
			"filter_type": {"CASH"}, "filter_amount": {"1.5"}, "include_archived": {"true"},
			"sort_by": {"TOTAL_PNL"}, "sort_direction": {"ASC"}, "start": {"0"}, "end": {"1700000000"},
		}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"pagination":{"limit":1000,"offset":0,"has_more":true,"next_cursor":"next+/="}}`, fixture)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetPositions(context.Background(), types.PositionsQuery{
		User: " " + v2User + " ", Condition: []string{v2Condition}, EventID: []int64{12, 34},
		Limit: v2Ptr(1000), Status: v2Ptr(types.PositionREDEEMABLE), Title: v2Ptr("测试%"),
		FilterType: v2Ptr(types.FilterCASH), FilterAmount: v2Ptr(1.5), IncludeArchived: v2Ptr(true),
		SortBy: v2Ptr(types.SortByTOTALPNL), SortDirection: v2Ptr(types.SortASC), Start: v2Ptr(int64(0)), End: v2Ptr(int64(1700000000)),
	})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assertV2JSON(t, fixture, page.Data[0])
	require.True(t, page.Pagination.HasMore) // A short page can still have another page.
	require.Equal(t, "next+/=", *page.Pagination.NextCursor)
}

func TestV2PositionsCursorPreservesAnchorWithoutDefaults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, url.Values{"condition": {v2Condition}, "cursor": {"next+/="}, "title": {"retained"}}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"pagination":{"limit":10,"offset":10,"has_more":false,"next_cursor":null}}`)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetPositions(context.Background(), types.PositionsQuery{
		Condition: []string{v2Condition}, Cursor: v2Ptr("next+/="), Limit: v2Ptr(1001), Title: v2Ptr("retained"),
	})
	require.NoError(t, err)
	require.Empty(t, page.Data)
	require.Nil(t, page.Pagination.NextCursor)
	require.Equal(t, 10, page.Pagination.Offset)
}

func TestV2PositionsValidation(t *testing.T) {
	many := make([]string, 21)
	for i := range many {
		many[i] = fmt.Sprintf("0x%064x", i)
	}
	for name, q := range map[string]types.PositionsQuery{
		"missing anchor":          {},
		"invalid condition":       {User: v2User, Condition: []string{"bad"}},
		"too many conditions":     {User: v2User, Condition: many},
		"multiple market anchors": {Condition: many[:2]},
		"event without user":      {Condition: []string{v2Condition}, EventID: []int64{1}},
		"large limit":             {User: v2User, Limit: v2Ptr(1001)},
		"long title":              {User: v2User, Title: v2Ptr(strings.Repeat("字", 201))},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := dataapi.NewClient("http://127.0.0.1:1", nil).GetPositions(context.Background(), q)
			require.Error(t, err)
		})
	}
}
