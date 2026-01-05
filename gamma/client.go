package gamma

import (
	"fmt"
	"github.com/override-coder/go-polymarket-sdk/gamma/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	client *http2.Client

	chainId *big.Int
}

func NewClient(host string, chainId *big.Int) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	return &Client{
		client:  http2.NewClient(host),
		chainId: chainId,
	}
}

func (c *Client) Search(p *types.SearchParams) (*types.SearchResponse, error) {
	if p == nil {
		return nil, fmt.Errorf("search params is nil")
	}
	q := strings.TrimSpace(p.Q)
	if q == "" {
		return nil, fmt.Errorf("q is required")
	}

	u := url.URL{Path: types.GET_PUBLIC_SEARCH}
	qs := url.Values{}
	qs.Set("q", q)

	if p.Cache != nil {
		qs.Set("cache", strconv.FormatBool(*p.Cache))
	}
	if p.EventsStatus != nil && strings.TrimSpace(*p.EventsStatus) != "" {
		qs.Set("events_status", strings.TrimSpace(*p.EventsStatus))
	}
	if p.LimitPerType != nil {
		qs.Set("limit_per_type", strconv.Itoa(*p.LimitPerType))
	}
	if p.Page != nil {
		qs.Set("page", strconv.Itoa(*p.Page))
	}
	if len(p.EventsTag) > 0 {
		for _, t := range p.EventsTag {
			t = strings.TrimSpace(t)
			if t != "" {
				qs.Add("events_tag", t)
			}
		}
	}
	if p.KeepClosedMarkets != nil {
		qs.Set("keep_closed_markets", strconv.Itoa(*p.KeepClosedMarkets))
	}
	if p.Sort != nil && strings.TrimSpace(*p.Sort) != "" {
		qs.Set("sort", strings.TrimSpace(*p.Sort))
	}
	if p.Ascending != nil {
		qs.Set("ascending", strconv.FormatBool(*p.Ascending))
	}
	if p.SearchTags != nil {
		qs.Set("search_tags", strconv.FormatBool(*p.SearchTags))
	}
	if p.SearchProfiles != nil {
		qs.Set("search_profiles", strconv.FormatBool(*p.SearchProfiles))
	}
	if p.Recurrence != nil && strings.TrimSpace(*p.Recurrence) != "" {
		qs.Set("recurrence", strings.TrimSpace(*p.Recurrence))
	}
	if len(p.ExcludeTagID) > 0 {
		for _, id := range p.ExcludeTagID {
			qs.Add("exclude_tag_id", strconv.Itoa(id))
		}
	}
	if p.Optimized != nil {
		qs.Set("optimized", strconv.FormatBool(*p.Optimized))
	}

	u.RawQuery = qs.Encode()

	var out types.SearchResponse
	resp, err := c.client.DoRequest(http.MethodGet, u.String(), &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}

	return &out, nil
}
