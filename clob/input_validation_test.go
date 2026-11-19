package clob

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"net/http"
	"testing"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
)

type forbiddenTransport struct{ t *testing.T }

func (f forbiddenTransport) RoundTrip(*http.Request) (*http.Response, error) {
	f.t.Error("invalid input reached HTTP")
	return nil, fmt.Errorf("unexpected HTTP")
}

func TestAuthenticatedEndpointsRejectNilOption(t *testing.T) {
	c := NewClient("http://invalid", big.NewInt(137), nil, nil, sdkhttp.WithTransport(forbiddenTransport{t}))
	ctx := context.Background()
	calls := map[string]func() error{
		"ensure":         func() error { _, e := c.EnsureAPIKey(ctx, nil, nil); return e },
		"create key":     func() error { _, e := c.CreateApiKey(ctx, nil, nil); return e },
		"derive key":     func() error { _, e := c.DeriveAPIKey(ctx, nil, nil); return e },
		"builder key":    func() error { _, e := c.CreateBuilderApiKey(ctx, nil); return e },
		"balance":        func() error { _, e := c.GetBalanceAllowance(nil); return e },
		"update balance": func() error { _, e := c.UpdateBalanceAllowance(nil, nil); return e },
		"order":          func() error { _, e := c.GetOrder(ctx, "id", ct.GetOrderRequest{}, nil); return e },
		"orders":         func() error { _, e := c.GetOrders(ctx, ct.GetActiveOrdersRequest{}, nil); return e },
		"trades":         func() error { _, e := c.GetTrades(ctx, ct.GetTradesRequest{}, nil); return e },
		"cancel":         func() error { _, e := c.CancelOrder(ctx, "id", nil); return e },
		"cancel batch":   func() error { _, e := c.CancelOrders(ctx, []string{"id"}, nil); return e },
		"cancel all":     func() error { _, e := c.CancelOrderAll(ctx, nil); return e },
		"cancel market":  func() error { _, e := c.CancelOrderByMarket(ctx, nil, nil); return e },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if call() == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestCreateOrdersRejectInvalidInputBeforeIO(t *testing.T) {
	calls := 0
	sign := func(string, []byte) ([]byte, error) { calls++; return make([]byte, 65), nil }
	c := NewClient("http://invalid", big.NewInt(137), sign, nil, sdkhttp.WithTransport(forbiddenTransport{t}))
	auth := &sdk.AuthOption{ApiKeyCreds: &sdk.ApiKeyCreds{ApiKey: "key"}}
	for _, tc := range []struct {
		name   string
		change func(*ct.UserOrderV2)
	}{
		{"nan price", func(o *ct.UserOrderV2) { o.Price = math.NaN() }},
		{"max finite size", func(o *ct.UserOrderV2) { o.Size = math.MaxFloat64 }},
		{"infinite size", func(o *ct.UserOrderV2) { o.Size = math.Inf(1) }},
		{"negative size", func(o *ct.UserOrderV2) { o.Size = -1 }},
		{"zero size", func(o *ct.UserOrderV2) { o.Size = 0 }},
		{"invalid side", func(o *ct.UserOrderV2) { o.Side = "sell" }},
		{"invalid asset", func(o *ct.UserOrderV2) { o.TokenId = "not-a-number" }},
		{"unsupported tick", func(o *ct.UserOrderV2) { v := "0.02"; o.TickSize = &v }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := ct.UserOrderV2{TokenId: "1", Price: .5, Size: 1, Side: ct.BUY}
			tc.change(&o)
			if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderTypeFOK, false, false, auth); e == nil {
				t.Fatal("V2 expected error")
			}
			v1 := ct.UserOrder{TokenID: o.TokenId, Price: o.Price, Size: o.Size, Side: o.Side, TickSize: o.TickSize}
			if _, e := c.CreateOrder(context.Background(), v1, ct.OrderTypeFOK, false, auth); e == nil {
				t.Fatal("legacy expected error")
			}
		})
	}
	o := ct.UserOrderV2{TokenId: "1", Price: .5, Size: 1, Side: ct.BUY}
	if _, e := c.CreateOrderV2(context.Background(), nil, ct.OrderTypeFOK, false, false, auth); e == nil {
		t.Fatal("nil order")
	}
	if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderType("invalid"), false, false, auth); e == nil {
		t.Fatal("invalid type")
	}
	if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderTypeFOK, true, false, auth); e == nil {
		t.Fatal("postOnly market order")
	}
	for _, v := range []float64{-1, math.NaN(), math.Inf(1)} {
		o.UserUsdcBalance = &v
		if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderTypeFOK, false, false, auth); e == nil {
			t.Fatal("invalid balance")
		}
		o.UserUsdcBalance = nil
		o.BuilderFeeRate = &v
		if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderTypeFOK, false, false, auth); e == nil {
			t.Fatal("invalid fee")
		}
		o.BuilderFeeRate = nil
	}
	c.signFn = nil
	if _, e := c.CreateOrderV2(context.Background(), &o, ct.OrderTypeFOK, false, false, auth); e == nil {
		t.Fatal("missing signer")
	}
	if _, e := c.DeriveAPIKey(context.Background(), nil, auth); e == nil {
		t.Fatal("missing L1 signer")
	}
	if _, e := c.CancelOrderByMarket(context.Background(), nil, auth); e == nil {
		t.Fatal("nil cancel request")
	}
	if calls != 0 {
		t.Fatalf("invalid input invoked signer %d times", calls)
	}
}

func TestRoundedOrderAmountBounds(t *testing.T) {
	// Input fits uint256 collateral, but BUY shares at 0.5 do not.
	size := maxOrderUnits.Shift(-6).Mul(decimal.RequireFromString("0.75")).InexactFloat64()
	if err := validateOrderInput("1", .5, size, ct.BUY, ct.OrderTypeFOK, nil); err != nil {
		t.Fatal(err)
	}
	if err := validateOrderAmounts(.5, size, ct.BUY, ct.OrderTypeFOK, "0.01"); err == nil {
		t.Fatal("overflowing taker amount accepted")
	}
	if err := validateOrderAmounts(.5, .001, ct.SELL, ct.OrderTypeFOK, "0.01"); err == nil {
		t.Fatal("zero rounded shares accepted")
	}
	if err := validateOrderAmounts(.5, 1.2, ct.BUY, ct.OrderTypeFOK, "0.01"); err != nil {
		t.Fatal(err)
	}
}
