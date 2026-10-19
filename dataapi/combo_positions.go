package dataapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
)

// GetComboPositions reads a page of enriched Combo holdings. Sorting is done by
// the server across all pages, not locally on the returned page.
func (c *Client) GetComboPositions(ctx context.Context, q types.ComboPositionsQuery) (*types.Page[types.ComboPosition], error) {
	user := strings.TrimSpace(q.User)
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	if err := validateConditions(q.Condition); err != nil {
		return nil, err
	}
	if q.SortBy != nil {
		switch *q.SortBy {
		case "FIRST_ENTRY", "ENTRY_COST", "CURRENT_VALUE", "UPDATED":
		default:
			return nil, fmt.Errorf("invalid combo sort_by")
		}
	}
	if q.SortDirection != nil && *q.SortDirection != types.SortASC && *q.SortDirection != types.SortDESC {
		return nil, fmt.Errorf("invalid sort_direction")
	}
	if q.Status != nil {
		values := strings.Split(*q.Status, ",")
		for _, value := range values {
			switch value {
			case "OPEN", "PARTIAL", "RESOLVED_WIN", "RESOLVED_LOSS", "RESOLVED_PARTIAL":
			case "REDEEMABLE":
				if len(values) != 1 {
					return nil, fmt.Errorf("REDEEMABLE must be the sole status")
				}
			default:
				return nil, fmt.Errorf("invalid combo status")
			}
		}
	}
	if (q.UpdatedAfter != nil && *q.UpdatedAfter < 0) || (q.UpdatedBefore != nil && *q.UpdatedBefore < 0) {
		return nil, fmt.Errorf("updated bounds must be non-negative")
	}
	if q.UpdatedAfter != nil && q.UpdatedBefore != nil && *q.UpdatedAfter > *q.UpdatedBefore {
		return nil, fmt.Errorf("updated_before precedes updated_after")
	}
	params := map[string]any{"user": user}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	if len(q.Condition) > 0 {
		params["condition"] = strings.Join(q.Condition, ",")
	}
	addParam(params, "status", q.Status)
	addParam(params, "sort_by", q.SortBy)
	addParam(params, "sort_direction", q.SortDirection)
	addParam(params, "updated_after", q.UpdatedAfter)
	addParam(params, "updated_before", q.UpdatedBefore)
	var out types.Page[types.ComboPosition]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_COMBO_POSITIONS, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}
