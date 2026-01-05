package gamma

import (
	"fmt"
	"github.com/override-coder/go-polymarket-sdk/gamma/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"net/http"
)

func (c *Client) GetMarketsBySlug(slug string) (*types.Market, error) {
	requestPath := fmt.Sprintf("%s%s", types.GET_MARKETS_SLUG, slug)

	var out types.Market
	resp, err := c.client.DoRequest(http.MethodGet, requestPath, &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) GetMarketsByID(id uint64) (*types.Market, error) {
	requestPath := fmt.Sprintf("%s%d", types.GET_MARKETS_ID, id)

	var out types.Market
	resp, err := c.client.DoRequest(http.MethodGet, requestPath, &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}
