# 批量报价方向

新增 GetMarketPricesBySide(ctx, requests)，返回 map[string]map[string]string：资产ID → BUY/SELL → 十进制价格。
同资产双方向请求使用该方法，例如 result[assetID]["BUY"] 和 result[assetID]["SELL"]；缺失的资产/方向不补零。

旧 GetMarketPrices 保留 map[string]string 返回类型，严格按请求的 Side 选价。
同资产请求不同方向时在HTTP调用前返回错误，提示使用新方法；不再静默以SELL覆盖BUY。
同资产同方向重复请求仍可用。旧双向调用需迁移，此错误行为收紧不改变方法签名。
HTTP重试策略保持原样。

依据（2026-11-19）：https://docs.polymarket.com/api-spec/clob-openapi.yaml 的 /prices 双层响应，及本地 ts-sdk PricesSchema。
