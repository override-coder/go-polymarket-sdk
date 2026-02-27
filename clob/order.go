package clob

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/override-coder/go-polymarket-sdk/clob/types"
	sdkheaders "github.com/override-coder/go-polymarket-sdk/headers"
	http2 "github.com/override-coder/go-polymarket-sdk/http"
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

	tickSize := ""
	if userOrder.TickSize != nil {
		tickSize = *userOrder.TickSize
	} else {
		tickSz, err := c.GetTickSize(tokenID)
		if err != nil {
			return nil, errors.WithMessage(err, "create order get tickSize")
		}
		tickSize = tickSz
	}

	feeRateBps := float64(0)
	if userOrder.FeeRateBps != nil {
		feeRateBps = *userOrder.FeeRateBps
	} else {
		bps, err := c.GetFeeRateBps(tokenID)
		if err != nil {
			return nil, errors.WithMessage(err, "create order get feeRateBps")
		}
		feeRateBps = bps
	}

	userOrder.FeeRateBps = &feeRateBps
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
		risk, err := c.GetNegRisk(tokenID)
		if err != nil {
			return nil, errors.WithMessage(err, "create order get negRisk")
		}
		negRisk = risk
	}

	signedOrder, err := c.orderBuilder.buildOrder(userOrder, orderType, types.CreateOrderOptions{
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
	resp, err := c.client.DoRequest(http.MethodPost, types.POST_ORDER, &http2.RequestOptions{
		Headers: headers,
		Data:    bodyStr,
	}, &out)
	if _, e := http2.ParseHTTPError(resp, err); e != nil {
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

func (c *Client) GetOrder(orderId string, req types.GetOrderRequest, option *sdktypes.AuthOption) (*types.OpenOrder, error) {
	requestPath := fmt.Sprintf("%s%s", types.GET_ORDER, orderId)

	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: requestPath,
	}

	params := make(map[string]any, 3)
	if req.ID != "" {
		params["id"] = req.ID
	}

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.OpenOrder
	res, err := c.client.DoRequest(http.MethodGet, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "get order buy id:%v", params)
	}
	return &resp, nil
}

func (c *Client) GetOrders(req types.GetActiveOrdersRequest, option *sdktypes.AuthOption) (*types.OpenOrders, error) {
	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodGet,
		RequestPath: types.GET_OPEN_ORDERS,
	}
	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	params := make(map[string]any, 3)
	if req.ID != "" {
		params["id"] = req.ID
	}
	if req.Market != "" {
		params["market"] = req.Market
	}
	if req.AssetID != "" {
		params["asset_id"] = req.AssetID
	}

	var resp types.OpenOrders
	res, err := c.client.DoRequest(http.MethodGet, types.GET_OPEN_ORDERS, &http2.RequestOptions{
		Headers: l2Headers,
		Params:  params,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "get orders")
	}
	return &resp, nil
}

func (c *Client) CancelOrder(orderId string, option *sdktypes.AuthOption) (*types.CancelOrder, error) {
	if orderId == "" {
		return nil, errors.New("order id is empty")
	}
	requestPath := fmt.Sprintf("%s", types.CANCEL_ORDER)

	ts := time.Now().Unix()

	params := make(map[string]any, 1)
	params["orderID"] = orderId

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.WithMessage(err, "cancel order post order marshal")
	}
	bodyStr := string(bodyBytes)

	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodDelete,
		RequestPath: requestPath,
		Body:        bodyStr,
	}

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.CancelOrder
	res, err := c.client.DoRequest(http.MethodDelete, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Data:    bodyStr,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "cancel order buy id:%v", orderId)
	}
	return &resp, nil
}

func (c *Client) CancelOrders(orderId []string, option *sdktypes.AuthOption) (*types.CancelOrder, error) {
	if len(orderId) == 0 {
		return nil, errors.New("orderId is empty")
	}
	requestPath := fmt.Sprintf("%s", types.CANCEL_ORDERS)

	ts := time.Now().Unix()

	bodyBytes, err := json.Marshal(orderId)
	if err != nil {
		return nil, errors.WithMessage(err, "cancel order post order marshal")
	}
	bodyStr := string(bodyBytes)

	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodDelete,
		RequestPath: requestPath,
		Body:        bodyStr,
	}

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.CancelOrder
	res, err := c.client.DoRequest(http.MethodDelete, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Data:    bodyStr,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "cancel orders buy ids:%v", orderId)
	}
	return &resp, nil
}

func (c *Client) CancelOrderAll(option *sdktypes.AuthOption) (*types.CancelOrder, error) {
	requestPath := fmt.Sprintf("%s", types.CANCEL_ALL)
	ts := time.Now().Unix()
	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodDelete,
		RequestPath: requestPath,
	}

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.CancelOrder
	res, err := c.client.DoRequest(http.MethodDelete, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrap(e, "cancel order all")
	}
	return &resp, nil
}

func (c *Client) CancelOrderByMarket(request *types.CancelOrderRequest, option *sdktypes.AuthOption) (*types.CancelOrder, error) {
	requestPath := fmt.Sprintf("%s", types.CANCEL_MARKET_ORDERS)
	ts := time.Now().Unix()

	params := make(map[string]any, 2)
	if request.ConditionID != nil && *request.ConditionID != "" {
		params["market"] = *request.ConditionID
	}
	if request.AssetID != nil && *request.AssetID != "" {
		params["asset_id"] = *request.AssetID
	}

	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.WithMessage(err, "cancel order post order marshal")
	}
	bodyStr := string(bodyBytes)

	l2HeaderArgs := types.L2HeaderArgs{
		Method:      http.MethodDelete,
		RequestPath: requestPath,
		Body:        bodyStr,
	}

	l2Headers, err := sdkheaders.CreateL2Headers(option.SingerAddress, option.ApiKeyCreds, l2HeaderArgs, &ts)
	if err != nil {
		return nil, errors.WithMessage(err, "create l2 headers")
	}

	var resp types.CancelOrder
	res, err := c.client.DoRequest(http.MethodDelete, requestPath, &http2.RequestOptions{
		Headers: l2Headers,
		Data:    bodyStr,
	}, &resp)
	if _, e := http2.ParseHTTPError(res, err); e != nil {
		return nil, errors.Wrapf(e, "cancel order by market:%v", params)
	}
	return &resp, nil
}
