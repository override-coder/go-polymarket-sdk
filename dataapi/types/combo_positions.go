package types

import "encoding/json"

// ComboPositionsQuery is independent of ordinary positions. Continuation needs
// User and Cursor; the cursor adopts omitted filters and sorting.
type ComboPositionsQuery struct {
	User          string
	Condition     []string
	Limit         *int
	Cursor        *string
	Status        *string
	SortBy        *string
	SortDirection *SortDirection
	UpdatedAfter  *int64
	UpdatedBefore *int64
}

// ComboPosition preserves decimal JSON numbers so shares and proceeds need not
// pass through binary floating point when preparing a cash-out request.
type ComboPosition struct {
	ComboConditionID   string      `json:"combo_condition_id"`
	ComboPositionID    string      `json:"combo_position_id"`
	ProxyWallet        string      `json:"proxy_wallet"`
	OutcomeIndex       int32       `json:"outcome_index"`
	OutcomeLabel       string      `json:"outcome_label"`
	CurrentSize        json.Number `json:"current_size"`
	EntryAvgPriceUSDC  json.Number `json:"entry_avg_price_usdc"`
	EntryCostUSDC      json.Number `json:"entry_cost_usdc"`
	GrossEntryCostUSDC json.Number `json:"gross_entry_cost_usdc"`
	EntryFeesUSDC      json.Number `json:"entry_fees_usdc"`
	RealizedPayoutUSDC json.Number `json:"realized_payout_usdc"`
	Status             string      `json:"status"`
	Redeemable         bool        `json:"redeemable"`
	FirstEntryAt       string      `json:"first_entry_at"`
	FirstEntryAtMicros *int64      `json:"first_entry_at_micros"`
	UpdatedAt          string      `json:"updated_at"`
	UpdatedAtMicros    int64       `json:"updated_at_micros"`
	ResolvedAt         *string     `json:"resolved_at"`
	LegsTotal          int32       `json:"legs_total"`
	LegsResolved       int32       `json:"legs_resolved"`
	LegsPending        int32       `json:"legs_pending"`
	Legs               []ComboLeg  `json:"legs"`
}

type ComboLeg struct {
	LegIndex        int32          `json:"leg_index"`
	LegPositionID   string         `json:"leg_position_id"`
	LegConditionID  string         `json:"leg_condition_id"`
	LegOutcomeIndex int32          `json:"leg_outcome_index"`
	LegOutcomeLabel string         `json:"leg_outcome_label"`
	LegStatus       string         `json:"leg_status"`
	LegCurrentPrice json.Number    `json:"leg_current_price"`
	LegResolvedAt   *string        `json:"leg_resolved_at"`
	Market          ComboLegMarket `json:"market"`
}

type ComboLegMarket struct {
	MarketID         string        `json:"market_id"`
	Slug             string        `json:"slug"`
	Title            string        `json:"title"`
	Outcome          string        `json:"outcome"`
	ImageURL         string        `json:"image_url"`
	IconURL          string        `json:"icon_url"`
	Category         string        `json:"category"`
	Subcategory      string        `json:"subcategory"`
	Tags             []string      `json:"tags"`
	EndDate          string        `json:"end_date"`
	Event            ComboLegEvent `json:"event"`
	GroupItemTitle   string        `json:"group_item_title"`
	Line             *json.Number  `json:"line"`
	Outcomes         []string      `json:"outcomes"`
	Question         string        `json:"question"`
	SportsMarketType string        `json:"sports_market_type"`
}

type ComboLegEvent struct {
	EventID    string `json:"event_id"`
	EventSlug  string `json:"event_slug"`
	EventTitle string `json:"event_title"`
	EventImage string `json:"event_image"`
}
