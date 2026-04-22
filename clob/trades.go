package clob

import (
	"context"
	"net/http"
	"time"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/pkg/errors"
)

func (c *Client) GetTrades(ctx context.Context, req clobtypes.GetTradesRequest, option *sdktypes.AuthOption) (*clobtypes.Trades, error) {
	ts := time.Now().Unix()
	l2HeaderArgs := clobtypes.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: clobtypes.GET_TRADES,
	}
	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	params := make(map[string]any, 7)
	if req.ID != nil && *req.ID != "" {
		params["id"] = *req.ID
	}
	if req.MakerAddress != "" {
		params["maker_address"] = req.MakerAddress
	}
	if req.Market != nil && *req.Market != "" {
		params["market"] = *req.Market
	}
	if req.AssetID != nil && *req.AssetID != "" {
		params["asset_id"] = *req.AssetID
	}
	if req.Before != nil && *req.Before != "" {
		params["before"] = *req.Before
	}
	if req.After != nil && *req.After != "" {
		params["after"] = *req.After
	}
	if req.NextCursor != nil && *req.NextCursor != "" {
		params["next_cursor"] = *req.NextCursor
	}

	var resp clobtypes.Trades
	res, err := c.client.DoRequest(ctx, http.MethodGet, clobtypes.GET_TRADES, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get trades")
	}
	return &resp, nil
}
