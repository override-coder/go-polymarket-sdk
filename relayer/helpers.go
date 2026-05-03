package relayer

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/override-coder/go-polymarket-sdk/relayer/types"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

var (
	erc1967Const1 = common.FromHex("0xcc3735a920a3ca505d382bbc545af43d6000803e6038573d6000fd5b3d6000f3")
	erc1967Const2 = common.FromHex("0x5155f3363d3d373d3d363d7f360894a13ba1a3210667c828492db98dca3e2076")
	erc1967Prefix = mustBigInt("0x61003d3d8160233d3973")
)

func buildSafeCreateTransactionRequest(
	chainId *big.Int,
	signatureFunc signing.SignatureFunc,
	safeCfg string,
	args types.SafeCreateTransactionArgs,
) (*types.TransactionRequest, error) {
	owner := args.From

	signature, err := signing.BuildSafeCreateTransactionEip712Signature(signatureFunc, chainId, owner, safeCfg, args.PaymentToken, args.Payment, args.PaymentReceiver)
	if err != nil {
		return nil, err
	}

	safeAddr, err := deriveSafe(owner, safeCfg)
	if err != nil {
		return nil, fmt.Errorf("buildSafeCreateTransactionRequest: deriveSafe failed: %w", err)
	}

	req := &types.TransactionRequest{
		Type:        string(types.TransactionTypeSAFECreate),
		From:        owner,
		To:          safeCfg,
		ProxyWallet: &safeAddr,
		Data:        "0x",
		Signature:   signature,
		SignatureParams: types.SignatureParams{
			PaymentToken:    &args.PaymentToken,
			Payment:         &args.Payment,
			PaymentReceiver: &args.PaymentReceiver,
		},
	}
	return req, nil
}

func deriveSafe(ownerAddress string, safeFactory string) (string, error) {
	if !common.IsHexAddress(ownerAddress) {
		return "", fmt.Errorf("deriveSafe: invalid owner address: %s", ownerAddress)
	}
	if !common.IsHexAddress(safeFactory) {
		return "", fmt.Errorf("deriveSafe: invalid factory address: %s", safeFactory)
	}

	initCodeHashBytes, err := hex.DecodeString(types.SafeInitCodeHashHex)
	if err != nil {
		return "", fmt.Errorf("deriveSafe: invalid SAFE_INIT_CODE_HASH_HEX: %w", err)
	}
	if len(initCodeHashBytes) != 32 {
		return "", fmt.Errorf("deriveSafe: unexpected init code hash length %d", len(initCodeHashBytes))
	}

	paddedOwner := common.LeftPadBytes(common.HexToAddress(ownerAddress).Bytes(), 32)
	saltBytes := crypto.Keccak256(paddedOwner)
	factoryAddr := common.HexToAddress(safeFactory)
	salt32 := [32]byte{}
	copy(salt32[:], saltBytes[:32])

	computedAddress := crypto.CreateAddress2(factoryAddr, salt32, initCodeHashBytes)
	return computedAddress.Hex(), nil
}

func deriveDepositWallet(ownerAddress, factoryAddress, implementationAddress string) (string, error) {
	if !common.IsHexAddress(ownerAddress) {
		return "", fmt.Errorf("deriveDepositWallet: invalid owner address: %s", ownerAddress)
	}
	if !common.IsHexAddress(factoryAddress) {
		return "", fmt.Errorf("deriveDepositWallet: invalid factory address: %s", factoryAddress)
	}
	if !common.IsHexAddress(implementationAddress) {
		return "", fmt.Errorf("deriveDepositWallet: invalid implementation address: %s", implementationAddress)
	}

	owner := common.HexToAddress(ownerAddress)
	factory := common.HexToAddress(factoryAddress)
	implementation := common.HexToAddress(implementationAddress)

	walletID := common.LeftPadBytes(owner.Bytes(), 32)

	args, err := abi.Arguments{
		{Type: mustABIType("address")},
		{Type: mustABIType("bytes32")},
	}.Pack(factory, [32]byte(walletID))
	if err != nil {
		return "", fmt.Errorf("deriveDepositWallet: abi pack args: %w", err)
	}

	salt := crypto.Keccak256(args)
	bytecodeHash := initCodeHashERC1967(implementation, args)

	var salt32 [32]byte
	copy(salt32[:], salt)

	computed := crypto.CreateAddress2(factory, salt32, bytecodeHash)
	return computed.Hex(), nil
}

func initCodeHashERC1967(implementation common.Address, args []byte) []byte {
	n := big.NewInt(int64(len(args)))
	shiftedN := new(big.Int).Lsh(n, 56)
	combined := new(big.Int).Add(new(big.Int).Set(erc1967Prefix), shiftedN)

	prefixBytes := leftPadBigInt(combined, 10)

	initCode := make([]byte, 0, 10+20+2+32+32+len(args))
	initCode = append(initCode, prefixBytes...)
	initCode = append(initCode, implementation.Bytes()...)
	initCode = append(initCode, 0x60, 0x09)
	initCode = append(initCode, erc1967Const2...)
	initCode = append(initCode, erc1967Const1...)
	initCode = append(initCode, args...)

	return crypto.Keccak256(initCode)
}

func mustBigInt(v string) *big.Int {
	n, ok := new(big.Int).SetString(v[2:], 16)
	if !ok {
		panic("invalid big int: " + v)
	}
	return n
}

func leftPadBigInt(v *big.Int, size int) []byte {
	out := make([]byte, size)
	b := v.Bytes()
	copy(out[size-len(b):], b)
	return out
}

func mustABIType(t string) abi.Type {
	typ, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}

func buildSafeTransactionRequest(signatureFunc signing.SignatureFunc, args types.SafeTransactionArgs, safeCfg *types.ContractConfig, metadata string) (*types.TransactionRequest, error) {
	safeAddr, err := deriveSafe(args.From, safeCfg.SafeFactory)
	if err != nil {
		return nil, fmt.Errorf("build safe transaction request: deriveSafe failed: %w", err)
	}

	transaction, err := aggregateTransaction(args.Transactions, safeCfg.SafeMultisend)
	if err != nil {
		return nil, fmt.Errorf("build safe transaction request: aggregateTransaction failed: %w", err)
	}

	sigBytes, err := signing.BuildSafeCreateSafeSignature(
		signatureFunc,
		args.ChainID,
		args.From,
		safeAddr,
		transaction.To,
		transaction.Value,
		transaction.Data,
		transaction.Operation,
		"0",
		"0",
		"0",
		sdktypes.ZeroAddress,
		sdktypes.ZeroAddress,
		args.Nonce,
	)
	if err != nil {
		return nil, fmt.Errorf("build safe transaction request: create safe signature failed: %w", err)
	}

	sig, err := types.SplitAndPackSignature("0x" + hex.EncodeToString(sigBytes))
	if err != nil {
		return nil, fmt.Errorf("build safe transaction request:  invalid signature: %w", err)
	}

	var (
		gasPrice       = "0"
		operation      = fmt.Sprintf("%d", transaction.Operation)
		safeTxnGas     = "0"
		baseGas        = "0"
		gasToken       = sdktypes.ZeroAddress
		refundReceiver = sdktypes.ZeroAddress
	)

	sigParams := types.SignatureParams{
		GasPrice:       &gasPrice,
		Operation:      &operation,
		SafeTxnGas:     &safeTxnGas,
		BaseGas:        &baseGas,
		GasToken:       &gasToken,
		RefundReceiver: &refundReceiver,
	}

	req := &types.TransactionRequest{
		Type:            string(types.TransactionTypeSAFE),
		From:            args.From,
		To:              transaction.To,
		ProxyWallet:     &safeAddr,
		Data:            transaction.Data,
		Nonce:           &args.Nonce,
		Signature:       sig,
		SignatureParams: sigParams,
		Metadata:        &metadata,
	}
	return req, nil
}

func buildDepositWalletBatchRequest(
	signatureFunc signing.SignatureFunc,
	args types.DepositWalletTransactionArgs,
	cfg *types.DepositWalletContractConfig,
) (*types.DepositWalletBatchRequest, error) {
	signature, err := signing.BuildDepositWalletBatchSignature(
		signatureFunc,
		args.ChainID,
		args.From,
		args.WalletAddress,
		args.Nonce,
		args.Deadline,
		args.Calls,
	)
	if err != nil {
		return nil, fmt.Errorf("build deposit wallet batch request: sign batch failed: %w", err)
	}

	req := &types.DepositWalletBatchRequest{
		Type:      string(types.TransactionTypeWALLET),
		From:      args.From,
		To:        cfg.DepositWalletFactory,
		Nonce:     args.Nonce,
		Signature: signature,
		DepositWalletParams: types.DepositWalletParams{
			DepositWallet: args.WalletAddress,
			Deadline:      args.Deadline,
			Calls:         args.Calls,
		},
	}
	return req, nil
}

func buildDepositWalletExecuteRequest(
	signatureFunc signing.SignatureFunc,
	args types.DepositWalletTransactionArgs,
) (*types.DepositWalletExecuteRequest, error) {
	signature, err := signing.BuildDepositWalletBatchSignature(
		signatureFunc,
		args.ChainID,
		args.From,
		args.WalletAddress,
		args.Nonce,
		args.Deadline,
		args.Calls,
	)
	if err != nil {
		return nil, fmt.Errorf("build deposit wallet execute request: sign batch failed: %w", err)
	}

	execCalls := make([]types.Call, 0, len(args.Calls))
	for _, call := range args.Calls {
		value, ok := new(big.Int).SetString(call.Value, 10)
		if !ok {
			return nil, fmt.Errorf("build deposit wallet execute request: invalid call value: %s", call.Value)
		}
		execCalls = append(execCalls, types.Call{
			Target: common.HexToAddress(call.Target),
			Value:  value,
			Data:   common.FromHex(call.Data),
		})
	}

	nonce, ok := new(big.Int).SetString(args.Nonce, 10)
	if !ok {
		return nil, fmt.Errorf("build deposit wallet execute request: invalid nonce: %s", args.Nonce)
	}
	deadline, ok := new(big.Int).SetString(args.Deadline, 10)
	if !ok {
		return nil, fmt.Errorf("build deposit wallet execute request: invalid deadline: %s", args.Deadline)
	}

	req := &types.DepositWalletExecuteRequest{
		WalletAddress: args.WalletAddress,
		Batch: types.Batch{
			Wallet:   common.HexToAddress(args.WalletAddress),
			Nonce:    nonce,
			Deadline: deadline,
			Calls:    execCalls,
		},
		Signature: signature,
	}
	return req, nil
}

// aggregateTransaction
func aggregateTransaction(txns []types.SafeTransaction, multisendAddr string) (types.SafeTransaction, error) {
	if len(txns) == 1 {
		return txns[0], nil
	}
	return types.CreateSafeMultiSendTransaction(txns, multisendAddr)
}
