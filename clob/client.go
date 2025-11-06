package clob

import (
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/override-coder/go-polymarket-sdk/signing"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"math/big"
	"net/http"
	"strings"
)

type Client struct {
	client *http2.Client

	chainId *big.Int
	signFn  signing.SignatureFunc

	builderApiKeyCreds *types.BuilderApiKeyCreds

	orderBuilder *OrderBuilder

	tickSizes types.TickSizes
	negRisk   types.NegRisks
	feeRates  types.FeeRates
}

func NewClient(host string, chainId *big.Int, signFn signing.SignatureFunc, builderApiKeyCreds *types.BuilderApiKeyCreds) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	return &Client{
		client:             http2.NewClient(host),
		chainId:            chainId,
		builderApiKeyCreds: builderApiKeyCreds,
		orderBuilder:       NewOrderBuilder(chainId, signFn),
		signFn:             signFn,
		tickSizes:          make(types.TickSizes, 500),
		negRisk:            make(types.NegRisks, 500),
		feeRates:           make(types.FeeRates, 500),
	}
}

func (c *Client) GetTickSize(tokenID string) (string, error) {
	if size, ok := c.tickSizes[tokenID]; ok {
		return string(size), nil
	}
	var resp map[string]float64
	res, err := c.client.DoRequest(http.MethodGet, types.GET_TICK_SIZE, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return "", e
	}
	if c.tickSizes == nil {
		c.tickSizes = make(types.TickSizes)
	}

	tickSize := utils.Float64ToDecimal(resp["minimum_tick_size"]).String()
	c.tickSizes[tokenID] = types.TickSize(tickSize)

	return tickSize, nil
}

// GetNegRisk
func (c *Client) GetNegRisk(tokenID string) (bool, error) {
	if neg, ok := c.negRisk[tokenID]; ok {
		return neg, nil
	}

	var resp map[string]bool
	res, err := c.client.DoRequest(http.MethodGet, types.GET_NEG_RISK, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return false, e
	}

	if c.negRisk == nil {
		c.negRisk = make(types.NegRisks)
	}

	negRisk := resp["neg_risk"]
	c.negRisk[tokenID] = negRisk

	return negRisk, nil
}

// GetFeeRateBps
func (c *Client) GetFeeRateBps(tokenID string) (float64, error) {
	if fee, ok := c.feeRates[tokenID]; ok {
		return fee, nil
	}

	var resp map[string]float64
	res, err := c.client.DoRequest(http.MethodGet, types.GET_FEE_RATE, &http2.RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return 0, e
	}

	if c.feeRates == nil {
		c.feeRates = make(types.FeeRates)
	}

	baseFee := resp["base_fee"]
	c.feeRates[tokenID] = baseFee

	return baseFee, nil
}
