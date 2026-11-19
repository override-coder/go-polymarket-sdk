# Polygon 主网交易验收

每次选择一个钱包、一个协议、一个具体市场和一个 outcome。验收矩阵有四组：Safe × v1、Safe × v2、Deposit Wallet × v1、Deposit Wallet × v2。两种市场都经过 Go SDK 的 `CreateOrderV2`。v1 使用 CTF token ID，v2 使用 position ID。Gamma 仅在测试准备阶段使用。

## 查找市场

```sh
curl --compressed -fsS 'https://gamma-api.polymarket.com/events?version=v2&closed=false&limit=100&offset=0'
POLY_LIVE_DISCOVER=1 go test -tags polymarket_live -count=1 -v -run '^TestDiscoverMarkets$' ./integration/live
```

第二条命令同时查询 v1 和 v2，各显示第一页最多五个正在接单的二元市场及 outcome 顺序。Gamma 对 `limit=100` 要求 gzip。实际测试使用显式选定的市场 ID 再读取市场详情。

## 本地环境

通过本地密钥管理器或交互式隐藏输入设置 `POLY_LIVE_PRIVATE_KEY`，填写所选钱包的 owner 私钥。每次切换账号都要更换对应的变量和 CLOB 凭据。以下变量使用 `POLY_LIVE_` 前缀。

| 变量后缀 | 使用方式 |
| --- | --- |
| `WALLET_TYPE`、`WALLET` | `safe` 或 `deposit`，以及实际持仓钱包的公开地址。代码核对 owner 派生地址，Deposit 支持 legacy 和 beacon。 |
| `PROTOCOL`、`MARKET_ID`、`OUTCOME_INDEX` | `v1` 或 `v2`，Gamma 市场 ID（不是 event ID），以及结果下标 `0` 或 `1`。 |
| `API_KEY`、`API_SECRET`、`API_PASSPHRASE` | 可选的现有 CLOB L2 凭据。未设置 API_KEY 时，用 owner 签名派生 nonce 0 已有凭据。 |
| `AMOUNT`、`MAX_PRICE` | 买单本金（pUSD，最多两位小数）及每份最高价格。按提交报文的整数金额检查。手续费另计。 |
| `MAX_SPEND` | 可选的总预算 pUSD，传入 SDK 的 UserUsdcBalance。启用时 AMOUNT 必须大于该值，且需配置 POLYGON_RPC。 |
| `ORDER_ID` | 只读复查本测试的 FOK 买单时填写。 |
| `POLYGON_RPC` | 赎回使用支持 eth_call、回执和历史区块余额查询的 Polygon 主网 RPC。 |
| `BUILDER_KEY`、`BUILDER_SECRET`、`BUILDER_PASSPHRASE` | Safe 和 Deposit 赎回所需的 Relayer Builder 凭据。 |

## 下单与复查

钱包需已部署、持有 pUSD，并完成对应协议的授权。以下命令在配置的市场发送一笔 FOK 买单，订单无法立即全额成交时可能被拒。成交后测试通过账户接口按订单 ID 查订单，再查与该 taker 订单 ID 关联的成交，核对资金钱包、BUY/FOK 和每笔成交哈希，并要求已确认成交份额之和等于订单的已成交份额。每次调用买单测试都可能产生一笔新的交易。

```sh
POLY_LIVE_ACTION=buy POLY_LIVE_ACK=I_UNDERSTAND_THIS_USES_REAL_FUNDS \
  go test -tags polymarket_live -count=1 -v -run '^TestBuyAndVerifyFill$' ./integration/live
```

保存输出的市场、资产及订单 ID。若提交返回后等待超时，设置 `POLY_LIVE_ORDER_ID`，用以下命令继续查询；该命令不下单。若网络错误导致提交结果未知，应先在账户订单及成交记录中核对。

```sh
POLY_LIVE_ACTION=query go test -tags polymarket_live -count=1 -v -run '^TestQueryExistingOrderAndFill$' ./integration/live
```

## 赎回

待同一市场结算后，选择实际有余额且 payout 为 1 的 outcome，再运行此项。如使用另一笔已有的获胜仓位，可单独验收赎回，并记录对应市场。新买入的仓位可能输掉，活跃市场不能立即提供获胜仓位用于验收。

**赎回范围**：V2 Router 赎回所选 outcome 的全部持仓。CTF collateral adapter 对所选条件的两个 index set 同时赎回，会涉及钱包在该市场已有的仓位。测试先核对持仓并从钱包地址模拟调用，再通过所选钱包的 Relayer 提交。成功需要链上回执成功、获胜仓位余额归零，并且同一笔回执包含向持仓钱包付款的正数 pUSD Transfer。保存 relayer ID 和交易哈希。测试不会自动创建钱包或发送授权。

```sh
POLY_LIVE_ACTION=redeem POLY_LIVE_ACK=I_UNDERSTAND_THIS_USES_REAL_FUNDS \
  go test -tags polymarket_live -count=1 -v -run '^TestRedeemResolvedPosition$' ./integration/live
```


## 预算保护验收

配置 MAX_SPEND 后，现有买单用例会要求 AMOUNT 大于 MAX_SPEND，以验证 SDK 缩量。
例如 AMOUNT=1.50、MAX_SPEND=1.20。测试在提交前校验签单本金不超过预算，
在成交确认后逐笔读取 Polygon 回执，按交易哈希去重，汇总资金钱包转出的 pUSD，
并要求该总额不超过 MAX_SPEND。汇总包含本金和费用，自转账忽略，退款不抵扣，
因此是净支出的保守上界。测试输出每笔回执和总转出金额，不单独推算手续费。
测试期间保持该钱包无其他交易，避免同一结算交易包含额外的资金活动。

先完成本页的密钥、RPC、市场及最高价配置。以下两个命令每次各提交一笔买单。
运行 Safe 前，切换到对应 owner 的 PRIVATE_KEY，并清空或更换 API_KEY、API_SECRET、API_PASSPHRASE。
本轮两个钱包的地址如下。

```sh
POLY_LIVE_WALLET_TYPE=deposit \
POLY_LIVE_WALLET=0x849eF86A5F186F0AdbC24C323340BAF67848C294 \
POLY_LIVE_AMOUNT=1.50 POLY_LIVE_MAX_SPEND=1.20 \
POLY_LIVE_ACTION=buy POLY_LIVE_ACK=I_UNDERSTAND_THIS_USES_REAL_FUNDS \
go test -tags polymarket_live -count=1 -v -run '^TestBuyAndVerifyFill$' ./integration/live

POLY_LIVE_WALLET_TYPE=safe \
POLY_LIVE_WALLET=0x963F832996F2Ab1352E846bA5209A6f8a91dDdD1 \
POLY_LIVE_AMOUNT=1.50 POLY_LIVE_MAX_SPEND=1.20 \
POLY_LIVE_ACTION=buy POLY_LIVE_ACK=I_UNDERSTAND_THIS_USES_REAL_FUNDS \
go test -tags polymarket_live -count=1 -v -run '^TestBuyAndVerifyFill$' ./integration/live
```

提交结果未知时，先核对已有订单，避免重复提交。已获得订单 ID 后，可以保留 MAX_SPEND、
配置 ORDER_ID，并使用上面的 query 命令恢复成交和预算核验。回执暂不可读也采用此路径。
FOK 买单不套用挂单的最小份额断言，最低成交金额由 CLOB 校验。
现有签单价格守卫使用 makerAmount / takerAmount 检查实际签名上限。
市价买单的份额截断可能使该比值略高于 MAX_PRICE，例如本金 1.20、价格 0.63。
此时测试会在发送前拒绝该报文，应检查最终金额和价格配置后再运行。
不足 1 pUSD 的缩量分支通过离线签单回归验证，真实市场可能拒绝此类小额订单。

离线验收命令为

```sh
go test -tags polymarket_live -race -count=1 -run '^(TestSignedBuyLimits|TestBudgetCollateralOutflow|TestBudgetReceiptsDeduplicateTransactions)$' ./integration/live
```


预算回执复查时显式传入同一预算和订单 ID

```sh
POLY_LIVE_ACTION=query POLY_LIVE_MAX_SPEND=1.20 POLY_LIVE_ORDER_ID='填入订单ID' go test -tags polymarket_live -count=1 -v -run '^TestQueryExistingOrderAndFill$' ./integration/live
```

## 统一入口与结果记录

从干净 checkout 执行，每次只执行一个 action，不自动遍历钱包或协议，也不自动重试提交。
脚本不设置资金操作 ACK，仍须由运行者显式设置。每次必须使用新的 JSONL 文件路径，
父目录须已存在。文件以 0600 权限独占创建，已有文件会在市场查询或交易前拒绝覆盖。

~~~sh
bash scripts/live-acceptance.sh buy /绝对路径/safe-v1-buy-001.jsonl
bash scripts/live-acceptance.sh position /绝对路径/safe-v1-position-001.jsonl
POLY_LIVE_ORDER_ID='买单ID' POLY_LIVE_SIDE=BUY   bash scripts/live-acceptance.sh query /绝对路径/safe-v1-query-buy-001.jsonl

# SHARES 为明确卖出的份数，最多两位小数；MIN_PRICE 为最低成交价。
POLY_LIVE_MAX_SPEND= POLY_LIVE_SHARES='2.00' POLY_LIVE_MIN_PRICE='0.45'   bash scripts/live-acceptance.sh sell /绝对路径/safe-v1-sell-001.jsonl
POLY_LIVE_ORDER_ID='卖单ID' POLY_LIVE_SIDE=SELL POLY_LIVE_MAX_SPEND=   bash scripts/live-acceptance.sh query /绝对路径/safe-v1-query-sell-001.jsonl

# 另选已结算且实际持有获胜仓位的市场。
bash scripts/live-acceptance.sh redeem /绝对路径/safe-v1-redeem-001.jsonl
~~~

金额仅为配置示例，运行前按实际确认的份额填写 SHARES。卖单为单笔 FOK，检查买盘深度、
链上可用份额、实际签单份额和最低价格，核对成交份额及链上持仓减少量，不默认卖出全部历史仓位。
position 和 sell 需要 POLYGON_RPC。position 同时记录链上及 CLOB 余额，两者可能存在缓存时差，
不强制即时相等。持仓读取对齐 ts-sdk 合约配置：V1 普通 CTF、V1 neg-risk adapter、V2 PositionManager；
这里只是 SDK 验收，不修改 watcher。

| 钱包 | 协议 | 验收步骤 |
| --- | --- | --- |
| safe | v1 | buy → position → query BUY → sell → query SELL；另行 redeem |
| safe | v2 | 同上 |
| deposit | v1 | 同上 |
| deposit | v2 | 同上 |

每步骤保存独立结果文件。切换矩阵行时更换钱包类型、地址、私钥及 API 凭据、协议、市场、outcome。
测试期间不要让其他程序操作该钱包。卖出后不再持有获胜仓位，就不能在同一份仓位上验收赎回。

JSONL 记录 started/submitting/submitted/confirmed/position 以及最终 passed/failed/skipped，
包含 SDK Git revision（有工作区修改时标记 dirty）、钱包、协议、市场、金额/价格、订单 ID、
Relayer ID、交易哈希、基础单位余额。不记录密钥、认证头或签名。
confirmed 是检查点，必须结合最后的 passed 判断整个 action 是否通过。程序被强制终止可能无最终行，
不表示没有提交。submitting 后出现错误或崩溃时先核对已有订单/Relayer 记录，不能直接重跑交易。

有订单 ID 可用 query 恢复成交核验，卖单必须设置 SIDE=SELL；query 不重跑卖出前后的持仓差值检查。
赎回尚无独立自动恢复用例，可按记录的 Relayer ID 使用 GetTransactionContext 和链上回执复查，
不能盲目再次 redeem。新结果文件只记录新尝试，不提供链上幂等性。
直接运行 go test 时 RESULT_FILE 可选，未设置则输出 LIVE_RESULT JSON 行。
离线测试通过不代表四组主网矩阵已经验收通过。

~~~sh
go test -tags polymarket_live -race -count=1   -run '^(TestSignedSellLimits|TestPositionContractRouting|TestLiveResultRefusesOverwrite)$' ./integration/live
bash -n scripts/live-acceptance.sh
~~~
