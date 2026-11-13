package dataapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
)

// GetComboActivity reads one V2 page. Keep the same user and condition filters
// when following NextCursor, including after a short or empty page.
func (c *Client) GetComboActivity(ctx context.Context, q types.ComboActivityQuery) (*types.Page[types.ComboActivity], error) {
	user := strings.TrimSpace(q.User)
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	if err := validateComboConditions(q.Condition); err != nil {
		return nil, err
	}
	params := map[string]any{"user": user}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	if len(q.Condition) > 0 {
		params["condition"] = strings.Join(q.Condition, ",")
	}
	var out types.Page[types.ComboActivity]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_COMBO_ACTIVITY, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}
