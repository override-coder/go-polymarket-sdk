package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/builder"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderV2RoutesByAssetID(t *testing.T) {
	tests := []struct {
		name, assetID string
		negRisk       bool
		expected      model.VerifyingContract
		orderType     types.OrderType
		sigType       model.SignatureType
	}{
		{"v2 hex with bits 104 and 39 set", "0x0000000000000000000000000000000000000100000000000000008000000001", true, model.ExchangeV3, types.OrderTypeGTC, model.POLY_1271},
		{"v2 decimal", new(big.Int).Lsh(big.NewInt(1), 248).String(), false, model.ExchangeV3, types.OrderTypeFAK, model.EOA},
		{"ctf bit 40", "1099511627776", false, model.CTFExchange, types.OrderTypeGTD, model.EOA},
		{"neg-risk ctf bit 103", new(big.Int).Lsh(big.NewInt(1), 103).String(), true, model.NegRiskCTFExchange, types.OrderTypeFOK, model.EOA},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posted := make(chan types.NewOrderV2, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method + " " + r.URL.Path {
				case "GET /markets-by-token/" + tt.assetID:
					_, _ = fmt.Fprint(w, `{"condition_id":"condition"}`)
				case "GET /clob-markets/condition":
					_, _ = fmt.Fprintf(w, `{"mts":0.01,"nr":%t,"t":[{"t":%q,"o":"yes"}]}`, tt.negRisk, tt.assetID)
				case "POST /order":
					var payload types.NewOrderV2
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Errorf("decode posted order: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					posted <- payload
					_, _ = fmt.Fprint(w, `{"success":true}`)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()

			chain := big.NewInt(137)
			innerSig := make([]byte, 65)
			var signedDigest []byte
			client := NewClient(server.URL, chain, func(_ string, digest []byte) ([]byte, error) {
				signedDigest = digest
				return innerSig, nil
			}, nil)
			auth := &sdktypes.AuthOption{
				SingerAddress: "0x0000000000000000000000000000000000000001",
				SignatureType: tt.sigType,
				ApiKeyCreds:   &sdktypes.ApiKeyCreds{ApiKey: "key", Secret: "c2VjcmV0"},
			}
			order := types.UserOrderV2{TokenId: tt.assetID, Side: types.BUY, Price: 0.5, Size: 10}
			resp, err := client.CreateOrderV2(context.Background(), &order, tt.orderType, false, true, auth)
			require.NoError(t, err)
			require.True(t, resp.Success)
			require.Len(t, posted, 1)
			wire := (<-posted).Order
			require.Equal(t, tt.assetID, wire.TokenID)

			asset, ok := new(big.Int).SetString(tt.assetID, 0)
			require.True(t, ok)
			reference := builder.NewExchangeOrderBuilderImplV2(chain, func() int64 { return wire.Salt })
			expected, err := reference.BuildOrder(&model.OrderDataV2{
				Maker: wire.Maker, Signer: wire.Signer, TokenID: asset.String(),
				MakerAmount: wire.MakerAmount, TakerAmount: wire.TakerAmount, Side: model.BUY,
				Timestamp: wire.Timestamp, Expiration: wire.Expiration, SignatureType: tt.sigType,
				Builder: wire.Builder, Metadata: wire.Metadata,
			})
			require.NoError(t, err)
			if tt.sigType == model.POLY_1271 {
				digest, err := reference.BuildPoly1271WrappedHash(expected, tt.expected)
				require.NoError(t, err)
				require.Equal(t, digest.Bytes(), signedDigest)
				sig, err := reference.BuildPoly1271FinalSignature(expected, tt.expected, innerSig)
				require.NoError(t, err)
				require.Equal(t, "0x"+common.Bytes2Hex(sig), wire.Signature)
			} else {
				digest, err := reference.BuildOrderHash(expected, tt.expected)
				require.NoError(t, err)
				require.Equal(t, digest.Bytes(), signedDigest)
				wantSig := common.Bytes2Hex(innerSig)
				if tt.expected == model.ExchangeV3 {
					wantSig = "0x" + wantSig
				}
				require.Equal(t, wantSig, wire.Signature)
			}
		})
	}
}

func TestParseOrderAssetID(t *testing.T) {
	for _, value := range []string{"", "-1", "0x", "0x0_1", "0x1" + strings.Repeat("0", 64)} {
		_, err := parseOrderAssetID(value)
		require.Error(t, err, value)
	}
	for _, value := range []string{"00010", "0xA", "0b1010", "0o12"} {
		id, err := parseOrderAssetID(value)
		require.NoError(t, err)
		require.Equal(t, int64(10), id.Int64())
	}
}
