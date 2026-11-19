# 未成交订单分页

GetOrders 保持单页返回，新增 GetActiveOrdersRequest.NextCursor（*string）。
首屏省略该字段；后续请求保留原来的 ID / Market / AssetID，传入响应 NextCursor。
响应 NextCursor 为 "" 或 "LTE=" 时停止。不要仅因 Data 为空就停止迭代。

~~~go
request := clobtypes.GetActiveOrdersRequest{Market: conditionID}
for {
    page, err := client.GetOrders(ctx, request, auth)
    if err != nil { return err }
    // 处理 page.Data。
    if page.NextCursor == "" || page.NextCursor == "LTE=" { break }
    request.NextCursor = &page.NextCursor
}
~~~

与 ts-sdk listOpenOrders 一样传递 next_cursor；保留 AssetID 必须同时指定 Market 的约束。
新增字段不改变原先不传游标的调用。调用方负责分页、去重和业务对账，不自动批量获取所有页。
