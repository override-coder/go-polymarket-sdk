//go:build polymarket_live

package live_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/stretchr/testify/require"
)

func positionContract(m selection) common.Address {
	if m.Version == "v2" {
		return common.HexToAddress(manager)
	}
	if m.NegRisk {
		return common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296")
	}
	return common.HexToAddress(ctf)
}

func readPosition(t *testing.T, ctx context.Context, rpc *ethclient.Client, a account, m selection, record recordResult) {
	t.Helper()
	contract := positionContract(m)
	balance := erc1155Balance(t, ctx, rpc, contract, common.HexToAddress(a.auth.FunderAddress), m.number, nil)
	record("position", liveResult{Asset: m.asset, Contract: contract.Hex(), Balance: balance.String()})
}

func TestQueryPosition(t *testing.T) {
	enabled(t, "position")
	record := beginResult(t)
	m := selected(t)
	a := wallet(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	rpc := budgetRPC(t, ctx)
	defer rpc.Close()
	readPosition(t, ctx, rpc, a, m, record)
	a.credentials(t, ctx)
	assetType := ct.AssetTypeConditional
	if m.Version == "v2" {
		assetType = ct.AssetTypeConditionalV2
	}
	balance, err := a.clob.FetchBalanceAllowance(ctx, ct.BalanceAllowanceRequest{AssetType: assetType, TokenID: m.asset}, a.auth)
	require.NoError(t, err)
	require.NotNil(t, balance)
	record("clob_balance", liveResult{Asset: m.asset, Balance: balance.Balance})
}

func TestPositionContractRouting(t *testing.T) {
	for _, tc := range []struct {
		version string
		neg     bool
		want    string
	}{
		{"v1", false, ctf}, {"v1", true, "0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"},
		{"v2", false, manager}, {"v2", true, manager},
	} {
		require.Equal(t, common.HexToAddress(tc.want), positionContract(selection{market: market{Version: tc.version, NegRisk: tc.neg}}))
	}
}
