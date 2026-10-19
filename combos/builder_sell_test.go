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
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestBuilderSellSafeAnd1271Lifecycle(t *testing.T) {
	for _, signatureType := range []model.SignatureType{model.POLY_GNOSIS_SAFE, model.POLY_1271} {
		t.Run(big.NewInt(int64(signatureType)).String(), func(t *testing.T) {
			option := builderTakerTestAuth()
			option.SignatureType = signatureType
			wantSigner := option.SingerAddress
			if signatureType == model.POLY_1271 {
				wantSigner = option.FunderAddress
			}
			accepts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case combostypes.CreateBuilderRFQ:
					assertAccountAndBuilderHeaders(t, r)
					var req combostypes.BuilderTakerRequest
					require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
					require.Equal(t, combostypes.DirectionSell, req.Direction)
					require.Equal(t, combostypes.RequestedSize{Unit: "shares", ValueE6: "10000000"}, req.RequestedSize)
					require.Equal(t, signatureType, req.SignatureType)
					require.Equal(t, wantSigner, req.SignerAddress)
					require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerRFQ{RFQID: "sell", Status: builderTakerStatusAwaitingRequesterAcceptance, ExpiresAt: time.Now().Add(time.Minute).UnixMilli(), BuilderCode: combostypes.Bytes32Zero, Request: combostypes.BuilderTakerRFQRequest{YesPositionID: "123", Direction: combostypes.DirectionSell, Side: combostypes.ComboSideYes}, Quote: &combostypes.BuilderTakerQuote{QuoteID: "quote", MakerAmountE6: "10000000", TakerAmountE6: "5000000", TotalRequiredE6: "10000000", NetReceiveE6: "4900000"}}))
				case combostypes.CreateBuilderRFQ + "/sell/accept":
					accepts++
					assertAccountAndBuilderHeaders(t, r)
					var req combostypes.BuilderTakerAcceptRequest
					require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
					require.Equal(t, 1, req.SignedOrder.Side)
					require.Equal(t, wantSigner, req.SignedOrder.Signer)
					require.Equal(t, option.FunderAddress, req.SignedOrder.Maker)
					require.Equal(t, "10000000", req.SignedOrder.MakerAmount)
					require.Equal(t, "5000000", req.SignedOrder.TakerAmount)
					require.Equal(t, int(signatureType), req.SignedOrder.SignatureType)
					if signatureType == model.POLY_GNOSIS_SAFE {
						require.Len(t, req.SignedOrder.Signature, 132)
					} else {
						require.Greater(t, len(req.SignedOrder.Signature), 132)
					}
					require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerStatus{RFQID: "sell", Status: "EXECUTING", TakerOrderHash: "0xorder"}))
				case combostypes.CreateBuilderRFQ + "/sell":
					require.Empty(t, r.Header.Get("POLY_BUILDER_API_KEY"))
					require.NoError(t, json.NewEncoder(w).Encode(combostypes.BuilderTakerStatus{RFQID: "sell", Status: "CONFIRMED", TxHash: "0xtx"}))
				default:
					t.Errorf("unexpected path %s", r.URL.Path)
				}
			}))
			defer server.Close()
			result, err := NewBuilderTakerClient(server.URL, big.NewInt(137), testSign).PlaceAndWait(context.Background(), TakerRequest{LegPositionIDs: []string{"111", "222"}, Direction: combostypes.DirectionSell, RequestedSize: combostypes.RequestedSize{Unit: "shares", ValueE6: "10000000"}}, option, time.Millisecond)
			require.NoError(t, err)
			require.Equal(t, "CONFIRMED", result.Status)
			require.Equal(t, "0xtx", result.TxHash)
			require.Equal(t, "0xorder", result.TakerOrderHash)
			require.Equal(t, 1, accepts)
		})
	}
}
