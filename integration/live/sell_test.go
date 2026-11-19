//go:build polymarket_live

package live_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// A sell is explicit and bounded by SHARES, never the wallet's whole balance.
func TestSellAndVerifyFill(t *testing.T) {
	enabled(t, "sell")
	record := beginResult(t)
	m := selected(t)
	require.True(t, m.Active && m.AcceptingOrders && !m.Closed)
	if m.Version == "v2" {
		require.Equal(t, "active", m.ResolutionStatus)
	}
	require.Empty(t, os.Getenv("POLY_LIVE_MAX_SPEND"), "clear BUY-only MAX_SPEND before selling")
	shares := positive(t, env(t, "SHARES"))
	require.True(t, shares.Equal(shares.Truncate(2)), "SHARES must have at most two decimals")
	price := positive(t, env(t, "MIN_PRICE"))
	a := wallet(t, sdkhttp.WithTransport(sellLimitTransport{m.asset, shares, price}))
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	rpc := budgetRPC(t, ctx)
	defer rpc.Close()
	contract := positionContract(m)
	holder := common.HexToAddress(a.auth.FunderAddress)
	before := erc1155Balance(t, ctx, rpc, contract, holder, m.number, nil)
	require.True(t, decimal.NewFromBigInt(before, -6).GreaterThanOrEqual(shares), "insufficient chain position for requested SHARES")
	record("position_before", liveResult{Asset: m.asset, Contract: contract.Hex(), Balance: before.String()})
	a.credentials(t, ctx)
	book, err := a.clob.GetOrderBook(m.asset)
	require.NoError(t, err)
	require.NotNil(t, book)
	require.Equal(t, m.asset, book.AssetID)
	tick := positive(t, book.TickSize)
	require.True(t, price.Mod(tick).IsZero() && price.GreaterThanOrEqual(tick) && price.LessThanOrEqual(decimal.NewFromInt(1).Sub(tick)), "MIN_PRICE must lie on the live price grid")
	depth := decimal.Zero
	for _, bid := range book.Bids {
		if positive(t, bid.Price).GreaterThanOrEqual(price) {
			depth = depth.Add(positive(t, bid.Size))
		}
	}
	require.True(t, depth.GreaterThanOrEqual(shares), "insufficient bids above MIN_PRICE")
	order := ct.UserOrderV2{TokenId: m.asset, ConditionId: &m.ConditionID, Side: ct.SELL, Price: price.InexactFloat64(), Size: shares.InexactFloat64(), TickSize: &book.TickSize}
	record("submitting", liveResult{Asset: m.asset, Side: string(ct.SELL)})
	resp, err := a.clob.CreateOrderV2(ctx, &order, ct.OrderTypeFOK, false, false, a.auth)
	if resp != nil {
		record("submitted", liveResult{Asset: m.asset, Side: string(ct.SELL), OrderID: resp.OrderID})
	}
	require.NoError(t, err, "submission may be unknown; reconcile before retrying")
	require.NotNil(t, resp)
	require.True(t, resp.Success, "%s", resp.ErrorMsg)
	require.NotEmpty(t, resp.OrderID)
	fills := verifyFill(t, ctx, a, m, resp.OrderID, ct.SELL)
	filled := decimal.Zero
	for _, fill := range fills {
		filled = filled.Add(positive(t, fill.Size))
	}
	require.True(t, filled.Equal(shares), "FOK sell must fill exactly the requested shares")
	recordFills(record, m, resp.OrderID, ct.SELL, fills)
	after := erc1155Balance(t, ctx, rpc, contract, holder, m.number, nil)
	record("position_after", liveResult{Asset: m.asset, Contract: contract.Hex(), Balance: after.String()})
	require.True(t, decimal.NewFromBigInt(before, -6).Sub(decimal.NewFromBigInt(after, -6)).Equal(shares), "position change differs from sold shares; keep the wallet idle during acceptance")
}

type sellLimitTransport struct {
	asset            string
	shares, minPrice decimal.Decimal
}

func (limit sellLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodPost && req.URL.Path == "/order" {
		body, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			return nil, err
		}
		var payload ct.NewOrderV2
		if err = json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}
		if err = limit.validate(payload); err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	return http.DefaultTransport.RoundTrip(req)
}
func (limit sellLimitTransport) validate(payload ct.NewOrderV2) error {
	o := payload.Order
	maker, e1 := decimal.NewFromString(o.MakerAmount)
	taker, e2 := decimal.NewFromString(o.TakerAmount)
	if payload.OrderType != ct.OrderTypeFOK || o.Side != ct.SELL || o.TokenID != limit.asset || e1 != nil || e2 != nil || !maker.IsPositive() || !taker.IsPositive() {
		return fmt.Errorf("unexpected signed sell order")
	}
	if !maker.Equal(limit.shares.Shift(6)) || taker.LessThan(maker.Mul(limit.minPrice)) {
		return fmt.Errorf("signed amounts differ from SHARES or violate MIN_PRICE; no order sent")
	}
	return nil
}
func TestSignedSellLimits(t *testing.T) {
	limit := sellLimitTransport{"asset", decimal.NewFromInt(2), decimal.RequireFromString("0.5")}
	valid := ct.NewOrderV2{OrderType: ct.OrderTypeFOK, Order: ct.OrderV2{Side: ct.SELL, TokenID: "asset", MakerAmount: "2000000", TakerAmount: "1000000"}}
	require.NoError(t, limit.validate(valid))
	for _, mutate := range []func(*ct.NewOrderV2){
		func(p *ct.NewOrderV2) { p.Order.MakerAmount = "2000001" },
		func(p *ct.NewOrderV2) { p.Order.MakerAmount = "1999999" },
		func(p *ct.NewOrderV2) { p.Order.TakerAmount = "999999" },
		func(p *ct.NewOrderV2) { p.Order.Side = ct.BUY },
		func(p *ct.NewOrderV2) { p.Order.TokenID = "wrong" },
		func(p *ct.NewOrderV2) { p.OrderType = ct.OrderTypeGTC },
		func(p *ct.NewOrderV2) { p.Order.TakerAmount = "invalid" },
	} {
		p := valid
		mutate(&p)
		require.Error(t, limit.validate(p))
	}
}
