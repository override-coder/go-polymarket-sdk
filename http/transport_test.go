package http_test

import (
	"context"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/clob"
	"github.com/override-coder/go-polymarket-sdk/combos"
	"github.com/override-coder/go-polymarket-sdk/dataapi"
	datatype "github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/override-coder/go-polymarket-sdk/gamma"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/override-coder/go-polymarket-sdk/relayer"
	relaytype "github.com/override-coder/go-polymarket-sdk/relayer/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/stretchr/testify/require"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type terminalError struct{}

func (terminalError) Error() string   { return "local admission busy" }
func (terminalError) Retryable() bool { return false }

func jsonResponse(r *http.Request, body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

// These calls use fake transports only; no credentials or live API requests.
func TestTransportReachesEverySlaveClient(t *testing.T) {
	ctx := context.Background()
	chain := big.NewInt(137)
	for _, tc := range []struct {
		name, host, body string
		call             func(sdkhttp.Option) error
	}{
		{"clob", "clob.invalid", `{"minimum_tick_size":0.01}`, func(o sdkhttp.Option) error {
			_, e := clob.NewClient("https://clob.invalid", chain, nil, nil, o).GetTickSize(ctx, "1")
			return e
		}},
		{"data", "data.invalid", `{"data":[],"pagination":{}}`, func(o sdkhttp.Option) error {
			_, e := dataapi.NewClient("https://data.invalid", chain, o).GetPositions(ctx, datatype.PositionsQuery{User: sdktypes.ZeroAddress})
			return e
		}},
		{"data internal clob", "clob.polymarket.com", `{}`, func(o sdkhttp.Option) error {
			_, e := dataapi.NewClient("https://data.invalid", chain, o).GetMarketByToken(ctx, "1")
			return e
		}},
		{"data custom clob", "custom.invalid", `{}`, func(o sdkhttp.Option) error {
			_, e := dataapi.NewClientWithClobHost("https://data.invalid", "https://custom.invalid", chain, o).GetMarketByToken(ctx, "1")
			return e
		}},
		{"gamma", "gamma.invalid", `{}`, func(o sdkhttp.Option) error {
			_, e := gamma.NewClient("https://gamma.invalid", chain, o).GetMarketsByID(ctx, "1")
			return e
		}},
		{"relayer", "relayer.invalid", `{"nonce":"1"}`, func(o sdkhttp.Option) error {
			_, e := relayer.NewClient("https://relayer.invalid", chain, nil, nil, o).GetNonce(sdktypes.ZeroAddress, relaytype.TransactionTypeSAFE)
			return e
		}},
		{"combo", "combos-rfq-gateway-builder.polymarket.com", `{"status":"FAILED"}`, func(o sdkhttp.Option) error {
			_, e := combos.NewBuilderTakerClient("", chain, nil, o).GetRFQStatus(ctx, "test-rfq", &sdktypes.AuthOption{SingerAddress: sdktypes.ZeroAddress, ApiKeyCreds: &sdktypes.ApiKeyCreds{ApiKey: "fixture", Secret: "dGVzdA==", Passphrase: "fixture"}})
			return e
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			option := sdkhttp.WithTransport(transportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, tc.host, r.URL.Host)
				return jsonResponse(r, tc.body), nil
			}))
			require.NoError(t, tc.call(option))
			require.Equal(t, 1, calls)
		})
	}
}

func TestTransportTerminalErrorDoesNotRetryOrCancelParent(t *testing.T) {
	calls := 0
	c := sdkhttp.NewClient("https://test.invalid", sdkhttp.WithTransport(transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return nil, terminalError{}
		}
		return jsonResponse(r, `{}`), nil
	})))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := c.DoRequest(ctx, http.MethodPost, "/order", nil, nil)
	require.ErrorIs(t, err, terminalError{})
	require.Equal(t, 1, calls)
	require.NoError(t, ctx.Err())
	_, err = c.DoRequest(ctx, http.MethodGet, "/next", nil, nil)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestTransportPreservesExistingRetryPolicy(t *testing.T) {
	for _, mode := range []string{"network", "429", "decode"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			c := sdkhttp.NewClient("https://test.invalid", sdkhttp.WithTransport(transportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				switch mode {
				case "network":
					return nil, errors.New("network failure")
				case "429":
					resp := jsonResponse(r, `{}`)
					resp.StatusCode = 429
					return resp, nil
				default:
					return jsonResponse(r, `{invalid`), nil
				}
			})))
			var out map[string]any
			_, err := c.DoRequest(context.Background(), http.MethodPost, "/order", nil, &out)
			if mode == "network" {
				require.Error(t, err)
				require.Equal(t, 3, calls)
			} else {
				require.Equal(t, 1, calls)
			}
			if mode == "decode" {
				require.Error(t, err)
			}
		})
	}
}

func TestTransportOmittedAndNilKeepDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok") }))
	defer server.Close()
	for _, c := range []*sdkhttp.Client{sdkhttp.NewClient(server.URL), sdkhttp.NewClient(server.URL, sdkhttp.WithTransport(nil))} {
		r, err := c.DoRequest(nil, http.MethodGet, "/", nil, nil)
		require.NoError(t, err)
		require.Equal(t, "ok", r.String())
	}
}
