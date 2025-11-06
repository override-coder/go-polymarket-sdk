package clob

import (
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/headers"
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
	resp, err := c.doRequest(http.MethodPost, types.CREATE_API_KEY, &RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := parseHTTPError(resp, err); e != nil {
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
	resp, err := c.doRequest(http.MethodGet, types.DERIVE_API_KEY, &RequestOptions{
		Headers: l1Headers,
	}, &raw)
	if _, e := parseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return raw, nil
}
