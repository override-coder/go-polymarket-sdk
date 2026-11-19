//go:build polymarket_live

package live_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func positive(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	require.NoError(t, err)
	require.True(t, d.IsPositive(), "expected positive decimal")
	return d
}

// One invocation can submit at most one FOK BUY; the query action never signs an order.
func TestBuyAndVerifyFill(t *testing.T) {
	enabled(t, "buy")
	m := selected(t)
	require.True(t, m.Active && m.AcceptingOrders && !m.Closed)
	if m.Version == "v2" {
		require.Equal(t, "active", m.ResolutionStatus)
	}
	amount := positive(t, env(t, "AMOUNT"))
	price := positive(t, env(t, "MAX_PRICE"))
	budget := amount
	var maxSpend *float64
	if value := os.Getenv("POLY_LIVE_MAX_SPEND"); value != "" {
		budget = positive(t, value)
		require.True(t, amount.GreaterThan(budget), "budget test requires AMOUNT > MAX_SPEND")
		value := budget.InexactFloat64()
		maxSpend = &value
	}
	a := wallet(t, sdkhttp.WithTransport(buyLimitTransport{m.asset, budget, price}))
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	a.credentials(t, ctx)
	if maxSpend != nil {
		rpc := budgetRPC(t, ctx)
		defer rpc.Close()
	}
	require.True(t, amount.Equal(amount.Truncate(2)), "amount must be representable in cents")
	book, err := a.clob.GetOrderBook(m.asset)
	require.NoError(t, err)
	require.Equal(t, m.asset, book.AssetID)
	tick := positive(t, book.TickSize)
	require.True(t, price.Mod(tick).IsZero() && price.GreaterThanOrEqual(tick) && price.LessThanOrEqual(decimal.NewFromInt(1).Sub(tick)), "max price must lie on the live price grid")
	var depth decimal.Decimal
	for _, ask := range book.Asks {
		p := positive(t, ask.Price)
		if p.LessThanOrEqual(price) {
			depth = depth.Add(p.Mul(positive(t, ask.Size)))
		}
	}
	require.True(t, depth.GreaterThanOrEqual(budget), "insufficient executable asks under the price limit")
	// FOK BUY Size is collateral. Pin the verified tick so normalization cannot raise the limit.
	amountFloat := amount.InexactFloat64()
	order := ct.UserOrderV2{TokenId: m.asset, ConditionId: &m.ConditionID, Side: ct.BUY, Price: price.InexactFloat64(), Size: amountFloat, TickSize: &book.TickSize, UserUsdcBalance: maxSpend}
	resp, err := a.clob.CreateOrderV2(ctx, &order, ct.OrderTypeFOK, false, false, a.auth)
	require.NoError(t, err, "submit outcome may be unknown on a transport error; reconcile before rerunning")
	require.NotNil(t, resp)
	t.Logf("market=%s asset=%s order=%s status=%s", m.ID, m.asset, resp.OrderID, resp.Status)
	require.True(t, resp.Success, "%s", resp.ErrorMsg)
	require.NotEmpty(t, resp.OrderID)
	fills := verifyFill(t, ctx, a, m, resp.OrderID)
	if maxSpend != nil {
		verifyBudgetSpend(t, ctx, a, fills, budget)
	}
}

func TestQueryExistingOrderAndFill(t *testing.T) {
	enabled(t, "query")
	m := selected(t)
	a := wallet(t)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	a.credentials(t, ctx)
	fills := verifyFill(t, ctx, a, m, env(t, "ORDER_ID"))
	if value := os.Getenv("POLY_LIVE_MAX_SPEND"); value != "" {
		verifyBudgetSpend(t, ctx, a, fills, positive(t, value))
	}
}

func verifyFill(t *testing.T, ctx context.Context, a account, m selection, id string) []ct.Trade {
	t.Helper()
	var last string
	for ctx.Err() == nil {
		order, err := a.clob.GetOrder(ctx, id, ct.GetOrderRequest{}, a.auth)
		if err == nil && order != nil && order.ID != "" {
			require.True(t, strings.EqualFold(id, order.ID))
			require.True(t, strings.EqualFold(order.MakerAddress, a.auth.FunderAddress), "order belongs to a different funding wallet")
			require.Equal(t, string(ct.BUY), order.Side)
			require.Equal(t, string(ct.OrderTypeFOK), order.OrderType)
			require.Equal(t, m.asset, order.AssetID)
			require.True(t, strings.EqualFold(m.ConditionID, order.Market))
			matched, err := decimal.NewFromString(order.SizeMatched)
			require.NoError(t, err)
			if matched.IsPositive() {
				filter := ct.GetTradesRequest{Market: &m.ConditionID, AssetID: &m.asset}
				seen := map[string]bool{}
				confirmedTrades := map[string]bool{}
				confirmed := decimal.Zero
				var fills []ct.Trade
				pending := false
				for ctx.Err() == nil {
					page, err := a.clob.GetTrades(ctx, filter, a.auth)
					if err != nil {
						last = err.Error()
						break
					}
					for _, fill := range page.Data {
						if !strings.EqualFold(fill.TakerOrderID, id) {
							continue
						}
						require.Equal(t, m.asset, fill.AssetID)
						require.True(t, strings.EqualFold(m.ConditionID, fill.Market))
						status := strings.TrimPrefix(strings.ToUpper(fill.Status), "TRADE_STATUS_")
						last = status
						require.NotEqual(t, "FAILED", status, "trade %s failed", fill.ID)
						if status != "CONFIRMED" {
							pending = true
							continue
						}
						if !confirmedTrades[fill.ID] {
							require.NotEmpty(t, fill.ID)
							confirmed = confirmed.Add(positive(t, fill.Size))
							confirmedTrades[fill.ID] = true
							fills = append(fills, fill)
							require.NotEmpty(t, fill.TransactionHash)
							t.Logf("confirmed order=%s trade=%s tx=%s", id, fill.ID, fill.TransactionHash)
						}
					}
					next := page.NextCursor
					if next == "" || next == "LTE=" {
						if confirmed.Equal(matched) && !pending {
							return fills
						}
						break
					}
					require.False(t, seen[next], "repeated trade cursor")
					seen[next] = true
					filter.NextCursor = &next
				}
			} else {
				last = order.Status
			}
		} else if err != nil {
			last = err.Error()
		}
		select {
		case <-ctx.Done():
		case <-time.After(3 * time.Second):
		}
	}
	t.Fatalf("order %s lacks a confirmed matching fill; last observation %s; resume using the query action", id, last)
	return nil
}

// Check the rounded integer amounts that will actually be posted. The CLOB
// builder can round the requested share amount down after dividing by price.
type buyLimitTransport struct {
	asset              string
	maxSpend, maxPrice decimal.Decimal
}

func (limit buyLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodPost && req.URL.Path == "/order" {
		body, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			return nil, err
		}
		var payload ct.NewOrderV2
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}
		if err := limit.validate(payload); err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	return http.DefaultTransport.RoundTrip(req)
}

func (limit buyLimitTransport) validate(payload ct.NewOrderV2) error {
	o := payload.Order
	maker, e1 := decimal.NewFromString(o.MakerAmount)
	taker, e2 := decimal.NewFromString(o.TakerAmount)
	if payload.OrderType != ct.OrderTypeFOK || o.Side != ct.BUY || o.TokenID != limit.asset || e1 != nil || e2 != nil || !maker.IsPositive() || !taker.IsPositive() {
		return fmt.Errorf("unexpected signed order")
	}
	if maker.GreaterThan(limit.maxSpend.Shift(6)) || maker.GreaterThan(taker.Mul(limit.maxPrice)) {
		return fmt.Errorf("signed amounts exceed the selected spend or price limit; no order sent")
	}
	return nil
}

func TestSignedBuyLimits(t *testing.T) {
	limit := buyLimitTransport{"asset", decimal.NewFromInt(5), decimal.RequireFromString("0.5")}
	p := ct.NewOrderV2{OrderType: ct.OrderTypeFOK, Order: ct.OrderV2{Side: ct.BUY, TokenID: "asset", MakerAmount: "5000000", TakerAmount: "10000000"}}
	require.NoError(t, limit.validate(p))
	p.Order.TakerAmount = "9999999"
	require.Error(t, limit.validate(p))
	p.Order.TakerAmount = "20000000"
	p.Order.MakerAmount = "5000001"
	require.Error(t, limit.validate(p))
}
