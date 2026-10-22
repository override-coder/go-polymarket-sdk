package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/override-coder/go-polymarket-sdk/gamma/types"
)

func TestMarketComboFields(t *testing.T) {
	const positionID = "548324018002783558052327573095137442627908907682778868240145789358304133121"
	for _, tc := range []struct {
		name   string
		body   string
		ids    []string
		status *string
	}{
		{name: "enabled", body: `{"positionIds":["` + positionID + `"],"comboStatus":"enabled"}`, ids: []string{positionID}, status: stringPointer("enabled")},
		{name: "unknown status preserved", body: `{"positionIds":[],"comboStatus":"future-status"}`, ids: []string{}, status: stringPointer("future-status")},
		{name: "missing", body: `{}`},
		{name: "null", body: `{"positionIds":null,"comboStatus":null}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var market types.Market
			if err := json.Unmarshal([]byte(tc.body), &market); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(market.PositionIDs, tc.ids) || !reflect.DeepEqual(market.ComboStatus, tc.status) {
				t.Fatalf("unexpected combo fields: %v, %v", market.PositionIDs, market.ComboStatus)
			}
			encoded, err := json.Marshal(market)
			if err != nil {
				t.Fatal(err)
			}
			var roundTrip types.Market
			if err := json.Unmarshal(encoded, &roundTrip); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(roundTrip.PositionIDs, tc.ids) || !reflect.DeepEqual(roundTrip.ComboStatus, tc.status) {
				t.Fatal("JSON round trip lost combo fields")
			}
		})
	}
}

func stringPointer(s string) *string { return &s }
