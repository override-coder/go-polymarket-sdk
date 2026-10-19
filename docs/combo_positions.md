# Combo positions and manual SELL

`dataapi.Client.GetComboPositions(ctx, types.ComboPositionsQuery)` calls the independent
`GET /v2/positions/combos` endpoint. It does not merge ordinary CLOB positions.
Set `User` to the holding wallet. To sort by entry amount, set `SortBy` to
`ENTRY_COST` and `SortDirection` to `DESC`. Follow `Pagination.NextCursor` with
the same `User`; omitted filters and sort are adopted from the cursor. Never
invent an offset request or locally sort just one page.

The complete enriched position/leg/market/event response is preserved. Decimal
numbers use `json.Number`; parse them as decimals for financial calculations.
`CurrentSize` is the held shares. `GrossEntryCostUSDC` is the exact fee-inclusive
entry basis; it is not the sum of rounded `EntryCostUSDC` and `EntryFeesUSDC`.
Neither the entry price nor the product of leg marks is an executable exit quote.

For manual SELL, use `BuilderTakerClient.CreateRFQ` with `DirectionSell`,
`ComboSideYes`, the original leg position IDs, and `RequestedSize{Unit:"shares",
ValueE6:"10000000"}` for ten shares. Present `Quote.NetReceiveE6` as the fee-net
proceeds, retain `ExpiresAt` and the RFQ, and call `AcceptQuote` only on explicit
confirmation before expiry. Safe and POLY_1271 use their existing signing paths.
Poll the same RFQ after acceptance; an ambiguous response is not permission to
create and accept another RFQ. Builder identity must remain the same for create
and accept. No live trades are part of the default tests.

Sources: https://docs.polymarket.com/api-reference/wallet/list-combo-positions
and https://docs.polymarket.com/trading/combos/builders .
