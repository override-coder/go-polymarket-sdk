package signing

import (
	"math/big"
	"testing"

	rt "github.com/override-coder/go-polymarket-sdk/relayer/types"
)

func TestMissingSigningDependencies(t *testing.T) {
	if _, err := BuildClobEip712Signature(nil, big.NewInt(137), "", "", nil); err == nil {
		t.Fatal("missing signer")
	}
	if _, err := BuildSafeCreateTransactionEip712Signature(nil, big.NewInt(137), "", "", "", "", ""); err == nil {
		t.Fatal("missing signer")
	}
	if _, err := BuildSafeCreateSafeSignature(nil, 137, "", "", "", "", "", rt.OperationCall, "", "", "", "", "", ""); err == nil {
		t.Fatal("missing signer")
	}
	if _, err := BuildDepositWalletBatchSignature(nil, 137, "", "", "", "", nil); err == nil {
		t.Fatal("missing signer")
	}
	sign := func(string, []byte) ([]byte, error) { t.Fatal("invalid input reached signer"); return nil, nil }
	if _, err := BuildClobEip712Signature(sign, nil, "", "", nil); err == nil {
		t.Fatal("missing chain")
	}
	if _, err := BuildSafeCreateTransactionEip712Signature(sign, nil, "", "", "", "", ""); err == nil {
		t.Fatal("missing chain")
	}
}
