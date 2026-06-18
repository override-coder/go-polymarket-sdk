package combos

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	"github.com/override-coder/go-polymarket-sdk/signing"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/builder"
	"github.com/polymarket/go-order-utils/pkg/eip712"
	"github.com/polymarket/go-order-utils/pkg/model"
)

const e6 = int64(1_000_000)

var (
	comboOrderStructureV2 = []abi.Type{
		eip712.Bytes32,
		eip712.Uint256,
		eip712.Address,
		eip712.Address,
		eip712.Uint256,
		eip712.Uint256,
		eip712.Uint256,
		eip712.Uint8,
		eip712.Uint8,
		eip712.Uint256,
		eip712.Bytes32,
		eip712.Bytes32,
	}
	comboTypedDataSignStructureV2 = []abi.Type{
		eip712.Bytes32,
		eip712.Bytes32,
		eip712.Bytes32,
		eip712.Bytes32,
		eip712.Uint256,
		eip712.Address,
		eip712.Bytes32,
	}
	comboOrderStructureHashV2 = crypto.Keccak256Hash(
		[]byte("Order(uint256 salt,address maker,address signer,uint256 tokenId,uint256 makerAmount,uint256 takerAmount,uint8 side,uint8 signatureType,uint256 timestamp,bytes32 metadata,bytes32 builder)"),
	)
	comboTypedDataSignStructureHashV2 = crypto.Keccak256Hash(
		[]byte("TypedDataSign(Order contents,string name,string version,uint256 chainId,address verifyingContract,bytes32 salt)Order(uint256 salt,address maker,address signer,uint256 tokenId,uint256 makerAmount,uint256 takerAmount,uint8 side,uint8 signatureType,uint256 timestamp,bytes32 metadata,bytes32 builder)"),
	)
	comboDepositWalletNameHashV2    = crypto.Keccak256Hash([]byte("DepositWallet"))
	comboDepositWalletVersionHashV2 = crypto.Keccak256Hash([]byte("1"))

	// Exchange.domainSeparator() on the combo matcher 0xe333...c00Aa.
	comboMatcherDomainSeparatorV2 = common.HexToHash("0x466c63910185bbd55e8679264200c4e0abdcbb0c6264eb3d41d13326022e095b")
	comboOrderTypeStringV2        = []byte("Order(uint256 salt,address maker,address signer,uint256 tokenId,uint256 makerAmount,uint256 takerAmount,uint8 side,uint8 signatureType,uint256 timestamp,bytes32 metadata,bytes32 builder)")
)

type OrderBuilder struct {
	chainID *big.Int
	signFn  signing.SignatureFunc
}

func NewOrderBuilder(chainID *big.Int, signFn signing.SignatureFunc) *OrderBuilder {
	return &OrderBuilder{
		chainID: chainID,
		signFn:  signFn,
	}
}

type BuildQuoteRequestInput struct {
	RFQID         string
	QuoteID       string
	Direction     combostypes.Direction
	YesPositionID string
	PriceE6       string
	SizeE6        string
	Source        combostypes.QuoteSource
	ValidUntil    int64
	Expiration    int64
	Metadata      string
	Builder       string
}

func (b *OrderBuilder) BuildQuoteRequest(input BuildQuoteRequestInput, option *sdktypes.AuthOption) (*combostypes.QuoteRequest, error) {
	if option == nil {
		return nil, fmt.Errorf("auth option is required")
	}
	if b.signFn == nil {
		return nil, fmt.Errorf("signature function is required")
	}
	if strings.TrimSpace(input.RFQID) == "" {
		return nil, fmt.Errorf("rfq id is required")
	}
	if strings.TrimSpace(input.QuoteID) == "" {
		return nil, fmt.Errorf("quote id is required")
	}
	if strings.TrimSpace(input.YesPositionID) == "" {
		return nil, fmt.Errorf("yes position id is required")
	}

	priceE6, err := parsePositiveE6(input.PriceE6, "price_e6")
	if err != nil {
		return nil, err
	}
	sizeE6, err := parsePositiveE6(input.SizeE6, "size_e6")
	if err != nil {
		return nil, err
	}

	maker := option.SingerAddress
	if option.FunderAddress != "" {
		maker = option.FunderAddress
	}
	if maker == "" {
		return nil, fmt.Errorf("maker address is required")
	}
	if option.SingerAddress == "" {
		return nil, fmt.Errorf("signer address is required")
	}
	rfqSigner := comboRFQSigner(option.SignatureType, option.SingerAddress, maker)

	side, makerAmount, takerAmount, err := quoteAmounts(input.Direction, priceE6, sizeE6)
	if err != nil {
		return nil, err
	}

	metadata := input.Metadata
	if metadata == "" {
		metadata = combostypes.Bytes32Zero
	}
	builderCode := input.Builder
	if builderCode == "" {
		builderCode = combostypes.Bytes32Zero
	}

	expiration := "0"
	if input.Expiration > 0 {
		expiration = fmt.Sprintf("%d", input.Expiration)
	}

	orderData := &model.OrderDataV2{
		Maker:         maker,
		TokenID:       input.YesPositionID,
		MakerAmount:   makerAmount.String(),
		TakerAmount:   takerAmount.String(),
		Side:          side,
		Signer:        rfqSigner,
		SignatureType: option.SignatureType,
		Timestamp:     fmt.Sprintf("%d", time.Now().UnixMilli()),
		Metadata:      metadata,
		Builder:       builderCode,
		Expiration:    expiration,
	}

	signedOrder, err := b.buildComboSignedOrder(option.SingerAddress, orderData)
	if err != nil {
		return nil, err
	}

	source := input.Source
	if source == "" {
		source = combostypes.QuoteSourceCollateral
	}

	return &combostypes.QuoteRequest{
		QuoteID:       input.QuoteID,
		RFQID:         input.RFQID,
		SignerAddress: rfqSigner,
		MakerAddress:  maker,
		SignatureType: option.SignatureType,
		PriceE6:       input.PriceE6,
		SizeE6:        input.SizeE6,
		Source:        source,
		SignedOrder:   orderToJSON(signedOrder),
		ValidUntil:    input.ValidUntil,
	}, nil
}

func comboRFQSigner(signatureType model.SignatureType, signer, maker string) string {
	if signatureType == model.POLY_1271 {
		return maker
	}
	return signer
}

func (b *OrderBuilder) buildComboSignedOrder(signer string, orderData *model.OrderDataV2) (*model.SignedOrderV2, error) {
	order, err := buildComboOrderV2(orderData)
	if err != nil {
		return nil, err
	}

	contentsHash, err := comboOrderContentsHash(order)
	if err != nil {
		return nil, err
	}

	var signature model.OrderSignatureV2
	if orderData.SignatureType != model.POLY_1271 {
		orderHash, err := eip712.HashTypedDataV4(comboMatcherDomainSeparatorV2, comboOrderStructureV2, comboOrderValues(order))
		if err != nil {
			return nil, err
		}
		signature, err = b.signFn(signer, orderHash.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		wrappedHash, err := comboPoly1271WrappedHash(order, contentsHash, b.chainID)
		if err != nil {
			return nil, err
		}
		innerSig, err := b.signFn(signer, wrappedHash.Bytes())
		if err != nil {
			return nil, err
		}
		signature = comboPoly1271FinalSignature(innerSig, contentsHash)
	}

	return &model.SignedOrderV2{
		OrderV2:   *order,
		Signature: signature,
	}, nil
}

func (b *OrderBuilder) buildSignedOrder(signer string, contract model.VerifyingContract, orderData *model.OrderDataV2) (*model.SignedOrderV2, error) {
	exchangeOrderBuilder := builder.NewExchangeOrderBuilderImplV2(b.chainID, nil)
	order, err := buildComboOrderV2(orderData)
	if err != nil {
		return nil, err
	}

	var signature model.OrderSignatureV2
	if orderData.SignatureType != model.POLY_1271 {
		orderHash, err := exchangeOrderBuilder.BuildOrderHash(order, contract)
		if err != nil {
			return nil, err
		}
		signature, err = b.signFn(signer, orderHash.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		wrappedHash, err := exchangeOrderBuilder.BuildPoly1271WrappedHash(order, contract)
		if err != nil {
			return nil, err
		}
		innerSig, err := b.signFn(signer, wrappedHash.Bytes())
		if err != nil {
			return nil, err
		}
		signature, err = exchangeOrderBuilder.BuildPoly1271FinalSignature(order, contract, innerSig)
		if err != nil {
			return nil, err
		}
	}

	return &model.SignedOrderV2{
		OrderV2:   *order,
		Signature: signature,
	}, nil
}

func buildComboOrderV2(orderData *model.OrderDataV2) (*model.OrderV2, error) {
	signer := orderData.Signer
	if strings.TrimSpace(signer) == "" {
		signer = orderData.Maker
	}

	tokenID, err := parseBigInt(orderData.TokenID, "TokenId")
	if err != nil {
		return nil, err
	}
	makerAmount, err := parseBigInt(orderData.MakerAmount, "MakerAmount")
	if err != nil {
		return nil, err
	}
	takerAmount, err := parseBigInt(orderData.TakerAmount, "TakerAmount")
	if err != nil {
		return nil, err
	}
	timestamp, err := parseBigInt(orderData.Timestamp, "Timestamp")
	if err != nil {
		return nil, err
	}

	expirationValue := orderData.Expiration
	if expirationValue == "" {
		expirationValue = "0"
	}
	expiration, err := parseBigInt(expirationValue, "Expiration")
	if err != nil {
		return nil, err
	}

	salt, err := randomPositiveInt63()
	if err != nil {
		return nil, err
	}

	return &model.OrderV2{
		Salt:          salt,
		Maker:         common.HexToAddress(orderData.Maker),
		Signer:        common.HexToAddress(signer),
		TokenID:       tokenID,
		MakerAmount:   makerAmount,
		TakerAmount:   takerAmount,
		Side:          new(big.Int).SetInt64(int64(orderData.Side)),
		SignatureType: new(big.Int).SetInt64(int64(orderData.SignatureType)),
		Timestamp:     timestamp,
		Metadata:      common.HexToHash(orderData.Metadata),
		Builder:       common.HexToHash(orderData.Builder),
		Expiration:    expiration,
	}, nil
}

func comboOrderContentsHash(order *model.OrderV2) (common.Hash, error) {
	encoded, err := eip712.Encode(comboOrderStructureV2, comboOrderValues(order))
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(encoded), nil
}

func comboOrderValues(order *model.OrderV2) []interface{} {
	return []interface{}{
		comboOrderStructureHashV2,
		order.Salt,
		order.Maker,
		order.Signer,
		order.TokenID,
		order.MakerAmount,
		order.TakerAmount,
		uint8(order.Side.Uint64()),
		uint8(order.SignatureType.Uint64()),
		order.Timestamp,
		order.Metadata,
		order.Builder,
	}
}

func comboPoly1271WrappedHash(order *model.OrderV2, contentsHash common.Hash, chainID *big.Int) (common.Hash, error) {
	values := []interface{}{
		comboTypedDataSignStructureHashV2,
		contentsHash,
		comboDepositWalletNameHashV2,
		comboDepositWalletVersionHashV2,
		chainID,
		order.Signer,
		common.Hash{},
	}
	return eip712.HashTypedDataV4(comboMatcherDomainSeparatorV2, comboTypedDataSignStructureV2, values)
}

func comboPoly1271FinalSignature(innerSig []byte, contentsHash common.Hash) model.OrderSignatureV2 {
	orderTypeLen := len(comboOrderTypeStringV2)
	finalSig := make([]byte, 0, len(innerSig)+32+32+len(comboOrderTypeStringV2)+2)
	finalSig = append(finalSig, innerSig...)
	finalSig = append(finalSig, comboMatcherDomainSeparatorV2.Bytes()...)
	finalSig = append(finalSig, contentsHash.Bytes()...)
	finalSig = append(finalSig, comboOrderTypeStringV2...)
	finalSig = append(finalSig, byte(orderTypeLen>>8), byte(orderTypeLen))
	return finalSig
}

func parseBigInt(value, field string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse %s: %s as valid *big.Int", field, value)
	}
	return n, nil
}

func randomPositiveInt63() (*big.Int, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 63)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	if n.Sign() == 0 {
		n.SetInt64(1)
	}
	return n, nil
}

func quoteAmounts(direction combostypes.Direction, priceE6, sizeE6 *big.Int) (model.Side, *big.Int, *big.Int, error) {
	collateral := mulDiv(priceE6, sizeE6, big.NewInt(e6))
	switch direction {
	case combostypes.DirectionBuy:
		return model.SELL, new(big.Int).Set(sizeE6), collateral, nil
	case combostypes.DirectionSell:
		collateral = mulDivCeil(priceE6, sizeE6, big.NewInt(e6))
		return model.BUY, collateral, new(big.Int).Set(sizeE6), nil
	default:
		return 0, nil, nil, fmt.Errorf("unsupported direction: %s", direction)
	}
}

func parsePositiveE6(value, field string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("%s must be an integer string", field)
	}
	if n.Sign() <= 0 {
		return nil, fmt.Errorf("%s must be positive", field)
	}
	return n, nil
}

func mulDiv(a, b, d *big.Int) *big.Int {
	out := new(big.Int).Mul(a, b)
	return out.Div(out, d)
}

func mulDivCeil(a, b, d *big.Int) *big.Int {
	product := new(big.Int).Mul(a, b)
	q, r := new(big.Int).QuoRem(product, d, new(big.Int))
	if r.Sign() > 0 {
		q.Add(q, big.NewInt(1))
	}
	return q
}

func orderToJSON(order *model.SignedOrderV2) combostypes.SignedOrderV2 {
	expiration := order.Expiration.String()
	if expiration == "0" {
		expiration = ""
	}
	return combostypes.SignedOrderV2{
		Salt:          order.Salt.String(),
		Maker:         order.Maker.String(),
		Signer:        order.Signer.String(),
		TokenID:       order.TokenID.String(),
		MakerAmount:   order.MakerAmount.String(),
		TakerAmount:   order.TakerAmount.String(),
		Side:          int(order.Side.Int64()),
		SignatureType: model.SignatureType(order.SignatureType.Uint64()),
		Timestamp:     order.Timestamp.String(),
		Expiration:    expiration,
		Metadata:      order.Metadata.String(),
		Builder:       order.Builder.String(),
		Signature:     "0x" + common.Bytes2Hex(order.Signature),
	}
}
