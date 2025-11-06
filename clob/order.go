package clob

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"github.com/pkg/errors"
	"github.com/polymarket/go-order-utils/pkg/model"
	"math/big"
	"net/http"
	"time"
)

func (c *Client) CreateOrder(userOrder types.UserOrder, orderType types.OrderType, deferExec bool, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	tokenID := userOrder.TokenID

	tickSize, err := c.GetTickSize(tokenID)
	if err != nil {
		return nil, errors.WithMessage(err, "create order get tickSize")
	}
	feeRateBps, err := c.GetFeeRateBps(tokenID)
	if err != nil {
		return nil, errors.WithMessage(err, "create order get feeRateBps")
	}
	userOrder.FeeRateBps = &feeRateBps

	tickSizeFloat64 := utils.StringToDecimal(tickSize).InexactFloat64()
	if !utils.PriceValid(userOrder.Price, tickSizeFloat64) {
		return nil, errors.Errorf("invalid price %f, min: %f - max: %f", userOrder.Price, tickSizeFloat64, 1-tickSizeFloat64)
	}

	negRisk, err := c.GetNegRisk(tokenID)
	if err != nil {
		return nil, errors.WithMessage(err, "create order get negRisk")
	}

	signedOrder, err := c.orderBuilder.buildOrder(userOrder, types.CreateOrderOptions{
		AuthOption: option,
		TickSize:   types.TickSize(tickSize),
		NegRisk:    negRisk,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "create order buildOrder")
	}

	return c.postOrder(signedOrder, orderType, deferExec, option)
}

func (c *Client) postOrder(order *model.SignedOrder, orderType types.OrderType, deferExec bool, option *sdktypes.AuthOption) (*types.OrderResponse, error) {
	orderPayload := orderToJson(order, option.ApiKeyCreds.ApiKey, &orderType, deferExec)
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

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}
	headers := l2Headers
	if c.builderApiKeyCreds != nil {
		builderHeaders, errBuilder := sdkheaders.CreateL2BuilderHeaders(c.builderApiKeyCreds, l2HeaderArgs, &ts)
		if errBuilder == nil && len(builderHeaders) != 0 {
			headers = sdkheaders.InjectBuilderHeaders(l2Headers, builderHeaders)
		}
	}

	var out types.OrderResponse
	resp, err := c.doRequest(http.MethodPost, types.POST_ORDER, &RequestOptions{
		Headers: headers,
		Data:    bodyStr,
	}, &out)
	if _, e := parseHTTPError(resp, err); e != nil {
		return nil, e
	}
	return &out, nil
}

func orderToJson(order *model.SignedOrder, owner string, orderType *types.OrderType, deferExec bool) types.NewOrder {
	side := types.BUY
	if order.Side.Cmp(big.NewInt(0)) == 1 {
		side = types.SELL
	}
	return types.NewOrder{
		Order: types.Order{
			Salt:          order.Salt.Int64(),
			Maker:         order.Maker.String(),
			Signer:        order.Signer.String(),
			Taker:         order.Taker.String(),
			TokenID:       order.TokenId.String(),
			MakerAmount:   order.MakerAmount.String(),
			TakerAmount:   order.TakerAmount.String(),
			Expiration:    order.Expiration.String(),
			Nonce:         order.Nonce.String(),
			FeeRateBps:    order.FeeRateBps.String(),
			Side:          side,
			SignatureType: model.SignatureType(order.SignatureType.Uint64()),
			Signature:     common.Bytes2Hex(order.Signature),
		},
		Owner:     owner,
		OrderType: *orderType,
		DeferExec: deferExec,
	}
}
