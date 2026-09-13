package combos

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilderTakerClientPlaceUsesOfficialGatewayFlow(t *testing.T) {
	const builderCode = "0x0000000000000000000000000000000000000000000000000000000000000001"
	var createCalls, acceptCalls, statusCalls int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case combostypes.CreateBuilderRFQ:
			createCalls++
			assertAccountAndBuilderHeaders(t, r)
			var body combostypes.BuilderTakerRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "0x2222222222222222222222222222222222222222", body.MakerAddress)
			assert.Equal(t, body.MakerAddress, body.SignerAddress)
			assert.Equal(t, model.POLY_1271, body.SignatureType)
			assert.Equal(t, combostypes.ComboSideYes, body.Side)
			assert.Equal(t, "notional", body.RequestedSize.Unit)
			require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerRFQ{
				RFQID:       "rfq-1",
				Status:      "AWAITING_REQUESTER_ACCEPTANCE",
				ExpiresAt:   time.Now().Add(time.Minute).UnixMilli(),
				BuilderCode: builderCode,
				Request: combostypes.BuilderTakerRFQRequest{
					RFQID:         "rfq-1",
					YesPositionID: "123456",
					Direction:     combostypes.DirectionBuy,
					Side:          combostypes.ComboSideYes,
				},
				Quote: &combostypes.BuilderTakerQuote{
					QuoteID:       "quote-1",
					MakerAmountE6: "450000",
					TakerAmountE6: "1000000",
				},
			}))
		case combostypes.CreateBuilderRFQ + "/rfq-1/accept":
			acceptCalls++
			assertAccountAndBuilderHeaders(t, r)
			var body combostypes.BuilderTakerAcceptRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "quote-1", body.QuoteID)
			assert.Equal(t, "123456", body.SignedOrder.TokenID)
			assert.Equal(t, "450000", body.SignedOrder.MakerAmount)
			assert.Equal(t, "1000000", body.SignedOrder.TakerAmount)
			assert.Equal(t, 0, body.SignedOrder.Side)
			assert.Equal(t, builderCode, body.SignedOrder.Builder)
			assert.Len(t, body.SignedOrder.Timestamp, 10)
			require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerStatus{
				RFQID:          "rfq-1",
				Status:         "EXECUTING",
				TakerOrderHash: "0xorder",
			}))
		case combostypes.CreateBuilderRFQ + "/rfq-1":
			statusCalls++
			assert.Equal(t, http.MethodGet, r.Method)
			assert.NotEmpty(t, r.Header.Get("POLY_ADDRESS"))
			assert.Empty(t, r.Header.Get("POLY_BUILDER_API_KEY"))
			require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerStatus{
				RFQID:  "rfq-1",
				Status: "FILLED",
				TxHash: "0xtx",
			}))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewBuilderTakerClient(server.URL, big.NewInt(137), testSign)
	option := builderTakerTestAuth()
	filled, err := client.PlaceAndWait(context.Background(), TakerRequest{
		LegPositionIDs: []string{"111", "222"},
		Direction:      combostypes.DirectionBuy,
		RequestedSize:  combostypes.RequestedSize{Unit: "notional", ValueE6: "1000000"},
	}, option, time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "FILLED", filled.Status)
	assert.Equal(t, "0xtx", filled.TxHash)
	assert.Equal(t, "0xorder", filled.TakerOrderHash)
	assert.Equal(t, 1, createCalls)
	assert.Equal(t, 1, acceptCalls)
	assert.Equal(t, 1, statusCalls)
}

func TestBuilderTakerClientReturnsNoQuoteBusinessOutcome(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, combostypes.CreateBuilderRFQ, r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerRFQ{
			RFQID:  "rfq-1",
			Status: "FAILED",
			Error:  &combostypes.BuilderTakerError{Code: "NO_QUOTES", Message: "no quote"},
		}))
	}))
	defer server.Close()

	client := NewBuilderTakerClient(server.URL, big.NewInt(137), testSign)
	status, err := client.PlaceAndWait(context.Background(), TakerRequest{
		LegPositionIDs: []string{"111", "222"},
		Direction:      combostypes.DirectionSell,
		Side:           combostypes.ComboSideYes,
		RequestedSize:  combostypes.RequestedSize{Unit: "shares", ValueE6: "1000000"},
	}, builderTakerTestAuth(), time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "FAILED", status.Status)
	require.NotNil(t, status.Error)
	assert.Equal(t, "NO_QUOTES", status.Error.Code)
	assert.Equal(t, 1, calls)
}

func TestBuilderTakerRequestRejectsUnsupportedOfficialInputs(t *testing.T) {
	_, err := builderTakerRequest(TakerRequest{
		LegPositionIDs: []string{"111", "222"},
		Direction:      combostypes.DirectionBuy,
		Side:           combostypes.ComboSideNo,
		RequestedSize:  combostypes.RequestedSize{Unit: "notional", ValueE6: "1000000"},
	}, builderTakerTestAuth())
	require.ErrorContains(t, err, "only YES")

	_, err = builderTakerRequest(TakerRequest{
		LegPositionIDs: []string{"111", "111"},
		Direction:      combostypes.DirectionBuy,
		RequestedSize:  combostypes.RequestedSize{Unit: "notional", ValueE6: "1000000"},
	}, builderTakerTestAuth())
	require.ErrorContains(t, err, "unique")
}

func TestBuilderTakerClientDefaultsGatewayHost(t *testing.T) {
	client := NewBuilderTakerClient("", big.NewInt(137), testSign)
	assert.Equal(t, DefaultBuilderTakerGatewayHost, client.host)
}

func assertAccountAndBuilderHeaders(t *testing.T, r *http.Request) {
	t.Helper()
	for _, header := range []string{
		"POLY_ADDRESS", "POLY_API_KEY", "POLY_PASSPHRASE", "POLY_TIMESTAMP", "POLY_SIGNATURE",
		"POLY_BUILDER_API_KEY", "POLY_BUILDER_PASSPHRASE", "POLY_BUILDER_TIMESTAMP", "POLY_BUILDER_SIGNATURE",
	} {
		assert.NotEmpty(t, r.Header.Get(header), "missing header %s", header)
	}
}

func builderTakerTestAuth() *sdktypes.AuthOption {
	return &sdktypes.AuthOption{
		SignatureType: model.POLY_1271,
		SingerAddress: "0x1111111111111111111111111111111111111111",
		FunderAddress: "0x2222222222222222222222222222222222222222",
		ApiKeyCreds: &sdktypes.ApiKeyCreds{
			ApiKey:     "account-key",
			Secret:     "c2VjcmV0",
			Passphrase: "account-passphrase",
		},
		BuilderApiKeyCreds: &sdktypes.BuilderApiKeyCreds{
			Key:        "builder-key",
			Secret:     "c2VjcmV0",
			Passphrase: "builder-passphrase",
		},
	}
}
