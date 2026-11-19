//go:build polymarket_live

package live_test

import (
	"context"
	"fmt"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/override-coder/go-polymarket-sdk/clob"
	"github.com/override-coder/go-polymarket-sdk/relayer"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

type account struct {
	auth   *sdk.AuthOption
	clob   *clob.Client
	relay  *relayer.Client
	beacon bool
}

func enabled(t *testing.T, action string) {
	t.Helper()
	if os.Getenv("POLY_LIVE_ACTION") != action {
		t.Skipf("requires POLY_LIVE_ACTION=%s", action)
	}
	if action != "query" && action != "position" {
		require.Equal(t, "I_UNDERSTAND_THIS_USES_REAL_FUNDS", os.Getenv("POLY_LIVE_ACK"))
	}
}

func wallet(t *testing.T, options ...sdkhttp.Option) account {
	t.Helper()
	key, err := crypto.HexToECDSA(strings.TrimPrefix(env(t, "PRIVATE_KEY"), "0x"))
	require.NoError(t, err)
	owner := crypto.PubkeyToAddress(key.PublicKey)
	sign := func(from string, digest []byte) ([]byte, error) {
		if !strings.EqualFold(from, owner.Hex()) {
			return nil, fmt.Errorf("unexpected signer address")
		}
		sig, err := crypto.Sign(digest, key)
		if err == nil {
			sig[64] += 27
		}
		return sig, err
	}
	chain := big.NewInt(137)
	once := func(c *resty.Client) { c.SetRetryCount(0) }
	r := relayer.NewClient("https://relayer-v2.polymarket.com", chain, sign, nil, once)
	addr := env(t, "WALLET")
	require.True(t, common.IsHexAddress(addr))
	auth := &sdk.AuthOption{SingerAddress: owner.Hex(), FunderAddress: common.HexToAddress(addr).Hex()}
	var expected string
	beacon := false
	switch env(t, "WALLET_TYPE") {
	case "safe":
		auth.SignatureType = model.POLY_GNOSIS_SAFE
		expected, err = r.GetExpectedSafe(owner.Hex())
	case "deposit":
		auth.SignatureType = model.POLY_1271
		expected, err = r.GetExpectedDepositWallet(owner.Hex())
		require.NoError(t, err)
		if !strings.EqualFold(expected, addr) {
			expected, err = r.GetExpectedBeaconDepositWallet(owner.Hex())
			beacon = true
		}
	default:
		t.Fatal("WALLET_TYPE must be deposit or safe")
	}
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress(expected), common.HexToAddress(addr), "wallet differs from owner-derived address")
	return account{auth, clob.NewClient("https://clob.polymarket.com", chain, sign, nil, append(options, once)...), r, beacon}
}

func (a account) credentials(t *testing.T, ctx context.Context) {
	if os.Getenv("POLY_LIVE_API_KEY") != "" {
		a.auth.ApiKeyCreds = &sdk.ApiKeyCreds{ApiKey: env(t, "API_KEY"), Secret: env(t, "API_SECRET"), Passphrase: env(t, "API_PASSPHRASE")}
		return
	}
	creds, err := a.clob.DeriveAPIKey(ctx, big.NewInt(0), a.auth)
	require.NoError(t, err)
	require.NotNil(t, creds)
	require.NotEmpty(t, creds.ApiKey)
	a.auth.ApiKeyCreds = creds
}
