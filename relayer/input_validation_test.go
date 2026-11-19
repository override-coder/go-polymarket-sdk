package relayer

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	rt "github.com/override-coder/go-polymarket-sdk/relayer/types"
	sdk "github.com/override-coder/go-polymarket-sdk/types"
)

type rejectTransport struct{ t *testing.T }

func (f rejectTransport) RoundTrip(*http.Request) (*http.Response, error) {
	f.t.Error("invalid input reached HTTP")
	return nil, fmt.Errorf("unexpected HTTP")
}

func TestRelayerRejectsInvalidInputBeforeIO(t *testing.T) {
	signCalls := 0
	sign := func(string, []byte) ([]byte, error) { signCalls++; return make([]byte, 65), nil }
	creds := &sdk.BuilderApiKeyCreds{Key: "key"}
	c := NewClient("http://invalid", big.NewInt(137), sign, creds, sdkhttp.WithTransport(rejectTransport{t}))
	auth := &sdk.AuthOption{}
	txs := []rt.SafeTransaction{{}}
	calls := map[string]func() error{
		"deploy nil":          func() error { _, e := c.Deploy(nil); return e },
		"execute nil":         func() error { _, e := c.Execute(txs, "", nil); return e },
		"build deploy nil":    func() error { _, e := c.BuildDeployTx(nil); return e },
		"build nil":           func() error { _, e := c.BuildTx(txs, big.NewInt(0), "", nil); return e },
		"nonce nil":           func() error { _, e := c.BuildTx(txs, nil, "", auth); return e },
		"nonce negative":      func() error { _, e := c.BuildTx(txs, big.NewInt(-1), "", auth); return e },
		"empty batch":         func() error { _, e := c.Execute(nil, "", auth); return e },
		"request nil":         func() error { _, e := c.ExecuteByTx(nil, auth); return e },
		"deposit nil":         func() error { _, e := c.BuildDepositWalletBatch(nil, "", nil, nil, nil); return e },
		"deposit request nil": func() error { _, e := c.ExecuteDepositWalletBatchByTx(nil, auth); return e },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if call() == nil {
				t.Fatal("expected error")
			}
		})
	}
	c.signFn = nil
	if _, e := c.BuildTx(txs, big.NewInt(0), "", auth); e == nil {
		t.Fatal("missing signer")
	}
	c.builderApiKeyCreds = nil
	if _, e := c.Execute(txs, "", auth); e == nil {
		t.Fatal("missing credentials")
	}
	if signCalls != 0 {
		t.Fatal("invalid input reached signer")
	}
}

func TestDepositBatchNumericBoundaries(t *testing.T) {
	auth := &sdk.AuthOption{SingerAddress: "0x0000000000000000000000000000000000000001"}
	calls := []rt.DepositWalletCall{{Target: "0x0000000000000000000000000000000000000002", Value: "0", Data: "0x"}}
	overflow := new(big.Int).Lsh(big.NewInt(1), 256).String()
	max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)).String()
	signCalls, requests := 0, 0
	sign := func(string, []byte) ([]byte, error) { signCalls++; return make([]byte, 65), nil }
	invalid := NewClient("http://invalid", big.NewInt(137), sign, nil, sdkhttp.WithTransport(rejectTransport{t}))
	for _, value := range []string{"", "abc", "1.5", "-1", overflow} {
		if _, err := invalid.BuildDepositWalletBatch(calls, "2000000000", &value, nil, auth); err == nil {
			t.Fatalf("nonce %q accepted", value)
		}
	}
	for _, value := range []string{"", "abc", "1.5", "-1", "0", overflow} {
		nonce := "1"
		if _, err := invalid.BuildDepositWalletBatch(calls, value, &nonce, nil, auth); err == nil {
			t.Fatalf("deadline %q accepted", value)
		}
	}
	if signCalls != 0 {
		t.Fatal("invalid numeric input reached signer")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != rt.GET_NONCE || !strings.EqualFold(r.URL.Query().Get("address"), auth.SingerAddress) || r.URL.Query().Get("type") != string(rt.TransactionTypeWALLET) {
			t.Errorf("unexpected request %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"nonce":"7"}`))
	}))
	defer server.Close()
	client := NewClient(server.URL, big.NewInt(137), sign, nil)
	for _, value := range []string{"1", max} {
		if _, err := client.BuildDepositWalletBatch(calls, max, &value, nil, auth); err != nil {
			t.Fatal(err)
		}
	}
	if requests != 0 {
		t.Fatal("explicit nonce triggered HTTP")
	}
	zero := "0"
	for _, nonce := range []*string{nil, &zero} {
		if _, err := client.BuildDepositWalletBatch(calls, "2000000000", nonce, nil, auth); err != nil {
			t.Fatal(err)
		}
	}
	if requests != 2 || signCalls != 4 {
		t.Fatalf("requests=%d signs=%d", requests, signCalls)
	}
}
