package dataapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/require"
)

func TestV2ValueSingleObjectAndZero(t *testing.T) {
	for _, filtered := range []bool{false, true} {
		t.Run(fmt.Sprint(filtered), func(t *testing.T) {
			expected := url.Values{"user": {v2User}}
			q := types.PositionValueQuery{User: " " + v2User + " "}
			value := 0.0
			if filtered {
				q.Condition = []string{v2Condition}
				expected.Set("condition", v2Condition)
				value = 12.3456
			}
			body := fmt.Sprintf(`{"data":{"proxy_wallet":%q,"value":%v}}`, v2User, value)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/v2/value", r.URL.Path)
				require.Equal(t, expected, r.URL.Query())
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, body)
			}))
			defer server.Close()
			result, err := dataapi.NewClient(server.URL, nil).GetPositionValue(context.Background(), q)
			require.NoError(t, err)
			require.Equal(t, v2User, result.Data.ProxyWallet)
			require.Equal(t, value, result.Data.Value)
			assertV2JSON(t, body, result)
		})
	}
}

func TestV2ValueValidation(t *testing.T) {
	many := make([]string, 21)
	for i := range many {
		many[i] = fmt.Sprintf("0x%064x", i)
	}
	for name, q := range map[string]types.PositionValueQuery{
		"missing user": {}, "invalid condition": {User: v2User, Condition: []string{"bad"}}, "too many conditions": {User: v2User, Condition: many},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := dataapi.NewClient("http://127.0.0.1:1", nil).GetPositionValue(context.Background(), q)
			require.Error(t, err)
		})
	}
}
