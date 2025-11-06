package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
	"strings"
)

const (
	ClobAuthDomain          = "ClobAuthDomain"
	ClobAuthDomainVersion   = "1"
	ClobAuthDomainMsgToSign = "This message attests that I control the given wallet"
)

func BuildClobEip712Signature(signatureFunc SignatureFunc, chainId *big.Int, singer, ts string, nonce *big.Int) (sigHex string, err error) {
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"ClobAuth": []apitypes.Type{
				{Name: "address", Type: "address"},
				{Name: "timestamp", Type: "string"},
				{Name: "nonce", Type: "uint256"},
				{Name: "message", Type: "string"},
			},
		},
		PrimaryType: "ClobAuth",
		Domain: apitypes.TypedDataDomain{
			Name:    ClobAuthDomain,
			Version: ClobAuthDomainVersion,
			ChainId: math.NewHexOrDecimal256(chainId.Int64()),
		},
		Message: map[string]interface{}{
			"address":   singer,
			"timestamp": ts,
			"nonce":     nonce,
			"message":   ClobAuthDomainMsgToSign,
		},
	}

	hash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return "", fmt.Errorf("TypedDataAndHash failed: %w", err)
	}

	sigBytes, err := signatureFunc(singer, hash)
	if err != nil {
		return "", fmt.Errorf("signature failed: %w", err)
	}

	return "0x" + hex.EncodeToString(sigBytes), nil
}

func BuildPolyHmacSignature(secret string, timestamp string, method string, requestPath string, body *string) (string, error) {
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret from base64 url: %w", err)
	}

	message := timestamp + method + requestPath
	if body != nil {
		processed := strings.ReplaceAll(*body, "'", `"`)
		message += processed
	}

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(message))
	digest := mac.Sum(nil)

	signature := base64.URLEncoding.EncodeToString(digest)
	return signature, nil
}
