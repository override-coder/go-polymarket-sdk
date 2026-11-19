# 成交诊断字段

Trade 新增 ErrMsg *string（err_msg）及 MatchTimeNano *string（match_time_nano）。
缺失/null保持nil，纳秒时间按原始字符串保留，避免浮点精度损失；原MatchTime保留。
纳秒时间不是唯一订单/成交标识，业务幂等仍使用相应ID。

依据（2026-11-19）：https://docs.polymarket.com/api-spec/clob-openapi.yaml 的 Trade schema。
