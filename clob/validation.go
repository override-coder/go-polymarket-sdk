package clob

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"

	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

var maxOrderUnits = decimal.NewFromBigInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)), 0)

// Keep authentication validation separate from order signing: read/cancel calls
// with existing API credentials do not require a private-key signer.
func validateAuthOption(option *sdktypes.AuthOption, requireCredentials bool) error {
	if option == nil {
		return fmt.Errorf("auth option is required")
	}
	if requireCredentials && option.ApiKeyCreds == nil {
		return fmt.Errorf("API credentials are required")
	}
	return nil
}

func (c *Client) validateOrderSigner(option *sdktypes.AuthOption) error {
	if err := validateAuthOption(option, true); err != nil {
		return err
	}
	if c.signFn == nil {
		return fmt.Errorf("signature function is required")
	}
	if c.chainId == nil || c.chainId.Sign() <= 0 {
		return fmt.Errorf("positive chain ID is required")
	}
	return nil
}

func validateOrderInput(asset string, price, size float64, side types.Side, orderType types.OrderType, tick *string) error {
	if _, err := parseOrderAssetID(asset); err != nil {
		return err
	}
	if !finite(price) || price <= 0 || price >= 1 {
		return fmt.Errorf("price must be finite and between 0 and 1")
	}
	if !finite(size) || size <= 0 {
		return fmt.Errorf("size must be finite and positive")
	}
	// Bound inputs before float arithmetic (e.g. market BUY size / price).
	// A share or collateral amount is represented by a uint256 with six decimals.
	if decimal.NewFromFloat(size).Shift(6).GreaterThan(maxOrderUnits) {
		return fmt.Errorf("size exceeds uint256 order amount range")
	}
	if side != types.BUY && side != types.SELL {
		return fmt.Errorf("side must be BUY or SELL")
	}
	switch orderType {
	case types.OrderTypeGTC, types.OrderTypeGTD, types.OrderTypeFOK, types.OrderTypeFAK:
	default:
		return fmt.Errorf("unsupported order type %q", orderType)
	}
	if tick != nil {
		return validateTick(*tick)
	}
	return nil
}

func validateTick(tick string) error {
	if _, ok := roundingConfig[types.TickSize(tick)]; !ok {
		return fmt.Errorf("unsupported tick size %q", tick)
	}
	return nil
}

func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func validateOptionalAmount(name string, value *float64) error {
	if value != nil && (!finite(*value) || *value < 0) {
		return fmt.Errorf("%s must be finite and non-negative", name)
	}
	return nil
}

// validateOrderAmounts checks the actual rounded values after price normalization
// and fee adjustment, before invoking the signer or serializing uint256 amounts.
func validateOrderAmounts(price, size float64, side types.Side, orderType types.OrderType, tick string) error {
	config := roundingConfig[types.TickSize(tick)]
	var maker, taker float64
	if orderType == types.OrderTypeFOK || orderType == types.OrderTypeFAK {
		_, maker, taker = getMarketOrderRawAmounts(side, size, price, config)
	} else {
		_, maker, taker = getOrderRawAmounts(side, size, price, config)
	}
	for _, amount := range []float64{maker, taker} {
		if !finite(amount) || amount <= 0 || decimal.NewFromFloat(amount).Shift(6).GreaterThan(maxOrderUnits) {
			return fmt.Errorf("rounded order amounts must be positive uint256 values")
		}
	}
	return nil
}
