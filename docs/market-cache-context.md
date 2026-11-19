# 市场缓存与调用取消

同一 token/condition 的冷缓存请求继续合并，避免重复查询 CLOB。
每个调用方独立等待：取消或 deadline 到达即退出，不再等待最先发起的调用方完成。

共享请求使用独立的 15 秒总超时，不继承首个调用方的取消或 deadline；
保留首个调用方的 context values，以便现有注入 Transport 使用路由信息。
即使所有调用方都已取消，该只读请求仍可完成并填充缓存，最长持续 15 秒。
首次调用已经取消时不启动查询。该机制仅作用于 CLOB 市场元数据，不延长订单提交的 context。
不会为 version 或 resolutionStatus 查询 Gamma。

这是共享缓存的生命周期约定；调用方仍应为每次交易设置自己的 deadline。
