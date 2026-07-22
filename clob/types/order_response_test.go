package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestOrderResponseUnmarshalPreservesTradeIDs(t *testing.T) {
	var response OrderResponse
	if err := json.Unmarshal([]byte(`{
		"success": true,
		"orderID": "order-1",
		"transactionsHashes": ["0xabc"],
		"tradeIDs": ["trade-1", "trade-2"],
		"status": "matched",
		"takingAmount": "10",
		"makingAmount": "5"
	}`), &response); err != nil {
		t.Fatalf("unmarshal order response: %v", err)
	}

	if want := []string{"trade-1", "trade-2"}; !reflect.DeepEqual(response.TradeIDs, want) {
		t.Fatalf("trade IDs = %v, want %v", response.TradeIDs, want)
	}
}
