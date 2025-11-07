package dataapi_test

import (
	"github.com/override-coder/go-polymarket-sdk/dataapi"
	"github.com/override-coder/go-polymarket-sdk/dataapi/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	PolymarketRelayURL = "https://data-api.polymarket.com"
	chaindId           = big.NewInt(137)
)

func TestGetPositions(t *testing.T) {
	client := dataapi.NewClient(PolymarketRelayURL, chaindId)

	positions, err := client.GetPositions(types.PositionsQuery{
		User: "0x0f863d92dd2b960e3eb6a23a35fd92a91981404e",
	})
	assert.Equal(t, nil, err)
	t.Logf("positions: %v", positions)

	activity, err := client.GetUserActivity(types.ActivityQuery{
		User: "0x0f863d92dd2b960e3eb6a23a35fd92a91981404e",
	})
	assert.Equal(t, nil, err)
	t.Logf("activitys: %v", activity)

}
