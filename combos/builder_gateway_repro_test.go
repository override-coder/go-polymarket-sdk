package combos_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestBuilderGatewayCreateRFQRepro is a minimal, standalone reproduction of
// the official TypeScript SDK's requestComboQuote call. It is intentionally
// opt-in: creating an RFQ is external state, although this test never accepts
// a quote and therefore cannot execute a trade.
//
// This test is suitable for sharing with Polymarket support. It never logs API
// credentials or request headers. Set the following environment variables:
//
//	POLY_COMBO_REPRO=1
//	POLY_COMBO_SIGNER_ADDRESS
//	POLY_COMBO_MAKER_ADDRESS
//	POLY_COMBO_SIGNATURE_TYPE           (for example, 0, 1, 2, or 3)
//	POLY_COMBO_ACCOUNT_API_KEY
//	POLY_COMBO_ACCOUNT_SECRET
//	POLY_COMBO_ACCOUNT_PASSPHRASE
//	POLY_COMBO_BUILDER_API_KEY
//	POLY_COMBO_BUILDER_SECRET
//	POLY_COMBO_BUILDER_PASSPHRASE
//	POLY_COMBO_LEG_POSITION_IDS         (comma-separated, 2-50 active legs)
//
// Optional:
//
//	POLY_COMBO_GATEWAY_HOST              (defaults to production)
//	POLY_COMBO_REQUEST_NOTIONAL_E6       (defaults to 1000000 = 1 pUSD)
func TestBuilderGatewayCreateRFQRepro(t *testing.T) {
	if os.Getenv("POLY_COMBO_REPRO") != "1" {
		t.Skip("set POLY_COMBO_REPRO=1 to call the production Builder Combo Gateway")
	}

	signatureType, err := strconv.ParseUint(requiredEnv(t, "POLY_COMBO_SIGNATURE_TYPE"), 10, 8)
	if err != nil {
		t.Fatalf("POLY_COMBO_SIGNATURE_TYPE must be an unsigned integer: %v", err)
	}

	legPositionIDs := splitRequiredCSV(t, "POLY_COMBO_LEG_POSITION_IDS")
	if len(legPositionIDs) < 2 || len(legPositionIDs) > 50 {
		t.Fatalf("POLY_COMBO_LEG_POSITION_IDS must contain 2 to 50 values; got %d", len(legPositionIDs))
	}

	requestBody := builderRfqCreateRequest{
		SignerAddress:  requiredEnv(t, "POLY_COMBO_SIGNER_ADDRESS"),
		MakerAddress:   requiredEnv(t, "POLY_COMBO_MAKER_ADDRESS"),
		SignatureType:  uint8(signatureType),
		LegPositionIDs: legPositionIDs,
		Direction:      "BUY",
		Side:           "YES",
		RequestedSize: requestedSize{
			Unit:    "notional",
			ValueE6: envOrDefault("POLY_COMBO_REQUEST_NOTIONAL_E6", "1000000"),
		},
	}

	// Marshal once. The exact bytes used to calculate both HMAC values must be
	// the same bytes sent in the HTTP request.
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	bodyString := string(body)

	accountCreds := &apiKeyCreds{
		ApiKey:     requiredEnv(t, "POLY_COMBO_ACCOUNT_API_KEY"),
		Secret:     requiredEnv(t, "POLY_COMBO_ACCOUNT_SECRET"),
		Passphrase: requiredEnv(t, "POLY_COMBO_ACCOUNT_PASSPHRASE"),
	}
	builderCreds := &builderApiKeyCreds{
		Key:        requiredEnv(t, "POLY_COMBO_BUILDER_API_KEY"),
		Secret:     requiredEnv(t, "POLY_COMBO_BUILDER_SECRET"),
		Passphrase: requiredEnv(t, "POLY_COMBO_BUILDER_PASSPHRASE"),
	}

	const requestPath = "/v1/builder/rfq/requests"
	timestamp := time.Now().Unix()
	hmacArgs := l2HeaderArgs{
		Method:      http.MethodPost,
		RequestPath: requestPath,
		Body:        bodyString,
	}

	// These are the same two authorizations resolved by the official TS SDK's
	// Builder Gateway client for a POST request.
	accountHeaders, err := createL2Headers(
		requestBody.SignerAddress,
		accountCreds,
		hmacArgs,
		&timestamp,
	)
	if err != nil {
		t.Fatal(err)
	}
	builderHeaders, err := createL2BuilderHeaders(
		builderCreds,
		hmacArgs,
		&timestamp,
	)
	if err != nil {
		t.Fatal(err)
	}
	requestHeaders := injectBuilderHeaders(accountHeaders, builderHeaders)

	host := strings.TrimRight(envOrDefault(
		"POLY_COMBO_GATEWAY_HOST",
		"https://combos-rfq-gateway-builder.polymarket.com",
	), "/")
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		host+requestPath,
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range requestHeaders {
		req.Header.Set(key, value)
	}

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Builder Gateway response: status=%d body=%s", resp.StatusCode, responseBody)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Builder Gateway rejected request: status=%d body=%s", resp.StatusCode, responseBody)
	}
}

func requiredEnv(t *testing.T, name string) string {
	t.Helper()
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		t.Skipf("set %s to run this reproduction", name)
	}
	return value
}

func splitRequiredCSV(t *testing.T, name string) []string {
	t.Helper()
	parts := strings.Split(requiredEnv(t, name), ",")
	values := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			t.Fatalf("%s must not contain an empty position ID", name)
		}
		for _, character := range value {
			if character < '0' || character > '9' {
				t.Fatalf("%s must contain decimal position IDs", name)
			}
		}
		if _, exists := seen[value]; exists {
			t.Fatalf("%s must not contain duplicate position IDs", name)
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return values
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

// l2HeaderArgs is deliberately local to this reproduction so the HMAC input is
// visible in one file. The signed message is timestamp + method + path + body.
type l2HeaderArgs struct {
	Method      string
	RequestPath string
	Body        string
}

// createL2Headers exactly mirrors this repository's account L2 header helper.
func createL2Headers(
	signerAddress string,
	creds *apiKeyCreds,
	args l2HeaderArgs,
	timestamp *int64,
) (map[string]string, error) {
	if creds == nil {
		return nil, fmt.Errorf("account API credentials are required")
	}
	signature, err := buildHMACSignature(creds.Secret, *timestamp, args)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"POLY_ADDRESS":    signerAddress,
		"POLY_API_KEY":    creds.ApiKey,
		"POLY_PASSPHRASE": creds.Passphrase,
		"POLY_SIGNATURE":  signature,
		"POLY_TIMESTAMP":  strconv.FormatInt(*timestamp, 10),
	}, nil
}

// createL2BuilderHeaders exactly mirrors this repository's Builder L2 header
// helper for non-GET Builder Gateway requests.
func createL2BuilderHeaders(
	creds *builderApiKeyCreds,
	args l2HeaderArgs,
	timestamp *int64,
) (map[string]string, error) {
	if creds == nil {
		return nil, fmt.Errorf("Builder API credentials are required")
	}
	signature, err := buildHMACSignature(creds.Secret, *timestamp, args)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"POLY_BUILDER_API_KEY":    creds.Key,
		"POLY_BUILDER_PASSPHRASE": creds.Passphrase,
		"POLY_BUILDER_SIGNATURE":  signature,
		"POLY_BUILDER_TIMESTAMP":  strconv.FormatInt(*timestamp, 10),
	}, nil
}

// injectBuilderHeaders combines the account L2 and Builder HMAC headers.
func injectBuilderHeaders(accountHeaders, builderHeaders map[string]string) map[string]string {
	result := make(map[string]string, len(accountHeaders)+len(builderHeaders))
	for key, value := range accountHeaders {
		result[key] = value
	}
	for key, value := range builderHeaders {
		result[key] = value
	}
	return result
}

// buildHMACSignature exactly mirrors this repository's sdkheaders dependency:
// URL-safe-base64 decode the secret, replace body apostrophes, HMAC-SHA256 the
// concatenated payload, then emit padded base64url.
func buildHMACSignature(secret string, timestamp int64, args l2HeaderArgs) (string, error) {
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("decode HMAC secret: %w", err)
	}

	body := strings.ReplaceAll(args.Body, "'", `"`)
	message := strconv.FormatInt(timestamp, 10) + args.Method + args.RequestPath + body
	mac := hmac.New(sha256.New, secretBytes)
	if _, err := mac.Write([]byte(message)); err != nil {
		return "", fmt.Errorf("write HMAC payload: %w", err)
	}
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), nil
}

// These request and credential shapes intentionally live in this file. The
// test can therefore be copied without importing this SDK's application types.
type builderRfqCreateRequest struct {
	SignerAddress  string        `json:"signer_address"`
	MakerAddress   string        `json:"maker_address"`
	SignatureType  uint8         `json:"signature_type"`
	LegPositionIDs []string      `json:"leg_position_ids"`
	Direction      string        `json:"direction"`
	Side           string        `json:"side"`
	RequestedSize  requestedSize `json:"requested_size"`
}

type requestedSize struct {
	Unit    string `json:"unit"`
	ValueE6 string `json:"value_e6"`
}

type apiKeyCreds struct {
	ApiKey     string
	Secret     string
	Passphrase string
}

type builderApiKeyCreds struct {
	Key        string
	Secret     string
	Passphrase string
}
