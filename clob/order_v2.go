package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// Protocol V2 reserves bits 40 through 103 of the uint256 position ID.
var v2PositionReservedBits = new(big.Int).Lsh(new(big.Int).SetUint64(math.MaxUint64), 40)

func (c *Client) CreateOrderV2(ctx context.Context, userOrder *types.UserOrderV2, orderType types.OrderType, postOnly, deferExec bool, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	if userOrder == nil {
		return nil, fmt.Errorf("order is required")
	}
	if err := c.validateOrderSigner(option); err != nil {
		return nil, err
	}
	if err := validateOrderInput(userOrder.TokenId, userOrder.Price, userOrder.Size, userOrder.Side, orderType, userOrder.TickSize); err != nil {
		return nil, err
	}
	if err := validateOptionalAmount("balance", userOrder.UserUsdcBalance); err != nil {
		return nil, err
	}
	if err := validateOptionalAmount("builder fee rate", userOrder.BuilderFeeRate); err != nil {
		return nil, err
	}
	if postOnly && (orderType == types.OrderTypeFOK || orderType == types.OrderTypeFAK) {
		return nil, fmt.Errorf("postOnly requires GTC or GTD")
	}

	tokenID := userOrder.TokenId

	market, err := c.orderMarketInfo(ctx, tokenID, userOrder.ConditionId)
	if err != nil {
		return nil, err
	}

	tickSize := utils.Float64ToDecimal(market.MinimumTickSize).String()
	if userOrder.TickSize != nil {
		tickSize = *userOrder.TickSize
	}

	if err := validateTick(tickSize); err != nil {
		return nil, err
	}
	tickSizeFloat64 := utils.StringToDecimal(tickSize).InexactFloat64()
	normalizedPrice := utils.NormalizePrice(userOrder.Price, tickSizeFloat64)
	if normalizedPrice != userOrder.Price {
		fmt.Printf(
			"price adjusted: origin=%f adjusted=%f (min=%f max=%f)",
			userOrder.Price,
			normalizedPrice,
			tickSizeFloat64,
			1-tickSizeFloat64,
		)

		if err = validateAdjustedLimitBuyUSDCBalance(userOrder, orderType, normalizedPrice); err != nil {
			return nil, err
		}
	}

	assetID, err := parseOrderAssetID(tokenID)
	if err != nil {
		return nil, err
	}
	contract := model.ExchangeV3
	if new(big.Int).And(assetID, v2PositionReservedBits).Sign() != 0 {
		negRisk := market.NegRisk
		if userOrder.NegRisk != nil {
			negRisk = userOrder.NegRisk
		}
		contract = model.CTFExchange
		if negRisk != nil && *negRisk {
			contract = model.NegRiskCTFExchange
		}
	}

	userOrder.Price = normalizedPrice
	c.adjustBuyMarketOrderAmountForFees(userOrder, market.FeeDetails, orderType)
	if !finite(userOrder.Size) || userOrder.Size <= 0 {
		return nil, fmt.Errorf("adjusted order size must be finite and positive")
	}
	if err := validateOrderAmounts(userOrder.Price, userOrder.Size, userOrder.Side, orderType, tickSize); err != nil {
		return nil, err
	}
	// The order-utils builder accepts decimal uint256s. CLOB uses the supplied asset ID.
	signingOrder := *userOrder
	signingOrder.TokenId = assetID.String()
	signedOrder, err := c.orderBuilder.buildOrderV2(signingOrder, orderType, types.TickSize(tickSize), contract, option)
	if err != nil {
		return nil, errors.WithMessage(err, "create order buildOrder")
	}

	payload := orderToJson2(signedOrder, tokenID, option.ApiKeyCreds.ApiKey, &orderType, postOnly, deferExec)
	if contract == model.ExchangeV3 {
		payload.Order.Signature = "0x" + payload.Order.Signature
	}
	return c.postOrderV2(ctx, payload, option)
}

func parseOrderAssetID(value string) (*big.Int, error) {
	value = strings.TrimSpace(value)
	// Keep unprefixed IDs decimal, including those with leading zeros.
	base := 10
	if len(value) >= 2 && value[0] == '0' && strings.ContainsAny(value[1:2], "xXbBoO") {
		base = 0
	}
	id, ok := new(big.Int).SetString(value, base)
	if !ok || id.Sign() < 0 || id.BitLen() > 256 || strings.Contains(value, "_") {
		return nil, fmt.Errorf("invalid uint256 asset ID %q", value)
	}
	return id, nil
}

func validateAdjustedLimitBuyUSDCBalance(userOrder *types.UserOrderV2, orderType types.OrderType, normalizedPrice float64) error {
	if userOrder.Price == normalizedPrice || userOrder.Side != types.BUY || userOrder.UserUsdcBalance == nil {
		return nil
	}
	if orderType != types.OrderTypeGTC && orderType != types.OrderTypeGTD {
		return nil
	}

	cost := utils.Float64ToDecimal(normalizedPrice).Mul(utils.Float64ToDecimal(userOrder.Size))
	budget := utils.Float64ToDecimal(*userOrder.UserUsdcBalance)
	if cost.GreaterThan(budget) {
		return fmt.Errorf("limit buy cost %s exceeds user USDC balance %s after price normalization", cost, budget)
	}
	return nil
}

func (c *Client) postOrderV2(ctx context.Context, orderPayload types.NewOrderV2, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	bodyBytes, err := json.Marshal(orderPayload)
	if err != nil {
		return nil, errors.WithMessage(err, "create order post order marshal")
	}
	bodyStr := string(bodyBytes)

	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodPost,
		RequestPath: types.POST_ORDER,
		Body:        bodyStr,
	}

	headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var out types.OrderResponse
	resp, err := c.client.DoRequest(ctx, http.MethodPost, types.POST_ORDER, &http2.RequestOptions{
		Headers: headers,
		Data:    bodyStr,
	}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
		return nil, e
	}
	if !orderPayload.DeferExec {
		c.resolveTransactionsHashes(ctx, &out, option)
	}
	return &out, nil
}

func orderToJson2(order *model.SignedOrderV2, tokenID, owner string, orderType *types.OrderType, postOnly, deferExec bool) types.NewOrderV2 {
	side := types.BUY
	if order.Side.Cmp(big.NewInt(0)) == 1 {
		side = types.SELL
	}
	return types.NewOrderV2{
		Order: types.OrderV2{
			Salt:          order.Salt.Int64(),
			Maker:         order.Maker.String(),
			Signer:        order.Signer.String(),
			TokenID:       tokenID,
			MakerAmount:   order.MakerAmount.String(),
			TakerAmount:   order.TakerAmount.String(),
			Side:          side,
			SignatureType: model.SignatureType(order.SignatureType.Uint64()),
			Timestamp:     order.Timestamp.String(),
			Expiration:    order.Expiration.String(),
			Metadata:      order.Metadata.String(),
			Builder:       order.Builder.String(),
			Signature:     common.Bytes2Hex(order.Signature),
		},
		Owner:     owner,
		OrderType: *orderType,
		DeferExec: deferExec,
		PostOnly:  postOnly,
	}
}

func (c *Client) adjustBuyMarketOrderAmountForFees(order *types.UserOrderV2, fees *types.CLOBMarketInfoFeeDetails, orderType types.OrderType) {
	if order == nil {
		return
	}

	if order.Side != types.BUY {
		return
	}

	if orderType != types.OrderTypeFAK && orderType != types.OrderTypeFOK {
		return
	}

	if order.UserUsdcBalance == nil {
		return
	}

	price := order.Price
	userUSDCBalance := *order.UserUsdcBalance

	builderTakerFeeRate := 0.0
	if c.isBuilderOrder(order.BuilderCode) && order.BuilderFeeRate != nil && *order.BuilderFeeRate > 0 {
		builderTakerFeeRate = *order.BuilderFeeRate
	}

	if fees == nil {
		fees = &types.CLOBMarketInfoFeeDetails{}
	}

	order.Size = adjustBuyAmountForFees(
		order.Size,
		price,
		userUSDCBalance,
		fees.Rate,
		fees.Exponent,
		builderTakerFeeRate,
	)
	return
}

func adjustBuyAmountForFees(
	amount float64,
	price float64,
	userUSDCBalance float64,
	feeRate float64,
	feeExponent float64,
	builderTakerFeeRate float64,
) float64 {
	platformFeeRate := feeRate * math.Pow(price*(1-price), feeExponent)
	platformFee := (amount / price) * platformFeeRate
	totalCost := amount + platformFee + amount*builderTakerFeeRate

	if userUSDCBalance <= totalCost {
		adjustAmount := userUSDCBalance / (1 + platformFeeRate/price + builderTakerFeeRate)
		if adjustAmount >= 1 {
			return adjustAmount
		}
		return userUSDCBalance / (1 + platformFeeRate/price + builderTakerFeeRate)
	}
	return amount
}

func (c *Client) isBuilderOrder(builderCode *string) bool {
	return builderCode != nil && *builderCode != types.Bytes32Zero
}
