package clob

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"math/big"
	"net/http"
	"time"
)

func (c *Client) EnsureAPIKey(nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	if creds, err := c.DeriveAPIKey(nonce, option); err == nil {
		option.ApiKeyCreds = creds
		return creds, nil
	}
	creds, err := c.CreateApiKey(nonce, option)
	if err != nil {
		return nil, err
	}
	option.ApiKeyCreds = creds
	return creds, nil
}

func (c *Client) CreateApiKey(nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	ts := time.Now().Unix()
	l1Headers, err := headers.CreateL1Headers(option.SingerAddress, c.signFn, c.chainId, nonce, &ts)
	if err != nil {
		return nil, err
	}
	var raw *sdktypes.ApiKeyCreds
	resp, err := c.client.DoRequest(http.MethodPost, types.CREATE_API_KEY, &http2.RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return raw, nil
}

func (c *Client) DeriveAPIKey(nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	ts := time.Now().Unix()
	l1Headers, err := headers.CreateL1Headers(option.SingerAddress, c.signFn, c.chainId, nonce, &ts)
	if err != nil {
		return nil, err
	}
	var raw *sdktypes.ApiKeyCreds
	resp, err := c.client.DoRequest(http.MethodGet, types.DERIVE_API_KEY, &http2.RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return raw, nil
}

func (c *Client) GetBalanceAllowance(option *sdktypes.AuthOption) (*types.BalanceAllowanceResponse, error) {
	requestPath := fmt.Sprintf("%s", types.GET_BALANCE_ALLOWANCE)

	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: requestPath,
	}

	params := make(map[string]any)
	params["asset_type"] = "COLLATERAL"
	params["signature_type"] = 2

	l2Headers, err := headers.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.BalanceAllowanceResponse
	res, err := c.client.DoRequest(http.MethodGet, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "get order buy id:%v", params)
	}
	return &resp, nil
}

func (c *Client) UpdateBalanceAllowance(option *sdktypes.AuthOption) (map[string]interface{}, error) {

	requestPath := fmt.Sprintf("%s", types.UPDATE_BALANCE_ALLOWANCE)

	ts := time.Now().Unix()

	params := make(map[string]any)
	params["asset_type"] = "COLLATERAL"
	params["signature_type"] = 2

	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: requestPath,
	}

	l2Headers, err := headers.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp map[string]interface{}
	res, err := c.client.DoRequest(http.MethodGet, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, e
	}
	return resp, nil
}
