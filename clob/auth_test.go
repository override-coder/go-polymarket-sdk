package clob_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/override-coder/go-polymarket-sdk/clob"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/assert"
	"math/big"
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
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, nil)
	key, err := client.DeriveAPIKey(big.NewInt(1), &sdktypes.AuthOption{
		SignatureType: model.POLY_GNOSIS_SAFE,
		SingerAddress: "0x8c5f23249462e20C4a202Ad35275562075F37e09",
		FunderAddress: "0x3BfD9C49E5B62cBc4b7DcE1b7a1f8123B515D278",
	})
	assert.Nil(t, err)
	t.Logf("%+v", key)
}

func TestGetTickSize(t *testing.T) {
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, nil)
	tokenID := "9271150303929886170307880275525448524362138570132971024043819830163602744465"
	size, err := client.GetTickSize(tokenID)
	assert.Nil(t, err)
	t.Logf("%+v", size)

	risk, err := client.GetNegRisk(tokenID)
	assert.Nil(t, err)
	t.Logf("%+v", risk)

	rateBps, err := client.GetFeeRateBps(tokenID)
	assert.Nil(t, err)
	t.Logf("%+v", rateBps)
}

func TestPostOrder(t *testing.T) {
	client := clob.NewClient(PolymarketClobURL, chaindId, signature, &sdktypes.BuilderApiKeyCreds{
		Key:        "019a4dec-fc6a-79ba-8937-d9bf3c2792ca",
		Secret:     "Q23ZHyR21V5_F8qVLvOvnXGhxtW6CmNCWDjHzFJQW7k=",
		Passphrase: "2a171196ddfe34aab62eea32ed63fe424fde8144413982dd90527c844cf2e8d3",
	})

	authOption := &sdktypes.AuthOption{
		SignatureType: model.POLY_GNOSIS_SAFE,
		SingerAddress: "0x8c5f23249462e20C4a202Ad35275562075F37e09",
		FunderAddress: "0x3BfD9C49E5B62cBc4b7DcE1b7a1f8123B515D278",
	}
	_, err := client.EnsureAPIKey(big.NewInt(0), authOption)
	assert.Nil(t, err)

	expiration := time.Now().Add(2 * time.Minute).Unix()
	userOrder := types.UserOrder{
		TokenID:    "29932229206038996544221694126815434341861961592336413071656609906503218641045",
		Price:      0.01,
		Size:       100,
		Side:       types.BUY,
		Expiration: &expiration,
	}

	order, err := client.CreateOrder(userOrder, types.OrderTypeGTD, false, authOption)
	assert.Nil(t, err)
	t.Logf("%+v", order)
}
