package combos

import (
	"context"
	"fmt"
	"strings"

	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

// TakerRequest is the input needed to request and accept one combo quote.
type TakerRequest struct {
	LegPositionIDs []string
	Direction      combostypes.Direction
	Side           combostypes.ComboSide
	RequestedSize  combostypes.RequestedSize
}

type TakerRunner struct {
	ws      *WSClient
	builder *OrderBuilder
	option  *sdktypes.AuthOption

	// MaxRetries is the number of re-quotes after an expired quote. Zero means 2.
	MaxRetries int
	onEvent    EventHandler
}

func NewTakerRunner(ws *WSClient, builder *OrderBuilder, option *sdktypes.AuthOption) *TakerRunner {
	return &TakerRunner{ws: ws, builder: builder, option: option, MaxRetries: 2}
}

func (r *TakerRunner) OnEvent(handler EventHandler) *TakerRunner {
	r.onEvent = handler
	return r
}

// Place performs auth -> RFQ_CREATE -> RFQ_QUOTE_READY -> RFQ_ACCEPT and
// returns after the gateway confirms settlement. Expired quotes are re-requested
// on the same authenticated socket, avoiding a second concurrent wallet prompt.
func (r *TakerRunner) Place(ctx context.Context, request TakerRequest) (*combostypes.ExecutionUpdate, error) {
	if r.ws == nil {
		return nil, fmt.Errorf("websocket client is required")
	}
	if r.builder == nil {
		return nil, fmt.Errorf("order builder is required")
	}
	if r.option == nil {
		return nil, fmt.Errorf("auth option is required")
	}
	if len(request.LegPositionIDs) < 2 {
		return nil, fmt.Errorf("at least two leg position ids are required")
	}
	if request.Direction != combostypes.DirectionBuy && request.Direction != combostypes.DirectionSell {
		return nil, fmt.Errorf("unsupported direction: %s", request.Direction)
	}
	if request.Side != combostypes.ComboSideYes && request.Side != combostypes.ComboSideNo {
		return nil, fmt.Errorf("unsupported combo side: %s", request.Side)
	}
	if strings.TrimSpace(request.RequestedSize.Unit) == "" || strings.TrimSpace(request.RequestedSize.ValueE6) == "" {
		return nil, fmt.Errorf("requested size is required")
	}

	if err := r.ws.Connect(ctx, r.option); err != nil {
		return nil, err
	}
	defer r.ws.Close()

	maxRetries := r.MaxRetries
	if maxRetries == 0 {
		maxRetries = 2
	}
	attempts := 0
	authenticated := false
	requestQuote := func() error {
		attempts++
		return r.ws.RequestTakerQuote(combostypes.TakerQuoteRequest{
			LegPositionIDs: request.LegPositionIDs,
			Direction:      request.Direction,
			Side:           request.Side,
			RequestedSize:  request.RequestedSize,
		})
	}

	for {
		event, err := r.ws.Read(ctx)
		if err != nil {
			return nil, err
		}
		if r.onEvent != nil {
			r.onEvent(ctx, *event)
		}

		switch {
		case event.Auth != nil:
			if !event.Auth.Success {
				return nil, fmt.Errorf("rfq auth failed: %s", event.Auth.Error)
			}
			if !authenticated {
				authenticated = true
				if err := requestQuote(); err != nil {
					return nil, err
				}
			}
		case event.TakerQuoteReady != nil:
			ready := event.TakerQuoteReady
			tokenID := ready.Request.YesPositionID
			if request.Side == combostypes.ComboSideNo {
				tokenID = ready.Request.NoPositionID
			}
			order, err := r.builder.BuildTakerOrder(BuildTakerOrderInput{
				TokenID:       tokenID,
				MakerAmountE6: ready.Quote.MakerAmountE6,
				TakerAmountE6: ready.Quote.TakerAmountE6,
				Direction:     request.Direction,
			}, r.option)
			if err != nil {
				return nil, err
			}
			if err := r.ws.AcceptTakerQuote(combostypes.TakerAcceptRequest{
				RFQID:       ready.Request.RFQID,
				QuoteID:     ready.Quote.QuoteID,
				SignedOrder: order,
			}); err != nil {
				return nil, err
			}
		case event.StatusUpdate != nil:
			if !isExpired(event.StatusUpdate.Status, event.StatusUpdate.Code) {
				return nil, fmt.Errorf("rfq status %s: %s", event.StatusUpdate.Status, event.StatusUpdate.Message)
			}
			if attempts > maxRetries {
				return nil, fmt.Errorf("rfq quote expired after %d attempts", attempts)
			}
			if err := requestQuote(); err != nil {
				return nil, err
			}
		case event.ExecutionUpdate != nil:
			if event.ExecutionUpdate.Status == "CONFIRMED" {
				return event.ExecutionUpdate, nil
			}
			if event.ExecutionUpdate.Status == "FAILED" || event.ExecutionUpdate.Status == "EXPIRED" {
				return nil, fmt.Errorf("rfq execution %s", event.ExecutionUpdate.Status)
			}
		case event.Error != nil:
			if isExpired(event.Error.Code, event.Error.Reason) && attempts <= maxRetries {
				if err := requestQuote(); err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("rfq error %s: %s", event.Error.Code, firstNonEmpty(event.Error.Error, event.Error.Reason, event.Error.Message))
		}
	}
}

func isExpired(values ...string) bool {
	for _, value := range values {
		if strings.EqualFold(value, "EXPIRED") || strings.Contains(strings.ToUpper(value), "EXPIRED") {
			return true
		}
	}
	return false
}
