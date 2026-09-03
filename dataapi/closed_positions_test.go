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

func TestGetClosedPositions(t *testing.T) {
	const user = "0x1111111111111111111111111111111111111111"
	market1 := "0x" + strings.Repeat("a", 64)
	market2 := "0x" + strings.Repeat("b", 64)
	title := "Election"
	limit := 50
	offset := 100000
	sortBy := types.ClosedPositionSortByTIMESTAMP
	sortDirection := types.SortASC

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/closed-positions", r.URL.Path)
		require.Equal(t, user, r.URL.Query().Get("user"))
		require.Equal(t, market1+","+market2, r.URL.Query().Get("market"))
		require.Equal(t, title, r.URL.Query().Get("title"))
		require.Equal(t, "50", r.URL.Query().Get("limit"))
		require.Equal(t, "100000", r.URL.Query().Get("offset"))
		require.Equal(t, "TIMESTAMP", r.URL.Query().Get("sortBy"))
		require.Equal(t, "ASC", r.URL.Query().Get("sortDirection"))

		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprint(w, `[{
			"proxyWallet":"0x1111111111111111111111111111111111111111",
			"asset":"123",
			"conditionId":"`+market1+`",
			"avgPrice":0.4,
			"totalBought":10.5,
			"realizedPnl":3.25,
			"curPrice":1,
			"timestamp":1700000000,
			"title":"Election",
			"slug":"election",
			"icon":"https://example.com/icon.png",
			"eventSlug":"election-event",
			"outcome":"Yes",
			"outcomeIndex":1,
			"oppositeOutcome":"No",
			"oppositeAsset":"456",
			"endDate":"2026-11-03"
		}]`)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := dataapi.NewClient(server.URL, big.NewInt(137))
	positions, err := client.GetClosedPositions(context.Background(), types.ClosedPositionsQuery{
		User:          user,
		Market:        []string{market1, market2},
		Title:         &title,
		Limit:         &limit,
		Offset:        &offset,
		SortBy:        &sortBy,
		SortDirection: &sortDirection,
	})

	require.NoError(t, err)
	require.Len(t, positions, 1)
	require.Equal(t, int64(1700000000), positions[0].Timestamp)
	require.Equal(t, 3.25, positions[0].RealizedPnl)
	require.Equal(t, "456", positions[0].OppositeAsset)
}

func TestGetClosedPositionsDefaultsAndEventIDs(t *testing.T) {
	const user = "0x1111111111111111111111111111111111111111"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "12,34", r.URL.Query().Get("eventId"))
		require.Equal(t, "10", r.URL.Query().Get("limit"))
		require.Equal(t, "0", r.URL.Query().Get("offset"))
		require.Equal(t, "TIMESTAMP", r.URL.Query().Get("sortBy"))
		require.Equal(t, "DESC", r.URL.Query().Get("sortDirection"))
		_, err := fmt.Fprint(w, `[]`)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := dataapi.NewClient(server.URL, big.NewInt(137))
	positions, err := client.GetClosedPositions(context.Background(), types.ClosedPositionsQuery{
		User:    user,
		EventID: []int64{12, 34},
	})

	require.NoError(t, err)
	require.Empty(t, positions)
}

func TestGetClosedPositionsValidation(t *testing.T) {
	validUser := "0x1111111111111111111111111111111111111111"
	validMarket := "0x" + strings.Repeat("a", 64)
	negative := -1
	limitTooHigh := 51
	offsetTooHigh := 100001
	longTitle := strings.Repeat("x", 101)

	tests := []struct {
		name  string
		query types.ClosedPositionsQuery
	}{
		{name: "missing user", query: types.ClosedPositionsQuery{}},
		{name: "invalid user", query: types.ClosedPositionsQuery{User: "0x1234"}},
		{name: "market and event", query: types.ClosedPositionsQuery{User: validUser, Market: []string{validMarket}, EventID: []int64{1}}},
		{name: "invalid market", query: types.ClosedPositionsQuery{User: validUser, Market: []string{"0x1234"}}},
		{name: "invalid event", query: types.ClosedPositionsQuery{User: validUser, EventID: []int64{0}}},
		{name: "negative limit", query: types.ClosedPositionsQuery{User: validUser, Limit: &negative}},
		{name: "limit too high", query: types.ClosedPositionsQuery{User: validUser, Limit: &limitTooHigh}},
		{name: "negative offset", query: types.ClosedPositionsQuery{User: validUser, Offset: &negative}},
		{name: "offset too high", query: types.ClosedPositionsQuery{User: validUser, Offset: &offsetTooHigh}},
		{name: "title too long", query: types.ClosedPositionsQuery{User: validUser, Title: &longTitle}},
	}

	client := dataapi.NewClient("http://127.0.0.1:1", big.NewInt(137))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetClosedPositions(context.Background(), tt.query)
			require.Error(t, err)
		})
	}
}
