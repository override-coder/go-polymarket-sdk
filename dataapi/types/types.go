package types

// Page is the V2 response envelope for cursor-paginated endpoints.
type Page[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"` // Display metadata only; never a request parameter.
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

type Envelope[T any] struct {
	Data T `json:"data"`
}

type SortBy string

const (
	SortByCURRENTVALUE  SortBy = "CURRENT_VALUE"
	SortByPRICE         SortBy = "PRICE"
	SortByTOKENS        SortBy = "TOKENS"
	SortByUNREALIZEDPNL SortBy = "UNREALIZED_PNL"
	SortByREALIZEDPNL   SortBy = "REALIZED_PNL"
	SortByTOTALPNL      SortBy = "TOTAL_PNL"
	SortByTIMESTAMP     SortBy = "TIMESTAMP"
)

type SortDirection string

const (
	SortASC  SortDirection = "ASC"
	SortDESC SortDirection = "DESC"
)

type PositionStatus string

const (
	PositionOPEN           PositionStatus = "OPEN"
	PositionREDEEMABLE     PositionStatus = "REDEEMABLE"
	PositionREDEEMABLELOST PositionStatus = "REDEEMABLE_LOST"
	PositionMERGEABLE      PositionStatus = "MERGEABLE"
	PositionCLOSED         PositionStatus = "CLOSED"
)

type FilterType string

const (
	FilterTOKENS FilterType = "TOKENS"
	FilterCASH   FilterType = "CASH"
)

// PositionsQuery follows /v2/positions. Re-send the user/condition anchor and
// filters on every page. Omitted status/sort/direction are adopted from Cursor.
type PositionsQuery struct {
	User            string   // User or Condition is required.
	Condition       []string // Up to 20 distinct IDs with User; exactly one without User.
	Limit           *int     // 0..1000 on the first page; ignored with Cursor.
	Cursor          *string
	Status          *PositionStatus // Default OPEN; CLOSED uses the same response shape.
	EventID         []int64         // Up to 20 distinct Gamma IDs; user-anchored only.
	Title           *string         // At most 200 characters; re-send on every page.
	FilterType      *FilterType     // Default TOKENS.
	FilterAmount    *float64        // Default token dust floor 0.1; CASH uses current value.
	IncludeArchived *bool           // Not supported with CLOSED.
	SortBy          *SortBy         // Server default depends on status.
	Start           *int64          // Inclusive last_event_at bound; omitted/0 is unbounded.
	End             *int64
	SortDirection   *SortDirection // Default DESC.
}

type Position struct {
	ProxyWallet        string         `json:"proxy_wallet"`
	TokenID            string         `json:"token_id"`
	ConditionID        string         `json:"condition_id"`
	CurrentSize        float64        `json:"current_size"`
	AvgPrice           float64        `json:"avg_price"`
	EntryCostUSDC      float64        `json:"entry_cost_usdc"` // Fee-exclusive basis.
	EntryFeesUSDC      float64        `json:"entry_fees_usdc"` // Disclosure; do not deduct again from PnL.
	TotalCostUSDC      float64        `json:"total_cost_usdc"`
	CurrentPrice       float64        `json:"current_price"`
	CurrentValue       float64        `json:"current_value"`
	TotalSize          float64        `json:"total_size"` // Lifetime bought shares, not current balance.
	RealizedPnl        float64        `json:"realized_pnl"`
	UnrealizedPnl      float64        `json:"unrealized_pnl"`
	TotalPnl           float64        `json:"total_pnl"`
	PercentPnl         float64        `json:"percent_pnl"`
	PercentRealizedPnl float64        `json:"percent_realized_pnl"` // Compatibility ratio, not realized_pnl / basis.
	Status             PositionStatus `json:"status"`
	Redeemable         bool           `json:"redeemable"` // Includes losing outcomes.
	Mergeable          bool           `json:"mergeable"`
	NegativeRisk       bool           `json:"negative_risk"`
	Archived           bool           `json:"archived"`
	Title              string         `json:"title"`
	Slug               string         `json:"slug"`
	Icon               string         `json:"icon"`
	EventID            string         `json:"event_id"`
	EventSlug          string         `json:"event_slug"`
	Outcome            string         `json:"outcome"`
	OutcomeIndex       int64          `json:"outcome_index"` // 999 means unlabeled.
	OppositeOutcome    string         `json:"opposite_outcome"`
	OppositeTokenID    string         `json:"opposite_token_id"`
	EndDate            string         `json:"end_date"`
	LastEventAt        int64          `json:"last_event_at"`
	FirstEntryAt       int64          `json:"first_entry_at"`
	Name               string         `json:"name"`
	ProfileImage       string         `json:"profile_image"`
	Verified           bool           `json:"verified"`
}

// ClosedPositionsQuery uses the unified V2 positions parameters. Status must
// be omitted or CLOSED; all other filters retain the /v2/positions contract.
type ClosedPositionsQuery = PositionsQuery

// ClosedPosition has the same complete V2 shape as an open position.
type ClosedPosition = Position

type ActivityType string

const (
	ActivityTIP            ActivityType = "TIP" // Opt-in user-to-user pUSD transfer.
	ActivityTRADE          ActivityType = "TRADE"
	ActivitySPLIT          ActivityType = "SPLIT"
	ActivityMERGE          ActivityType = "MERGE"
	ActivityREDEEM         ActivityType = "REDEEM"
	ActivityREWARD         ActivityType = "REWARD"
	ActivityCONVERSION     ActivityType = "CONVERSION"
	ActivityDEPOSIT        ActivityType = "DEPOSIT"
	ActivityWITHDRAWAL     ActivityType = "WITHDRAWAL"
	ActivityYIELD          ActivityType = "YIELD"
	ActivityMAKERREBATE    ActivityType = "MAKER_REBATE"
	ActivityTAKERREBATE    ActivityType = "TAKER_REBATE"
	ActivityREFERRALREWARD ActivityType = "REFERRAL_REWARD"
)

type ActivitySortBy string

const (
	ActivitySortTIMESTAMP ActivitySortBy = "TIMESTAMP" // default
)

type Side string

const (
	SideBUY  Side = "BUY"
	SideSELL Side = "SELL"
)

// ActivityQuery follows /v2/activity. Re-send identical filters on every page.
type ActivityQuery struct {
	User                       string // Required wallet.
	Limit                      *int   // Default 100, range 0..1000; cursor page size wins.
	Cursor                     *string
	Condition                  []string        // At most 20 distinct IDs; mutually exclusive with EventID.
	EventID                    []int64         // At most 20 distinct Gamma event IDs.
	Type                       []ActivityType  // TIP must be explicitly included.
	ExcludeDepositsWithdrawals *bool           // Default true; false includes deposits/withdrawals.
	Start                      *int64          // Omitted/0 means three years back; 1 requests full history.
	End                        *int64          // Omitted/0 means now plus one day.
	SortBy                     *ActivitySortBy // Only TIMESTAMP is supported.
	SortDirection              *SortDirection  // ASC/DESC; omitted direction is adopted from Cursor.
	Side                       *Side           // Query supports BUY/SELL; TIP response sides are IN/OUT.
}

type UserActivity struct {
	ProxyWallet           string  `json:"proxy_wallet"`
	Timestamp             int64   `json:"timestamp"`
	ConditionID           string  `json:"condition_id"`
	Type                  string  `json:"type"`
	Size                  float64 `json:"size"`
	USDCSize              float64 `json:"usdc_size"`
	TransactionHash       string  `json:"transaction_hash"`
	Price                 float64 `json:"price"`
	TokenID               string  `json:"token_id"`
	Side                  string  `json:"side"`
	OutcomeIndex          int64   `json:"outcome_index"`
	Title                 string  `json:"title"`
	Slug                  string  `json:"slug"`
	Icon                  string  `json:"icon"`
	EventSlug             string  `json:"event_slug"`
	Outcome               string  `json:"outcome"`
	Name                  string  `json:"name"`
	Pseudonym             string  `json:"pseudonym"`
	Bio                   string  `json:"bio"`
	ProfileImage          string  `json:"profile_image"`
	ProfileImageOptimized string  `json:"profile_image_optimized"`
	IsCombo               *bool   `json:"is_combo,omitempty"`
}

type PositionValueQuery struct {
	User      string   // Required wallet.
	Condition []string // Up to 20 distinct IDs; excludes the portfolio-level combo term.
}

type PositionValue struct {
	ProxyWallet string  `json:"proxy_wallet"`
	Value       float64 `json:"value"` // USDC, rounded by the service to four decimals.
}

type LeaderboardCategory string

const (
	LeaderboardCategoryCOMBOS    LeaderboardCategory = "combos" // PNL board only.
	LeaderboardCategoryESPORTS   LeaderboardCategory = "esports"
	LeaderboardCategoryOVERALL   LeaderboardCategory = "overall"
	LeaderboardCategoryPOLITICS  LeaderboardCategory = "politics"
	LeaderboardCategorySPORTS    LeaderboardCategory = "sports"
	LeaderboardCategoryCRYPTO    LeaderboardCategory = "crypto"
	LeaderboardCategoryCULTURE   LeaderboardCategory = "culture"
	LeaderboardCategoryMENTIONS  LeaderboardCategory = "mentions"
	LeaderboardCategoryWEATHER   LeaderboardCategory = "weather"
	LeaderboardCategoryECONOMICS LeaderboardCategory = "economics"
	LeaderboardCategoryTECH      LeaderboardCategory = "tech"
	LeaderboardCategoryFINANCE   LeaderboardCategory = "finance"
)

type LeaderboardTimePeriod string

const (
	LeaderboardTimeDAY   LeaderboardTimePeriod = "day"
	LeaderboardTimeWEEK  LeaderboardTimePeriod = "week"
	LeaderboardTimeMONTH LeaderboardTimePeriod = "month"
	LeaderboardTimeALL   LeaderboardTimePeriod = "all"
)

type LeaderboardSortBy string

const (
	LeaderboardSortByPNL    LeaderboardSortBy = "PNL"
	LeaderboardSortByVOLUME LeaderboardSortBy = "VOLUME"
)

// TraderLeaderboardQuery selects a ranked board. For the V2 user branch use
// TraderLeaderboardUserQuery and GetTraderLeaderboardUser instead.
type TraderLeaderboardQuery struct {
	Category   *LeaderboardCategory
	TimePeriod *LeaderboardTimePeriod
	SortBy     *LeaderboardSortBy
	Limit      *int    // 0..1000; ignored with Cursor.
	Cursor     *string // Pins category, time period and sort; omit defaults on resume.
}

type TraderLeaderboard struct {
	Rank         int     `json:"rank"` // Ties share a rank; never use rank to derive a cursor.
	UserID       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	XUsername    string  `json:"x_username"`
	Verified     bool    `json:"verified"`
	Volume       float64 `json:"volume"` // Both-sides shares, not USDC.
	PnL          float64 `json:"pnl"`    // Finite windows: marked equity change; all: realized-only.
	ProfileImage string  `json:"profile_image"`
}

// The user branch ignores sort_by, limit and cursor, so it exposes only its
// applicable parameters and returns one nullable object, not a page.
type TraderLeaderboardUserQuery struct {
	User       string
	Category   *LeaderboardCategory
	TimePeriod *LeaderboardTimePeriod
}

type TraderLeaderboardUser struct {
	UserID       string  `json:"user_id"`
	PnL          float64 `json:"pnl"`
	Volume       float64 `json:"volume"`
	RankPnL      *int    `json:"rank_pnl"` // nil means unranked.
	RankVolume   *int    `json:"rank_volume"`
	UserName     string  `json:"user_name"`
	ProfileImage string  `json:"profile_image"`
	XUsername    string  `json:"x_username"`
	Verified     bool    `json:"verified"`
}

type TopHoldersQuery struct {
	Condition  []string // Required; at most 20 distinct IDs, exactly one with IncludePnL.
	Limit      *int     // Per token: default 100, max 1000 (100 with IncludePnL).
	Cursor     *string  // Keep sending Condition and filters on every page.
	MinBalance *float64 // Default 0; net shares, or gross shares with IncludePnL.
	IncludePnL *bool    // Default false; changes Amount to per-side gross balances.
}

type TopHoldersForMarket struct {
	TokenID string      `json:"token_id"`
	Holders []TopHolder `json:"holders"`
}

type TopHolder struct {
	ProxyWallet           string  `json:"proxy_wallet"`
	Bio                   string  `json:"bio"`
	TokenID               string  `json:"token_id"`
	Pseudonym             string  `json:"pseudonym"`
	Amount                float64 `json:"amount"` // Net shares by default; per-side gross with include_pnl.
	DisplayUsernamePublic bool    `json:"display_username_public"`
	OutcomeIndex          int64   `json:"outcome_index"`
	Name                  string  `json:"name"`
	ProfileImage          string  `json:"profile_image"`
	ProfileImageOptimized string  `json:"profile_image_optimized"` // V2 currently returns empty.
	Verified              bool    `json:"verified"`
	// Economics are optional and nullable, present only with include_pnl=true.
	AvgPrice      *float64 `json:"avg_price,omitempty"`
	EntryCostUSDC *float64 `json:"entry_cost_usdc,omitempty"`
	CurrentPrice  *float64 `json:"current_price,omitempty"`
	CurrentValue  *float64 `json:"current_value,omitempty"`
	RealizedPnl   *float64 `json:"realized_pnl,omitempty"`
	UnrealizedPnl *float64 `json:"unrealized_pnl,omitempty"`
	TotalPnl      *float64 `json:"total_pnl,omitempty"`
}

type MarketByTokenResponse struct {
	ConditionID      string `json:"condition_id"`
	PrimaryTokenID   string `json:"primary_token_id"`
	SecondaryTokenID string `json:"secondary_token_id"`
}
