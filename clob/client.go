package clob

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/signing"
	"github.com/pkg/errors"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client *resty.Client

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
		client:             resty.New().SetBaseURL(host),
		chainId:            chainId,
		builderApiKeyCreds: builderApiKeyCreds,
		orderBuilder:       NewOrderBuilder(chainId, signFn),
		signFn:             signFn,
		tickSizes:          make(types.TickSizes, 500),
		negRisk:            make(types.NegRisks, 500),
		feeRates:           make(types.FeeRates, 500),
	}
}

type RequestOptions struct {
	Headers map[string]string
	Data    any
	Params  map[string]any
}

func (c *Client) withDefaults() *resty.Client {
	return c.client.
		SetTimeout(15*time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(300*time.Millisecond).
		SetHeader("Accept", "*/*").
		SetHeader("Connection", "keep-alive").
		SetHeader("User-Agent", "@polymarket/go-polymarket-sdk").
		SetHeader("Content-Type", "application/json")
}

func (c *Client) doRequest(method, endpoint string, opt *RequestOptions, out any) (*resty.Response, error) {
	rc := c.withDefaults().R()
	if opt != nil {
		if opt.Headers != nil {
			for k, v := range opt.Headers {
				rc.SetHeader(k, v)
			}
		}
		if opt.Params != nil {
			rc.SetQueryParamsFromValues(toValues(opt.Params))
		}
		if opt.Data != nil {
			switch b := opt.Data.(type) {
			case string:
				rc.SetBody(b)
			case []byte:
				rc.SetBody(b)
			default:
				rc.SetBody(opt.Data)
			}
		}
	}
	if out != nil {
		rc.SetResult(out)
	}

	switch strings.ToUpper(method) {
	case http.MethodGet:
		return rc.Get(endpoint)
	case http.MethodPost:
		return rc.Post(endpoint)
	case http.MethodDelete:
		return rc.Delete(endpoint)
	case http.MethodPut:
		return rc.Put(endpoint)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

func toValues(m map[string]any) map[string][]string {
	v := make(map[string][]string, len(m))
	for k, val := range m {
		switch t := val.(type) {
		case []string:
			v[k] = t
		default:
			v[k] = []string{fmt.Sprint(val)}
		}
	}
	return v
}

func parseHTTPError(resp *resty.Response, err error) (any, error) {
	if err != nil {
		return map[string]any{"error": err.Error()}, err
	}
	if resp.IsSuccess() {
		return resp, nil
	}
	var body any
	b := resp.Body()
	_ = json.Unmarshal(b, &body)
	if body == nil {
		body = string(b)
	}
	return map[string]any{
		"status":      resp.StatusCode(),
		"status_text": resp.Status(),
		"error":       body,
	}, errors.Errorf("http non-2xx: %s", body)
}
