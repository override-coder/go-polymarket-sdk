package clob

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestMarketBuyBudgetBelowOne(t *testing.T) {
	for _, tc := range []struct{ balance, rate, builder float64 }{
		{0.5, 0, 0}, {1, 0.05, 0.01}, {0, 0, 0}, {2, 0.05, 0.01},
	} {
		for _, orderType := range []ct.OrderType{ct.OrderTypeFOK, ct.OrderTypeFAK} {
			code := "builder"
			order := ct.UserOrderV2{Side: ct.BUY, Price: 0.5, Size: 1.2, UserUsdcBalance: &tc.balance, BuilderCode: &code, BuilderFeeRate: &tc.builder}
			client := &Client{}
			client.adjustBuyMarketOrderAmountForFees(&order, &ct.CLOBMarketInfoFeeDetails{Rate: tc.rate, Exponent: 1}, orderType)
			cost := order.Size * (1 + tc.rate*math.Pow(0.25, 1)/0.5 + tc.builder)
			if cost > tc.balance+1e-12 || order.Size > 1.2 {
				t.Errorf("budget ignored: order=%+v cost=%g", order, cost)
			}
		}
	}
}

func TestCreateOrderV2HonorsSmallMarketBuyBudget(t *testing.T) {
	for _, orderType := range []ct.OrderType{ct.OrderTypeFOK, ct.OrderTypeFAK} {
		var posted ct.NewOrderV2
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order", r.URL.Path)
			require.NoError(t, json.NewDecoder(r.Body).Decode(&posted))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"success":true}`))
		}))
		client := NewClient(server.URL, big.NewInt(137), func(string, []byte) ([]byte, error) {
			return make([]byte, 65), nil
		}, nil)
		condition := "condition"
		asset := new(big.Int).Lsh(big.NewInt(1), 248).String()
		client.storeMarketInfo(condition, &ct.CLOBMarketInfo{
			MinimumTickSize: 0.01, Tokens: []ct.CLOBMarketInfoToken{{TokenID: asset}},
		})
		balance := 0.5
		order := ct.UserOrderV2{TokenId: asset, ConditionId: &condition, Side: ct.BUY, Price: 0.5, Size: 1.2, UserUsdcBalance: &balance}
		auth := testAuthOption()
		auth.SingerAddress = "0x0000000000000000000000000000000000000001"
		auth.SignatureType = model.EOA
		response, err := client.CreateOrderV2(context.Background(), &order, orderType, false, true, auth)
		require.NoError(t, err)
		require.True(t, response.Success)
		require.Equal(t, "500000", posted.Order.MakerAmount)
		require.Equal(t, "1000000", posted.Order.TakerAmount)
		server.Close()
	}
}
