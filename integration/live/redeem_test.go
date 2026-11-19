//go:build polymarket_live

package live_test

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	rt "github.com/override-coder/go-polymarket-sdk/relayer/types"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

// Polygon addresses and redeem entry points from the pinned ts-sdk environments.ts and abis.ts.
const (
	pusd       = "0xC011a7E12a19f7B1f670d46F03B03f3342E82DFB"
	ctf        = "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
	ctfAdapter = "0xAdA100Db00Ca00073811820692005400218FcE1f"
	negAdapter = "0xadA2005600Dec949baf300f4C6120000bDB6eAab"
	manager    = "0x006F54F7f9A22e0000CC2AB60031000000ae9fEF"
	router     = "0x12121212006e4CD160D18e3f00711DA5c3372600"
)

func calldata(t *testing.T, signature string, values ...any) []byte {
	t.Helper()
	_, params, _ := strings.Cut(signature, "(")
	var args abi.Arguments
	for _, kind := range strings.Split(strings.TrimSuffix(params, ")"), ",") {
		typ, err := abi.NewType(kind, "", nil)
		require.NoError(t, err)
		args = append(args, abi.Argument{Type: typ})
	}
	packed, err := args.Pack(values...)
	require.NoError(t, err)
	return append(crypto.Keccak256([]byte(signature))[:4], packed...)
}

func erc1155Balance(t *testing.T, ctx context.Context, rpc *ethclient.Client, contract, holder common.Address, asset, block *big.Int) *big.Int {
	t.Helper()
	out, err := rpc.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: calldata(t, "balanceOf(address,uint256)", holder, asset)}, block)
	require.NoError(t, err)
	require.Len(t, out, 32)
	return new(big.Int).SetBytes(out)
}

func TestRedeemResolvedPosition(t *testing.T) {
	enabled(t, "redeem")
	m := selected(t)
	require.True(t, m.Closed, "use a resolved market with an existing winning position")
	if m.Version == "v2" {
		require.Equal(t, "resolved", m.ResolutionStatus)
	}
	var prices []string
	require.NoError(t, json.Unmarshal([]byte(m.OutcomePrices), &prices))
	require.Len(t, prices, 2)
	require.True(t, positive(t, prices[m.index]).Equal(positive(t, "1")), "selected outcome must pay out")
	a := wallet(t)
	a.auth.BuilderApiKeyCreds = &sdk.BuilderApiKeyCreds{Key: env(t, "BUILDER_KEY"), Secret: env(t, "BUILDER_SECRET"), Passphrase: env(t, "BUILDER_PASSPHRASE")}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	rpc, err := ethclient.DialContext(ctx, env(t, "POLYGON_RPC"))
	require.NoError(t, err)
	defer rpc.Close()
	chain, err := rpc.ChainID(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(137), chain)
	token := common.HexToAddress(ctf)
	if m.NegRisk && m.Version == "v1" {
		token = common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296")
	}
	if m.Version == "v2" {
		token = common.HexToAddress(manager)
	}
	holder := common.HexToAddress(a.auth.FunderAddress)
	before := erc1155Balance(t, ctx, rpc, token, holder, m.number, nil)
	require.Positive(t, before.Sign(), "wallet has no selected winning shares")
	var target common.Address
	var data []byte
	if m.Version == "v2" {
		target = common.HexToAddress(router)
		var id [31]byte
		copy(id[:], m.condition)
		data = calldata(t, "redeem(bytes31,uint256,uint256)", id, big.NewInt(int64(m.index)), before)
	} else {
		target = common.HexToAddress(ctfAdapter)
		if m.NegRisk {
			target = common.HexToAddress(negAdapter)
		}
		data = calldata(t, "redeemPositions(address,bytes32,bytes32,uint256[])", common.HexToAddress(pusd), [32]byte{}, common.BytesToHash(m.condition), []*big.Int{big.NewInt(1), big.NewInt(2)})
	}
	_, err = rpc.CallContract(ctx, ethereum.CallMsg{From: holder, To: &target, Data: data}, nil)
	require.NoError(t, err, "redeem simulation from the funding wallet")
	var tx *rt.RelayerTransactionResponse
	if a.auth.SignatureType == model.POLY_1271 {
		tx, err = a.relay.ExecuteDepositWalletBatch([]rt.DepositWalletCall{{Target: target.Hex(), Value: "0", Data: hexutil.Encode(data)}}, strconv.FormatInt(time.Now().Add(5*time.Minute).Unix(), 10), &a.beacon, a.auth)
	} else {
		tx, err = a.relay.Execute([]rt.SafeTransaction{{To: target.Hex(), Value: "0", Operation: rt.OperationCall, Data: hexutil.Encode(data)}}, "redeem market "+m.ID, a.auth)
	}
	require.NoError(t, err)
	require.NotNil(t, tx)
	require.NotEmpty(t, tx.TransactionID)
	t.Logf("redeem market=%s wallet=%s relayer_id=%s", m.ID, holder, tx.TransactionID)
	var receipt *ethtypes.Receipt
	for receipt == nil && ctx.Err() == nil {
		rows, err := a.relay.GetTransaction(tx.TransactionID)
		if err == nil {
			for _, row := range rows {
				if row.TransactionID != tx.TransactionID {
					continue
				}
				require.NotEqual(t, rt.RelayerStateFailed, row.State)
				require.NotEqual(t, rt.RelayerStateInvalid, row.State)
				if row.TransactionHash != "" && (row.State == rt.RelayerStateMined || row.State == rt.RelayerStateConfirmed) {
					r, err := rpc.TransactionReceipt(ctx, common.HexToHash(row.TransactionHash))
					if err == nil {
						receipt = r
						break
					}
				}
			}
		}
		if receipt == nil {
			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Second):
			}
		}
	}
	require.NotNil(t, receipt, "no mined receipt; reconcile relayer ID before retry")
	require.Equal(t, uint64(ethtypes.ReceiptStatusSuccessful), receipt.Status)
	after := erc1155Balance(t, ctx, rpc, token, holder, m.number, receipt.BlockNumber)
	require.Zero(t, after.Sign(), "winning position remains")
	received := new(big.Int)
	transfer := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	for _, log := range receipt.Logs {
		if log.Address == common.HexToAddress(pusd) && len(log.Topics) == 3 && log.Topics[0] == transfer && log.Topics[2] == common.BytesToHash(holder.Bytes()) {
			received.Add(received, new(big.Int).SetBytes(log.Data))
		}
	}
	require.Positive(t, received.Sign(), "receipt has no pUSD payment to funding wallet")
	t.Logf("redeemed tx=%s position_units=%s pUSD_units=%s", receipt.TxHash, before, received)
}
