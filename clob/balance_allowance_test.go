package clob

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/polymarket/go-order-utils/pkg/model"
)

func TestFetchBalanceAllowanceWalletAndAsset(t *testing.T) {
	for _, signature := range []model.SignatureType{model.EOA, model.POLY_GNOSIS_SAFE, model.POLY_1271} {
		for _, asset := range []ct.AssetType{ct.AssetTypeCollateral, ct.AssetTypeConditional, ct.AssetTypeConditionalV2} {
			auth := testAuthOption()
			auth.SignatureType = signature
			request := ct.BalanceAllowanceRequest{AssetType: asset}
			if asset != ct.AssetTypeCollateral {
				request.TokenID = "123"
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" || r.URL.Path != "/balance-allowance" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				var got model.SignatureType
				if err := json.Unmarshal([]byte(r.URL.Query().Get("signature_type")), &got); err != nil || got != signature {
					t.Errorf("wrong signature type %v, error %v", got, err)
				}
				if r.URL.Query().Get("asset_type") != string(asset) || r.URL.Query().Get("token_id") != request.TokenID {
					t.Errorf("unexpected query %s", r.URL.RawQuery)
				}
				if r.Header.Get("POLY_API_KEY") != auth.ApiKeyCreds.ApiKey || r.Header.Get("POLY_SIGNATURE") == "" {
					t.Error("missing authentication")
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"balance":"1234567","allowances":{"0xe3333700cA9d93003F00f0F71f8515005F6c00Aa":"987654321","0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296":"42"}}`))
			}))
			client := NewClient(server.URL, nil, nil, nil)
			result, err := client.FetchBalanceAllowance(context.Background(), request, auth)
			if err != nil {
				t.Fatal(err)
			}
			if result.Balance != "1234567" || result.Allowances.ExchangeV3 != "987654321" || result.Allowances.NegRiskAdapter != "42" {
				t.Fatalf("lost balance or allowance: %+v", result)
			}
			if asset == ct.AssetTypeCollateral {
				if _, err := client.GetBalanceAllowance(auth); err != nil {
					t.Fatal(err)
				}
			}
			server.Close()
		}
	}
}

func TestFetchBalanceAllowanceRejectsInvalidAndCancelledRequests(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
	}))
	defer server.Close()
	client := NewClient(server.URL, nil, nil, nil)
	for _, request := range []ct.BalanceAllowanceRequest{
		{},
		{AssetType: ct.AssetType("INVALID")},
	} {
		if _, err := client.FetchBalanceAllowance(context.Background(), request, testAuthOption()); err == nil {
			t.Errorf("accepted invalid request %+v", request)
		}
	}
	request := ct.BalanceAllowanceRequest{AssetType: ct.AssetTypeCollateral}
	if _, err := client.FetchBalanceAllowance(context.Background(), request, nil); err == nil {
		t.Error("accepted missing auth")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.FetchBalanceAllowance(ctx, request, testAuthOption()); err == nil {
		t.Error("ignored cancellation")
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid requests reached network: %d", calls.Load())
	}
}

// TS keeps token IDs optional at this low-level endpoint and delegates token
// validation to CLOB. Preserve the supplied string when one is present.
func TestFetchBalanceAllowanceOptionalToken(t *testing.T) {
	for _, asset := range []ct.AssetType{ct.AssetTypeCollateral, ct.AssetTypeConditional, ct.AssetTypeConditionalV2} {
		for _, token := range []string{"", "123", "0x123", "invalid"} {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("token_id") != token {
					t.Errorf("token changed: %q", r.URL.Query().Get("token_id"))
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"balance":"0","allowances":{}}`))
			}))
			client := NewClient(server.URL, nil, nil, nil)
			_, err := client.FetchBalanceAllowance(context.Background(), ct.BalanceAllowanceRequest{AssetType: asset, TokenID: token}, testAuthOption())
			server.Close()
			if err != nil {
				t.Fatalf("asset=%s token=%q: %v", asset, token, err)
			}
		}
	}
}
