# CLOB query deprecations (2026-11-13)

Official notice: https://docs.polymarket.com/changelog/predictions#nov-11-2026

The December 14 query changes deprecate `associate_trades` and require `market`
when filtering `GET /data/orders` by `asset_id`.

- `GetOrders` rejects an asset filter with an empty or whitespace-only market
  before authentication or HTTP I/O. The remaining query behavior is unchanged.
- `OpenOrder.AssociateTrades` remains source/wire compatible, but is deprecated.
  Missing or empty values must not be interpreted as evidence of no fills.
- `GetTrades` retains its existing filters, pagination and response model. Use
  `taker_order_id` and `maker_orders[].order_id` to discover order relationships;
  the caller must validate ownership, quantities and settlement independently.
- The `asset_id`/`market` constraint does not apply to trade reads or cancellation
  based on this notice. Order submission `tradeIDs` is a distinct field.

## Offline validation

`clob/order_query_deprecation_test.go` covers rejected asset-only queries without
network calls, accepted query combinations, single/list orders without the
removed field, two-page trade queries with encoded cursors and maker/taker
relationships, and HTTP/decode errors that must not appear as empty successes.
Existing trade-ID-to-hash tests remain part of the focused verification.

Commands (unset a stale GOROOT when the installed Go executable differs):

```sh
env -u GOROOT go test ./clob -run 'Test(OrderQuer|TradeQuery|QueryErrors|ResolveTransactionsHashes)' -count=1
env -u GOROOT go test -race ./clob -run 'Test(OrderQuer|TradeQuery|QueryErrors|ResolveTransactionsHashes)' -count=1
env -u GOROOT go vet ./clob/...
env -u GOROOT go build ./...
```

No live credentials, submissions or cancellations are used by these tests.
Consumer fill discovery and reconciliation changes are separate downstream work.

## Canonical trade route correction

GetTrades and POST trade-ID hash resolution use the documented `/data/trades`
path for both the HTTP request and HMAC signature. The query pagination test
asserts the literal path and independently computes the HMAC with the standard
library; it does not derive expected behavior from the SDK endpoint constant.
The by-ID resolution regression also asserts the literal path.
