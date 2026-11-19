//go:build polymarket_live

package live_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	ct "github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/stretchr/testify/require"
)

// Only public identifiers and acceptance observations belong in this allowlist.
type liveResult struct {
	SDKRevision string   `json:"sdk_revision,omitempty"`
	Amount      string   `json:"amount,omitempty"`
	MaxSpend    string   `json:"max_spend,omitempty"`
	Shares      string   `json:"shares,omitempty"`
	MaxPrice    string   `json:"max_price,omitempty"`
	MinPrice    string   `json:"min_price,omitempty"`
	Time        string   `json:"time"`
	Action      string   `json:"action"`
	WalletType  string   `json:"wallet_type"`
	Wallet      string   `json:"wallet"`
	Protocol    string   `json:"protocol"`
	Market      string   `json:"market"`
	Outcome     string   `json:"outcome_index"`
	Stage       string   `json:"stage"`
	Asset       string   `json:"asset,omitempty"`
	Side        string   `json:"side,omitempty"`
	OrderID     string   `json:"order_id,omitempty"`
	RelayerID   string   `json:"relayer_id,omitempty"`
	Hashes      []string `json:"transaction_hashes,omitempty"`
	Contract    string   `json:"position_contract,omitempty"`
	Balance     string   `json:"balance_units,omitempty"`
}
type recordResult func(string, liveResult)

func createResultFile(path string) (*os.File, error) {
	// Refuse overwriting a previous attempt, including failed/unknown submissions.
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
}

func beginResult(t *testing.T) recordResult {
	t.Helper()
	var file *os.File
	if path := os.Getenv("POLY_LIVE_RESULT_FILE"); path != "" {
		var err error
		file, err = createResultFile(path)
		require.NoError(t, err, "use a fresh result path for each attempt; reconcile old results before another submission")
	}
	record := func(stage string, r liveResult) {
		t.Helper()
		r.Time = time.Now().UTC().Format(time.RFC3339Nano)
		r.Stage = stage
		r.SDKRevision = os.Getenv("POLY_LIVE_SDK_REVISION")
		r.Amount = os.Getenv("POLY_LIVE_AMOUNT")
		r.MaxSpend = os.Getenv("POLY_LIVE_MAX_SPEND")
		r.Shares = os.Getenv("POLY_LIVE_SHARES")
		r.MaxPrice = os.Getenv("POLY_LIVE_MAX_PRICE")
		r.MinPrice = os.Getenv("POLY_LIVE_MIN_PRICE")
		r.Action = os.Getenv("POLY_LIVE_ACTION")
		if r.Action == "query" {
			if r.OrderID == "" {
				r.OrderID = os.Getenv("POLY_LIVE_ORDER_ID")
			}
			if r.Side == "" {
				r.Side = os.Getenv("POLY_LIVE_SIDE")
				if r.Side == "" {
					r.Side = string(ct.BUY)
				}
			}
		}
		r.WalletType = os.Getenv("POLY_LIVE_WALLET_TYPE")
		r.Wallet = os.Getenv("POLY_LIVE_WALLET")
		r.Protocol = os.Getenv("POLY_LIVE_PROTOCOL")
		r.Market = os.Getenv("POLY_LIVE_MARKET_ID")
		r.Outcome = os.Getenv("POLY_LIVE_OUTCOME_INDEX")
		data, err := json.Marshal(r)
		require.NoError(t, err)
		t.Logf("LIVE_RESULT %s", data)
		if file != nil {
			_, err = file.Write(append(data, '\n'))
			require.NoError(t, err)
			require.NoError(t, file.Sync())
		}
	}
	t.Cleanup(func() {
		if file != nil {
			defer file.Close()
		}
		stage := "passed"
		if t.Failed() {
			stage = "failed"
		} else if t.Skipped() {
			stage = "skipped"
		}
		record(stage, liveResult{})
	})
	record("started", liveResult{})
	return record
}

func recordFills(record recordResult, m selection, id string, side ct.Side, fills []ct.Trade) {
	seen := map[string]bool{}
	var hashes []string
	for _, fill := range fills {
		if !seen[fill.TransactionHash] {
			hashes = append(hashes, fill.TransactionHash)
			seen[fill.TransactionHash] = true
		}
	}
	record("confirmed", liveResult{Asset: m.asset, OrderID: id, Side: string(side), Hashes: hashes})
}

func TestLiveResultRefusesOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "result.jsonl")
	f, err := createResultFile(path)
	require.NoError(t, err)
	_, err = f.WriteString("previous attempt\n")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	_, err = createResultFile(path)
	require.ErrorIs(t, err, os.ErrExist)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "previous attempt\n", string(data))
}

// Query identity is durable in the first event, before market/API/fill checks can fail.
func TestLiveQueryResultIdentity(t *testing.T) {
	for _, side := range []string{"", "SELL"} {
		t.Run("side="+side, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "query.jsonl")
			t.Setenv("POLY_LIVE_RESULT_FILE", path)
			t.Setenv("POLY_LIVE_ACTION", "query")
			t.Setenv("POLY_LIVE_ORDER_ID", "existing-order")
			t.Setenv("POLY_LIVE_SIDE", side)
			t.Setenv("POLY_LIVE_PRIVATE_KEY", "must-not-appear")
			beginResult(t)
			data, err := os.ReadFile(path)
			require.NoError(t, err)
			require.NotContains(t, string(data), "must-not-appear")
			var row liveResult
			require.NoError(t, json.Unmarshal(data, &row))
			require.Equal(t, "started", row.Stage)
			require.Equal(t, "existing-order", row.OrderID)
			want := side
			if want == "" {
				want = "BUY"
			}
			require.Equal(t, want, row.Side)
		})
	}
}
