package combos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

const (
	DefaultRFQAPIHost  = "https://combos-rfq-api.polymarket.com"
	DefaultDataAPIHost = "https://data-api.polymarket.com"
)

type Client struct {
	rfqClient  *http2.Client
	dataClient *http2.Client

	chainID *big.Int
	signFn  signing.SignatureFunc
	builder *OrderBuilder
}

func NewClient(rfqHost string, chainID *big.Int, signFn signing.SignatureFunc) *Client {
	return NewClientWithDataHost(rfqHost, DefaultDataAPIHost, chainID, signFn)
}

func NewClientWithDataHost(rfqHost, dataHost string, chainID *big.Int, signFn signing.SignatureFunc) *Client {
	if strings.TrimSpace(rfqHost) == "" {
		rfqHost = DefaultRFQAPIHost
	}
	if strings.TrimSpace(dataHost) == "" {
		dataHost = DefaultDataAPIHost
	}
	rfqHost = strings.TrimRight(rfqHost, "/")
	dataHost = strings.TrimRight(dataHost, "/")

	return &Client{
		rfqClient:  http2.NewClient(rfqHost),
		dataClient: http2.NewClient(dataHost),
		chainID:    chainID,
		signFn:     signFn,
		builder:    NewOrderBuilder(chainID, signFn),
	}
}

func (c *Client) OrderBuilder() *OrderBuilder {
	return c.builder
}

func (c *Client) GetComboMarkets(ctx context.Context, q combostypes.ComboMarketsQuery) (*combostypes.ComboMarketsResponse, error) {
	params := make(map[string]any)
	if q.Limit != nil {
		if *q.Limit < 1 || *q.Limit > 100 {
			return nil, fmt.Errorf("limit out of range (1..100)")
		}
		params["limit"] = *q.Limit
	}
	if q.Cursor != nil && *q.Cursor != "" {
		params["cursor"] = *q.Cursor
	}
	if len(q.Exclude) > 0 {
		params["exclude"] = strings.Join(q.Exclude, ",")
	}

	var out combostypes.ComboMarketsResponse
	resp, err := c.rfqClient.DoRequest(ctx, http.MethodGet, combostypes.GetComboMarkets, &http2.RequestOptions{
		Params: params,
	}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) SubmitQuote(ctx context.Context, req combostypes.QuoteRequest, option *sdktypes.AuthOption) (*combostypes.RFQSnapshot, error) {
	req.Type = ""
	var out combostypes.RFQSnapshot
	err := c.doL2(ctx, http.MethodPost, combostypes.SubmitQuote, req, option, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CancelQuote(ctx context.Context, req combostypes.QuoteCancelRequest, option *sdktypes.AuthOption) (*combostypes.RFQSnapshot, error) {
	req.Type = ""
	var out combostypes.RFQSnapshot
	err := c.doL2(ctx, http.MethodPost, combostypes.CancelQuote, req, option, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ConfirmQuote(ctx context.Context, req combostypes.ConfirmationRequest, option *sdktypes.AuthOption) (*combostypes.ConfirmationResponse, error) {
	req.Type = ""
	req.Decision = combostypes.ConfirmationDecisionConfirm
	var out combostypes.ConfirmationResponse
	err := c.doL2(ctx, http.MethodPost, combostypes.ConfirmQuote, req, option, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeclineQuote(ctx context.Context, req combostypes.ConfirmationRequest, option *sdktypes.AuthOption) (*combostypes.ConfirmationResponse, error) {
	req.Type = ""
	req.Decision = combostypes.ConfirmationDecisionDecline
	var out combostypes.ConfirmationResponse
	err := c.doL2(ctx, http.MethodPost, combostypes.ConfirmQuote, req, option, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetComboPositions(ctx context.Context, q combostypes.ComboPositionsQuery) (*combostypes.ComboPositionsResponse, error) {
	if strings.TrimSpace(q.User) == "" {
		return nil, fmt.Errorf("user is required")
	}
	params := map[string]any{"user": q.User}
	if q.Status != nil && *q.Status != "" {
		params["status"] = *q.Status
	}
	if q.Sort != nil && *q.Sort != "" {
		params["sort"] = *q.Sort
	}
	if len(q.MarketID) > 0 {
		params["market_id"] = strings.Join(q.MarketID, ",")
	}
	if q.Limit != nil {
		if *q.Limit < 1 || *q.Limit > 100 {
			return nil, fmt.Errorf("limit out of range (1..100)")
		}
		params["limit"] = *q.Limit
	}
	if q.Offset != nil {
		if *q.Offset < 0 {
			return nil, fmt.Errorf("offset must be >= 0")
		}
		params["offset"] = *q.Offset
	}

	var out combostypes.ComboPositionsResponse
	resp, err := c.dataClient.DoRequest(ctx, http.MethodGet, combostypes.GetComboPositions, &http2.RequestOptions{
		Params: params,
	}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) GetComboActivity(ctx context.Context, q combostypes.ComboActivityQuery) (*combostypes.ComboActivityResponse, error) {
	if strings.TrimSpace(q.User) == "" {
		return nil, fmt.Errorf("user is required")
	}
	params := map[string]any{"user": q.User}
	if len(q.MarketID) > 0 {
		params["market_id"] = strings.Join(q.MarketID, ",")
	}
	if q.Limit != nil {
		if *q.Limit < 1 || *q.Limit > 500 {
			return nil, fmt.Errorf("limit out of range (1..500)")
		}
		params["limit"] = *q.Limit
	}
	if q.Offset != nil {
		if *q.Offset < 0 {
			return nil, fmt.Errorf("offset must be >= 0")
		}
		params["offset"] = *q.Offset
	}

	var out combostypes.ComboActivityResponse
	resp, err := c.dataClient.DoRequest(ctx, http.MethodGet, combostypes.GetComboActivity, &http2.RequestOptions{
		Params: params,
	}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) doL2(ctx context.Context, method, path string, body any, option *sdktypes.AuthOption, out any) error {
	if option == nil {
		return fmt.Errorf("auth option is required")
	}
	if option.ApiKeyCreds == nil {
		return fmt.Errorf("api key creds are required")
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyStr := string(bodyBytes)
	ts := time.Now().Unix()
	headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, clobtypes.L2HeaderArgs{
		Method:      method,
		RequestPath: path,
		Body:        bodyStr,
	}, &ts)
	if err != nil {
		return err
	}

	resp, err := c.rfqClient.DoRequest(ctx, method, path, &http2.RequestOptions{
		Headers: headers,
		Data:    bodyStr,
	}, out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return e
	}
	return nil
}
