//go:build polymarket_live

package live_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const gammaHost = "https://gamma-api.polymarket.com"

type market struct {
	ID                                       string
	Version                                  string
	ConditionID                              string
	PositionIDs                              []string
	ClobTokenIDs                             string
	Outcomes                                 string
	OutcomePrices                            string
	ResolutionStatus                         string
	Active, Closed, AcceptingOrders, NegRisk bool
}

func gammaGet(t *testing.T, path string, result any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gammaHost+path, nil)
	require.NoError(t, err)
	// Go's default transport requests and decompresses gzip for large event pages.
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode, path)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(result))
}

func assetIDs(m market) ([]string, error) {
	switch m.Version {
	case "v2":
		return m.PositionIDs, nil
	case "v1":
		var ids []string
		err := json.Unmarshal([]byte(m.ClobTokenIDs), &ids)
		return ids, err
	default:
		return nil, fmt.Errorf("unknown market version %q", m.Version)
	}
}

func TestDiscoverMarkets(t *testing.T) {
	if os.Getenv("POLY_LIVE_DISCOVER") != "1" {
		t.Skip("set POLY_LIVE_DISCOVER=1 for public Gamma discovery")
	}
	for _, version := range []string{"v1", "v2"} {
		t.Run(version, func(t *testing.T) {
			var events []struct{ Markets []market }
			gammaGet(t, "/events?version="+version+"&closed=false&limit=100&offset=0", &events)
			count := 0
			for _, e := range events {
				for _, m := range e.Markets {
					if m.Version != version || !m.Active || m.Closed || !m.AcceptingOrders || (version == "v2" && m.ResolutionStatus != "active") {
						continue
					}
					ids, err := assetIDs(m)
					require.NoError(t, err)
					if len(ids) != 2 {
						continue
					}
					if count < 5 {
						t.Logf("market=%s outcomes=%s assets=%v", m.ID, m.Outcomes, ids)
					}
					count++
				}
			}
			require.Positive(t, count, "no eligible binary market on the first page")
			t.Logf("%d open binary markets on page", count)
		})
	}
}

func env(t *testing.T, key string) string {
	t.Helper()
	v := strings.TrimSpace(os.Getenv("POLY_LIVE_" + key))
	if v == "" {
		t.Fatalf("set POLY_LIVE_%s", key)
	}
	return v
}

type selection struct {
	market
	index     int
	asset     string
	number    *big.Int
	condition []byte
}

func selected(t *testing.T) selection {
	t.Helper()
	id := env(t, "MARKET_ID")
	_, err := strconv.ParseUint(id, 10, 64)
	require.NoError(t, err)
	var m market
	gammaGet(t, "/markets/"+url.PathEscape(id), &m)
	require.Equal(t, id, m.ID)
	require.Equal(t, env(t, "PROTOCOL"), m.Version)
	ids, err := assetIDs(m)
	require.NoError(t, err)
	require.Len(t, ids, 2)
	i, err := strconv.Atoi(env(t, "OUTCOME_INDEX"))
	require.NoError(t, err)
	require.Contains(t, []int{0, 1}, i)
	n, ok := new(big.Int).SetString(ids[i], 10)
	if !ok && strings.HasPrefix(ids[i], "0x") {
		n, ok = new(big.Int).SetString(ids[i][2:], 16)
	}
	require.True(t, ok)
	require.Positive(t, n.Sign())
	require.LessOrEqual(t, n.BitLen(), 256)
	condition, err := hex.DecodeString(strings.TrimPrefix(m.ConditionID, "0x"))
	require.NoError(t, err)
	raw := n.FillBytes(make([]byte, 32))
	isV2 := new(big.Int).SetBytes(raw[19:27]).Sign() == 0 // reserved bits 40..103
	require.Equal(t, m.Version == "v2", isV2, "asset namespace disagrees with Gamma")
	if isV2 {
		if len(condition) == 32 {
			require.LessOrEqual(t, condition[31], byte(1))
			condition = condition[:31]
		}
		require.Equal(t, raw[:31], condition)
		require.Equal(t, byte(i), raw[31])
	} else {
		require.Len(t, condition, 32)
	}
	return selection{m, i, ids[i], n, condition}
}
