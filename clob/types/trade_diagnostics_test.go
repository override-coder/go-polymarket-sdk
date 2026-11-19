package types

import (
	"encoding/json"
	"testing"
)

func TestTradeDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		body      string
		populated bool
	}{
		{`{}`, false},
		{`{"err_msg":null,"match_time_nano":null}`, false},
		{`{"err_msg":"settlement failed","match_time_nano":"1795090000123456789"}`, true},
	} {
		var trades Trades
		if e := json.Unmarshal([]byte(`{"data":[`+tc.body+`],"next_cursor":"LTE="}`), &trades); e != nil {
			t.Fatal(e)
		}
		got := trades.Data[0]
		if tc.populated {
			if got.ErrMsg == nil || *got.ErrMsg != "settlement failed" || got.MatchTimeNano == nil || *got.MatchTimeNano != "1795090000123456789" {
				t.Fatal(got)
			}
		} else if got.ErrMsg != nil || got.MatchTimeNano != nil {
			t.Fatal("missing fields defaulted")
		}
		encoded, e := json.Marshal(got)
		if e != nil {
			t.Fatal(e)
		}
		var roundTrip Trade
		if e = json.Unmarshal(encoded, &roundTrip); e != nil {
			t.Fatal(e)
		}
		if got.ErrMsg == nil {
			if roundTrip.ErrMsg != nil || roundTrip.MatchTimeNano != nil {
				t.Fatal("missing fields defaulted")
			}
		} else if *roundTrip.ErrMsg != *got.ErrMsg || *roundTrip.MatchTimeNano != *got.MatchTimeNano {
			t.Fatal("diagnostics changed")
		}
	}
}
