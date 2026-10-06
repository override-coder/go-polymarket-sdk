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

func TestV2LeaderboardAllParametersAndFields(t *testing.T) {
	fixture := v2Fixture(t, "leaderboard")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v2/leaderboard", r.URL.Path)
		require.Equal(t, url.Values{"category": {"sports"}, "time_period": {"all"}, "sort_by": {"VOLUME"}, "limit": {"1000"}}, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"pagination":{"limit":1000,"offset":0,"has_more":true,"next_cursor":"board+/="}}`, fixture)
	}))
	defer server.Close()
	page, err := dataapi.NewClient(server.URL, nil).GetTraderLeaderboardRankings(context.Background(), types.TraderLeaderboardQuery{
		Category: v2Ptr(types.LeaderboardCategorySPORTS), TimePeriod: v2Ptr(types.LeaderboardTimeALL), SortBy: v2Ptr(types.LeaderboardSortByVOLUME), Limit: v2Ptr(1000),
	})
	require.NoError(t, err)
	require.Len(t, page.Data, 1)
	assertV2JSON(t, fixture, page.Data[0])
	require.Equal(t, "board+/=", *page.Pagination.NextCursor)
}

func TestV2LeaderboardCursorAndDefaults(t *testing.T) {
	for _, cursor := range []string{"", "board+/="} {
		t.Run(cursor, func(t *testing.T) {
			expected := url.Values{}
			q := types.TraderLeaderboardQuery{}
			if cursor != "" {
				expected.Set("cursor", cursor)
				q.Cursor = &cursor
				q.Limit = v2Ptr(1001)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, expected, r.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[],"pagination":{"limit":100,"offset":100,"has_more":false,"next_cursor":null}}`)
			}))
			defer server.Close()
			page, err := dataapi.NewClient(server.URL, nil).GetTraderLeaderboardRankings(context.Background(), q)
			require.NoError(t, err)
			require.Nil(t, page.Pagination.NextCursor)
		})
	}
}

func TestV2LeaderboardUserFieldsAndNull(t *testing.T) {
	for _, fixture := range []string{v2Fixture(t, "leaderboard_user"), "null"} {
		t.Run(fixture, func(t *testing.T) {
			body := `{"data":` + fixture + `}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/v2/leaderboard", r.URL.Path)
				require.Equal(t, url.Values{"user": {v2User}, "category": {"combos"}, "time_period": {"week"}}, r.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, body)
			}))
			defer server.Close()
			result, err := dataapi.NewClient(server.URL, nil).GetTraderLeaderboardUser(context.Background(), types.TraderLeaderboardUserQuery{
				User: v2User, Category: v2Ptr(types.LeaderboardCategoryCOMBOS), TimePeriod: v2Ptr(types.LeaderboardTimeWEEK),
			})
			require.NoError(t, err)
			assertV2JSON(t, body, result)
			if fixture == "null" {
				require.Nil(t, result.Data)
			} else {
				require.Equal(t, 2, *result.Data.RankPnL)
				require.Nil(t, result.Data.RankVolume)
			}
		})
	}
}

func TestV2LeaderboardValidation(t *testing.T) {
	client := dataapi.NewClient("http://127.0.0.1:1", nil)
	_, err := client.GetTraderLeaderboardRankings(context.Background(), types.TraderLeaderboardQuery{Limit: v2Ptr(1001)})
	require.Error(t, err)
	_, err = client.GetTraderLeaderboardUser(context.Background(), types.TraderLeaderboardUserQuery{})
	require.Error(t, err)
}
