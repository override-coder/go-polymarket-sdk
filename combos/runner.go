package combos

import (
	"context"
	"fmt"

	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

type QuoteHandler func(ctx context.Context, rfq combostypes.RFQRequest) (*combostypes.QuoteRequest, error)
type LastLookHandler func(ctx context.Context, req combostypes.ConfirmationRequest) (combostypes.ConfirmationDecision, error)
type EventHandler func(ctx context.Context, event Event)
type ErrorHandler func(ctx context.Context, err error)

type MakerRunner struct {
	ws     *WSClient
	option *sdktypes.AuthOption

	onRFQ      QuoteHandler
	onLastLook LastLookHandler
	onEvent    EventHandler
	onError    ErrorHandler
}

func NewMakerRunner(ws *WSClient, option *sdktypes.AuthOption) *MakerRunner {
	return &MakerRunner{
		ws:     ws,
		option: option,
	}
}

func (r *MakerRunner) OnRFQ(handler QuoteHandler) *MakerRunner {
	r.onRFQ = handler
	return r
}

func (r *MakerRunner) OnLastLook(handler LastLookHandler) *MakerRunner {
	r.onLastLook = handler
	return r
}

func (r *MakerRunner) OnEvent(handler EventHandler) *MakerRunner {
	r.onEvent = handler
	return r
}

func (r *MakerRunner) OnError(handler ErrorHandler) *MakerRunner {
	r.onError = handler
	return r
}

func (r *MakerRunner) Run(ctx context.Context) error {
	if r.ws == nil {
		return fmt.Errorf("websocket client is required")
	}
	if err := r.ws.Connect(ctx, r.option); err != nil {
		return err
	}
	defer r.ws.Close()

	authenticated := false
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		event, err := r.ws.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			r.handleError(ctx, err)
			return err
		}
		if r.onEvent != nil {
			r.onEvent(ctx, *event)
		}

		switch {
		case event.Auth != nil:
			if !event.Auth.Success {
				r.handleError(ctx, fmt.Errorf("rfq auth failed: %s", event.Auth.Error))
				continue
			}
			authenticated = true
		case event.RFQRequest != nil:
			if !authenticated {
				continue
			}
			if err := r.handleRFQ(ctx, *event.RFQRequest); err != nil {
				r.handleError(ctx, err)
			}
		case event.ConfirmationRequest != nil:
			if !authenticated {
				continue
			}
			if err := r.handleLastLook(ctx, *event.ConfirmationRequest); err != nil {
				r.handleError(ctx, err)
			}
		case event.Error != nil:
			r.handleError(ctx, fmt.Errorf("rfq error %s: %s", event.Error.Code, firstNonEmpty(event.Error.Error, event.Error.Reason, event.Error.Message)))
		}
	}
}

func (r *MakerRunner) handleRFQ(ctx context.Context, rfq combostypes.RFQRequest) error {
	if r.onRFQ == nil {
		return nil
	}
	quote, err := r.onRFQ(ctx, rfq)
	if err != nil {
		return err
	}
	if quote == nil {
		return nil
	}
	return r.ws.SendQuote(*quote)
}

func (r *MakerRunner) handleLastLook(ctx context.Context, req combostypes.ConfirmationRequest) error {
	if r.onLastLook == nil {
		req.Decision = combostypes.ConfirmationDecisionDecline
		return r.ws.DeclineQuote(req)
	}
	decision, err := r.onLastLook(ctx, req)
	if err != nil {
		return err
	}
	switch decision {
	case combostypes.ConfirmationDecisionConfirm:
		return r.ws.ConfirmQuote(req)
	case combostypes.ConfirmationDecisionDecline:
		return r.ws.DeclineQuote(req)
	default:
		return fmt.Errorf("unsupported last-look decision: %s", decision)
	}
}

func (r *MakerRunner) handleError(ctx context.Context, err error) {
	if r.onError != nil {
		r.onError(ctx, err)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
