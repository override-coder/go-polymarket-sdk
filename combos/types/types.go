package types

import "github.com/polymarket/go-order-utils/pkg/model"

const (
	Bytes32Zero = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type Direction string

const (
	DirectionBuy  Direction = "BUY"
	DirectionSell Direction = "SELL"
)

type ComboSide string

const (
	ComboSideYes ComboSide = "YES"
)

type QuoteSource string

const (
	QuoteSourceCollateral QuoteSource = "collateral"
	QuoteSourceInventory  QuoteSource = "inventory"
)

type ConfirmationDecision string

const (
	ConfirmationDecisionConfirm ConfirmationDecision = "CONFIRM"
	ConfirmationDecisionDecline ConfirmationDecision = "DECLINE"
)

type MessageType string

const (
	MessageTypeAuth                       MessageType = "auth"
	MessageTypeRFQRequest                 MessageType = "RFQ_REQUEST"
	MessageTypeRFQQuote                   MessageType = "RFQ_QUOTE"
	MessageTypeAckRFQQuote                MessageType = "ACK_RFQ_QUOTE"
	MessageTypeRFQQuoteCancel             MessageType = "RFQ_QUOTE_CANCEL"
	MessageTypeAckRFQQuoteCancel          MessageType = "ACK_RFQ_QUOTE_CANCEL"
	MessageTypeRFQConfirmationRequest     MessageType = "RFQ_CONFIRMATION_REQUEST"
	MessageTypeRFQConfirmationResponse    MessageType = "RFQ_CONFIRMATION_RESPONSE"
	MessageTypeAckRFQConfirmationResponse MessageType = "ACK_RFQ_CONFIRMATION_RESPONSE"
	MessageTypeRFQExecutionUpdate         MessageType = "RFQ_EXECUTION_UPDATE"
	MessageTypeRFQError                   MessageType = "RFQ_ERROR"
)

type AuthMessage struct {
	Type     MessageType `json:"type"`
	Auth     AuthPayload `json:"auth"`
	Identity RFQIdentity `json:"identity"`
}

type AuthPayload struct {
	APIKey     string `json:"apiKey"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

type RFQIdentity struct {
	SignerAddress string              `json:"signer_address"`
	MakerAddress  string              `json:"maker_address"`
	SignatureType model.SignatureType `json:"signature_type"`
}

type AuthResponse struct {
	Type    MessageType `json:"type"`
	Success bool        `json:"success"`
	Address string      `json:"address"`
	Error   string      `json:"error,omitempty"`
}

type ComboMarketsQuery struct {
	Limit   *int
	Cursor  *string
	Exclude []string
}

type ComboMarketsResponse struct {
	Markets    []ComboMarket `json:"markets"`
	NextCursor string        `json:"next_cursor"`
}

type ComboMarket struct {
	ID            string   `json:"id"`
	ConditionID   string   `json:"condition_id"`
	PositionIDs   []string `json:"position_ids"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	Outcomes      []string `json:"outcomes"`
	OutcomePrices []string `json:"outcome_prices"`
	Image         string   `json:"image"`
	Volume        float64  `json:"volume"`
	Tags          []string `json:"tags"`
}

type RequestedSize struct {
	Unit    string `json:"unit"`
	ValueE6 string `json:"value_e6"`
}

type RFQRequest struct {
	Type               MessageType   `json:"type"`
	RFQID              string        `json:"rfq_id"`
	RequestorPublicID  string        `json:"requestor_public_id"`
	LegPositionIDs     []string      `json:"leg_position_ids"`
	ConditionID        string        `json:"condition_id"`
	YesPositionID      string        `json:"yes_position_id"`
	NoPositionID       string        `json:"no_position_id"`
	Direction          Direction     `json:"direction"`
	Side               ComboSide     `json:"side"`
	RequestedSize      RequestedSize `json:"requested_size"`
	SubmissionDeadline int64         `json:"submission_deadline"`
}

type SignedOrderV2 struct {
	Salt          string              `json:"salt"`
	Maker         string              `json:"maker"`
	Signer        string              `json:"signer"`
	TokenID       string              `json:"tokenId"`
	MakerAmount   string              `json:"makerAmount"`
	TakerAmount   string              `json:"takerAmount"`
	Side          int                 `json:"side"`
	SignatureType model.SignatureType `json:"signatureType"`
	Timestamp     string              `json:"timestamp"`
	Expiration    string              `json:"expiration,omitempty"`
	Metadata      string              `json:"metadata"`
	Builder       string              `json:"builder"`
	Signature     string              `json:"signature"`
}

type QuoteRequest struct {
	Type          MessageType         `json:"type,omitempty"`
	QuoteID       string              `json:"quote_id"`
	RFQID         string              `json:"rfq_id"`
	SignerAddress string              `json:"signer_address"`
	MakerAddress  string              `json:"maker_address"`
	SignatureType model.SignatureType `json:"signature_type"`
	PriceE6       string              `json:"price_e6"`
	SizeE6        string              `json:"size_e6"`
	Source        QuoteSource         `json:"source,omitempty"`
	SignedOrder   SignedOrderV2       `json:"signed_order"`
	ValidUntil    int64               `json:"valid_until,omitempty"`
}

type WSQuoteRequest struct {
	Type        MessageType   `json:"type"`
	RFQID       string        `json:"rfq_id"`
	PriceE6     string        `json:"price_e6"`
	SizeE6      string        `json:"size_e6"`
	SignedOrder SignedOrderV2 `json:"signed_order"`
}

type QuoteCancelRequest struct {
	Type          MessageType         `json:"type,omitempty"`
	RFQID         string              `json:"rfq_id"`
	QuoteID       string              `json:"quote_id"`
	SignerAddress string              `json:"signer_address"`
	MakerAddress  string              `json:"maker_address"`
	SignatureType model.SignatureType `json:"signature_type"`
}

type ConfirmationRequest struct {
	Type           MessageType          `json:"type,omitempty"`
	RFQID          string               `json:"rfq_id"`
	QuoteID        string               `json:"quote_id"`
	SignerAddress  string               `json:"signer_address"`
	MakerAddress   string               `json:"maker_address"`
	SignatureType  model.SignatureType  `json:"signature_type"`
	Decision       ConfirmationDecision `json:"decision,omitempty"`
	LegPositionIDs []string             `json:"leg_position_ids,omitempty"`
	ConditionID    string               `json:"condition_id,omitempty"`
	YesPositionID  string               `json:"yes_position_id,omitempty"`
	NoPositionID   string               `json:"no_position_id,omitempty"`
	Direction      Direction            `json:"direction,omitempty"`
	Side           ComboSide            `json:"side,omitempty"`
	FillSizeE6     string               `json:"fill_size_e6,omitempty"`
	PriceE6        string               `json:"price_e6,omitempty"`
	ConfirmBy      int64                `json:"confirm_by,omitempty"`
}

type ConfirmationResponse struct {
	Snapshot  *RFQSnapshot `json:"snapshot,omitempty"`
	Execution any          `json:"execution,omitempty"`
}

type QuoteAck struct {
	Type     MessageType  `json:"type"`
	RFQID    string       `json:"rfq_id"`
	QuoteID  string       `json:"quote_id"`
	Snapshot *RFQSnapshot `json:"snapshot,omitempty"`
}

type ExecutionUpdate struct {
	Type      MessageType  `json:"type"`
	RFQID     string       `json:"rfq_id"`
	QuoteID   string       `json:"quote_id,omitempty"`
	Status    string       `json:"status"`
	TxHash    string       `json:"tx_hash,omitempty"`
	Snapshot  *RFQSnapshot `json:"snapshot,omitempty"`
	Execution any          `json:"execution,omitempty"`
}

type RFQError struct {
	Type        MessageType `json:"type"`
	RequestType string      `json:"request_type,omitempty"`
	RFQID       string      `json:"rfq_id,omitempty"`
	QuoteID     string      `json:"quote_id,omitempty"`
	Code        string      `json:"code,omitempty"`
	Reason      string      `json:"reason,omitempty"`
	Message     string      `json:"message,omitempty"`
	Error       string      `json:"error,omitempty"`
}

type RFQSnapshot struct {
	Request                  *RFQStoredRequest   `json:"request,omitempty"`
	Status                   string              `json:"status"`
	CompetitionStartedAt     int64               `json:"competition_started_at,omitempty"`
	CompetitionEndsAt        int64               `json:"competition_ends_at,omitempty"`
	ConfirmationStartedAt    int64               `json:"confirmation_started_at,omitempty"`
	ConfirmationEndsAt       int64               `json:"confirmation_ends_at,omitempty"`
	QuoteID                  string              `json:"quote_id,omitempty"`
	Bundle                   *RFQBundle          `json:"bundle,omitempty"`
	MakerConfirmations       []MakerConfirmation `json:"maker_confirmations,omitempty"`
	SelectedBundle           any                 `json:"selected_bundle,omitempty"`
	CompetitionWindowEndsAt  int64               `json:"competition_window_ends_at,omitempty"`
	ConfirmationWindowEndsAt int64               `json:"confirmation_window_ends_at,omitempty"`
}

type RFQStoredRequest struct {
	RFQID             string              `json:"rfq_id"`
	LegPositionIDs    []string            `json:"leg_position_ids"`
	AuthAddress       string              `json:"auth_address,omitempty"`
	SignerAddress     string              `json:"signer_address,omitempty"`
	MakerAddress      string              `json:"maker_address,omitempty"`
	SignatureType     model.SignatureType `json:"signature_type,omitempty"`
	RequestorPublicID string              `json:"requestor_public_id,omitempty"`
	ConditionID       string              `json:"condition_id,omitempty"`
	YesPositionID     string              `json:"yes_position_id,omitempty"`
	NoPositionID      string              `json:"no_position_id,omitempty"`
	RequestedSize     RequestedSize       `json:"requested_size"`
	CreatedAt         int64               `json:"created_at,omitempty"`
}

type RFQBundle struct {
	RequestedSharesE6   string          `json:"requested_shares_e6"`
	BlendedPriceE6      string          `json:"blended_price_e6"`
	Allocations         []RFQAllocation `json:"allocations"`
	RequestedNotionalE6 string          `json:"requested_notional_e6"`
}

type RFQAllocation struct {
	MakerQuoteID  string `json:"maker_quote_id"`
	SignerAddress string `json:"signer_address"`
	MakerAddress  string `json:"maker_address"`
	SizeE6        string `json:"size_e6"`
	PriceE6       string `json:"price_e6"`
	ReceivedAt    int64  `json:"received_at"`
}

type MakerConfirmation struct {
	QuoteID       string `json:"quote_id"`
	SignerAddress string `json:"signer_address"`
	MakerAddress  string `json:"maker_address"`
	Reason        string `json:"reason,omitempty"`
	RespondedAt   int64  `json:"responded_at,omitempty"`
}

type ComboPositionsQuery struct {
	User     string
	Status   *string
	Sort     *string
	MarketID []string
	Limit    *int
	Offset   *int
}

type ComboPositionsResponse struct {
	Combos     []ComboPosition `json:"combos"`
	Pagination Pagination      `json:"pagination"`
}

type ComboPosition struct {
	ComboConditionID   string     `json:"combo_condition_id"`
	ComboPositionID    string     `json:"combo_position_id"`
	ModuleID           string     `json:"module_id"`
	UserAddress        string     `json:"user_address"`
	SharesBalance      string     `json:"shares_balance"`
	EntryAvgPriceUSDC  string     `json:"entry_avg_price_usdc"`
	EntryCostUSDC      string     `json:"entry_cost_usdc"`
	RealizedPayoutUSDC string     `json:"realized_payout_usdc"`
	TotalCostUSDC      string     `json:"total_cost_usdc"`
	Legs               []ComboLeg `json:"legs"`
}

type ComboActivityQuery struct {
	User     string
	MarketID []string
	Limit    *int
	Offset   *int
}

type ComboActivityResponse struct {
	Activity   []ComboActivity `json:"activity"`
	Pagination Pagination      `json:"pagination"`
}

type ComboActivity struct {
	ID               string     `json:"id"`
	EventKind        string     `json:"event_kind"`
	ModuleKind       string     `json:"module_kind"`
	UserAddress      string     `json:"user_address"`
	ComboConditionID string     `json:"combo_condition_id"`
	ComboPositionID  string     `json:"combo_position_id"`
	AmountUSDC       string     `json:"amount_usdc"`
	PayoutUSDC       string     `json:"payout_usdc"`
	Timestamp        int64      `json:"timestamp"`
	TxHash           string     `json:"tx_hash"`
	Legs             []ComboLeg `json:"legs"`
}

type ComboLeg struct {
	ConditionID string `json:"condition_id"`
	PositionID  string `json:"position_id"`
	Outcome     string `json:"outcome"`
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
