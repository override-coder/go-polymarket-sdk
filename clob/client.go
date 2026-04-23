package clob

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/override-coder/go-polymarket-sdk/clob/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

const (
	defaultMarketInfoCacheTTL     = 30 * time.Minute
	defaultMarketInfoCacheIdleTTL = 72 * time.Hour
)

type Client struct {
	client *http2.Client

	chainId *big.Int
	signFn  signing.SignatureFunc

	builderApiKeyCreds *sdktypes.BuilderApiKeyCreds

	orderBuilder *OrderBuilder

	tickSizes types.TickSizes
	negRisk   types.NegRisks
	feeRates  types.FeeRates

	// v2
	feeInfos          types.FeeInfos
	tokenConditionMap map[string]string

	cacheMu                   sync.RWMutex
	marketInfoCacheTTL        time.Duration
	marketInfoCacheIdleTTL    time.Duration
	marketInfoCachedAt        map[string]time.Time
	marketInfoLastUsedAt      map[string]time.Time
	tokenConditionCachedAt    map[string]time.Time
	tokenConditionLastUsedAt  map[string]time.Time
	conditionTokenMap         map[string][]string
	marketInfoSingleflight    singleflight.Group
	marketByTokenSingleflight singleflight.Group
}

func NewClient(host string, chainId *big.Int, signFn signing.SignatureFunc, builderApiKeyCreds *sdktypes.BuilderApiKeyCreds) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	return &Client{
		client:                   http2.NewClient(host),
		chainId:                  chainId,
		builderApiKeyCreds:       builderApiKeyCreds,
		orderBuilder:             NewOrderBuilder(chainId, signFn),
		signFn:                   signFn,
		tickSizes:                make(types.TickSizes, 500),
		negRisk:                  make(types.NegRisks, 500),
		feeRates:                 make(types.FeeRates, 500),
		feeInfos:                 make(types.FeeInfos, 500),
		tokenConditionMap:        make(map[string]string, 500),
		marketInfoCacheTTL:       defaultMarketInfoCacheTTL,
		marketInfoCacheIdleTTL:   defaultMarketInfoCacheIdleTTL,
		marketInfoCachedAt:       make(map[string]time.Time, 500),
		marketInfoLastUsedAt:     make(map[string]time.Time, 500),
		tokenConditionCachedAt:   make(map[string]time.Time, 500),
		tokenConditionLastUsedAt: make(map[string]time.Time, 500),
		conditionTokenMap:        make(map[string][]string, 500),
	}
}

func (c *Client) WithSignatureFunc(signFn signing.SignatureFunc) error {
	if c.signFn != nil {
		return errors.New("signFn already set")
	}
	c.signFn = signFn
	return c.orderBuilder.WithSignatureFunc(signFn)
}

func (c *Client) WithMarketInfoCacheTTL(ttl time.Duration) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.marketInfoCacheTTL = ttl
}

func (c *Client) WithMarketInfoCacheIdleTTL(ttl time.Duration) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.marketInfoCacheIdleTTL = ttl
}

func (c *Client) GetTickSize(ctx context.Context, tokenID string) (string, error) {
	c.cacheMu.Lock()
	if size, ok := c.tickSizes[tokenID]; ok {
		c.touchTokenLocked(tokenID, time.Now())
		c.cacheMu.Unlock()
		return string(size), nil
	}
	c.cacheMu.Unlock()

	var resp map[string]float64
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_TICK_SIZE, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return "", e
	}
	tickSize := utils.Float64ToDecimal(resp["minimum_tick_size"]).String()

	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if c.tickSizes == nil {
		c.tickSizes = make(types.TickSizes)
	}
	c.tickSizes[tokenID] = types.TickSize(tickSize)
	c.touchTokenLocked(tokenID, time.Now())

	return tickSize, nil
}

// GetNegRisk
func (c *Client) GetNegRisk(ctx context.Context, tokenID string) (bool, error) {
	c.cacheMu.Lock()
	if neg, ok := c.negRisk[tokenID]; ok {
		c.touchTokenLocked(tokenID, time.Now())
		c.cacheMu.Unlock()
		return neg, nil
	}
	c.cacheMu.Unlock()

	var resp map[string]bool
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_NEG_RISK, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return false, e
	}
	negRisk := resp["neg_risk"]

	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if c.negRisk == nil {
		c.negRisk = make(types.NegRisks)
	}
	c.negRisk[tokenID] = negRisk
	c.touchTokenLocked(tokenID, time.Now())

	return negRisk, nil
}

// GetFeeRateBps
func (c *Client) GetFeeRateBps(ctx context.Context, tokenID string) (float64, error) {
	c.cacheMu.Lock()
	if fee, ok := c.feeRates[tokenID]; ok {
		c.touchTokenLocked(tokenID, time.Now())
		c.cacheMu.Unlock()
		return fee, nil
	}
	c.cacheMu.Unlock()

	var resp map[string]float64
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_FEE_RATE, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return 0, e
	}
	baseFee := resp["base_fee"]

	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if c.feeRates == nil {
		c.feeRates = make(types.FeeRates)
	}
	c.feeRates[tokenID] = baseFee
	c.touchTokenLocked(tokenID, time.Now())

	return baseFee, nil
}

func (c *Client) GetCLOBMarketInfo(ctx context.Context, conditionID string) (*types.CLOBMarketInfo, error) {
	var resp types.CLOBMarketInfo
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_CLOB_MARKET+conditionID, nil, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get clob market info")
	}
	return &resp, nil
}

func (c *Client) GetMarketByToken(ctx context.Context, tokenID string) (*types.MarketByTokenResponse, error) {
	var resp types.MarketByTokenResponse
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_MARKET_BY_TOKEN+tokenID, nil, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get market by token")
	}
	return &resp, nil
}

func (c *Client) GetOrderBook(tokenID string) (*types.OrderBookSummary, error) {
	var resp types.OrderBookSummary
	res, err := c.client.DoRequest(context.Background(), http.MethodGet, types.GET_ORDER_BOOK, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get order book")
	}
	return &resp, nil
}

func (c *Client) GetMarketPrice(ctx context.Context, tokenID, side string) (string, error) {
	var resp map[string]string
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_PRICE, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID, "side": side},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return "", errors.Wrap(e, "get market price")
	}
	return resp["price"], nil
}

func (c *Client) GetBuilderFee(ctx context.Context, builderCode string) (string, string, error) {
	var resp map[string]string
	res, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_BUILDER_FEES, &http2.RequestOptions{
		Params: map[string]any{"builder_code": builderCode},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return "", "", errors.Wrap(e, "get market price")
	}
	return resp["builder_maker_fee_rate_bps"], resp["builder_taker_fee_rate_bps"], nil
}

func (c *Client) GetMarketPrices(ctx context.Context, prices []types.PricesRequest) (map[string]string, error) {
	var resp map[string]map[string]string
	res, err := c.client.DoRequest(ctx, http.MethodPost, types.GET_PRICES, &http2.RequestOptions{
		Data: prices,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get market price")
	}
	out := make(map[string]string, len(prices))
	for _, id := range prices {
		if sideMap, ok := resp[id.TokenId]; ok {
			if price, found := sideMap["BUY"]; found {
				out[id.TokenId] = price
			}
			if price, found := sideMap["SELL"]; found {
				out[id.TokenId] = price
			}
		}
	}
	return out, nil
}

func (c *Client) ensureMarketInfoCached(ctx context.Context, tokenID string, conditionId *string) error {
	if c.hasFreshMarketInfoForToken(tokenID, conditionId) {
		return nil
	}

	conditionID := ""
	if conditionId != nil {
		conditionID = *conditionId
	}

	if conditionID == "" {
		cachedConditionID, ok := c.getCachedConditionID(tokenID)
		if !ok {
			value, err, _ := c.marketByTokenSingleflight.Do(tokenID, func() (any, error) {
				if cachedConditionIDSF, okSF := c.getCachedConditionID(tokenID); okSF {
					return cachedConditionIDSF, nil
				}

				result, err := c.GetMarketByToken(ctx, tokenID)
				if err != nil {
					return "", err
				}
				if result.ConditionID == "" {
					return "", fmt.Errorf("failed to resolve condition id for token %s", tokenID)
				}

				c.storeTokenConditionMap(result.ConditionID, tokenID, result.PrimaryTokenID, result.SecondaryTokenID)
				return result.ConditionID, nil
			})
			if err != nil {
				return err
			}
			cachedConditionID = value.(string)
		}
		conditionID = cachedConditionID
	}

	_, err := c.getClobMarketInfoCached(ctx, conditionID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getClobMarketInfo(ctx context.Context, conditionID string) (*types.CLOBMarketInfo, error) {
	marketInfo, err := c.GetCLOBMarketInfo(ctx, conditionID)
	if err != nil {
		return nil, errors.Wrapf(err, "get clob market info")
	}
	if len(marketInfo.Tokens) == 0 {
		return nil, errors.Errorf("no token found for condition %s", conditionID)
	}
	c.storeMarketInfo(conditionID, marketInfo)
	return marketInfo, nil
}

func (c *Client) getClobMarketInfoCached(ctx context.Context, conditionID string) (*types.CLOBMarketInfo, error) {
	if c.hasFreshMarketInfoForCondition(conditionID) {
		return nil, nil
	}

	value, err, _ := c.marketInfoSingleflight.Do(conditionID, func() (any, error) {
		if c.hasFreshMarketInfoForCondition(conditionID) {
			return nil, nil
		}
		return c.getClobMarketInfo(ctx, conditionID)
	})
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return value.(*types.CLOBMarketInfo), nil
}

func (c *Client) storeMarketInfo(conditionID string, marketInfo *types.CLOBMarketInfo) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.tokenConditionMap == nil {
		c.tokenConditionMap = make(map[string]string)
	}
	if c.tickSizes == nil {
		c.tickSizes = make(types.TickSizes)
	}
	if c.negRisk == nil {
		c.negRisk = make(types.NegRisks)
	}
	if c.feeInfos == nil {
		c.feeInfos = make(types.FeeInfos)
	}
	if c.marketInfoCachedAt == nil {
		c.marketInfoCachedAt = make(map[string]time.Time)
	}
	if c.marketInfoLastUsedAt == nil {
		c.marketInfoLastUsedAt = make(map[string]time.Time)
	}
	if c.tokenConditionCachedAt == nil {
		c.tokenConditionCachedAt = make(map[string]time.Time)
	}
	if c.tokenConditionLastUsedAt == nil {
		c.tokenConditionLastUsedAt = make(map[string]time.Time)
	}
	if c.conditionTokenMap == nil {
		c.conditionTokenMap = make(map[string][]string)
	}

	now := time.Now()
	c.marketInfoCachedAt[conditionID] = now
	c.marketInfoLastUsedAt[conditionID] = now
	c.conditionTokenMap[conditionID] = c.conditionTokenMap[conditionID][:0]
	for _, token := range marketInfo.Tokens {
		tokenID := token.TokenID
		c.tokenConditionMap[tokenID] = conditionID
		c.tokenConditionCachedAt[tokenID] = now
		c.tokenConditionLastUsedAt[tokenID] = now
		c.conditionTokenMap[conditionID] = append(c.conditionTokenMap[conditionID], tokenID)
		c.tickSizes[tokenID] = types.TickSize(utils.Float64ToDecimal(marketInfo.MinimumTickSize).String())

		negRisk := false
		if marketInfo.NegRisk != nil {
			negRisk = *marketInfo.NegRisk
		}
		c.negRisk[tokenID] = negRisk

		rate := 0.0
		exponent := 0.0
		if marketInfo.FeeDetails != nil {
			rate = marketInfo.FeeDetails.Rate
			exponent = marketInfo.FeeDetails.Exponent
		}

		c.feeInfos[tokenID] = types.FeeInfo{
			Rate:     rate,
			Exponent: exponent,
		}
	}
}

func (c *Client) storeTokenConditionMap(conditionID string, tokenIDs ...string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.tokenConditionMap == nil {
		c.tokenConditionMap = make(map[string]string)
	}
	if c.tokenConditionCachedAt == nil {
		c.tokenConditionCachedAt = make(map[string]time.Time)
	}
	if c.tokenConditionLastUsedAt == nil {
		c.tokenConditionLastUsedAt = make(map[string]time.Time)
	}

	now := time.Now()
	for _, tokenID := range tokenIDs {
		if tokenID == "" {
			continue
		}
		c.tokenConditionMap[tokenID] = conditionID
		c.tokenConditionCachedAt[tokenID] = now
		c.tokenConditionLastUsedAt[tokenID] = now
	}
}

func (c *Client) hasFreshMarketInfoForToken(tokenID string, conditionId *string) bool {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if _, ok := c.feeInfos[tokenID]; !ok {
		return false
	}

	conditionID := ""
	if conditionId != nil {
		conditionID = *conditionId
	}
	if conditionID == "" {
		conditionID = c.tokenConditionMap[tokenID]
	}
	if conditionID == "" {
		return false
	}
	if !c.isFreshLocked(c.marketInfoCachedAt[conditionID]) {
		return false
	}
	c.touchConditionLocked(conditionID, time.Now())
	return true
}

func (c *Client) hasFreshMarketInfoForCondition(conditionID string) bool {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if !c.isFreshLocked(c.marketInfoCachedAt[conditionID]) {
		return false
	}
	c.touchConditionLocked(conditionID, time.Now())
	return true
}

func (c *Client) getCachedConditionID(tokenID string) (string, bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	conditionID := c.tokenConditionMap[tokenID]
	if conditionID == "" {
		return "", false
	}
	if !c.isFreshLocked(c.tokenConditionCachedAt[tokenID]) {
		return "", false
	}
	c.touchTokenLocked(tokenID, time.Now())
	return conditionID, true
}

func (c *Client) getCachedFeeInfo(tokenID string) types.FeeInfo {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	c.touchTokenLocked(tokenID, time.Now())
	return c.feeInfos[tokenID]
}

func (c *Client) ClearIdleMarketInfoCache() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	now := time.Now()
	c.clearIdleMarketInfoCacheLocked(now)
}

func (c *Client) StartMarketInfoCacheJanitor(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Hour
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.ClearIdleMarketInfoCache()
		}
	}
}

func (c *Client) clearIdleMarketInfoCacheLocked(now time.Time) {
	if c.marketInfoCacheIdleTTL <= 0 {
		return
	}

	for conditionID, lastUsedAt := range c.marketInfoLastUsedAt {
		if lastUsedAt.IsZero() || now.Sub(lastUsedAt) <= c.marketInfoCacheIdleTTL {
			continue
		}
		c.deleteConditionCacheLocked(conditionID)
	}

	for tokenID, lastUsedAt := range c.tokenConditionLastUsedAt {
		if lastUsedAt.IsZero() || now.Sub(lastUsedAt) <= c.marketInfoCacheIdleTTL {
			continue
		}
		c.deleteTokenCacheLocked(tokenID)
	}
}

func (c *Client) deleteConditionCacheLocked(conditionID string) {
	tokenIDs := c.conditionTokenMap[conditionID]
	for _, tokenID := range tokenIDs {
		c.deleteTokenCacheLocked(tokenID)
	}
	delete(c.conditionTokenMap, conditionID)
	delete(c.marketInfoCachedAt, conditionID)
	delete(c.marketInfoLastUsedAt, conditionID)
}

func (c *Client) deleteTokenCacheLocked(tokenID string) {
	delete(c.tickSizes, tokenID)
	delete(c.negRisk, tokenID)
	delete(c.feeRates, tokenID)
	delete(c.feeInfos, tokenID)
	delete(c.tokenConditionMap, tokenID)
	delete(c.tokenConditionCachedAt, tokenID)
	delete(c.tokenConditionLastUsedAt, tokenID)
}

func (c *Client) touchTokenLocked(tokenID string, now time.Time) {
	if tokenID == "" {
		return
	}
	if c.tokenConditionLastUsedAt == nil {
		c.tokenConditionLastUsedAt = make(map[string]time.Time)
	}
	c.tokenConditionLastUsedAt[tokenID] = now
	if conditionID := c.tokenConditionMap[tokenID]; conditionID != "" {
		c.touchConditionLocked(conditionID, now)
	}
}

func (c *Client) touchConditionLocked(conditionID string, now time.Time) {
	if conditionID == "" {
		return
	}
	if c.marketInfoLastUsedAt == nil {
		c.marketInfoLastUsedAt = make(map[string]time.Time)
	}
	c.marketInfoLastUsedAt[conditionID] = now
	for _, tokenID := range c.conditionTokenMap[conditionID] {
		if tokenID == "" {
			continue
		}
		if c.tokenConditionLastUsedAt == nil {
			c.tokenConditionLastUsedAt = make(map[string]time.Time)
		}
		c.tokenConditionLastUsedAt[tokenID] = now
	}
}

func (c *Client) isFreshLocked(cachedAt time.Time) bool {
	if cachedAt.IsZero() {
		return false
	}
	if c.marketInfoCacheTTL <= 0 {
		return true
	}
	return time.Since(cachedAt) <= c.marketInfoCacheTTL
}
