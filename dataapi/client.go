package dataapi

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
)

type Client struct {
	client     *http2.Client
	clobClient *http2.Client

	chainId *big.Int
}

const defaultClobHost = "https://clob.polymarket.com"

func NewClient(host string, chainId *big.Int) *Client {
	return NewClientWithClobHost(host, defaultClobHost, chainId)
}

func NewClientWithClobHost(host, clobHost string, chainId *big.Int) *Client {
	if strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	if strings.HasSuffix(clobHost, "/") {
		clobHost = clobHost[:len(clobHost)-1]
	}
	return &Client{
		client:     http2.NewClient(host),
		clobClient: http2.NewClient(clobHost),
		chainId:    chainId,
	}
}

// GetPositions returns one V2 page. Follow Pagination.NextCursor and keep the
// same user/condition anchor and filters; a short page is not an end marker.
func (c *Client) GetPositions(ctx context.Context, q types.PositionsQuery) (*types.Page[types.Position], error) {
	user := strings.TrimSpace(q.User)
	if user == "" && len(q.Condition) == 0 {
		return nil, fmt.Errorf("user or condition is required")
	}
	if err := validateConditions(q.Condition); err != nil {
		return nil, err
	}
	if user == "" && distinctCount(q.Condition) != 1 {
		return nil, fmt.Errorf("market-anchored positions require exactly one condition")
	}
	if user == "" && len(q.EventID) > 0 {
		return nil, fmt.Errorf("event_id requires user")
	}
	if err := validateEventIDs(q.EventID); err != nil {
		return nil, err
	}
	if q.Title != nil && utf8.RuneCountInString(*q.Title) > 200 {
		return nil, fmt.Errorf("title too long (max 200 characters)")
	}
	params := map[string]any{}
	if user != "" {
		params["user"] = user
	}
	if len(q.Condition) > 0 {
		params["condition"] = strings.Join(q.Condition, ",")
	}
	if len(q.EventID) > 0 {
		params["event_id"] = joinEventIDs(q.EventID)
	}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	addParam(params, "status", q.Status)
	addParam(params, "title", q.Title)
	addParam(params, "filter_type", q.FilterType)
	addParam(params, "filter_amount", q.FilterAmount)
	addParam(params, "include_archived", q.IncludeArchived)
	addParam(params, "sort_by", q.SortBy)
	addParam(params, "start", q.Start)
	addParam(params, "end", q.End)
	addParam(params, "sort_direction", q.SortDirection)
	var out types.Page[types.Position]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_POSITIONS, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

// GetClosedPositions is GetPositions with status=CLOSED. The service chooses
// REALIZED_PNL DESC by default; callers may explicitly select another V2 sort.
func (c *Client) GetClosedPositions(ctx context.Context, q types.ClosedPositionsQuery) (*types.Page[types.ClosedPosition], error) {
	if q.Status != nil && *q.Status != types.PositionCLOSED {
		return nil, fmt.Errorf("closed positions require status=CLOSED")
	}
	status := types.PositionCLOSED
	q.Status = &status
	return c.GetPositions(ctx, q)
}

// GetUserActivity returns one V2 activity page. Keep all filters unchanged
// when following NextCursor; the cursor alone does not bind the feed filters.
func (c *Client) GetUserActivity(ctx context.Context, q types.ActivityQuery) (*types.Page[types.UserActivity], error) {
	user := strings.TrimSpace(q.User)
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	if len(q.Condition) > 0 && len(q.EventID) > 0 {
		return nil, fmt.Errorf("condition and event_id are mutually exclusive")
	}
	if err := validateConditions(q.Condition); err != nil {
		return nil, err
	}
	if err := validateEventIDs(q.EventID); err != nil {
		return nil, err
	}
	if q.SortBy != nil && *q.SortBy != types.ActivitySortTIMESTAMP {
		return nil, fmt.Errorf("activity only supports sort_by=TIMESTAMP")
	}
	params := map[string]any{"user": user}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	if len(q.Condition) > 0 {
		params["condition"] = strings.Join(q.Condition, ",")
	}
	if len(q.EventID) > 0 {
		params["event_id"] = joinEventIDs(q.EventID)
	}
	if len(q.Type) > 0 {
		values := make([]string, len(q.Type))
		for i, t := range q.Type {
			values[i] = string(t)
		}
		params["type"] = strings.Join(values, ",")
	}
	addParam(params, "exclude_deposits_withdrawals", q.ExcludeDepositsWithdrawals)
	addParam(params, "start", q.Start)
	addParam(params, "end", q.End)
	addParam(params, "sort_by", q.SortBy)
	addParam(params, "sort_direction", q.SortDirection)
	addParam(params, "side", q.Side)
	var out types.Page[types.UserActivity]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_Activity, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

// GetPositionValue returns single-market mark value plus unresolved combos at
// cost basis. A condition filter excludes combos. An empty wallet has value 0.
func (c *Client) GetPositionValue(ctx context.Context, q types.PositionValueQuery) (*types.Envelope[types.PositionValue], error) {
	user := strings.TrimSpace(q.User)
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	if err := validateConditions(q.Condition); err != nil {
		return nil, err
	}
	params := map[string]any{"user": user}
	if len(q.Condition) > 0 {
		params["condition"] = strings.Join(q.Condition, ",")
	}
	var out types.Envelope[types.PositionValue]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_VALUE, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

// GetTraderLeaderboardRankings returns one page of the V2 board. A cursor pins
// the board; leave omitted parameters unset to resume without overriding it.
func (c *Client) GetTraderLeaderboardRankings(ctx context.Context, q types.TraderLeaderboardQuery) (*types.Page[types.TraderLeaderboard], error) {
	params := map[string]any{}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	addParam(params, "category", q.Category)
	addParam(params, "time_period", q.TimePeriod)
	addParam(params, "sort_by", q.SortBy)
	var out types.Page[types.TraderLeaderboard]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_LEADERBOARD, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

// GetTraderLeaderboardUser reads the user branch of /v2/leaderboard. Data is
// nil for an unknown user; either rank may be nil for an existing user.
func (c *Client) GetTraderLeaderboardUser(ctx context.Context, q types.TraderLeaderboardUserQuery) (*types.Envelope[*types.TraderLeaderboardUser], error) {
	user := strings.TrimSpace(q.User)
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	params := map[string]any{"user": user}
	addParam(params, "category", q.Category)
	addParam(params, "time_period", q.TimePeriod)
	var out types.Envelope[*types.TraderLeaderboardUser]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_LEADERBOARD, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

// GetTopHoldersForMarkets returns a page per outcome token. Merge pages by
// TokenID, not array position: exhausted token groups disappear on later pages.
func (c *Client) GetTopHoldersForMarkets(ctx context.Context, q types.TopHoldersQuery) (*types.Page[types.TopHoldersForMarket], error) {
	if len(q.Condition) == 0 {
		return nil, fmt.Errorf("condition is required")
	}
	if err := validateConditions(q.Condition); err != nil {
		return nil, err
	}
	if q.IncludePnL != nil && *q.IncludePnL {
		if distinctCount(q.Condition) != 1 {
			return nil, fmt.Errorf("include_pnl requires exactly one condition")
		}
		if (q.Cursor == nil || *q.Cursor == "") && q.Limit != nil && *q.Limit > 100 {
			return nil, fmt.Errorf("include_pnl limit must be <= 100")
		}
	}
	params := map[string]any{"condition": strings.Join(q.Condition, ",")}
	if err := addPageParams(params, q.Limit, q.Cursor); err != nil {
		return nil, err
	}
	addParam(params, "min_balance", q.MinBalance)
	addParam(params, "include_pnl", q.IncludePnL)
	var out types.Page[types.TopHoldersForMarket]
	resp, err := c.client.DoRequest(ctx, http.MethodGet, types.GET_HOLDERS, &http2.RequestOptions{Params: params}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func (c *Client) GetMarketByToken(ctx context.Context, tokenID string) (*types.MarketByTokenResponse, error) {
	tokenID = strings.TrimSpace(tokenID)
	if tokenID == "" {
		return nil, fmt.Errorf("tokenID is required")
	}

	var out types.MarketByTokenResponse
	resp, err := c.clobClient.DoRequest(ctx, http.MethodGet, types.GET_MARKET_BY_TOKEN+url.PathEscape(tokenID), nil, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}
