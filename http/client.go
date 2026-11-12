package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client            *resty.Client
	terminalTransport bool
}

// Option configures a client before it is used concurrently.
type Option func(*resty.Client)

// WithTransport injects an HTTP transport without changing other client defaults.
// Errors implementing Retryable() bool with false stop this request's retries.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *resty.Client) {
		if transport != nil {
			c.SetTransport(terminalErrorTransport{transport})
		}
	}
}

type retryCancelKey struct{}
type terminalErrorTransport struct{ http.RoundTripper }

func (t terminalErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(req)
	var terminal interface{ Retryable() bool }
	if errors.As(err, &terminal) && !terminal.Retryable() {
		if cancel, ok := req.Context().Value(retryCancelKey{}).(context.CancelFunc); ok {
			// Cancel only the SDK-owned child context. Resty checks it before
			// retrying, preserving its existing policy for every other error.
			cancel()
		}
	}
	return resp, err
}

func NewClient(host string, options ...Option) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	c := resty.New().
		SetBaseURL(host).
		SetTimeout(15 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(300 * time.Millisecond).
		SetRetryMaxWaitTime(3 * time.Second)
	for _, option := range options {
		option(c)
	}
	_, terminalTransport := c.GetClient().Transport.(terminalErrorTransport)
	return &Client{client: c, terminalTransport: terminalTransport}
}

type RequestOptions struct {
	Headers map[string]string
	Data    any
	Params  map[string]any
}

// 仅设置本次请求的默认 Header（不要再改 client 级 Header）
func (c *Client) newRequest(ctx context.Context) *resty.Request {
	r := c.client.R()
	if ctx != nil {
		r.SetContext(ctx)
	}
	r.SetHeader("Accept", "*/*")
	r.SetHeader("Connection", "keep-alive")
	r.SetHeader("User-Agent", "@polymarket/go-polymarket-sdk")
	return r
}

func (c *Client) DoRequest(ctx context.Context, method, endpoint string, opt *RequestOptions, out any) (*resty.Response, error) {
	if c.terminalTransport {
		if ctx == nil {
			ctx = context.Background()
		}
		requestCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		ctx = context.WithValue(requestCtx, retryCancelKey{}, cancel)
	}
	rc := c.newRequest(ctx)
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
				rc.SetHeader("Content-Type", "application/json")
				rc.SetBody(b)
			case []byte:
				rc.SetHeader("Content-Type", "application/json")
				rc.SetBody(b)
			default:
				rc.SetHeader("Content-Type", "application/json")
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

func ParseHTTPError(resp *resty.Response, err error) (any, error) {
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

	errMsg := extractErrorMessage(body)
	if errMsg == "" {
		errMsg = string(b)
	}

	return map[string]any{
			"status":      resp.StatusCode(),
			"status_text": resp.Status(),
			"error":       body,
		}, &UpstreamServiceError{
			Service:    "polymarket",
			StatusCode: resp.StatusCode(),
			Message:    errMsg,
			Body:       body,
		}
}

func extractErrorMessage(body any) string {
	switch v := body.(type) {
	case map[string]any:
		if msg, ok := v["error"].(string); ok {
			return msg
		}
		if msg, ok := v["message"].(string); ok {
			return msg
		}
		if msg, ok := v["msg"].(string); ok {
			return msg
		}
		b, _ := json.Marshal(v)
		return string(b)
	case string:
		return v
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
