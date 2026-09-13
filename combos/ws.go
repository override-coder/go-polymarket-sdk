package combos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
)

const DefaultRFQWSHost = "wss://combos-rfq-gateway-quoter.polymarket.sh"

// DefaultTakerRFQWSHost is separate from the quoter gateway used by WSClient.
const DefaultTakerRFQWSHost = "wss://combos-rfq-gateway-requester.polymarket.sh/ws"

type WSClient struct {
	url    string
	dialer *websocket.Dialer

	conn *websocket.Conn
	mu   sync.Mutex
}

func NewWSClient(host string) *WSClient {
	if strings.TrimSpace(host) == "" {
		host = DefaultRFQWSHost
	}
	return &WSClient{
		url:    buildWSURL(host),
		dialer: websocket.DefaultDialer,
	}
}

// NewTakerWSClient connects to the requester gateway used to place combos.
// It intentionally returns WSClient so authentication, lifecycle, and event
// consumption are identical to the existing maker client.
func NewTakerWSClient(host string) *WSClient {
	if strings.TrimSpace(host) == "" {
		host = DefaultTakerRFQWSHost
	}
	return &WSClient{
		url:    buildTakerWSURL(host),
		dialer: websocket.DefaultDialer,
	}
}

func (c *WSClient) Connect(ctx context.Context, option *sdktypes.AuthOption) error {
	if option == nil {
		return fmt.Errorf("auth option is required")
	}
	if option.ApiKeyCreds == nil {
		return fmt.Errorf("api key creds are required")
	}
	conn, _, err := c.dialer.DialContext(ctx, c.url, http.Header{})
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	maker := option.SingerAddress
	if option.FunderAddress != "" {
		maker = option.FunderAddress
	}
	signer := option.SingerAddress
	if option.SignatureType == model.POLY_1271 {
		signer = maker
	}
	return c.writeJSON(combostypes.AuthMessage{
		Type: combostypes.MessageTypeAuth,
		Auth: combostypes.AuthPayload{
			APIKey:     option.ApiKeyCreds.ApiKey,
			Secret:     option.ApiKeyCreds.Secret,
			Passphrase: option.ApiKeyCreds.Passphrase,
		},
		Identity: combostypes.RFQIdentity{
			SignerAddress: signer,
			MakerAddress:  maker,
			SignatureType: option.SignatureType,
		},
	})
}

func (c *WSClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *WSClient) Read(ctx context.Context) (*Event, error) {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return nil, fmt.Errorf("websocket is not connected")
	}

	type result struct {
		messageType int
		data        []byte
		err         error
	}
	ch := make(chan result, 1)
	go func() {
		mt, data, err := conn.ReadMessage()
		ch <- result{messageType: mt, data: data, err: err}
	}()
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = c.Close()
		case <-done:
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.err != nil {
			return nil, res.err
		}
		if res.messageType != websocket.TextMessage && res.messageType != websocket.BinaryMessage {
			return &Event{Raw: res.data}, nil
		}
		return decodeEvent(res.data)
	}
}

func (c *WSClient) SendQuote(req combostypes.QuoteRequest) error {
	return c.writeJSON(combostypes.WSQuoteRequest{
		Type:        combostypes.MessageTypeRFQQuote,
		RFQID:       req.RFQID,
		PriceE6:     req.PriceE6,
		SizeE6:      req.SizeE6,
		SignedOrder: req.SignedOrder,
	})
}

func (c *WSClient) RequestTakerQuote(req combostypes.TakerQuoteRequest) error {
	req.Type = combostypes.MessageTypeRFQCreate
	return c.writeJSON(req)
}

func (c *WSClient) AcceptTakerQuote(req combostypes.TakerAcceptRequest) error {
	req.Type = combostypes.MessageTypeRFQAccept
	return c.writeJSON(req)
}

func (c *WSClient) CancelQuote(req combostypes.QuoteCancelRequest) error {
	req.Type = combostypes.MessageTypeRFQQuoteCancel
	return c.writeJSON(req)
}

func (c *WSClient) ConfirmQuote(req combostypes.ConfirmationRequest) error {
	req.Type = combostypes.MessageTypeRFQConfirmationResponse
	req.Decision = combostypes.ConfirmationDecisionConfirm
	return c.writeJSON(req)
}

func (c *WSClient) DeclineQuote(req combostypes.ConfirmationRequest) error {
	req.Type = combostypes.MessageTypeRFQConfirmationResponse
	req.Decision = combostypes.ConfirmationDecisionDecline
	return c.writeJSON(req)
}

func (c *WSClient) writeJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("websocket is not connected")
	}
	return c.conn.WriteJSON(v)
}

type Event struct {
	Type                    combostypes.MessageType
	Auth                    *combostypes.AuthResponse
	RFQRequest              *combostypes.RFQRequest
	QuoteAck                *combostypes.QuoteAck
	QuoteCancelAck          *combostypes.QuoteAck
	ConfirmationRequest     *combostypes.ConfirmationRequest
	ConfirmationResponseAck *combostypes.ConfirmationRequest
	ExecutionUpdate         *combostypes.ExecutionUpdate
	TakerQuoteReady         *combostypes.TakerQuoteReady
	StatusUpdate            *combostypes.RFQStatusUpdate
	Error                   *combostypes.RFQError
	Raw                     json.RawMessage
	ReceivedAt              time.Time
}

func decodeEvent(data []byte) (*Event, error) {
	var envelope struct {
		Type combostypes.MessageType `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}

	ev := &Event{
		Type:       envelope.Type,
		Raw:        append(json.RawMessage(nil), data...),
		ReceivedAt: time.Now(),
	}

	switch envelope.Type {
	case combostypes.MessageTypeAuth:
		var msg combostypes.AuthResponse
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.Auth = &msg
	case combostypes.MessageTypeRFQRequest:
		var msg combostypes.RFQRequest
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.RFQRequest = &msg
	case combostypes.MessageTypeAckRFQQuote:
		var msg combostypes.QuoteAck
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.QuoteAck = &msg
	case combostypes.MessageTypeAckRFQQuoteCancel:
		var msg combostypes.QuoteAck
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.QuoteCancelAck = &msg
	case combostypes.MessageTypeRFQConfirmationRequest:
		var msg combostypes.ConfirmationRequest
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.ConfirmationRequest = &msg
	case combostypes.MessageTypeAckRFQConfirmationResponse:
		var msg combostypes.ConfirmationRequest
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.ConfirmationResponseAck = &msg
	case combostypes.MessageTypeRFQExecutionUpdate:
		var msg combostypes.ExecutionUpdate
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.ExecutionUpdate = &msg
	case combostypes.MessageTypeRFQQuoteReady:
		var msg combostypes.TakerQuoteReady
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.TakerQuoteReady = &msg
	case combostypes.MessageTypeRFQStatusUpdate:
		var msg combostypes.RFQStatusUpdate
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.StatusUpdate = &msg
	case combostypes.MessageTypeRFQError:
		var msg combostypes.RFQError
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		ev.Error = &msg
	}

	return ev, nil
}

func buildWSURL(host string) string {
	host = strings.TrimRight(host, "/")
	u, err := url.Parse(host)
	if err != nil {
		return strings.TrimRight(host, "/") + combostypes.RFQWebSocketPath
	}
	if u.Path == "" {
		u.Path = combostypes.RFQWebSocketPath
	}
	return u.String()
}

func buildTakerWSURL(host string) string {
	host = strings.TrimRight(host, "/")
	u, err := url.Parse(host)
	if err != nil {
		return strings.TrimRight(host, "/") + "/ws"
	}
	if u.Path == "" {
		u.Path = "/ws"
	}
	return u.String()
}
