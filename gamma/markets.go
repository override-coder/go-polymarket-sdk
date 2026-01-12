package gamma

import (
	"fmt"
	"github.com/override-coder/go-polymarket-sdk/gamma/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (c *Client) GetMarkets(p *types.GetMarketsParams) ([]*types.Market, error) {
	if p == nil {
		return nil, fmt.Errorf("get markets params is nil")
	}
	if p.Limit < 0 || p.Offset < 0 {
		return nil, fmt.Errorf("limit/offset must be >= 0")
	}

	u := url.URL{Path: types.GET_MARKETS}
	qs := url.Values{}
	qs.Set("limit", strconv.Itoa(p.Limit))
	qs.Set("offset", strconv.Itoa(p.Offset))

	if p.Order != nil && strings.TrimSpace(*p.Order) != "" {
		qs.Set("order", strings.TrimSpace(*p.Order))
	}
	if p.Ascending != nil {
		qs.Set("ascending", strconv.FormatBool(*p.Ascending))
	}

	if p.LiquidityNumMin != nil {
		qs.Set("liquidity_num_min", p.LiquidityNumMin.String())
	}
	if p.LiquidityNumMax != nil {
		qs.Set("liquidity_num_max", p.LiquidityNumMax.String())
	}
	if p.VolumeNumMin != nil {
		qs.Set("volume_num_min", p.VolumeNumMin.String())
	}
	if p.VolumeNumMax != nil {
		qs.Set("volume_num_max", p.VolumeNumMax.String())
	}
	if p.StartDateMin != nil && strings.TrimSpace(*p.StartDateMin) != "" {
		qs.Set("start_date_min", strings.TrimSpace(*p.StartDateMin))
	}
	if p.StartDateMax != nil && strings.TrimSpace(*p.StartDateMax) != "" {
		qs.Set("start_date_max", strings.TrimSpace(*p.StartDateMax))
	}
	if p.EndDateMin != nil && strings.TrimSpace(*p.EndDateMin) != "" {
		qs.Set("end_date_min", strings.TrimSpace(*p.EndDateMin))
	}
	if p.EndDateMax != nil && strings.TrimSpace(*p.EndDateMax) != "" {
		qs.Set("end_date_max", strings.TrimSpace(*p.EndDateMax))
	}
	if p.TagID != nil {
		qs.Set("tag_id", strconv.Itoa(int(*p.TagID)))
	}
	if p.RelatedTags != nil {
		qs.Set("related_tags", strconv.FormatBool(*p.RelatedTags))
	}
	if p.CYOM != nil {
		qs.Set("cyom", strconv.FormatBool(*p.CYOM))
	}
	if p.UMAResolution != nil && strings.TrimSpace(*p.UMAResolution) != "" {
		qs.Set("uma_resolution_status", strings.TrimSpace(*p.UMAResolution))
	}
	if p.GameID != nil && strings.TrimSpace(*p.GameID) != "" {
		qs.Set("game_id", strings.TrimSpace(*p.GameID))
	}
	if p.RewardsMinSize != nil {
		qs.Set("rewards_min_size", p.RewardsMinSize.String())
	}
	if p.IncludeTag != nil {
		qs.Set("include_tag", strconv.FormatBool(*p.IncludeTag))
	}
	if p.Closed != nil {
		qs.Set("closed", strconv.FormatBool(*p.Closed))
	}

	for _, id := range p.ID {
		qs.Add("id", strconv.Itoa(int(id)))
	}
	for _, s := range p.Slug {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("slug", s)
		}
	}
	for _, s := range p.CLOBTokenIDs {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("clob_token_ids", s)
		}
	}
	for _, s := range p.ConditionIDs {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("condition_ids", s)
		}
	}
	for _, s := range p.MarketMakerAddr {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("market_maker_address", s)
		}
	}
	for _, s := range p.SportsMarketTypes {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("sports_market_types", s)
		}
	}
	for _, s := range p.QuestionIDs {
		s = strings.TrimSpace(s)
		if s != "" {
			qs.Add("question_ids", s)
		}
	}

	u.RawQuery = qs.Encode()

	var out []*types.Market
	resp, err := c.client.DoRequest(http.MethodGet, u.String(), &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return out, nil
}
