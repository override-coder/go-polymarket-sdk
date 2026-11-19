# go-polymarket-sdk

## 上游依赖配置

SDK 使用维护的 go-order-utils fork 提供 Protocol V2 和 POLY_1271 签名。
Go 依赖模块内的 replace 不会传递给调用方，上游主模块必须配置下面的替换，再引入 SDK。

```sh
go mod edit -replace=github.com/polymarket/go-order-utils=github.com/override-coder/go-order-utils@v1.22.7-0.20261119031642-91a0e8d83f5f
go get github.com/override-coder/go-polymarket-sdk@main_v2
go build ./...
```

go get 会在上游 go.mod 中记录解析到的版本。发布时保留该版本及 go.sum，以便复现构建。
SDK 的公开签名类型仍使用 github.com/polymarket/go-order-utils/pkg/model 导入路径。

## 余额与额度

GetBalanceAllowance(auth) 保留为抵押资产查询，并使用 auth.SignatureType。
FetchBalanceAllowance(ctx, request, auth) 支持资产参数和请求取消。

```go
result, err := client.FetchBalanceAllowance(ctx, clobtypes.BalanceAllowanceRequest{
    AssetType: clobtypes.AssetTypeConditionalV2,
    TokenID: positionID,
}, auth)
if err != nil {
    return err
}
exchangeAllowance := result.Allowances.ExchangeV3
_ = exchangeAllowance
```

示例中的 clobtypes 指 github.com/override-coder/go-polymarket-sdk/clob/types。
抵押资产使用 AssetTypeCollateral 且省略 TokenID，CTF 持仓使用 AssetTypeConditional，
Protocol V2 持仓使用 AssetTypeConditionalV2。余额和额度以十进制基础单位字符串返回。
ExchangeV3 字段对应 0xe3333700cA9d93003F00f0F71f8515005F6c00Aa，现有额度字段继续可用。

## 市价买单金额

FOK/FAK BUY 的 Size 为抵押资产金额。传入 UserUsdcBalance 时，SDK 预留已知手续费，
并在金额不足时缩减 Size，包括缩减后小于 1 的情况，与 ts-sdk 的 maxSpend 处理一致。
调用方应核对平台最低订单要求与实际响应，继续通过订单和成交查询确认执行结果。
