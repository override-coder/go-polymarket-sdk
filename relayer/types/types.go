package types

import (
	"fmt"
	"time"
)

type TransactionType string

const (
	TransactionTypeSAFE       TransactionType = "SAFE"
	TransactionTypeSAFECreate TransactionType = "SAFE-CREATE"

	TransactionTypeWALLET       TransactionType = "WALLET"
	TransactionTypeWALLETCreate TransactionType = "WALLET-CREATE"
)

type SignatureParams struct {
	GasPrice *string `json:"gasPrice,omitempty"`

	// SAFE sig parameters
	Operation      *string `json:"operation,omitempty"`
	SafeTxnGas     *string `json:"safeTxnGas,omitempty"`
	BaseGas        *string `json:"baseGas,omitempty"`
	GasToken       *string `json:"gasToken,omitempty"`
	RefundReceiver *string `json:"refundReceiver,omitempty"`

	// SAFE CREATE sig parameters
	PaymentToken    *string `json:"paymentToken,omitempty"`
	Payment         *string `json:"payment,omitempty"`
	PaymentReceiver *string `json:"paymentReceiver,omitempty"`
}

type NoncePayload struct {
	Nonce string `json:"nonce"`
}

type TransactionRequest struct {
	Type            string          `json:"type"`
	From            string          `json:"from"`
	To              string          `json:"to"`
	ProxyWallet     *string         `json:"proxyWallet,omitempty"`
	Data            string          `json:"data"`
	Nonce           *string         `json:"nonce,omitempty"`
	Signature       string          `json:"signature"`
	SignatureParams SignatureParams `json:"signatureParams"`
	Metadata        *string         `json:"metadata,omitempty"`
}

type OperationType uint8

const (
	OperationCall         OperationType = iota // 0
	OperationDelegateCall                      // 1
)

type SafeTransaction struct {
	To        string        `json:"to"`
	Operation OperationType `json:"operation"` // 0/1
	Data      string        `json:"data"`
	Value     string        `json:"value"`
}

type SafeTransactionArgs struct {
	From         string            `json:"from"`
	Nonce        string            `json:"nonce"`
	ChainID      int64             `json:"chainId"`
	Transactions []SafeTransaction `json:"transactions"`
}

type SafeCreateTransactionArgs struct {
	From            string `json:"from"`
	ChainID         int64  `json:"chainId"`
	PaymentToken    string `json:"paymentToken"`
	Payment         string `json:"payment"`
	PaymentReceiver string `json:"paymentReceiver"`
}

type DepositWalletCall struct {
	Target string `json:"target"`
	Value  string `json:"value"`
	Data   string `json:"data"`
}

type DepositWalletTransactionArgs struct {
	From          string              `json:"from"`
	ChainID       int64               `json:"chainId"`
	WalletAddress string              `json:"walletAddress"`
	Nonce         string              `json:"nonce"`
	Deadline      string              `json:"deadline"`
	Calls         []DepositWalletCall `json:"calls"`
}

type DepositWalletParams struct {
	DepositWallet string              `json:"depositWallet"`
	Deadline      string              `json:"deadline"`
	Calls         []DepositWalletCall `json:"calls"`
}

type DepositWalletBatchRequest struct {
	Type                string              `json:"type"`
	From                string              `json:"from"`
	To                  string              `json:"to"`
	Nonce               string              `json:"nonce"`
	Signature           string              `json:"signature"`
	DepositWalletParams DepositWalletParams `json:"depositWalletParams"`
}

type DepositWalletExecuteRequest struct {
	WalletAddress string `json:"walletAddress"`
	Batch         Batch  `json:"batch"`
	Signature     string `json:"signature"`
}

func DepositWalletCallFromSafeTransaction(txn SafeTransaction) (DepositWalletCall, error) {
	if txn.Operation != OperationCall {
		return DepositWalletCall{}, fmt.Errorf("unsupported safe transaction operation for deposit wallet call: %d", txn.Operation)
	}

	return DepositWalletCall{
		Target: txn.To,
		Value:  txn.Value,
		Data:   txn.Data,
	}, nil
}

func DepositWalletCallsFromSafeTransactions(txns []SafeTransaction) ([]DepositWalletCall, error) {
	calls := make([]DepositWalletCall, 0, len(txns))
	for _, txn := range txns {
		call, err := DepositWalletCallFromSafeTransaction(txn)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	return calls, nil
}

type RelayerTransactionState string

const (
	RelayerStateNew       RelayerTransactionState = "STATE_NEW"
	RelayerStateExecuted  RelayerTransactionState = "STATE_EXECUTED"
	RelayerStateMined     RelayerTransactionState = "STATE_MINED"
	RelayerStateInvalid   RelayerTransactionState = "STATE_INVALID"
	RelayerStateConfirmed RelayerTransactionState = "STATE_CONFIRMED"
	RelayerStateFailed    RelayerTransactionState = "STATE_FAILED"
)

type RelayerTransaction struct {
	TransactionID   string                  `json:"transactionID"`
	TransactionHash string                  `json:"transactionHash"`
	From            string                  `json:"from"`
	To              string                  `json:"to"`
	ProxyAddress    string                  `json:"proxyAddress"`
	Data            string                  `json:"data"`
	Nonce           string                  `json:"nonce"`
	Value           string                  `json:"value"`
	State           RelayerTransactionState `json:"state"`
	Type            string                  `json:"type"`
	Metadata        string                  `json:"metadata"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
}

type RelayerTransactionResponse struct {
	TransactionID   string `json:"transactionID"`
	State           string `json:"state"`
	Hash            string `json:"hash"`
	TransactionHash string `json:"transactionHash"`
}

type GetDeployedResponse struct {
	Deployed bool `json:"deployed"`
}

type DepositWalletCreateRequest struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}
