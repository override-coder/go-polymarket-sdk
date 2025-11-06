package clob

import (
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"net/http"
)

func (c *Client) GetTickSize(tokenID string) (string, error) {
	if size, ok := c.tickSizes[tokenID]; ok {
		return string(size), nil
	}
	var resp map[string]float64
	res, err := c.doRequest(http.MethodGet, types.GET_TICK_SIZE, &RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := parseHTTPError(res, err); e != nil {
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
	res, err := c.doRequest(http.MethodGet, types.GET_NEG_RISK, &RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := parseHTTPError(res, err); e != nil {
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
	res, err := c.doRequest(http.MethodGet, types.GET_FEE_RATE, &RequestOptions{
		Params: map[string]any{"token_id": tokenID},
	}, &resp)
	if _, e := parseHTTPError(res, err); e != nil {
		return 0, e
	}

	if c.feeRates == nil {
		c.feeRates = make(types.FeeRates)
	}

	baseFee := resp["base_fee"]
	c.feeRates[tokenID] = baseFee

	return baseFee, nil
}
