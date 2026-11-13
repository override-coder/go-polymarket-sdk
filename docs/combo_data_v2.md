# Combo Data API V2 migration

Both `combos.Client.GetComboPositions` and `combos.Client.GetComboActivity` now
use `dataapi.Client` on the configured Data API host. RFQ, maker, Builder Gateway
and WebSocket routes are unchanged. No authentication is added to these reads.

This is a source-breaking migration of the retired V1 read contract, consistent
with the existing `dataapi` V2 migration. Method and constructor names remain;
`combos/types` query, row and page types alias the canonical `dataapi/types` types.
Existing callers of `dataapi.GetComboPositions` need no changes.

| V1 caller | V2 caller |
| --- | --- |
| `MarketID` | `Condition` containing Combo condition IDs, not Gamma market IDs |
| `Offset` | `Cursor` from `Pagination.NextCursor` |
| Position `Sort` | `SortBy` and `SortDirection` |
| Position `Combos` / activity `Activity` | `Data` |
| `Pagination.Total` | Not available; follow `NextCursor` until nil |
| Position `UserAddress` / `SharesBalance` | `ProxyWallet` / `CurrentSize` |
| Position `TotalCostUSDC` | `GrossEntryCostUSDC` (fee-inclusive entry basis) |
| Activity `UserAddress` / `EventKind` / `TxHash` | `ProxyWallet` / `Type` / `TransactionHash` |
| Legacy leg fields | Enriched `ComboLeg` fields and nested `Market` / `Event` |

The removed `ModuleID` and `ModuleKind` have no V2 field; do not synthesize them.
Position amounts are `json.Number`. Activity `AmountUSDC` and `PayoutUSDC` are
`*json.Number`: missing/null remains nil, while a confirmed zero remains a number.
Parse numbers using a decimal type, never through float64 for accounting.

## Requests and pagination

- Positions: `GET /v2/positions/combos`. Reuse `dataapi.GetComboPositions` validation,
  filters, lifecycle statuses and server-side sorting. See [positions](combo_positions.md).
- Activity: `GET /v2/activity/combos`. Set `User`; optionally set `Condition`, `Limit`
  and `Cursor`. Keep the same user and conditions on every continuation request.
- Combo condition IDs are bytes31 (`0x` plus 62 hexadecimal characters), with at
  most 20 distinct IDs. They are not ordinary bytes32 CTF condition IDs.
- First-page limit is 0..1000. With a cursor, its page size wins and `Limit` is not sent.
- Short or empty pages can still carry a next cursor. Neither row count nor the
  response-only `Offset` is a pagination control. V1 query parameters are not sent.

```go
q := types.ComboActivityQuery{User: wallet, Condition: conditions}
for {
    page, err := client.GetComboActivity(ctx, q)
    if err != nil { return err }
    // Consume page.Data, preserving nullable payout semantics.
    if page.Pagination.NextCursor == nil { break }
    q.Cursor = page.Pagination.NextCursor
}
```

Known consumers: poly-trader's Combo holdings, Cash Out and auto-claim already
use the canonical V2 positions client. Its redemption PnL reader already uses
`/v2/activity/combos` with nullable decimal payouts; it does not use the old
`combos` read models. poly-openapi's public routes are application contracts and
do not need `/v2` prefixes. Upgrade the SDK dependency and regress these consumers.

Sources checked 2026-11-13:
- https://docs.polymarket.com/migrate/data-api-v1-to-v2
- https://docs.polymarket.com/api-reference/wallet/list-combo-positions
- https://docs.polymarket.com/api-reference/feeds/list-combo-activity
