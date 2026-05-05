package clob

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/override-coder/go-polymarket-sdk/clob/types"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/override-coder/go-polymarket-sdk/types/utils"
	"github.com/polymarket/go-order-utils/pkg/builder"
	"github.com/polymarket/go-order-utils/pkg/model"
)

// RoundConfig
type RoundConfig struct {
	Price  int
	Size   int
	Amount int
}

// roundingConfig
var roundingConfig = map[types.TickSize]RoundConfig{
	types.TickSize01: {
		Price:  1,
		Size:   2,
		Amount: 3,
	},
	types.TickSize001: {
		Price:  2,
		Size:   2,
		Amount: 4,
	},
	types.TickSize0001: {
		Price:  3,
		Size:   2,
		Amount: 5,
	},
	types.TickSize00001: {
		Price:  4,
		Size:   2,
		Amount: 6,
	},
}

type OrderBuilder struct {
	chaindId *big.Int
	signFn   signing.SignatureFunc
}

func NewOrderBuilder(chaindId *big.Int, signFn signing.SignatureFunc) *OrderBuilder {
	return &OrderBuilder{
		chaindId: chaindId,
		signFn:   signFn,
	}
}

func (o *OrderBuilder) WithSignatureFunc(signFn signing.SignatureFunc) error {
	if o.signFn != nil {
		return errors.New("signFn already set")
	}
	o.signFn = signFn
	return nil
}

func (o *OrderBuilder) buildOrder(order types.UserOrder, orderType types.OrderType, options types.CreateOrderOptions) (*model.SignedOrder, error) {
	return o.createOrder(order, orderType, options)
}

func (o *OrderBuilder) buildOrderV2(order types.UserOrderV2, orderType types.OrderType, options types.CreateOrderOptions) (*model.SignedOrderV2, error) {
	return o.createOrderV2(order, orderType, options)
}

func (o *OrderBuilder) createOrder(order types.UserOrder, orderType types.OrderType, options types.CreateOrderOptions) (*model.SignedOrder, error) {
	orderData := o.buildOrderCreationArgs(order, orderType, roundingConfig[options.TickSize], options.AuthOption)
	exchangeContract := model.CTFExchange
	if options.NegRisk {
		exchangeContract = model.NegRiskCTFExchange
	}
	return buildOrder(o.signFn, exchangeContract, o.chaindId, orderData)
}

func (o *OrderBuilder) createOrderV2(order types.UserOrderV2, orderType types.OrderType, options types.CreateOrderOptions) (*model.SignedOrderV2, error) {
	orderData := o.buildOrderCreationArgsV2(order, orderType, roundingConfig[options.TickSize], options.AuthOption)
	exchangeContract := model.CTFExchange
	if options.NegRisk {
		exchangeContract = model.NegRiskCTFExchange
	}
	return buildOrderV2(o.signFn, exchangeContract, o.chaindId, orderData)
}

func buildOrder(signFn signing.SignatureFunc, exchangeAddress model.VerifyingContract, chainId *big.Int, orderData *model.OrderData) (*model.SignedOrder, error) {
	cTFExchangeOrderBuilder := builder.NewExchangeOrderBuilderImpl(chainId, nil)
	order, err := cTFExchangeOrderBuilder.BuildOrder(orderData)
	if err != nil {
		return nil, err
	}
	orderHash, err := cTFExchangeOrderBuilder.BuildOrderHash(order, exchangeAddress)
	if err != nil {
		return nil, err
	}
	signature, err := signFn(order.Signer.String(), func(key *ecdsa.PrivateKey) ([]byte, error) {
		return cTFExchangeOrderBuilder.BuildOrderSignature(key, orderHash)
	})
	if err != nil {
		return nil, err
	}
	return &model.SignedOrder{
		Order:     *order,
		Signature: signature,
	}, nil
}

func buildOrderV2(signFn signing.SignatureFunc, exchangeAddress model.VerifyingContract, chainId *big.Int, orderData *model.OrderDataV2) (*model.SignedOrderV2, error) {
	cTFExchangeOrderBuilder := builder.NewExchangeOrderBuilderImplV2(chainId, nil)
	signedOrderV2 := new(model.SignedOrderV2)
	_, err := signFn(orderData.Signer, func(key *ecdsa.PrivateKey) ([]byte, error) {
		order, err := cTFExchangeOrderBuilder.BuildSignedOrder(key, orderData, exchangeAddress)
		if err != nil {
			return nil, err
		}
		signedOrderV2 = order
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return signedOrderV2, nil
}

func (o *OrderBuilder) buildOrderCreationArgs(order types.UserOrder, orderType types.OrderType, roundConfig RoundConfig, option *sdktypes.AuthOption) *model.OrderData {
	var (
		side     model.Side
		makerAmt float64
		takerAmt float64
	)
	if orderType == types.OrderTypeFAK || orderType == types.OrderTypeFOK {
		side, makerAmt, takerAmt = getMarketOrderRawAmounts(order.Side, order.Size, order.Price, roundConfig)
	} else {
		side, makerAmt, takerAmt = getOrderRawAmounts(order.Side, order.Size, order.Price, roundConfig)
	}
	makerAmount := utils.Pow(utils.Float64ToDecimal(makerAmt), types.CollateralTokenDecimals)
	takerAmount := utils.Pow(utils.Float64ToDecimal(takerAmt), types.ConditionalTokenDecimals)

	maker := option.SingerAddress
	if !strings.EqualFold(option.FunderAddress, "") {
		maker = option.FunderAddress
	}
	var (
		taker      = sdktypes.ZeroAddress
		feeRateBps = "0"
		nonce      = "0"
		expiration = "0"
	)
	if order.Taker != nil && *order.Taker != "" {
		taker = *order.Taker
	}
	if order.FeeRateBps != nil && *order.FeeRateBps > 0 {
		feeRateBps = utils.Float64ToDecimal(*order.FeeRateBps).String()
	}
	if order.Nonce != nil && *order.Nonce > 0 {
		nonce = fmt.Sprintf("%d", *order.Nonce)
	}
	if order.Expiration != nil && *order.Expiration > 0 {
		expiration = fmt.Sprintf("%d", *order.Expiration)
	}

	return &model.OrderData{
		Maker:         maker,
		Taker:         taker,
		TokenId:       order.TokenID,
		MakerAmount:   makerAmount.String(),
		TakerAmount:   takerAmount.String(),
		FeeRateBps:    feeRateBps,
		Nonce:         nonce,
		Signer:        option.SingerAddress,
		Expiration:    expiration,
		Side:          side,
		SignatureType: option.SignatureType,
	}
}

func (o *OrderBuilder) buildOrderCreationArgsV2(order types.UserOrderV2, orderType types.OrderType, roundConfig RoundConfig, option *sdktypes.AuthOption) *model.OrderDataV2 {
	var (
		side     model.Side
		makerAmt float64
		takerAmt float64
	)
	if orderType == types.OrderTypeFAK || orderType == types.OrderTypeFOK {
		side, makerAmt, takerAmt = getMarketOrderRawAmounts(order.Side, order.Size, order.Price, roundConfig)
	} else {
		side, makerAmt, takerAmt = getOrderRawAmounts(order.Side, order.Size, order.Price, roundConfig)
	}
	makerAmount := utils.Pow(utils.Float64ToDecimal(makerAmt), types.CollateralTokenDecimals)
	takerAmount := utils.Pow(utils.Float64ToDecimal(takerAmt), types.ConditionalTokenDecimals)

	maker := option.SingerAddress
	if !strings.EqualFold(option.FunderAddress, "") {
		maker = option.FunderAddress
	}
	var (
		builder    = types.Bytes32Zero
		metadata   = types.Bytes32Zero
		expiration = "0"
	)

	if order.BuilderCode != nil && !strings.EqualFold(*order.BuilderCode, types.Bytes32Zero) {
		builder = *order.BuilderCode
	}

	if order.Metadata != nil && !strings.EqualFold(*order.Metadata, types.Bytes32Zero) {
		metadata = *order.Metadata
	}

	if order.Expiration != nil && *order.Expiration > 0 {
		expiration = fmt.Sprintf("%d", *order.Expiration)
	}

	return &model.OrderDataV2{
		Maker:         maker,
		TokenID:       order.TokenId,
		MakerAmount:   makerAmount.String(),
		TakerAmount:   takerAmount.String(),
		Side:          side,
		Signer:        option.SingerAddress,
		SignatureType: option.SignatureType,
		Timestamp:     fmt.Sprintf("%d", time.Now().UnixMilli()),
		Metadata:      metadata,
		Builder:       builder,
		Expiration:    expiration,
	}
}

func getOrderRawAmounts(side types.Side, size float64, price float64, roundConfig RoundConfig) (model.Side, float64, float64) {
	rawPrice := utils.RoundNormal(price, roundConfig.Price)
	if side == types.BUY {
		rawTakerAmt := utils.RoundDown(size, roundConfig.Size)
		rawMakerAmt := rawTakerAmt * rawPrice
		if utils.DecimalPlaces(rawMakerAmt) > roundConfig.Amount {
			rawMakerAmt = utils.RoundUp(rawMakerAmt, roundConfig.Amount+4)
			if utils.DecimalPlaces(rawMakerAmt) > roundConfig.Amount {
				rawMakerAmt = utils.RoundDown(rawMakerAmt, roundConfig.Amount)
			}
		}
		return model.BUY, rawMakerAmt, rawTakerAmt
	}
	rawMakerAmt := utils.RoundDown(size, roundConfig.Size)
	rawTakerAmt := rawMakerAmt * rawPrice
	if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
		rawTakerAmt = utils.RoundUp(rawTakerAmt, roundConfig.Amount+4)
		if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = utils.RoundDown(rawTakerAmt, roundConfig.Amount)
		}
	}
	return model.SELL, rawMakerAmt, rawTakerAmt
}

func getMarketOrderRawAmounts(side types.Side, size float64, price float64, roundConfig RoundConfig) (model.Side, float64, float64) {
	rawPrice := utils.RoundDown(price, roundConfig.Price)
	if side == types.BUY {
		rawMakerAmt := utils.RoundDown(size, roundConfig.Size)
		rawTakerAmt := rawMakerAmt / rawPrice
		if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = utils.RoundUp(rawTakerAmt, roundConfig.Amount+4)
			if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
				rawTakerAmt = utils.RoundDown(rawTakerAmt, roundConfig.Amount)
			}
		}
		return model.BUY, rawMakerAmt, rawTakerAmt
	}
	rawMakerAmt := utils.RoundDown(size, roundConfig.Size)
	rawTakerAmt := rawMakerAmt * rawPrice
	if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
		rawTakerAmt = utils.RoundUp(rawTakerAmt, roundConfig.Amount+4)
		if utils.DecimalPlaces(rawTakerAmt) > roundConfig.Amount {
			rawTakerAmt = utils.RoundDown(rawTakerAmt, roundConfig.Amount)
		}
	}
	return model.SELL, rawMakerAmt, rawTakerAmt
}
