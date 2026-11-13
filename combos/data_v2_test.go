package combos_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/combos"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dataWallet = "0x1111111111111111111111111111111111111111"
const dataCondition = "0x03a4918ddf5a3fc35d65ef396fb54ce8f90000000000000000000000000000"

func dataPtr[T any](v T) *T { return &v }

type comboDataReader interface {
	GetComboPositions(context.Context, types.ComboPositionsQuery) (*types.Page[types.ComboPosition], error)
	GetComboActivity(context.Context, types.ComboActivityQuery) (*types.Page[types.ComboActivity], error)
}

var dataReaders = map[string]func(string, string) comboDataReader{
	"dataapi": func(_ string, dataHost string) comboDataReader { return dataapi.NewClient(dataHost, nil) },
	"combos": func(rfqHost, dataHost string) comboDataReader {
		return combos.NewClientWithDataHost(rfqHost, dataHost, nil, nil)
	},
}

func TestComboDataV2ActivityContract(t *testing.T) {
	for name, factory := range dataReaders {
		t.Run(name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v2/activity/combos", r.URL.Path)
				assert.Empty(t, r.Header.Get("POLY_SIGNATURE"))
				want := url.Values{"user": {dataWallet}, "condition": {dataCondition}}
				w.Header().Set("Content-Type", "application/json")
				if calls == 1 {
					want.Set("limit", "2")
					// An empty page with a cursor must still be exposed as resumable.
					fmt.Fprint(w, `{"data":[],"pagination":{"limit":2,"offset":0,"has_more":true,"next_cursor":"next+/="}}`)
				} else {
					want.Set("cursor", "next+/=")
					fmt.Fprint(w, `{"data":[{"id":"event-1","block_number":123456,"combo_condition_id":"condition","combo_position_id":"123456789012345678901234567890","proxy_wallet":"wallet","timestamp":1792368000,"transaction_hash":"0xtx","type":"REDEEM","amount_usdc":0.123456789012345678,"payout_usdc":0,"legs":[{"leg_index":1,"leg_position_id":"234567890123456789012345678901","leg_condition_id":"leg-condition","leg_outcome_index":0,"leg_outcome_label":"YES","leg_status":"RESOLVED_WIN","leg_current_price":1,"leg_resolved_at":null,"market":{"market_id":"42","title":"Final","line":null,"event":{"event_id":"7","event_title":"Tournament"}}}]},{"id":"event-2","amount_usdc":null}],"pagination":{"limit":2,"offset":2,"has_more":false,"next_cursor":null}}`)
				}
				assert.Equal(t, want, r.URL.Query(), "must not leak V1 market_id or offset")
			}))
			defer server.Close()
			client := factory("http://127.0.0.1:1", server.URL+"/")
			q := combostypes.ComboActivityQuery{User: " " + dataWallet + " ", Condition: []string{dataCondition}, Limit: dataPtr(2)}
			page, err := client.GetComboActivity(context.Background(), q)
			require.NoError(t, err)
			require.Empty(t, page.Data)
			require.True(t, page.Pagination.HasMore)
			require.NotNil(t, page.Pagination.NextCursor)
			q.Cursor, q.Limit = page.Pagination.NextCursor, dataPtr(-1)
			page, err = client.GetComboActivity(context.Background(), q)
			require.NoError(t, err)
			require.Len(t, page.Data, 2)
			item := page.Data[0]
			assert.Equal(t, "event-1", item.ID)
			assert.EqualValues(t, 123456, item.BlockNumber)
			assert.Equal(t, "condition", item.ComboConditionID)
			assert.Equal(t, "123456789012345678901234567890", item.ComboPositionID)
			assert.Equal(t, "wallet", item.ProxyWallet)
			assert.EqualValues(t, 1792368000, item.Timestamp)
			assert.Equal(t, "0xtx", item.TransactionHash)
			assert.Equal(t, "REDEEM", item.Type)
			require.NotNil(t, item.AmountUSDC)
			assert.Equal(t, "0.123456789012345678", item.AmountUSDC.String())
			require.NotNil(t, item.PayoutUSDC)
			assert.Equal(t, "0", item.PayoutUSDC.String())
			require.Len(t, item.Legs, 1)
			assert.Equal(t, "234567890123456789012345678901", item.Legs[0].LegPositionID)
			assert.Equal(t, "Tournament", item.Legs[0].Market.Event.EventTitle)
			assert.Nil(t, item.Legs[0].LegResolvedAt)
			assert.Nil(t, item.Legs[0].Market.Line)
			assert.Nil(t, page.Data[1].AmountUSDC)
			assert.Nil(t, page.Data[1].PayoutUSDC)
			assert.Nil(t, page.Pagination.NextCursor)
			assert.False(t, page.Pagination.HasMore)
			assert.Equal(t, 2, calls)
		})
	}
}

func TestComboDataV2PositionsDelegation(t *testing.T) {
	for name, factory := range dataReaders {
		t.Run(name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, "/v2/positions/combos", r.URL.Path)
				want := url.Values{"user": {dataWallet}, "condition": {dataCondition}, "status": {"OPEN"}, "sort_by": {"ENTRY_COST"}, "sort_direction": {"DESC"}, "updated_after": {"0"}, "updated_before": {"1792368000"}}
				if calls == 1 {
					want.Set("limit", "1000")
				} else {
					want.Set("cursor", "positions-next")
				}
				assert.Equal(t, want, r.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[{"proxy_wallet":"wallet","current_size":123456789012345678.123456,"gross_entry_cost_usdc":0.100001,"status":"OPEN","legs":[{"leg_position_id":"999999999999999999999999","market":{"event":{"event_title":"Final"}}}]}],"pagination":{"limit":1000,"has_more":true,"next_cursor":"positions-next"}}`)
			}))
			defer server.Close()
			client := factory("http://127.0.0.1:1", server.URL)
			q := combostypes.ComboPositionsQuery{User: dataWallet, Condition: []string{dataCondition}, Limit: dataPtr(1000), Status: dataPtr("OPEN"), SortBy: dataPtr("ENTRY_COST"), SortDirection: dataPtr(types.SortDESC), UpdatedAfter: dataPtr(int64(0)), UpdatedBefore: dataPtr(int64(1792368000))}
			page, err := client.GetComboPositions(context.Background(), q)
			require.NoError(t, err)
			require.Len(t, page.Data, 1)
			assert.Equal(t, "wallet", page.Data[0].ProxyWallet)
			assert.Equal(t, "123456789012345678.123456", page.Data[0].CurrentSize.String())
			assert.Equal(t, "0.100001", page.Data[0].GrossEntryCostUSDC.String())
			require.Len(t, page.Data[0].Legs, 1)
			assert.Equal(t, "Final", page.Data[0].Legs[0].Market.Event.EventTitle)
			q.Cursor, q.Limit = page.Pagination.NextCursor, dataPtr(-1)
			_, err = client.GetComboPositions(context.Background(), q)
			require.NoError(t, err)
			assert.Equal(t, 2, calls)
		})
	}
}

func TestComboDataV2ValidationAndErrors(t *testing.T) {
	for name, factory := range dataReaders {
		t.Run(name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":"invalid cursor"}`, http.StatusBadRequest)
			}))
			defer server.Close()
			client := factory("http://127.0.0.1:1", server.URL)
			ids := make([]string, 21)
			for i := range ids {
				ids[i] = fmt.Sprintf("0x%062x", i)
			}
			for _, q := range []types.ComboActivityQuery{
				{}, {User: " "}, {User: dataWallet, Limit: dataPtr(-1)}, {User: dataWallet, Limit: dataPtr(1001)},
				{User: dataWallet, Condition: []string{dataCondition + "00"}},
				{User: dataWallet, Condition: []string{"0x" + strings.Repeat("g", 62)}},
				{User: dataWallet, Condition: ids},
			} {
				_, err := client.GetComboActivity(context.Background(), q)
				require.Error(t, err)
			}
			_, err := client.GetComboPositions(context.Background(), types.ComboPositionsQuery{User: dataWallet, SortBy: dataPtr("PRICE")})
			require.Error(t, err)
			assert.Zero(t, calls, "invalid queries must not reach HTTP")
			for _, limit := range []int{0, 1000} {
				_, err = client.GetComboActivity(context.Background(), types.ComboActivityQuery{User: dataWallet, Condition: append(ids[:20:20], ids[0]), Limit: dataPtr(limit)})
				require.ErrorContains(t, err, "400", "valid query must propagate upstream error")
			}
			_, err = client.GetComboPositions(context.Background(), types.ComboPositionsQuery{User: dataWallet})
			require.ErrorContains(t, err, "400")
			assert.Equal(t, 3, calls)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err = client.GetComboActivity(ctx, types.ComboActivityQuery{User: dataWallet})
			require.ErrorIs(t, err, context.Canceled)
		})
	}
}
