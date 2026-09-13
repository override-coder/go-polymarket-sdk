# Combo 主动下单与跟单

SDK 同时保留两条 Combo 路径：

- `TakerRunner`：已有的 requester WebSocket 路径。
- `BuilderTakerClient`：官方 Builder Gateway REST 路径，适合上游订单服务和信号跟单。

Builder Gateway 的一次主动下单是：

```text
legs + direction + requested_size
  -> 创建 RFQ
  -> 获得实时 quote
  -> 跟单账户签名并接受 quote
  -> 轮询终态，取得 tx_hash
```

不要复用被跟单地址的订单签名、报价或交易对手。跟单只复制交易意图和 legs；价格、报价、订单签名与链上交易由跟单账户在下单时重新生成。

## 给上游订单服务的调用方式

`PlaceAndWait` 是与普通下单最接近的入口。它在 `context` 的期限内完成 RFQ 创建、接受和状态轮询，终态结果包含 `RFQID`、`TakerOrderHash` 与 `TxHash`。

```go
client := combos.NewBuilderTakerClient(
	"", // 使用官方 Builder Gateway production host
	big.NewInt(137),
	signFn,
)

result, err := client.PlaceAndWait(ctx, combos.TakerRequest{
	LegPositionIDs: []string{
		"first-yes-position-id",
		"second-yes-position-id",
	},
	Direction: combostypes.DirectionBuy,
	Side:      combostypes.ComboSideYes,
	RequestedSize: combostypes.RequestedSize{
		Unit:    "notional", // BUY 使用 pUSD 名义金额
		ValueE6: "1000000", // 1 pUSD
	},
}, authOption, time.Second)
if err != nil {
	return err // 网络、认证、签名或 context 超时
}

switch result.Status {
case "CONFIRMED", "FILLED":
	return result.TxHash
case "FAILED", "EXPIRED", "CANCELED":
	return fmt.Errorf("combo order terminal status: %s", result.Status)
default:
	return fmt.Errorf("unexpected terminal status: %s", result.Status)
}
```

`authOption` 必须同时含有账户 L2 凭证与 Builder 凭证：

```go
authOption := &sdktypes.AuthOption{
	SignatureType:      signatureType,
	SingerAddress:      signerAddress,
	FunderAddress:      makerAddress,
	ApiKeyCreds:        accountL2Credentials,
	BuilderApiKeyCreds: builderCredentials,
}
```

SDK 会在创建和接受 RFQ 时自动注入账户 L2 与 Builder HMAC 标头；查询状态只使用账户 L2 标头。

## 请求约束

- `LegPositionIDs` 必须是 2 至 50 个不重复的 position ID。
- 官方 Builder Gateway 当前只支持 `YES` 侧。
- `BUY` 必须使用 `requested_size.unit = "notional"`，`value_e6` 为 pUSD 的 6 位精度整数。
- `SELL` 必须使用 `requested_size.unit = "shares"`。
- 创建 RFQ 或接受报价成功不等于成交。上游只能把 `CONFIRMED` 或 `FILLED` 当作成功，并以 `TxHash` 作为交易标识。

## 跟单信号处理

对一个目标地址的已确认 Combo 交易：

1. 从链上 calldata 或 Combo Activity 提取 `combinatorialLegs`、方向和侧别。
2. 校验方向、legs 数量、风险上限与去重键。
3. 按跟单配置重新填写 `RequestedSize`，例如固定 `1 pUSD`。
4. 调用 `PlaceAndWait` 并保存 RFQ ID、订单哈希和交易哈希。

行情变化、legs 状态变化、没有可执行报价或余额不足时，订单可能进入 `FAILED` 或 `EXPIRED`，这属于业务结果，不应重试或重放原订单。
