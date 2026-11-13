package types

import datatypes "github.com/override-coder/go-polymarket-sdk/dataapi/types"

const (
	GetComboMarkets = "/v1/rfq/combo-markets"

	SubmitQuote  = "/v1/maker/quotes"
	CancelQuote  = "/v1/maker/quotes/cancel"
	ConfirmQuote = "/v1/maker/confirmations"

	GetComboPositions = datatypes.GET_COMBO_POSITIONS
	GetComboActivity  = datatypes.GET_COMBO_ACTIVITY

	CreateBuilderRFQ = "/v1/builder/rfq/requests"

	RFQWebSocketPath = "/ws/rfq"
)
