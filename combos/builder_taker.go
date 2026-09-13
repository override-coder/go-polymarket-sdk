package combos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

// DefaultBuilderTakerGatewayHost is the production Builder Gateway configured
// by Polymarket's official TypeScript client.
const DefaultBuilderTakerGatewayHost = "https://combos-rfq-gateway-builder.polymarket.com"

// BuilderTakerClient implements the official Builder Gateway Combo requester
// API. It is independent from TakerRunner, which retains compatibility with
// the existing requester WebSocket gateway.
type BuilderTakerClient struct {
	host    string
	client  *http2.Client
	builder *OrderBuilder
}

// NewBuilderTakerClient creates a client for the Builder Gateway host supplied
// during Polymarket builder onboarding. An empty host uses the official
// production gateway. The host must not include /v1/builder/rfq, which is
// added by the client.
func NewBuilderTakerClient(host string, chainID *big.Int, signFn signing.SignatureFunc) *BuilderTakerClient {
	host = strings.TrimRight(strings.TrimSpace(host), "/")
	if host == "" {
		host = DefaultBuilderTakerGatewayHost
	}
	return &BuilderTakerClient{
		host:    host,
		client:  http2.NewClient(host),
		builder: NewOrderBuilder(chainID, signFn),
	}
}

// CreateRFQ requests the best executable Combo quote. It returns terminal
// business outcomes (for example status FAILED with no quote) without an error.
func (c *BuilderTakerClient) CreateRFQ(ctx context.Context, request TakerRequest, option *sdktypes.AuthOption) (*combostypes.BuilderTakerRFQ, error) {
	body, err := builderTakerRequest(request, option)
	if err != nil {
		return nil, err
	}

	var out combostypes.BuilderTakerRFQ
	if err := c.do(ctx, http.MethodPost, combostypes.CreateBuilderRFQ, body, option, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AcceptRFQ accepts a previously returned quote. Use GetRFQStatus or
// WaitForTerminalStatus afterwards because EXECUTING is not a confirmed fill.
func (c *BuilderTakerClient) AcceptRFQ(ctx context.Context, rfqID string, request combostypes.BuilderTakerAcceptRequest, option *sdktypes.AuthOption) (*combostypes.BuilderTakerStatus, error) {
	if strings.TrimSpace(rfqID) == "" {
		return nil, fmt.Errorf("rfq id is required")
	}
	if strings.TrimSpace(request.QuoteID) == "" {
		return nil, fmt.Errorf("quote id is required")
	}

	var out combostypes.BuilderTakerStatus
	path := combostypes.CreateBuilderRFQ + "/" + rfqID + "/accept"
	if err := c.do(ctx, http.MethodPost, path, request, option, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRFQStatus reads the durable status after acceptance. It deliberately uses
// only account authentication, as required by Builder Gateway.
func (c *BuilderTakerClient) GetRFQStatus(ctx context.Context, rfqID string, option *sdktypes.AuthOption) (*combostypes.BuilderTakerStatus, error) {
	if strings.TrimSpace(rfqID) == "" {
		return nil, fmt.Errorf("rfq id is required")
	}

	var out combostypes.BuilderTakerStatus
	path := combostypes.CreateBuilderRFQ + "/" + rfqID
	if err := c.do(ctx, http.MethodGet, path, nil, option, false, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Place creates an RFQ, builds the official requester order from the returned
// quote, and submits acceptance. It does not wait for a fill; call
// WaitForTerminalStatus with the returned RFQ ID to track execution.
func (c *BuilderTakerClient) Place(ctx context.Context, request TakerRequest, option *sdktypes.AuthOption) (*combostypes.BuilderTakerStatus, error) {
	rfq, err := c.CreateRFQ(ctx, request, option)
	if err != nil {
		return nil, err
	}
	if rfq.Quote == nil {
		return &combostypes.BuilderTakerStatus{RFQID: rfq.RFQID, Status: rfq.Status, Error: rfq.Error}, nil
	}
	if strings.TrimSpace(rfq.BuilderCode) == "" {
		return nil, fmt.Errorf("builder gateway response missing builder_code")
	}
	if rfq.Request.Side != combostypes.ComboSideYes {
		return nil, fmt.Errorf("unsupported builder gateway combo side: %s", rfq.Request.Side)
	}

	order, err := c.builder.BuildTakerOrder(BuildTakerOrderInput{
		TokenID:       rfq.Request.YesPositionID,
		MakerAmountE6: rfq.Quote.MakerAmountE6,
		TakerAmountE6: rfq.Quote.TakerAmountE6,
		Direction:     rfq.Request.Direction,
		Timestamp:     time.Now().Unix(),
		Builder:       rfq.BuilderCode,
	}, option)
	if err != nil {
		return nil, err
	}

	return c.AcceptRFQ(ctx, rfq.RFQID, combostypes.BuilderTakerAcceptRequest{
		QuoteID:     rfq.Quote.QuoteID,
		SignedOrder: order,
	}, option)
}

// PlaceAndWait creates and accepts a Combo RFQ, then waits for its terminal
// status. It is the convenience entry point for callers that need the same
// execution result shape as a normal order flow, including TxHash when the
// Combo Gateway confirms settlement. The context controls the total wait.
//
// A terminal no-quote outcome is returned without polling. For an accepted
// quote, pollInterval defaults to one second when zero or negative.
func (c *BuilderTakerClient) PlaceAndWait(ctx context.Context, request TakerRequest, option *sdktypes.AuthOption, pollInterval time.Duration) (*combostypes.BuilderTakerStatus, error) {
	status, err := c.Place(ctx, request, option)
	if err != nil || status == nil || isBuilderTakerTerminal(status.Status) {
		return status, err
	}
	terminal, err := c.WaitForTerminalStatus(ctx, status.RFQID, pollInterval, option)
	if terminal == nil {
		return nil, err
	}
	if terminal.RFQID == "" {
		terminal.RFQID = status.RFQID
	}
	if terminal.TakerOrderHash == "" {
		terminal.TakerOrderHash = status.TakerOrderHash
	}
	return terminal, err
}

// WaitForTerminalStatus polls until Builder Gateway reports a terminal state
// or ctx expires. pollInterval defaults to one second when zero or negative.
func (c *BuilderTakerClient) WaitForTerminalStatus(ctx context.Context, rfqID string, pollInterval time.Duration, option *sdktypes.AuthOption) (*combostypes.BuilderTakerStatus, error) {
	if pollInterval <= 0 {
		pollInterval = time.Second
	}
	for {
		status, err := c.GetRFQStatus(ctx, rfqID, option)
		if err != nil || isBuilderTakerTerminal(status.Status) {
			return status, err
		}

		timer := time.NewTimer(pollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func builderTakerRequest(request TakerRequest, option *sdktypes.AuthOption) (combostypes.BuilderTakerRequest, error) {
	if option == nil {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("auth option is required")
	}
	if len(request.LegPositionIDs) < 2 || len(request.LegPositionIDs) > 50 {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("leg position ids must contain 2 to 50 entries")
	}
	unique := make(map[string]struct{}, len(request.LegPositionIDs))
	for _, id := range request.LegPositionIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			return combostypes.BuilderTakerRequest{}, fmt.Errorf("leg position ids must not contain an empty value")
		}
		if _, exists := unique[id]; exists {
			return combostypes.BuilderTakerRequest{}, fmt.Errorf("leg position ids must be unique")
		}
		unique[id] = struct{}{}
	}
	if request.Direction != combostypes.DirectionBuy && request.Direction != combostypes.DirectionSell {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("unsupported direction: %s", request.Direction)
	}
	if request.Side == "" {
		request.Side = combostypes.ComboSideYes
	}
	if request.Side != combostypes.ComboSideYes {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("builder gateway currently supports only YES combo side")
	}
	if request.Direction == combostypes.DirectionBuy && request.RequestedSize.Unit != "notional" {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("BUY requested size unit must be notional")
	}
	if request.Direction == combostypes.DirectionSell && request.RequestedSize.Unit != "shares" {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("SELL requested size unit must be shares")
	}
	if _, err := parsePositiveE6(request.RequestedSize.ValueE6, "requested_size.value_e6"); err != nil {
		return combostypes.BuilderTakerRequest{}, err
	}
	if strings.TrimSpace(option.SingerAddress) == "" {
		return combostypes.BuilderTakerRequest{}, fmt.Errorf("signer address is required")
	}

	maker := option.SingerAddress
	if option.FunderAddress != "" {
		maker = option.FunderAddress
	}
	return combostypes.BuilderTakerRequest{
		SignerAddress:  comboRFQSigner(option.SignatureType, option.SingerAddress, maker),
		MakerAddress:   maker,
		SignatureType:  option.SignatureType,
		LegPositionIDs: request.LegPositionIDs,
		Direction:      request.Direction,
		Side:           request.Side,
		RequestedSize:  request.RequestedSize,
	}, nil
}

func (c *BuilderTakerClient) do(ctx context.Context, method, path string, body any, option *sdktypes.AuthOption, needsBuilderAuth bool, out any) error {
	if c == nil || c.client == nil || c.host == "" {
		return fmt.Errorf("builder taker client is required")
	}
	if option == nil || option.ApiKeyCreds == nil {
		return fmt.Errorf("account api key creds are required")
	}

	bodyString := ""
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyString = string(bodyBytes)
	}
	ts := time.Now().Unix()
	args := clobtypes.L2HeaderArgs{Method: method, RequestPath: path, Body: bodyString}
	headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, args, &ts)
	if err != nil {
		return err
	}
	if needsBuilderAuth {
		if option.BuilderApiKeyCreds == nil {
			return fmt.Errorf("builder api key creds are required")
		}
		builderHeaders, err := sdkheaders.CreateL2BuilderHeaders(option.BuilderApiKeyCreds, args, &ts)
		if err != nil {
			return err
		}
		headers = sdkheaders.InjectBuilderHeaders(headers, builderHeaders)
	}

	resp, err := c.client.DoRequest(ctx, method, path, &http2.RequestOptions{Headers: headers, Data: bodyString}, out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return e
	}
	return nil
}

func isBuilderTakerTerminal(status string) bool {
	switch strings.ToUpper(status) {
	case "CONFIRMED", "FILLED", "FAILED", "EXPIRED", "CANCELED":
		return true
	default:
		return false
	}
}
