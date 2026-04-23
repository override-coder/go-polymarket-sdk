package clob

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
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

func (c *Client) CreateOrderV2(ctx context.Context, userOrder *types.UserOrderV2, orderType types.OrderType, postOnly, deferExec bool, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	tokenID := userOrder.TokenId

	err := c.ensureMarketInfoCached(ctx, tokenID, userOrder.ConditionId)
	if err != nil {
		return nil, err
	}

	tickSize := ""
	if userOrder.TickSize != nil {
		tickSize = *userOrder.TickSize
	} else {
		tickSz, err := c.GetTickSize(ctx, tokenID)
		if err != nil {
			return nil, errors.WithMessage(err, "create order get tickSize")
		}
		tickSize = tickSz
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
	}

	userOrder.Price = normalizedPrice

	negRisk := false
	if userOrder.NegRisk != nil {
		negRisk = *userOrder.NegRisk
	} else {
		risk, err := c.GetNegRisk(ctx, tokenID)
		if err != nil {
			return nil, errors.WithMessage(err, "create order get negRisk")
		}
		negRisk = risk
	}

	c.adjustBuyMarketOrderAmountForFees(userOrder, tokenID, orderType)

	signedOrder, err := c.orderBuilder.buildOrderV2(*userOrder, orderType, types.CreateOrderOptions{
		AuthOption: option,
		TickSize:   types.TickSize(tickSize),
		NegRisk:    negRisk,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "create order buildOrder")
	}

	return c.postOrderV2(ctx, signedOrder, orderType, postOnly, deferExec, option)

}

func (c *Client) postOrderV2(ctx context.Context, order *model.SignedOrderV2, orderType types.OrderType, postOnly, deferExec bool, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	orderPayload := orderToJson2(order, option.ApiKeyCreds.ApiKey, &orderType, postOnly, deferExec)
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
	return &out, nil
}

func orderToJson2(order *model.SignedOrderV2, owner string, orderType *types.OrderType, postOnly, deferExec bool) types.NewOrderV2 {
	side := types.BUY
	if order.Side.Cmp(big.NewInt(0)) == 1 {
		side = types.SELL
	}
	return types.NewOrderV2{
		Order: types.OrderV2{
			Salt:          order.Salt.Int64(),
			Maker:         order.Maker.String(),
			Signer:        order.Signer.String(),
			TokenID:       order.TokenID.String(),
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

func (c *Client) adjustBuyMarketOrderAmountForFees(order *types.UserOrderV2, tokenID string, orderType types.OrderType) {
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

	feeInfo := c.getCachedFeeInfo(tokenID)

	order.Size = adjustBuyAmountForFees(
		order.Size,
		price,
		userUSDCBalance,
		feeInfo.Rate,
		feeInfo.Exponent,
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
	}
	return amount
}

func (c *Client) isBuilderOrder(builderCode *string) bool {
	return builderCode != nil && *builderCode != types.Bytes32Zero
}
