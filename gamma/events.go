package gamma

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/override-coder/go-polymarket-sdk/gamma/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"net/http"
)

func (c *Client) GetEventsBySlug(ctx context.Context, slug string) (*types.Event, error) {
	requestPath := fmt.Sprintf("%s%s", types.GET_EVENTS_SLUG, slug)

	var out types.Event
	resp, err := c.client.DoRequest(ctx, http.MethodGet, requestPath, &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) GetEventsByID(ctx context.Context, id uint64) (*types.Event, error) {
	requestPath := fmt.Sprintf("%s%d", types.GET_EVENTS_ID, id)

	var out types.Event
	resp, err := c.client.DoRequest(ctx, http.MethodGet, requestPath, &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) GetEventsByKeyset(ctx context.Context, params *types.GetEventsKeysetParams) (*types.GetEventsKeysetResponse, error) {
	u := url.URL{Path: types.GET_EVENTS_KEYSET}
	qs := url.Values{}

	if params.Active != nil {
		qs.Set("active", strconv.FormatBool(*params.Active))
	}
	if params.Closed != nil {
		qs.Set("closed", strconv.FormatBool(*params.Closed))
	}
	if params.TagSlug != nil && strings.TrimSpace(*params.TagSlug) != "" {
		qs.Set("tag_slug", strings.TrimSpace(*params.TagSlug))
	}
	if params.TitleSearch != nil && strings.TrimSpace(*params.TitleSearch) != "" {
		qs.Set("title_search", strings.TrimSpace(*params.TitleSearch))
	}
	if params.Limit != nil {
		qs.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Ascending != nil {
		qs.Set("ascending", strconv.FormatBool(*params.Ascending))
	}
	if params.Live != nil {
		qs.Set("live", strconv.FormatBool(*params.Live))
	}
	if params.AfterCursor != nil && strings.TrimSpace(*params.AfterCursor) != "" {
		qs.Set("after_cursor", strings.TrimSpace(*params.AfterCursor))
	}
	if params.Order != nil && strings.TrimSpace(*params.Order) != "" {
		qs.Set("order", strings.TrimSpace(*params.Order))
	}
	if params.EventWeek != nil {
		qs.Set("event_week", strconv.Itoa(*params.EventWeek))
	}
	if params.EventDate != nil && strings.TrimSpace(*params.EventDate) != "" {
		qs.Set("event_date", strings.TrimSpace(*params.EventDate))
	}
	if params.TagMatch != nil && strings.TrimSpace(*params.TagMatch) != "" {
		qs.Set("tag_match", strings.TrimSpace(*params.TagMatch))
	}
	if len(params.ExcludeTagID) > 0 {
		for _, id := range params.ExcludeTagID {
			qs.Add("exclude_tag_id", strconv.FormatUint(id, 10))
		}
	}
	if len(params.TagID) > 0 {
		for _, id := range params.TagID {
			qs.Add("tag_id", strconv.FormatUint(id, 10))
		}
	}
	if params.StartTimeMax != nil && strings.TrimSpace(*params.StartTimeMax) != "" {
		qs.Set("start_time_max", strings.TrimSpace(*params.StartTimeMax))
	}
	if params.StartTimeMin != nil && strings.TrimSpace(*params.StartTimeMin) != "" {
		qs.Set("start_time_min", strings.TrimSpace(*params.StartTimeMin))
	}
	if params.EndDateMax != nil && strings.TrimSpace(*params.EndDateMax) != "" {
		qs.Set("end_date_max", strings.TrimSpace(*params.EndDateMax))
	}
	if params.EndDateMin != nil && strings.TrimSpace(*params.EndDateMin) != "" {
		qs.Set("end_date_min", strings.TrimSpace(*params.EndDateMin))
	}
	if params.StartDateMax != nil && strings.TrimSpace(*params.StartDateMax) != "" {
		qs.Set("start_date_max", strings.TrimSpace(*params.StartDateMax))
	}
	if params.StartDateMin != nil && strings.TrimSpace(*params.StartDateMin) != "" {
		qs.Set("start_date_min", strings.TrimSpace(*params.StartDateMin))
	}
	if params.VolumeMax != nil {
		qs.Set("volume_max", strconv.FormatFloat(*params.VolumeMax, 'f', -1, 64))
	}
	if params.VolumeMin != nil {
		qs.Set("volume_min", strconv.FormatFloat(*params.VolumeMin, 'f', -1, 64))
	}
	if params.LiquidityMax != nil {
		qs.Set("liquidity_max", strconv.FormatFloat(*params.LiquidityMax, 'f', -1, 64))
	}
	if params.LiquidityMin != nil {
		qs.Set("liquidity_min", strconv.FormatFloat(*params.LiquidityMin, 'f', -1, 64))
	}
	if len(params.ID) > 0 {
		for _, id := range params.ID {
			qs.Add("id", strconv.FormatUint(id, 10))
		}
	}
	if len(params.Slug) > 0 {
		for _, slug := range params.Slug {
			if strings.TrimSpace(slug) != "" {
				qs.Add("slug", strings.TrimSpace(slug))
			}
		}
	}

	u.RawQuery = qs.Encode()

	var out types.GetEventsKeysetResponse
	resp, err := c.client.DoRequest(ctx, http.MethodGet, u.String(), &http2.RequestOptions{}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}
