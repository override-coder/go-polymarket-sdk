//go:build polymarket_live

package live_test

import (
	"context"
	"encoding/json"
	"fmt"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func budgetRPC(t *testing.T, ctx context.Context) *ethclient.Client {
	t.Helper()
	rpc, err := ethclient.DialContext(ctx, env(t, "POLYGON_RPC"))
	require.NoError(t, err)
	chain, err := rpc.ChainID(ctx)
	if err != nil || chain.Cmp(big.NewInt(137)) != 0 {
		rpc.Close()
		t.Fatalf("budget verification requires Polygon chain 137: chain=%v err=%v", chain, err)
	}
	return rpc
}

func verifyBudgetSpend(t *testing.T, ctx context.Context, a account, fills []ct.Trade, budget decimal.Decimal) {
	t.Helper()
	rpc := budgetRPC(t, ctx)
	defer rpc.Close()
	seen := map[common.Hash]bool{}
	total := new(big.Int)
	for _, fill := range fills {
		raw := common.FromHex(fill.TransactionHash)
		require.Len(t, raw, common.HashLength, "invalid confirmed transaction hash")
		hash := common.BytesToHash(raw)
		require.NotEqual(t, common.Hash{}, hash)
		if seen[hash] {
			continue
		}
		seen[hash] = true
		receipt, err := rpc.TransactionReceipt(ctx, hash)
		require.NoError(t, err, "use query action to retry receipt verification")
		require.Equal(t, hash, receipt.TxHash)
		require.Equal(t, uint64(ethtypes.ReceiptStatusSuccessful), receipt.Status)
		spent, err := collateralOutflow(receipt.Logs, common.HexToAddress(a.auth.FunderAddress))
		require.NoError(t, err)
		require.Positive(t, spent.Sign(), "no pUSD outflow from funding wallet")
		total.Add(total, spent)
		t.Logf("budget receipt tx=%s block=%s pUSD_outflow_units=%s", hash, receipt.BlockNumber, spent)
	}
	require.NotEmpty(t, seen, "no confirmed fill receipts")
	actual := decimal.NewFromBigInt(total, -6)
	t.Logf("pUSD_outflow_including_fees=%s budget=%s", actual, budget)
	require.True(t, actual.LessThanOrEqual(budget), "pUSD outflow %s exceeds budget %s", actual, budget)
}

// Count gross collateral sent to other addresses. Including fees and ignoring
// refunds makes this a conservative upper bound on net spend for the receipts.
func collateralOutflow(logs []*ethtypes.Log, holder common.Address) (*big.Int, error) {
	total := new(big.Int)
	token := common.HexToAddress(pusd)
	transfer := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	owner := common.BytesToHash(holder.Bytes())
	for _, log := range logs {
		if log == nil || log.Address != token || len(log.Topics) != 3 || log.Topics[0] != transfer || log.Topics[1] != owner || log.Topics[2] == owner {
			continue
		}
		if log.Removed || len(log.Data) != 32 {
			return nil, fmt.Errorf("invalid collateral transfer log")
		}
		total.Add(total, new(big.Int).SetBytes(log.Data))
	}
	return total, nil
}

func TestBudgetCollateralOutflow(t *testing.T) {
	holder := common.HexToAddress("0x1234")
	other := common.HexToAddress("0x5678")
	makeLog := func(from, to common.Address, units int64) *ethtypes.Log {
		return &ethtypes.Log{
			Address: common.HexToAddress(pusd),
			Topics: []common.Hash{
				crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
				common.BytesToHash(from.Bytes()), common.BytesToHash(to.Bytes()),
			},
			Data: big.NewInt(units).FillBytes(make([]byte, 32)),
		}
	}
	unrelated := makeLog(holder, other, 9000000)
	unrelated.Address = other
	logs := []*ethtypes.Log{nil, makeLog(holder, other, 1150000), makeLog(holder, other, 10000), makeLog(other, holder, 20000), makeLog(holder, holder, 100), unrelated}
	amount, err := collateralOutflow(logs, holder)
	require.NoError(t, err)
	require.Equal(t, int64(1160000), amount.Int64())
	bad := makeLog(holder, other, 1)
	bad.Data = []byte{1}
	_, err = collateralOutflow([]*ethtypes.Log{bad}, holder)
	require.Error(t, err)
}

func TestBudgetReceiptsDeduplicateTransactions(t *testing.T) {
	holder := common.HexToAddress("0x1234")
	hash := common.HexToHash("0xabcd")
	receipt := &ethtypes.Receipt{
		TxHash: hash, BlockHash: common.HexToHash("0x123"), BlockNumber: big.NewInt(100),
		Status: ethtypes.ReceiptStatusSuccessful,
		Logs: []*ethtypes.Log{{
			Address: common.HexToAddress(pusd),
			Topics:  []common.Hash{crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")), common.BytesToHash(holder.Bytes()), {}},
			Data:    big.NewInt(1160000).FillBytes(make([]byte, 32)),
		}},
	}
	var receiptCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ID     json.RawMessage
			Method string
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		var result any
		switch request.Method {
		case "eth_chainId":
			result = "0x89"
		case "eth_getTransactionReceipt":
			receiptCalls.Add(1)
			result = receipt
		default:
			t.Errorf("unexpected RPC %s", request.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": result})
	}))
	defer server.Close()
	t.Setenv("POLY_LIVE_POLYGON_RPC", server.URL)
	a := account{auth: &sdk.AuthOption{FunderAddress: holder.Hex()}}
	fills := []ct.Trade{{ID: "one", TransactionHash: hash.Hex()}, {ID: "two", TransactionHash: hash.Hex()}}
	verifyBudgetSpend(t, context.Background(), a, fills, decimal.RequireFromString("1.20"))
	require.Equal(t, int32(1), receiptCalls.Load())
}
