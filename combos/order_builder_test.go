package combos

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildQuoteRequestUserBuyBuildsMakerSell(t *testing.T) {
	builder := NewOrderBuilder(big.NewInt(137), testSign)
	option := testAuthOption(model.EOA)

	quote, err := builder.BuildQuoteRequest(BuildQuoteRequestInput{
		RFQID:         "rfq_1",
		QuoteID:       "quote_1",
		Direction:     combostypes.DirectionBuy,
		YesPositionID: "123456",
		PriceE6:       "450000",
		SizeE6:        "2000000",
	}, option)

	require.NoError(t, err)
	assert.Equal(t, model.SELL, quote.SignedOrder.Side)
	assert.Equal(t, "2000000", quote.SignedOrder.MakerAmount)
	assert.Equal(t, "900000", quote.SignedOrder.TakerAmount)
	assert.Equal(t, combostypes.QuoteSourceCollateral, quote.Source)
	assert.Contains(t, quote.SignedOrder.Signature, "0x")
}

func TestBuildQuoteRequestUserSellBuildsMakerBuy(t *testing.T) {
	builder := NewOrderBuilder(big.NewInt(137), testSign)
	option := testAuthOption(model.POLY_GNOSIS_SAFE)

	quote, err := builder.BuildQuoteRequest(BuildQuoteRequestInput{
		RFQID:         "rfq_1",
		QuoteID:       "quote_1",
		Direction:     combostypes.DirectionSell,
		YesPositionID: "123456",
		PriceE6:       "333333",
		SizeE6:        "1000001",
		Source:        combostypes.QuoteSourceInventory,
	}, option)

	require.NoError(t, err)
	assert.Equal(t, model.BUY, quote.SignedOrder.Side)
	assert.Equal(t, "333334", quote.SignedOrder.MakerAmount)
	assert.Equal(t, "1000001", quote.SignedOrder.TakerAmount)
	assert.Equal(t, combostypes.QuoteSourceInventory, quote.Source)
}

func TestBuildQuoteRequestPOLYGnosisSafeSignsOrderHash(t *testing.T) {
	var signedBy string
	var signedDigest []byte
	builder := NewOrderBuilder(big.NewInt(137), func(signer string, digest []byte) ([]byte, error) {
		signedBy = signer
		signedDigest = append([]byte(nil), digest...)
		return make([]byte, 65), nil
	})
	option := testAuthOption(model.POLY_GNOSIS_SAFE)

	quote, err := builder.BuildQuoteRequest(BuildQuoteRequestInput{
		RFQID:         "rfq_1",
		QuoteID:       "quote_1",
		Direction:     combostypes.DirectionSell,
		YesPositionID: "123456",
		PriceE6:       "333333",
		SizeE6:        "1000001",
	}, option)

	require.NoError(t, err)
	assert.Equal(t, option.FunderAddress, quote.SignedOrder.Maker)
	assert.Equal(t, option.SingerAddress, quote.SignedOrder.Signer)
	assert.Equal(t, option.SingerAddress, quote.SignerAddress)
	assert.Equal(t, model.POLY_GNOSIS_SAFE, quote.SignatureType)
	assert.Equal(t, model.POLY_GNOSIS_SAFE, quote.SignedOrder.SignatureType)
	assert.Equal(t, option.SingerAddress, signedBy)
	assert.Len(t, signedDigest, common.HashLength)
	assert.Len(t, quote.SignedOrder.Signature, 132)
}

func TestBuildQuoteRequestPOLY1271SignsWrappedHash(t *testing.T) {
	builder := NewOrderBuilder(big.NewInt(137), testSign)
	option := testAuthOption(model.POLY_1271)

	quote, err := builder.BuildQuoteRequest(BuildQuoteRequestInput{
		RFQID:         "rfq_1",
		QuoteID:       "quote_1",
		Direction:     combostypes.DirectionBuy,
		YesPositionID: "123456",
		PriceE6:       "450000",
		SizeE6:        "1000000",
	}, option)

	require.NoError(t, err)
	assert.Equal(t, option.FunderAddress, quote.SignedOrder.Maker)
	assert.Equal(t, option.FunderAddress, quote.SignedOrder.Signer)
	assert.Equal(t, option.FunderAddress, quote.SignerAddress)
	assert.Greater(t, len(quote.SignedOrder.Signature), 132)
}

func testAuthOption(signatureType model.SignatureType) *sdktypes.AuthOption {
	return &sdktypes.AuthOption{
		SignatureType: signatureType,
		SingerAddress: "0x1111111111111111111111111111111111111111",
		FunderAddress: "0x2222222222222222222222222222222222222222",
	}
}

func testSign(_ string, _ []byte) ([]byte, error) {
	return make([]byte, 65), nil
}

func TestMakerRunnerCanBeConstructed(t *testing.T) {
	runner := NewMakerRunner(NewWSClient("wss://example.com"), testAuthOption(model.EOA))
	runner.OnError(func(context.Context, error) {})
	assert.NotNil(t, runner)
}
