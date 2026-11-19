package clob

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/pkg/errors"
)

func (c *Client) EnsureAPIKey(ctx context.Context, nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	if creds, err := c.DeriveAPIKey(ctx, nonce, option); err == nil {
		option.ApiKeyCreds = creds
		return creds, nil
	}
	creds, err := c.CreateApiKey(ctx, nonce, option)
	if err != nil {
		return nil, err
	}
	option.ApiKeyCreds = creds
	return creds, nil
}

func (c *Client) CreateApiKey(ctx context.Context, nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	ts := time.Now().Unix()
	l1Headers, err := headers.CreateL1Headers(option.SingerAddress, c.signFn, c.chainId, nonce, &ts)
	if err != nil {
		return nil, err
	}
	var raw *sdktypes.ApiKeyCreds
	resp, err := c.client.DoRequest(ctx, http.MethodPost, types.CREATE_API_KEY, &http2.RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return raw, nil
}

func (c *Client) DeriveAPIKey(ctx context.Context, nonce *big.Int, option *sdktypes.AuthOption) (*sdktypes.ApiKeyCreds, error) {
	ts := time.Now().Unix()
	l1Headers, err := headers.CreateL1Headers(option.SingerAddress, c.signFn, c.chainId, nonce, &ts)
	if err != nil {
		return nil, err
	}
	var raw *sdktypes.ApiKeyCreds
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.DERIVE_API_KEY, &http2.RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return raw, nil
}

func (c *Client) CreateBuilderApiKey(ctx context.Context, option *sdktypes.AuthOption) (*sdktypes.BuilderApiKeyCreds, error) {
	if option == nil {
		return nil, errors.New("auth option is required")
	}
	if option.ApiKeyCreds == nil {
		return nil, errors.New("api key creds are required")
	}
	if option.SingerAddress == "" {
		return nil, errors.New("signer address is required")
	}

	requestPath := types.CREATE_BUILDER_API_KEY
	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodPost,
		RequestPath: requestPath,
	}

	l2Headers, err := headers.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var raw sdktypes.BuilderApiKeyCreds
	resp, err := c.client.DoRequest(ctx, http.MethodPost, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
	}, &raw)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &raw, nil
}

func (c *Client) GetBalanceAllowance(option *sdktypes.AuthOption) (*types.BalanceAllowanceResponse, error) {
	return c.FetchBalanceAllowance(context.Background(), types.BalanceAllowanceRequest{
		AssetType: types.AssetTypeCollateral,
	}, option)
}

// FetchBalanceAllowance reads collateral or position balances for the funding wallet.
// Use CONDITIONAL for CTF tokens and CONDITIONAL-V2 for Protocol V2 positions.
func (c *Client) FetchBalanceAllowance(ctx context.Context, request types.BalanceAllowanceRequest, option *sdktypes.AuthOption) (*types.BalanceAllowanceResponse, error) {
	if option == nil || option.ApiKeyCreds == nil {
		return nil, fmt.Errorf("balance allowance requires CLOB credentials")
	}
	switch request.AssetType {
	case types.AssetTypeCollateral, types.AssetTypeConditional, types.AssetTypeConditionalV2:
	default:
		return nil, fmt.Errorf("unsupported balance allowance asset type %q", request.AssetType)
	}
	requestPath := types.GET_BALANCE_ALLOWANCE
	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{Method: http.MethodGet, RequestPath: requestPath}
	params := map[string]any{
		"asset_type":     request.AssetType,
		"signature_type": option.SignatureType,
	}
	if request.TokenID != "" {
		params["token_id"] = request.TokenID
	}

	l2Headers, err := headers.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.BalanceAllowanceResponse
	res, err := c.client.DoRequest(ctx, http.MethodGet, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "get order buy id:%v", params)
	}
	return &resp, nil
}

func (c *Client) UpdateBalanceAllowance(params map[string]any, option *sdktypes.AuthOption) (map[string]interface{}, error) {

	requestPath := fmt.Sprintf("%s", types.UPDATE_BALANCE_ALLOWANCE)

	ts := time.Now().Unix()

	//params := make(map[string]any)
	//params["asset_type"] = "CONDITIONAL"
	//params["token_id"] = "45763018441764333771124945243746174684578244015331389396782339063349542289693"
	//params["signature_type"] = 2

	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: requestPath,
	}

	l2Headers, err := headers.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp map[string]interface{}
	res, err := c.client.DoRequest(context.Background(), http.MethodGet, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, e
	}
	return resp, nil
}
