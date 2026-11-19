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

## Relayer 查询与错误处理

GetTransaction(id) 保持原接口；需要取消或超时控制时使用
GetTransactionContext(ctx, id)。PollUntilState 的 ctx 同时控制正在执行的 HTTP 查询和轮询等待。
取消轮询不会取消已提交的链上交易；保留 transaction ID，之后继续查询结果。

CreateL2Headers 和 CreateL2BuilderHeaders 在凭据指针为 nil 时返回 error。
这不是凭据有效性验证，过期、撤销或不完整凭据仍应按接口错误处理。
Safe/Deposit 提交入口仍须传入有效 AuthOption、签名函数及所需凭据。

## SDK 交付验证边界

生产构建及测试包编译可在干净 checkout 执行：

```sh
go build -mod=mod ./...
go test -mod=mod ./... -run '^$'
go vet -mod=mod ./...
go test -mod=mod -race ./headers ./relayer -run '^(TestMissingCredentials|TestCredentialHeaders|TestPollUntilStateCancelsInFlightRequest|TestGetTransactionCompatibility)$'
go test -mod=mod -race ./clob -run '^(TestFetchBalanceAllowance.*|TestMarketBuyBudgetBelowOne|TestCreateOrderV2HonorsSmallMarketBuyBudget)$'
```

上述命令不发起真实交易。不要将全仓无筛选的 go test 作为离线验收：其他包仍有历史联网测试。
旧 clob/auth_test.go 需 polymarket_legacy_live 标签才参与编译；无签名密钥的旧认证/下单用例跳过，
有效的显式实盘入口见 [integration/live](integration/live/README.md)，凭据仅通过环境变量提供。

本次交付限 SDK：订单路由维持从资产 ID 识别交易所，不新增 version/resolutionStatus 入参或交易期间的 Gamma 查询。
价格取整维持现有兼容行为；调用方应按市场步长传价格，并核对最终签单/成交金额。
watcher 的链上持仓合约选择不属于本次 SDK 修改。仅离线构建和回归不能证明真实赎回已成功，
实盘验收结果应另行记录订单 ID、交易哈希与回执。

## 公开交易入口的参数边界

CLOB 认证、查询和撤单入口拒绝 nil AuthOption；L2 操作拒绝缺失 API 凭据。
已有 API 凭据的查询/撤单不要求签名函数。CreateOrder/CreateOrderV2 在访问市场前检查签名函数、
正 chain ID、uint256 资产 ID、BUY/SELL、订单类型、有限正数量及 0 < price < 1。
V2 额外拒绝 nil 订单、非有限或负余额/Builder 费率，以及 FOK/FAK + postOnly。

tick 字符串支持 0.1、0.01、0.005、0.0025、0.001、0.0001；其他值报错，不再套用默认精度。
这是非法输入行为的收紧。有效输入保留现有价格归一化和取整，未整体改为 ts-sdk 严格步长拒绝策略。

Safe 构建/签名入口拒绝 nil AuthOption、缺失签名函数、空交易列表及 nil/负 nonce；
BuildTx 仍用 nonce=0 查询当前 nonce。Deposit batch 拒绝非法 uint256 nonce/deadline；
提交入口提前检查 Builder 凭据是否存在。已签名的 ByTx 提交不要求签名函数。
签名工具函数也返回缺失签名函数或 chain ID 的错误。

校验不验证凭据是否撤销、订单是否可成交或余额是否足够，也不替代服务端校验。
构造函数签名未改变，Relayer 构造仍要求受支持的 Polygon chain ID 137。
统一实盘验收入口及 Safe/Deposit × V1/V2 矩阵见 [验收说明](integration/live/README.md)。
