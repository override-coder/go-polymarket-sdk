package clob

import (
	"math/big"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"github.com/polymarket/go-order-utils/pkg/model"
)

func TestLimitBuyRoundingDoesNotExceedNotional(t *testing.T) {
	tests := []struct {
		name                 string
		price, size          float64
		tick                 types.TickSize
		wantMaker, wantTaker float64
	}{
		{"reported midpoint keeps shares", 0.295, 33.89, types.TickSize001, 9.8281, 33.89},
		{"above midpoint rounds down", 0.296, 33.89, types.TickSize001, 9.8281, 33.89},
		{"round down keeps shares", 0.294, 33.89, types.TickSize001, 9.8281, 33.89},
		{"aligned price keeps shares", 0.30, 33.89, types.TickSize001, 10.167, 33.89},
		{"fractional shares", 0.295, 33.899, types.TickSize001, 9.8281, 33.89},
		{"low price", 0.015, 100, types.TickSize001, 1, 100},
		{"low price preserves requested fill target", 0.085, 100, types.TickSize001, 8, 100},
		{"minimum five shares preserved", 0.295, 5, types.TickSize001, 1.45, 5},
		{"one decimal", 0.25, 10, types.TickSize01, 2, 10},
		{"three decimals", 0.2955, 10, types.TickSize0001, 2.95, 10},
		{"four decimals", 0.29555, 10, types.TickSize00001, 2.955, 10},
		{"zero shares", 0.295, 0, types.TickSize001, 0, 0},
		{"small positive shares preserved", 0.295, 0.01, types.TickSize001, 0.0029, 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			side, maker, taker := getOrderRawAmounts(types.BUY, tt.size, tt.price, roundConfigForTickSize(tt.tick))
			if side != model.BUY || maker != tt.wantMaker || taker != tt.wantTaker {
				t.Fatalf("got side=%v maker=%v taker=%v, want BUY maker=%v taker=%v", side, maker, taker, tt.wantMaker, tt.wantTaker)
			}
			budget := utils.Float64ToDecimal(tt.size).Mul(utils.Float64ToDecimal(tt.price))
			if utils.Float64ToDecimal(maker).GreaterThan(budget) {
				t.Fatalf("signed spend %v exceeds original notional %s", maker, budget)
			}
		})
	}
}

func TestLimitBuyRoundingSignedAmounts(t *testing.T) {
	for _, orderType := range []types.OrderType{types.OrderTypeGTC, types.OrderTypeGTD} {
		t.Run(string(orderType), func(t *testing.T) {
			signCalls := 0
			builder := NewOrderBuilder(big.NewInt(137), func(_ string, _ []byte) ([]byte, error) {
				signCalls++
				return make([]byte, 65), nil
			})
			opts := types.CreateOrderOptions{
				TickSize:   types.TickSize001,
				AuthOption: &sdktypes.AuthOption{SingerAddress: "0x0000000000000000000000000000000000000001"},
			}
			expiration := int64(2000000000)
			v1, err := builder.buildOrder(types.UserOrder{
				TokenID: "1", Side: types.BUY, Price: 0.295, Size: 33.89, Expiration: &expiration,
			}, orderType, opts)
			if err != nil {
				t.Fatal(err)
			}
			v2, err := builder.buildOrderV2(types.UserOrderV2{
				TokenId: "1", Side: types.BUY, Price: 0.295, Size: 33.89, Expiration: &expiration,
			}, orderType, opts)
			if err != nil {
				t.Fatal(err)
			}
			for version, amounts := range map[string][2]*big.Int{
				"v1": {v1.MakerAmount, v1.TakerAmount},
				"v2": {v2.MakerAmount, v2.TakerAmount},
			} {
				if amounts[0].String() != "9828100" || amounts[1].String() != "33890000" {
					t.Errorf("%s signed maker=%s taker=%s, want 9828100 / 33890000", version, amounts[0], amounts[1])
				}
			}
			if signCalls != 2 {
				t.Fatalf("sign calls=%d, want 2", signCalls)
			}
		})
	}
}

func TestLimitSellRoundingKeepsShares(t *testing.T) {
	for _, tc := range []struct {
		name              string
		price, wantAmount float64
	}{
		{"below midpoint", 0.294, 9.8281},
		{"at midpoint", 0.295, 10.167},
		{"above midpoint", 0.296, 10.167},
		{"aligned price", 0.30, 10.167},
	} {
		t.Run(tc.name, func(t *testing.T) {
			side, maker, taker := getOrderRawAmounts(types.SELL, 33.89, tc.price, roundConfigForTickSize(types.TickSize001))
			if side != model.SELL || maker != 33.89 || taker != tc.wantAmount {
				t.Fatalf("SELL changed: side=%v maker=%v taker=%v, want SELL maker=33.89 taker=%v", side, maker, taker, tc.wantAmount)
			}
		})
	}
}

func TestLimitBudgetFixLeavesMarketAmountsUnchanged(t *testing.T) {
	_, maker, taker := getMarketOrderRawAmounts(types.BUY, 9.9999, 0.295, roundConfigForTickSize(types.TickSize001))
	if maker != 9.99 || taker != 34.4482 {
		t.Fatalf("market BUY changed: maker=%v taker=%v", maker, taker)
	}
}
