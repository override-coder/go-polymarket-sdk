package types

const (
	GetComboMarkets = "/v1/rfq/combo-markets"

	SubmitQuote  = "/v1/maker/quotes"
	CancelQuote  = "/v1/maker/quotes/cancel"
	ConfirmQuote = "/v1/maker/confirmations"

	GetComboPositions = "/v1/positions/combos"
	GetComboActivity  = "/v1/activity/combos"

	RFQWebSocketPath = "/ws/rfq"
)
