#!/usr/bin/env bash
# One explicit action per invocation; credentials are read only from the environment.
set -euo pipefail
case "${1:-}" in
  buy) test_name=TestBuyAndVerifyFill ;;
  sell) test_name=TestSellAndVerifyFill ;;
  query) test_name=TestQueryExistingOrderAndFill ;;
  position) test_name=TestQueryPosition ;;
  redeem) test_name=TestRedeemResolvedPosition ;;
  *) printf '%s\n' 'Usage: bash scripts/live-acceptance.sh {buy|sell|query|position|redeem} NEW_RESULT_FILE' >&2; exit 2 ;;
esac
: "${2:?A new result file path is required}"
case "${POLY_LIVE_WALLET_TYPE:-}" in safe|deposit) ;; *) printf '%s\n' 'Set WALLET_TYPE to safe or deposit' >&2; exit 2 ;; esac
case "${POLY_LIVE_PROTOCOL:-}" in v1|v2) ;; *) printf '%s\n' 'Set PROTOCOL to v1 or v2' >&2; exit 2 ;; esac
# Resolve the output relative to the caller, before entering the SDK directory.
result_path="$2"
if [[ "$result_path" != /* ]]; then result_path="$PWD/$result_path"; fi
export POLY_LIVE_RESULT_FILE="$result_path"
export POLY_LIVE_ACTION="$1"
cd "$(dirname "${BASH_SOURCE[0]}")/.."
export POLY_LIVE_SDK_REVISION="$(git describe --always --dirty --abbrev=40)"
exec go test -mod=mod -tags polymarket_live -count=1 -v -run "^$test_name$" ./integration/live
