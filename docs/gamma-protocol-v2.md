# Gamma Protocol V2 字段

Market 新增可空 Version、ResolutionStatus；Event 新增可空 Version。单查、列表、嵌套事件和搜索复用相同类型，未知响应值原文保留，缺失/null 不默认成 v1。
根据 Market.Version 选择 outcome 对应的 positionIds（v2）或 JSON 文本 clobTokenIds（v1）；v2结算看 ResolutionStatus，v1保留 UmaResolutionStatus。不要仅根据 positionIds 存在判断协议。

GetMarketsParams、GetMarketsKeysetParams、GetEventsKeysetParams 新增 Versions []string。
不传保留服务端默认；传 []string{"v2"} 选择V2，[]string{"v1","v2"} 使用重复的 version 查询键。
GetMarketsKeysetParams 另支持 ResolutionStatus *string（inactive/active/resolved），编码为 resolution_status。
查询参数拒绝未知版本、超过两个版本和未知状态；响应未知值仍透传，供上游判断。

仅用于显式市场发现和状态读取，不新增下单入参，也不在交易中查询 Gamma。
本次不扩展费用/奖励等其他市场元数据。

依据（2026-11-19）：https://docs.polymarket.com/migrate/polymarket-v2/api-integrations
及本地 ts-sdk 的 markets/events 筛选和 Gamma 类型。
