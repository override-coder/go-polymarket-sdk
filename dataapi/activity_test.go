package dataapi_test

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/require"
)

func TestGetUserActivityAllParameters(t *testing.T) {
	const user = "0x1111111111111111111111111111111111111111"
	market1 := "0x" + strings.Repeat("a", 64)
	market2 := "0x" + strings.Repeat("b", 64)
	limit := 501
	offset := 5000
	excludeDepositsWithdrawals := false
	start := int64(1)
	end := int64(1700000000)
	sortBy := types.ActivitySortCASH
	sortDirection := types.SortASC
	side := types.SideSELL

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/activity", r.URL.Path)
		require.Equal(t, user, r.URL.Query().Get("user"))
		require.Equal(t, "500", r.URL.Query().Get("limit"))
		require.Equal(t, "5000", r.URL.Query().Get("offset"))
		require.Equal(t, market1+","+market2, r.URL.Query().Get("market"))
		require.Equal(t, "DEPOSIT,WITHDRAWAL,YIELD,MAKER_REBATE,TAKER_REBATE,REFERRAL_REWARD", r.URL.Query().Get("type"))
		require.Equal(t, "false", r.URL.Query().Get("excludeDepositsWithdrawals"))
		require.Equal(t, "1", r.URL.Query().Get("start"))
		require.Equal(t, "1700000000", r.URL.Query().Get("end"))
		require.Equal(t, "CASH", r.URL.Query().Get("sortBy"))
		require.Equal(t, "ASC", r.URL.Query().Get("sortDirection"))
		require.Equal(t, "SELL", r.URL.Query().Get("side"))

		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprint(w, `[{
			"proxyWallet":"`+user+`",
			"timestamp":1700000000,
			"conditionId":"`+market1+`",
			"type":"TRADE",
			"size":12.5,
			"usdcSize":5,
			"transactionHash":"0xabc",
			"price":0.4,
			"asset":"123",
			"side":"SELL",
			"outcomeIndex":1,
			"title":"Election",
			"slug":"election",
			"icon":"https://example.com/icon.png",
			"eventSlug":"election-event",
			"outcome":"Yes",
			"name":"Trader",
			"pseudonym":"trader",
			"bio":"bio",
			"profileImage":"https://example.com/profile.png",
			"profileImageOptimized":"https://example.com/profile-optimized.png",
			"isCombo":true
		}]`)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := dataapi.NewClient(server.URL, big.NewInt(137))
	activities, err := client.GetUserActivity(context.Background(), types.ActivityQuery{
		User:                       "  " + user + "  ",
		Limit:                      &limit,
		Offset:                     &offset,
		Market:                     []string{market1, market2},
		Type:                       []types.ActivityType{types.ActivityDEPOSIT, types.ActivityWITHDRAWAL, types.ActivityYIELD, types.ActivityMAKERREBATE, types.ActivityTAKERREBATE, types.ActivityREFERRALREWARD},
		ExcludeDepositsWithdrawals: &excludeDepositsWithdrawals,
		Start:                      &start,
		End:                        &end,
		SortBy:                     &sortBy,
		SortDirection:              &sortDirection,
		Side:                       &side,
	})

	require.NoError(t, err)
	require.Len(t, activities, 1)
	require.Equal(t, types.ActivityTRADE, activities[0].Type)
	require.Equal(t, types.SideSELL, activities[0].Side)
	require.True(t, activities[0].IsCombo)
}

func TestGetUserActivityDefaultsAndEventIDs(t *testing.T) {
	const user = "0x1111111111111111111111111111111111111111"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "12,34", r.URL.Query().Get("eventId"))
		require.Equal(t, "100", r.URL.Query().Get("limit"))
		require.Equal(t, "0", r.URL.Query().Get("offset"))
		require.Equal(t, "true", r.URL.Query().Get("excludeDepositsWithdrawals"))
		require.Equal(t, "TIMESTAMP", r.URL.Query().Get("sortBy"))
		require.Equal(t, "DESC", r.URL.Query().Get("sortDirection"))
		_, err := fmt.Fprint(w, `[]`)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := dataapi.NewClient(server.URL, big.NewInt(137))
	activities, err := client.GetUserActivity(context.Background(), types.ActivityQuery{
		User:    user,
		EventID: []int64{12, 34},
	})

	require.NoError(t, err)
	require.Empty(t, activities)
}

func TestGetUserActivityValidation(t *testing.T) {
	validUser := "0x1111111111111111111111111111111111111111"
	validMarket := "0x" + strings.Repeat("a", 64)
	negative := -1
	offsetTooHigh := 5001
	negativeTimestamp := int64(-1)

	tests := []struct {
		name  string
		query types.ActivityQuery
	}{
		{name: "missing user", query: types.ActivityQuery{}},
		{name: "invalid user", query: types.ActivityQuery{User: "0x1234"}},
		{name: "market and event", query: types.ActivityQuery{User: validUser, Market: []string{validMarket}, EventID: []int64{1}}},
		{name: "invalid market", query: types.ActivityQuery{User: validUser, Market: []string{"0x1234"}}},
		{name: "invalid event", query: types.ActivityQuery{User: validUser, EventID: []int64{0}}},
		{name: "negative limit", query: types.ActivityQuery{User: validUser, Limit: &negative}},
		{name: "negative offset", query: types.ActivityQuery{User: validUser, Offset: &negative}},
		{name: "offset too high", query: types.ActivityQuery{User: validUser, Offset: &offsetTooHigh}},
		{name: "negative start", query: types.ActivityQuery{User: validUser, Start: &negativeTimestamp}},
		{name: "negative end", query: types.ActivityQuery{User: validUser, End: &negativeTimestamp}},
	}

	client := dataapi.NewClient("http://127.0.0.1:1", big.NewInt(137))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetUserActivity(context.Background(), tt.query)
			require.Error(t, err)
		})
	}
}
