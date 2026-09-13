package combos

import (
	"math/big"
	"testing"

	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTakerOrderUsesGatewayAmountsAndRequesterDirection(t *testing.T) {
	builder := NewOrderBuilder(big.NewInt(137), testSign)

	for _, tc := range []struct {
		name      string
		direction combostypes.Direction
		side      int
	}{
		{name: "buy", direction: combostypes.DirectionBuy, side: 0},
		{name: "sell", direction: combostypes.DirectionSell, side: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			order, err := builder.BuildTakerOrder(BuildTakerOrderInput{
				TokenID:       "123456",
				MakerAmountE6: "450000",
				TakerAmountE6: "1000000",
				Direction:     tc.direction,
			}, testAuthOption(model.POLY_GNOSIS_SAFE))

			require.NoError(t, err)
			assert.Equal(t, tc.side, order.Side)
			assert.Equal(t, "450000", order.MakerAmount)
			assert.Equal(t, "1000000", order.TakerAmount)
			assert.Equal(t, "0", order.Expiration)
		})
	}
}

func TestDecodeTakerGatewayEvents(t *testing.T) {
	event, err := decodeEvent([]byte(`{"type":"RFQ_QUOTE_READY","request":{"rfq_id":"rfq_1","yes_position_id":"111","no_position_id":"222"},"quote":{"quote_id":"quote_1","maker_amount_e6":"450000","taker_amount_e6":"1000000"}}`))

	require.NoError(t, err)
	require.NotNil(t, event.TakerQuoteReady)
	assert.Equal(t, "rfq_1", event.TakerQuoteReady.Request.RFQID)
	assert.Equal(t, "quote_1", event.TakerQuoteReady.Quote.QuoteID)

	event, err = decodeEvent([]byte(`{"type":"RFQ_STATUS_UPDATE","status":"EXPIRED","code":"EXPIRED_RFQ"}`))
	require.NoError(t, err)
	assert.True(t, isExpired(event.StatusUpdate.Status, event.StatusUpdate.Code))
}

func TestTakerWSURL(t *testing.T) {
	assert.Equal(t, "wss://combos-rfq-gateway-requester.polymarket.sh/ws", NewTakerWSClient("").url)
	assert.Equal(t, "wss://example.com/ws", NewTakerWSClient("wss://example.com").url)
}
