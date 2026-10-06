// Package dataapi implements the migrated Polymarket Data API V2 reads.
//
// Contract: https://data-api.polymarket.com/v2/openapi.json (2026-10-06).
// Migration: https://docs.polymarket.com/migrate/data-api-v1-to-v2.
//
// V2 is a breaking change: paginated methods return types.Page[T], with rows in
// Data and the next opaque cursor in Pagination.NextCursor. Offset is response
// display metadata only. Query Condition replaces Market; unsupported V1 filters
// and sorts are no longer exposed. Nil query fields leave defaults to the server.
//
// Re-send positions anchors and filters, activity filters, and holders conditions
// and filters on every page. Stop when NextCursor is nil, not when a page is short.
// A cursor owns its page size, so Limit is ignored on subsequent pages.
//
// GetClosedPositions uses the complete Position shape with status=CLOSED.
// GetPositionValue returns Envelope[PositionValue], a single object rather than
// an array. Leaderboard board pages use GetTraderLeaderboardRankings; the user
// branch uses GetTraderLeaderboardUser and returns nullable Data and ranks.
// Holder economics are optional pointers, because absent/null does not mean zero.
//
// GetMarketByToken remains a CLOB request through the separate CLOB client.
package dataapi
