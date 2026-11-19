//go:build polymarket_legacy_live

package clob_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/override-coder/go-polymarket-sdk/clob"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
	"time"
)

var (
	PolymarketClobURL = "https://clob.polymarket.com"
	chaindId          = big.NewInt(137)
	privateKey, _     = crypto.ToECDSA(common.Hex2Bytes(""))
)

func signature(signer string, digest []byte) ([]byte, error) {
	sig, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

func TestClient(t *testing.T) {
	t.Skip("legacy authentication fixture has no signing key; use integration/live")
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, nil)
	key, err := client.DeriveAPIKey(context.Background(), big.NewInt(1), &sdktypes.AuthOption{
		SignatureType: model.POLY_GNOSIS_SAFE,
		SingerAddress: "0x8c5f23249462e20C4a202Ad35275562075F37e09",
		FunderAddress: "0x3BfD9C49E5B62cBc4b7DcE1b7a1f8123B515D278",
	})
	assert.Nil(t, err)
	t.Logf("%+v", key)
}

func TestGetTickSize(t *testing.T) {
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, nil)
	tokenID := "108743709732442130739073851488597967747030701044009651663118921104082786836017"
	size, err := client.GetTickSize(context.Background(), tokenID)
	assert.Nil(t, err)
	t.Logf("%+v", size)

	risk, err := client.GetOrderBook(tokenID)
	assert.Nil(t, err)
	t.Logf("%+v", risk)

	rateBps, err := client.GetMarketPrice(context.Background(), tokenID, "BUY")
	assert.Nil(t, err)
	t.Logf("%+v", rateBps)
}

func TestPostOrder(t *testing.T) {
	t.Skip("legacy authentication fixture has no signing key; use integration/live")
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, &sdktypes.BuilderApiKeyCreds{
		Key:        os.Getenv("POLY_BUILDER_API_KEY"),
		Secret:     os.Getenv("POLY_BUILDER_SECRET"),
		Passphrase: os.Getenv("POLY_BUILDER_PASSPHRASE"),
	})

	authOption := &sdktypes.AuthOption{
		SignatureType: model.POLY_GNOSIS_SAFE,
		SingerAddress: "0x8c5f23249462e20C4a202Ad35275562075F37e09",
		FunderAddress: "0x3BfD9C49E5B62cBc4b7DcE1b7a1f8123B515D278",
	}
	_, err := client.EnsureAPIKey(context.Background(), big.NewInt(0), authOption)
	assert.Nil(t, err)

	orders, err := client.GetOrders(context.Background(), types.GetActiveOrdersRequest{}, authOption)
	assert.Nil(t, err)

	if len(orders.Data) > 0 {
		getOrder, err := client.GetOrder(context.Background(), orders.Data[0].ID, types.GetOrderRequest{}, authOption)
		assert.Nil(t, err)
		t.Logf("%+v", getOrder)
	}

	t.Skip("legacy order submission disabled; use integration/live")
	expiration := time.Now().Add(2 * time.Minute).Unix()
	userOrder := types.UserOrder{
		TokenID:    "29932229206038996544221694126815434341861961592336413071656609906503218641045",
		Price:      0.01,
		Size:       100,
		Side:       types.BUY,
		Expiration: &expiration,
	}

	order, err := client.CreateOrder(context.Background(), userOrder, types.OrderTypeGTD, false, authOption)
	assert.Nil(t, err)
	t.Logf("%+v", order)
}
