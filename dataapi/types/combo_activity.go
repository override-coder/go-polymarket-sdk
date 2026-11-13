package types

import "encoding/json"

// ComboActivityQuery selects the user-anchored Combo feed. Re-send User and
// Condition alongside Cursor on every page; the cursor only carries the anchor.
type ComboActivityQuery struct {
	User      string
	Condition []string
	Limit     *int
	Cursor    *string
}

// ComboActivity is a lifecycle or redemption event. Optional amounts remain
// nil when unavailable; decimal numbers never pass through float64.
type ComboActivity struct {
	ID               string       `json:"id"`
	BlockNumber      int64        `json:"block_number"`
	ComboConditionID string       `json:"combo_condition_id"`
	ComboPositionID  string       `json:"combo_position_id"`
	ProxyWallet      string       `json:"proxy_wallet"`
	Timestamp        int64        `json:"timestamp"`
	TransactionHash  string       `json:"transaction_hash"`
	Type             string       `json:"type"`
	AmountUSDC       *json.Number `json:"amount_usdc,omitempty"`
	PayoutUSDC       *json.Number `json:"payout_usdc,omitempty"`
	Legs             []ComboLeg   `json:"legs"`
}
