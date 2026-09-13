package combos_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/override-coder/go-polymarket-sdk/combos"
	combostypes "github.com/override-coder/go-polymarket-sdk/combos/types"
	sdkhttp "github.com/override-coder/go-polymarket-sdk/http"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
	"github.com/polymarket/go-order-utils/pkg/model"
	"github.com/stretchr/testify/require"
)

// Fill these values before running live tests. Environment variables with the
// same names take precedence, so secrets do not need to be committed locally.
var liveCombo = struct {
	PrivateKey    string
	SignerAddress string
	MakerAddress  string
	SignatureType model.SignatureType

	APIKey     string
	APISecret  string
	Passphrase string

	RFQHost  string
	DataHost string
	WSHost   string

	RFQID         string
	YesPositionID string
	Direction     combostypes.Direction
	PriceE6       string
	SizeE6        string
}{
	PrivateKey:    "", // COMBOS_PRIVATE_KEY
	SignerAddress: "", // COMBOS_SIGNER_ADDRESS
	MakerAddress:  "", // COMBOS_MAKER_ADDRESS / funder / quoter wallet
	SignatureType: model.POLY_GNOSIS_SAFE,

	APIKey:     "", // COMBOS_API_KEY
	APISecret:  "", // COMBOS_API_SECRET
	Passphrase: "", // COMBOS_API_PASSPHRASE

	RFQHost:  combos.DefaultRFQAPIHost,
	DataHost: combos.DefaultDataAPIHost,
	WSHost:   combos.DefaultRFQWSHost,

	RFQID:         "", // COMBOS_RFQ_ID, required for REST submit test
	YesPositionID: "", // COMBOS_YES_POSITION_ID, required for quote build/submit
	Direction:     combostypes.DirectionBuy,
	PriceE6:       "100000",  // 0.10
	SizeE6:        "1000000", // 1 share
}

func TestComboLiveGetMarkets(t *testing.T) {
	requireLiveReadOnly(t)

	client := newLiveComboClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	markets, err := client.GetComboMarkets(ctx, combostypes.ComboMarketsQuery{})
	require.NoError(t, err)
	require.NotNil(t, markets)
	t.Logf("combo markets=%d next_cursor=%q", len(markets.Markets), markets.NextCursor)
	if len(markets.Markets) > 0 {
		t.Logf("first market: condition_id=%s title=%s yes_position=%s", markets.Markets[0].ConditionID, markets.Markets[0].Title, first(markets.Markets[0].PositionIDs))
	}
}

func TestComboLiveBuildQuote(t *testing.T) {
	requireLiveReadOnly(t)

	cfg := loadLiveComboConfig()
	client := newLiveComboClient(t)
	yesPositionID := cfg.YesPositionID
	if yesPositionID == "" {
		yesPositionID = firstLiveComboYesPositionID(t, client)
	}

	quote, err := client.OrderBuilder().BuildQuoteRequest(combos.BuildQuoteRequestInput{
		RFQID:         firstNonEmpty(cfg.RFQID, "rfq_local_build_only"),
		QuoteID:       newQuoteID(),
		Direction:     cfg.Direction,
		YesPositionID: yesPositionID,
		PriceE6:       cfg.PriceE6,
		SizeE6:        cfg.SizeE6,
		ValidUntil:    time.Now().Add(2 * time.Second).UnixMilli(),
	}, liveAuthOption(cfg))

	require.NoError(t, err)
	t.Logf("quote_id=%s side=%d maker=%s signer=%s makerAmount=%s takerAmount=%s timestamp=%s",
		quote.QuoteID,
		quote.SignedOrder.Side,
		quote.SignedOrder.Maker,
		quote.SignedOrder.Signer,
		quote.SignedOrder.MakerAmount,
		quote.SignedOrder.TakerAmount,
		quote.SignedOrder.Timestamp,
	)
}

func TestComboLiveSubmitAndCancelQuoteREST(t *testing.T) {
	requireLiveSubmit(t)

	cfg := loadLiveComboConfig()
	require.NotEmpty(t, cfg.RFQID, "fill COMBOS_RFQ_ID or liveCombo.RFQID from an active RFQ_REQUEST")
	require.NotEmpty(t, cfg.YesPositionID, "fill COMBOS_YES_POSITION_ID or liveCombo.YesPositionID")

	client := newLiveComboClient(t)
	option := liveAuthOption(cfg)
	quote, err := client.OrderBuilder().BuildQuoteRequest(combos.BuildQuoteRequestInput{
		RFQID:         cfg.RFQID,
		QuoteID:       newQuoteID(),
		Direction:     cfg.Direction,
		YesPositionID: cfg.YesPositionID,
		PriceE6:       cfg.PriceE6,
		SizeE6:        cfg.SizeE6,
		ValidUntil:    time.Now().Add(1500 * time.Millisecond).UnixMilli(),
	}, option)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	snapshot, err := client.SubmitQuote(ctx, *quote, option)
	require.NoError(t, err)
	require.NotNil(t, snapshot)
	t.Logf("submitted quote_id=%s status=%s", quote.QuoteID, snapshot.Status)

	if os.Getenv("COMBOS_SKIP_CANCEL") == "1" {
		return
	}

	cancelReq := combostypes.QuoteCancelRequest{
		RFQID:         quote.RFQID,
		QuoteID:       quote.QuoteID,
		SignerAddress: quote.SignerAddress,
		MakerAddress:  quote.MakerAddress,
		SignatureType: quote.SignatureType,
	}
	cancelSnapshot, err := client.CancelQuote(ctx, cancelReq, option)
	require.NoError(t, err)
	require.NotNil(t, cancelSnapshot)
	t.Logf("cancel requested quote_id=%s status=%s", quote.QuoteID, cancelSnapshot.Status)
}

func TestComboLiveWebSocketReadOnly(t *testing.T) {
	requireLiveWS(t)

	cfg := loadLiveComboConfig()
	ws := combos.NewWSClient(cfg.WSHost)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	require.NoError(t, ws.Connect(ctx, liveAuthOption(cfg)))
	defer ws.Close()

	for {
		event, err := ws.Read(ctx)
		if ctx.Err() != nil {
			t.Log("websocket read window elapsed")
			return
		}
		require.NoError(t, err)
		t.Logf("event type=%s raw=%s", event.Type, string(event.Raw))
	}
}

func TestComboLiveWebSocketAutoQuoteDeclineLastLook(t *testing.T) {
	requireLiveWSAutoQuote(t)

	runComboLiveWebSocketAutoQuote(t, false)
}

func TestComboLiveWebSocketAutoQuoteAcceptLastLook(t *testing.T) {
	requireLiveWSAutoQuote(t)
	requireLiveAcceptLastLook(t)

	runComboLiveWebSocketAutoQuote(t, true)
}

func TestComboLiveTakerPlace(t *testing.T) {
	requireLiveTaker(t)

	cfg := loadLiveComboConfig()
	client := newLiveComboClient(t)
	legPositionIDs := liveTakerLegPositionIDs(t, client)
	requestedSizeE6 := envOr("COMBOS_TAKER_SIZE_E6", "1000000")
	requireMaxLiveE6(t, "COMBOS_TAKER_SIZE_E6", requestedSizeE6, 1000000)

	runner := combos.NewTakerRunner(
		combos.NewTakerWSClient(envOr("COMBOS_TAKER_WS_HOST", combos.DefaultTakerRFQWSHost)),
		client.OrderBuilder(),
		liveAuthOption(cfg),
	)
	runner.OnEvent(func(_ context.Context, event combos.Event) {
		if event.ExecutionUpdate != nil {
			t.Logf("taker execution status=%s tx_hash=%s", event.ExecutionUpdate.Status, event.ExecutionUpdate.TxHash)
			return
		}
		if event.StatusUpdate != nil {
			t.Logf("taker status=%s code=%s", event.StatusUpdate.Status, event.StatusUpdate.Code)
			return
		}
		if event.Error != nil {
			t.Logf("taker error code=%s message=%s", event.Error.Code, firstNonEmpty(event.Error.Error, event.Error.Reason, event.Error.Message))
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), liveTakerTimeout())
	defer cancel()
	execution, err := runner.Place(ctx, combos.TakerRequest{
		LegPositionIDs: legPositionIDs,
		Direction:      cfg.Direction,
		Side:           liveTakerSide(t),
		RequestedSize: combostypes.RequestedSize{
			Unit:    "notional",
			ValueE6: requestedSizeE6,
		},
	})
	require.NoError(t, err)
	require.Equal(t, "CONFIRMED", execution.Status)
	require.NotEmpty(t, execution.TxHash)
	t.Logf("combo taker confirmed tx_hash=%s", execution.TxHash)
}

// TestComboLiveBuilderTakerStatus validates the official Builder Gateway path
// without creating or accepting an RFQ. The gateway host is provisioned during
// builder onboarding and the RFQ ID must be from a previously accepted request.
func TestComboLiveBuilderTakerStatus(t *testing.T) {
	requireLiveReadOnly(t)
	host := envOr("COMBOS_BUILDER_GATEWAY_HOST", combos.DefaultBuilderTakerGatewayHost)
	rfqID := os.Getenv("COMBOS_BUILDER_RFQ_ID")
	if rfqID == "" {
		t.Skip("set COMBOS_BUILDER_RFQ_ID to check Builder Gateway status")
	}

	cfg := loadLiveComboConfig()
	client := combos.NewBuilderTakerClient(host, big.NewInt(137), liveSignatureFunc(t, cfg))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	status, err := client.GetRFQStatus(ctx, rfqID, liveAuthOption(cfg))
	require.NoError(t, err)
	require.Equal(t, rfqID, status.RFQID)
	require.NotEmpty(t, status.Status)
	t.Logf("builder taker status=%s tx_hash=%s", status.Status, status.TxHash)
}

// TestComboLiveBuilderTakerUnknownStatus is a non-mutating production smoke
// test equivalent to the official TypeScript SDK integration test. A 404 proves
// that the Builder Gateway host and account L2 authentication were accepted.
func TestComboLiveBuilderTakerUnknownStatus(t *testing.T) {
	requireLiveReadOnly(t)
	cfg := loadLiveComboConfig()
	client := combos.NewBuilderTakerClient("", big.NewInt(137), liveSignatureFunc(t, cfg))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.GetRFQStatus(ctx, "rfq-00000000-0000-0000-0000-000000000000", liveAuthOption(cfg))
	var upstreamErr *sdkhttp.UpstreamServiceError
	require.ErrorAs(t, err, &upstreamErr)
	require.Equal(t, 404, upstreamErr.StatusCode)
}

// TestComboLiveBuilderTakerCreateRFQ validates both account and Builder HMAC
// authentication. It does not accept the quote and therefore cannot trade.
func TestComboLiveBuilderTakerCreateRFQ(t *testing.T) {
	requireLiveBuilderTakerRequest(t)
	cfg := loadLiveComboConfig()
	client := combos.NewBuilderTakerClient("", big.NewInt(137), liveSignatureFunc(t, cfg))
	marketClient := newLiveComboClient(t)
	legPositionIDs := liveTakerLegPositionIDs(t, marketClient)
	requestedSizeE6 := envOr("COMBOS_BUILDER_TAKER_REQUEST_SIZE_E6", "1000000")
	requireMaxLiveE6(t, "COMBOS_BUILDER_TAKER_REQUEST_SIZE_E6", requestedSizeE6, 1000000)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rfq, err := client.CreateRFQ(ctx, combos.TakerRequest{
		LegPositionIDs: legPositionIDs,
		Direction:      combostypes.DirectionBuy,
		Side:           combostypes.ComboSideYes,
		RequestedSize: combostypes.RequestedSize{
			Unit:    "notional",
			ValueE6: requestedSizeE6,
		},
	}, liveBuilderTakerAuth(t, cfg))
	require.NoError(t, err)
	require.NotEmpty(t, rfq.RFQID)
	require.NotEmpty(t, rfq.Status)
	t.Logf("builder RFQ created rfq_id=%s status=%s has_quote=%t", rfq.RFQID, rfq.Status, rfq.Quote != nil)
}

// TestComboLiveBuilderTakerPlace executes the official Builder Gateway flow.
// It is intentionally gated because accepting the resulting RFQ can trade.
func TestComboLiveBuilderTakerPlace(t *testing.T) {
	requireLiveBuilderTaker(t)
	host := envOr("COMBOS_BUILDER_GATEWAY_HOST", combos.DefaultBuilderTakerGatewayHost)

	cfg := loadLiveComboConfig()
	client := combos.NewBuilderTakerClient(host, big.NewInt(137), liveSignatureFunc(t, cfg))
	marketClient := newLiveComboClient(t)
	legPositionIDs := liveTakerLegPositionIDs(t, marketClient)
	requestedSizeE6 := envOr("COMBOS_BUILDER_TAKER_SIZE_E6", "1000000")
	requireMaxLiveE6(t, "COMBOS_BUILDER_TAKER_SIZE_E6", requestedSizeE6, 1000000)

	ctx, cancel := context.WithTimeout(context.Background(), liveTakerTimeout())
	defer cancel()
	fill, err := client.PlaceAndWait(ctx, combos.TakerRequest{
		LegPositionIDs: legPositionIDs,
		Direction:      combostypes.DirectionBuy,
		Side:           combostypes.ComboSideYes,
		RequestedSize: combostypes.RequestedSize{
			Unit:    "notional",
			ValueE6: requestedSizeE6,
		},
	}, liveBuilderTakerAuth(t, cfg), time.Second)
	require.NoError(t, err)
	if fill.Status == "FAILED" {
		t.Skipf("Builder Gateway returned no executable quote: %+v", fill.Error)
	}
	require.NotEqual(t, "EXPIRED", fill.Status)
	require.Contains(t, []string{"CONFIRMED", "FILLED"}, fill.Status)
	require.NotEmpty(t, fill.TxHash)
	t.Logf("builder taker confirmed tx_hash=%s", fill.TxHash)
}

func runComboLiveWebSocketAutoQuote(t *testing.T, acceptLastLook bool) {
	t.Helper()

	cfg := loadLiveComboConfig()
	client := newLiveComboClient(t)
	ws := combos.NewWSClient(cfg.WSHost)
	runner := combos.NewMakerRunner(ws, liveAuthOption(cfg))
	quoteCount := 0
	maxQuotes := liveMaxQuotes()
	minDeadlineMs := liveMinDeadlineMs()
	ackReceived := false
	lastLookReceived := false
	confirmationAckReceived := false
	matchedReceived := false
	executionUpdateReceived := false
	quotedRFQIDs := map[string]bool{}
	quotedQuoteIDs := map[string]bool{}
	ackedRFQIDs := map[string]bool{}
	rfqError := ""
	prerequisiteError := ""
	executionError := ""
	ctx, cancel := context.WithTimeout(context.Background(), liveWSTimeout())
	defer cancel()
	var maxQuoteTimer *time.Timer
	defer func() {
		if maxQuoteTimer != nil {
			maxQuoteTimer.Stop()
		}
	}()

	runner.OnRFQ(func(ctx context.Context, rfq combostypes.RFQRequest) (*combostypes.QuoteRequest, error) {
		if quoteCount >= maxQuotes {
			return nil, nil
		}
		if directionFilter := liveRFQDirectionFilter(); directionFilter != "" && rfq.Direction != directionFilter {
			return nil, nil
		}
		nowMs := time.Now().UnixMilli()
		if rfq.SubmissionDeadline > 0 && rfq.SubmissionDeadline-nowMs < minDeadlineMs {
			if os.Getenv("COMBOS_LOG_STALE") == "1" {
				t.Logf("skipping stale rfq_id=%s now_ms=%d deadline_ms=%d", rfq.RFQID, nowMs, rfq.SubmissionDeadline)
			}
			return nil, nil
		}
		quoteCount++
		quotedRFQIDs[rfq.RFQID] = true
		quote, err := client.OrderBuilder().BuildQuoteRequest(combos.BuildQuoteRequestInput{
			RFQID:         rfq.RFQID,
			QuoteID:       newQuoteID(),
			Direction:     rfq.Direction,
			YesPositionID: rfq.YesPositionID,
			PriceE6:       cfg.PriceE6,
			SizeE6:        cfg.SizeE6,
			ValidUntil:    time.Now().Add(1500 * time.Millisecond).UnixMilli(),
		}, liveAuthOption(cfg))
		if err != nil {
			return nil, err
		}
		quotedQuoteIDs[quote.QuoteID] = true
		if os.Getenv("COMBOS_LOG_QUOTE_PAYLOAD") == "1" {
			wsPayload := combostypes.WSQuoteRequest{
				Type:        combostypes.MessageTypeRFQQuote,
				RFQID:       quote.RFQID,
				PriceE6:     quote.PriceE6,
				SizeE6:      quote.SizeE6,
				SignedOrder: quote.SignedOrder,
			}
			payloadBytes, _ := json.Marshal(wsPayload)
			t.Logf("sending ws quote #%d/%d deadline_ms=%d now_ms=%d payload=%s", quoteCount, maxQuotes, rfq.SubmissionDeadline, nowMs, payloadBytes)
		} else {
			t.Logf("sending ws quote #%d/%d rfq_id=%s deadline_ms=%d now_ms=%d side=%d makerAmount=%s takerAmount=%s",
				quoteCount,
				maxQuotes,
				quote.RFQID,
				rfq.SubmissionDeadline,
				nowMs,
				quote.SignedOrder.Side,
				quote.SignedOrder.MakerAmount,
				quote.SignedOrder.TakerAmount,
			)
		}
		if quoteCount == maxQuotes && maxQuoteTimer == nil {
			wait := livePostMaxQuoteWait()
			t.Logf("max quotes reached; waiting %s for last-look/execution updates", wait)
			maxQuoteTimer = time.AfterFunc(wait, cancel)
		}
		return quote, nil
	})
	runner.OnLastLook(func(ctx context.Context, req combostypes.ConfirmationRequest) (combostypes.ConfirmationDecision, error) {
		if !quotedRFQIDs[req.RFQID] && !quotedQuoteIDs[req.QuoteID] {
			t.Logf("declining unrelated last-look rfq_id=%s quote_id=%s fill_size_e6=%s price_e6=%s", req.RFQID, req.QuoteID, req.FillSizeE6, req.PriceE6)
			return combostypes.ConfirmationDecisionDecline, nil
		}
		lastLookReceived = true
		if acceptLastLook {
			t.Logf("accepting last-look rfq_id=%s quote_id=%s fill_size_e6=%s price_e6=%s confirm_by=%d",
				req.RFQID,
				req.QuoteID,
				req.FillSizeE6,
				req.PriceE6,
				req.ConfirmBy,
			)
			return combostypes.ConfirmationDecisionConfirm, nil
		}
		t.Logf("declining last-look rfq_id=%s quote_id=%s fill_size_e6=%s price_e6=%s", req.RFQID, req.QuoteID, req.FillSizeE6, req.PriceE6)
		return combostypes.ConfirmationDecisionDecline, nil
	})
	runner.OnEvent(func(ctx context.Context, event combos.Event) {
		if event.RFQRequest != nil {
			if directionFilter := liveRFQDirectionFilter(); directionFilter == "" || event.RFQRequest.Direction == directionFilter {
				t.Logf("event type=%s rfq_id=%s condition_id=%s yes_position_id=%s direction=%s size=%s deadline_ms=%d",
					event.Type,
					event.RFQRequest.RFQID,
					event.RFQRequest.ConditionID,
					event.RFQRequest.YesPositionID,
					event.RFQRequest.Direction,
					event.RFQRequest.RequestedSize.ValueE6,
					event.RFQRequest.SubmissionDeadline,
				)
			}
		} else {
			t.Logf("event type=%s raw=%s", event.Type, string(event.Raw))
		}
		if event.QuoteAck != nil {
			if !quotedRFQIDs[event.QuoteAck.RFQID] && !quotedQuoteIDs[event.QuoteAck.QuoteID] {
				t.Logf("ignoring unrelated ACK_RFQ_QUOTE rfq_id=%s quote_id=%s", event.QuoteAck.RFQID, event.QuoteAck.QuoteID)
				return
			}
			ackReceived = true
			ackedRFQIDs[event.QuoteAck.RFQID] = true
			t.Logf("ACK_RFQ_QUOTE rfq_id=%s quote_id=%s", event.QuoteAck.RFQID, event.QuoteAck.QuoteID)
			if !acceptLastLook {
				cancel()
			}
		}
		if event.ConfirmationRequest != nil {
			if !quotedRFQIDs[event.ConfirmationRequest.RFQID] && !quotedQuoteIDs[event.ConfirmationRequest.QuoteID] {
				t.Logf("unrelated RFQ_CONFIRMATION_REQUEST rfq_id=%s quote_id=%s fill_size_e6=%s price_e6=%s",
					event.ConfirmationRequest.RFQID,
					event.ConfirmationRequest.QuoteID,
					event.ConfirmationRequest.FillSizeE6,
					event.ConfirmationRequest.PriceE6,
				)
				return
			}
			t.Logf("RFQ_CONFIRMATION_REQUEST rfq_id=%s quote_id=%s fill_size_e6=%s price_e6=%s confirm_by=%d",
				event.ConfirmationRequest.RFQID,
				event.ConfirmationRequest.QuoteID,
				event.ConfirmationRequest.FillSizeE6,
				event.ConfirmationRequest.PriceE6,
				event.ConfirmationRequest.ConfirmBy,
			)
		}
		if event.ConfirmationResponseAck != nil {
			if !quotedRFQIDs[event.ConfirmationResponseAck.RFQID] && !quotedQuoteIDs[event.ConfirmationResponseAck.QuoteID] {
				t.Logf("ignoring unrelated ACK_RFQ_CONFIRMATION_RESPONSE rfq_id=%s quote_id=%s decision=%s",
					event.ConfirmationResponseAck.RFQID,
					event.ConfirmationResponseAck.QuoteID,
					event.ConfirmationResponseAck.Decision,
				)
				return
			}
			confirmationAckReceived = true
			t.Logf("ACK_RFQ_CONFIRMATION_RESPONSE rfq_id=%s quote_id=%s decision=%s",
				event.ConfirmationResponseAck.RFQID,
				event.ConfirmationResponseAck.QuoteID,
				event.ConfirmationResponseAck.Decision,
			)
		}
		if event.ExecutionUpdate != nil {
			if !quotedRFQIDs[event.ExecutionUpdate.RFQID] && !quotedQuoteIDs[event.ExecutionUpdate.QuoteID] {
				t.Logf("ignoring unrelated RFQ_EXECUTION_UPDATE rfq_id=%s quote_id=%s status=%s", event.ExecutionUpdate.RFQID, event.ExecutionUpdate.QuoteID, event.ExecutionUpdate.Status)
				return
			}
			if !ackedRFQIDs[event.ExecutionUpdate.RFQID] && event.ExecutionUpdate.QuoteID == "" {
				t.Logf("ignoring RFQ_EXECUTION_UPDATE without quote_id for unacked rfq_id=%s status=%s", event.ExecutionUpdate.RFQID, event.ExecutionUpdate.Status)
				return
			}
			t.Logf("RFQ_EXECUTION_UPDATE rfq_id=%s quote_id=%s status=%s tx_hash=%s raw=%s",
				event.ExecutionUpdate.RFQID,
				event.ExecutionUpdate.QuoteID,
				event.ExecutionUpdate.Status,
				event.ExecutionUpdate.TxHash,
				string(event.Raw),
			)
			switch event.ExecutionUpdate.Status {
			case "MATCHED":
				matchedReceived = true
			case "MINED", "CONFIRMED":
				executionUpdateReceived = true
				cancel()
			case "FAILED", "REVERTED", "CANCELLED", "EXPIRED":
				executionError = string(event.Raw)
				cancel()
			}
		}
		if event.Error != nil && quotedRFQIDs[event.Error.RFQID] {
			message := firstNonEmpty(event.Error.Error, event.Error.Reason, event.Error.Message, event.Error.Code)
			if event.Error.Code == "SUBMISSION_WINDOW_CLOSED" {
				t.Logf("retryable RFQ_ERROR rfq_id=%s code=%s message=%s", event.Error.RFQID, event.Error.Code, message)
				return
			}
			if event.Error.Code == "BALANCE_VALIDATION_FAILED" || event.Error.Code == "ALLOWANCE_VALIDATION_FAILED" {
				prerequisiteError = message
				cancel()
				return
			}
			rfqError = message
			cancel()
		}
	})
	runner.OnError(func(ctx context.Context, err error) {
		t.Logf("runner error: %v", err)
	})

	err := runner.Run(ctx)
	if !acceptLastLook && ackReceived {
		return
	}
	if acceptLastLook && executionUpdateReceived {
		return
	}
	require.Empty(t, executionError, "quoted RFQ execution failed")
	if acceptLastLook && matchedReceived {
		t.Skipf("RFQ reached MATCHED, but no mined/confirmed execution update arrived after %d quote attempt(s)", quoteCount)
	}
	if acceptLastLook && confirmationAckReceived {
		t.Skip("last-look was accepted and acknowledged, but no execution update arrived before websocket closed")
		return
	}
	if prerequisiteError != "" {
		if os.Getenv("COMBOS_REQUIRE_ACK") == "1" {
			require.Empty(t, prerequisiteError, "live wallet prerequisite failed")
		}
		t.Skipf("live wallet prerequisite failed before ACK_RFQ_QUOTE: %s", prerequisiteError)
	}
	require.Empty(t, rfqError, "quoted RFQ returned an error")
	if os.Getenv("COMBOS_REQUIRE_ACK") != "1" && (err == context.DeadlineExceeded || err == context.Canceled) {
		if quoteCount == 0 {
			t.Skip("no matching RFQ arrived before timeout")
		}
		if acceptLastLook && ackReceived && !lastLookReceived {
			t.Skipf("ACK_RFQ_QUOTE received after %d quote attempt(s), but no last-look arrived before websocket ended", quoteCount)
		}
		if acceptLastLook && lastLookReceived {
			t.Skipf("last-look was accepted, but no mined/confirmed execution update arrived before websocket ended after %d quote attempt(s)", quoteCount)
		}
		if acceptLastLook && matchedReceived {
			t.Skipf("RFQ reached MATCHED, but no mined/confirmed execution update arrived before websocket ended after %d quote attempt(s)", quoteCount)
		}
		t.Skipf("ACK_RFQ_QUOTE was not received before websocket ended after %d quote attempt(s)", quoteCount)
	}
	if os.Getenv("COMBOS_REQUIRE_ACK") != "1" && isRetryableLiveWSError(err) {
		if quoteCount == 0 {
			t.Skipf("websocket closed before a matching RFQ arrived: %v", err)
		}
		if acceptLastLook && ackReceived && !lastLookReceived {
			t.Skipf("ACK_RFQ_QUOTE received after %d quote attempt(s), but websocket closed before last-look: %v", quoteCount, err)
		}
		t.Skipf("websocket closed after %d quote attempt(s): %v", quoteCount, err)
	}
	require.NoError(t, err)
	require.Fail(t, "ACK_RFQ_QUOTE was not received before timeout")
}

func requireLiveReadOnly(t *testing.T) {
	t.Helper()
	if os.Getenv("COMBOS_LIVE") != "1" {
		t.Skip("set COMBOS_LIVE=1 to run live combo tests")
	}
}

func requireLiveSubmit(t *testing.T) {
	t.Helper()
	requireLiveReadOnly(t)
	if os.Getenv("COMBOS_LIVE_SUBMIT") != "1" {
		t.Skip("set COMBOS_LIVE_SUBMIT=1 to submit a live REST quote")
	}
}

func requireLiveWS(t *testing.T) {
	t.Helper()
	requireLiveReadOnly(t)
	if os.Getenv("COMBOS_LIVE_WS") != "1" {
		t.Skip("set COMBOS_LIVE_WS=1 to connect to live RFQ websocket")
	}
}

func requireLiveWSAutoQuote(t *testing.T) {
	t.Helper()
	requireLiveWS(t)
	if os.Getenv("COMBOS_LIVE_WS_AUTO_QUOTE") != "1" {
		t.Skip("set COMBOS_LIVE_WS_AUTO_QUOTE=1 to auto-submit quotes over websocket")
	}
}

func requireLiveAcceptLastLook(t *testing.T) {
	t.Helper()
	if os.Getenv("COMBOS_LIVE_ACCEPT_LAST_LOOK") != "1" {
		t.Skip("set COMBOS_LIVE_ACCEPT_LAST_LOOK=1 to accept last-look and allow a live fill")
	}
	require.Equal(t, string(combostypes.DirectionSell), string(liveRFQDirectionFilter()), "set COMBOS_RFQ_DIRECTION_FILTER=SELL to buy YES with USDC")
	requireMaxLiveE6(t, "COMBOS_PRICE_E6", envOr("COMBOS_PRICE_E6", liveCombo.PriceE6), 1000000)
	requireMaxLiveE6(t, "COMBOS_SIZE_E6", envOr("COMBOS_SIZE_E6", liveCombo.SizeE6), 1000000)
}

func requireLiveTaker(t *testing.T) {
	t.Helper()
	requireLiveReadOnly(t)
	if os.Getenv("COMBOS_LIVE_TAKER") != "1" {
		t.Skip("set COMBOS_LIVE_TAKER=1 to run a live combo taker order")
	}
	if os.Getenv("COMBOS_LIVE_TAKER_ACK") != "I_UNDERSTAND_THIS_CAN_TRADE" {
		t.Skip("set COMBOS_LIVE_TAKER_ACK=I_UNDERSTAND_THIS_CAN_TRADE to allow a live fill")
	}
}

func requireLiveBuilderTaker(t *testing.T) {
	t.Helper()
	requireLiveReadOnly(t)
	if os.Getenv("COMBOS_LIVE_BUILDER_TAKER") != "1" {
		t.Skip("set COMBOS_LIVE_BUILDER_TAKER=1 to run an official Builder Gateway taker order")
	}
	if os.Getenv("COMBOS_LIVE_BUILDER_TAKER_ACK") != "I_UNDERSTAND_THIS_CAN_TRADE" {
		t.Skip("set COMBOS_LIVE_BUILDER_TAKER_ACK=I_UNDERSTAND_THIS_CAN_TRADE to allow a live fill")
	}
	for _, name := range []string{"COMBOS_BUILDER_API_KEY", "COMBOS_BUILDER_SECRET", "COMBOS_BUILDER_PASSPHRASE"} {
		require.NotEmpty(t, os.Getenv(name), "%s is required", name)
	}
}

func requireLiveBuilderTakerRequest(t *testing.T) {
	t.Helper()
	requireLiveReadOnly(t)
	if os.Getenv("COMBOS_LIVE_BUILDER_TAKER_REQUEST") != "1" {
		t.Skip("set COMBOS_LIVE_BUILDER_TAKER_REQUEST=1 to request a live Builder Gateway quote without accepting it")
	}
	for _, name := range []string{"COMBOS_BUILDER_API_KEY", "COMBOS_BUILDER_SECRET", "COMBOS_BUILDER_PASSPHRASE"} {
		require.NotEmpty(t, os.Getenv(name), "%s is required", name)
	}
}

func requireMaxLiveE6(t *testing.T, name string, value string, max int64) {
	t.Helper()
	n, err := strconv.ParseInt(value, 10, 64)
	require.NoError(t, err, "invalid %s", name)
	require.LessOrEqual(t, n, max, "%s is capped for live accept-last-look tests", name)
}

func newLiveComboClient(t *testing.T) *combos.Client {
	t.Helper()
	cfg := loadLiveComboConfig()
	return combos.NewClientWithDataHost(cfg.RFQHost, cfg.DataHost, big.NewInt(137), liveSignatureFunc(t, cfg))
}

func firstLiveComboYesPositionID(t *testing.T, client *combos.Client) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	markets, err := client.GetComboMarkets(ctx, combostypes.ComboMarketsQuery{})
	require.NoError(t, err)
	require.NotEmpty(t, markets.Markets, "no combo markets returned; fill COMBOS_YES_POSITION_ID manually")
	require.NotEmpty(t, markets.Markets[0].PositionIDs, "first combo market has no position_ids; fill COMBOS_YES_POSITION_ID manually")

	yesPositionID := markets.Markets[0].PositionIDs[0]
	t.Logf("using first combo market yes_position_id=%s condition_id=%s title=%s", yesPositionID, markets.Markets[0].ConditionID, markets.Markets[0].Title)
	return yesPositionID
}

func liveAuthOption(cfg liveComboConfig) *sdktypes.AuthOption {
	return &sdktypes.AuthOption{
		SignatureType: cfg.SignatureType,
		SingerAddress: cfg.SignerAddress,
		FunderAddress: cfg.MakerAddress,
		ApiKeyCreds: &sdktypes.ApiKeyCreds{
			ApiKey:     cfg.APIKey,
			Secret:     cfg.APISecret,
			Passphrase: cfg.Passphrase,
		},
	}
}

func liveBuilderTakerAuth(t *testing.T, cfg liveComboConfig) *sdktypes.AuthOption {
	t.Helper()
	option := liveAuthOption(cfg)
	option.BuilderApiKeyCreds = &sdktypes.BuilderApiKeyCreds{
		Key:        os.Getenv("COMBOS_BUILDER_API_KEY"),
		Secret:     os.Getenv("COMBOS_BUILDER_SECRET"),
		Passphrase: os.Getenv("COMBOS_BUILDER_PASSPHRASE"),
	}
	return option
}

type liveComboConfig struct {
	PrivateKey    string
	SignerAddress string
	MakerAddress  string
	SignatureType model.SignatureType
	APIKey        string
	APISecret     string
	Passphrase    string
	RFQHost       string
	DataHost      string
	WSHost        string
	RFQID         string
	YesPositionID string
	Direction     combostypes.Direction
	PriceE6       string
	SizeE6        string
}

func loadLiveComboConfig() liveComboConfig {
	return liveComboConfig{
		PrivateKey:    envOr("COMBOS_PRIVATE_KEY", liveCombo.PrivateKey),
		SignerAddress: envOr("COMBOS_SIGNER_ADDRESS", liveCombo.SignerAddress),
		MakerAddress:  envOr("COMBOS_MAKER_ADDRESS", liveCombo.MakerAddress),
		SignatureType: liveSignatureType(),
		APIKey:        envOr("COMBOS_API_KEY", liveCombo.APIKey),
		APISecret:     envOr("COMBOS_API_SECRET", liveCombo.APISecret),
		Passphrase:    envOr("COMBOS_API_PASSPHRASE", liveCombo.Passphrase),
		RFQHost:       envOr("COMBOS_RFQ_HOST", liveCombo.RFQHost),
		DataHost:      envOr("COMBOS_DATA_HOST", liveCombo.DataHost),
		WSHost:        envOr("COMBOS_WS_HOST", liveCombo.WSHost),
		RFQID:         envOr("COMBOS_RFQ_ID", liveCombo.RFQID),
		YesPositionID: envOr("COMBOS_YES_POSITION_ID", liveCombo.YesPositionID),
		Direction:     combostypes.Direction(envOr("COMBOS_DIRECTION", string(liveCombo.Direction))),
		PriceE6:       envOr("COMBOS_PRICE_E6", liveCombo.PriceE6),
		SizeE6:        envOr("COMBOS_SIZE_E6", liveCombo.SizeE6),
	}
}

func liveSignatureType() model.SignatureType {
	switch envOr("COMBOS_SIGNATURE_TYPE", "") {
	case "0":
		return model.EOA
	case "1":
		return model.POLY_PROXY
	case "2":
		return model.POLY_GNOSIS_SAFE
	case "3":
		return model.POLY_1271
	default:
		return liveCombo.SignatureType
	}
}

func liveMaxQuotes() int {
	value := envOr("COMBOS_MAX_QUOTES", "3")
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 3
	}
	return n
}

func liveMinDeadlineMs() int64 {
	value := envOr("COMBOS_MIN_DEADLINE_MS", "0")
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

func liveWSTimeout() time.Duration {
	value := envOr("COMBOS_LIVE_WS_TIMEOUT_SECONDS", "60")
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 60 * time.Second
	}
	return time.Duration(n) * time.Second
}

func livePostMaxQuoteWait() time.Duration {
	value := envOr("COMBOS_POST_MAX_QUOTE_WAIT_SECONDS", "45")
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 45 * time.Second
	}
	return time.Duration(n) * time.Second
}

func liveRFQDirectionFilter() combostypes.Direction {
	return combostypes.Direction(envOr("COMBOS_RFQ_DIRECTION_FILTER", ""))
}

func liveTakerLegPositionIDs(t *testing.T, client *combos.Client) []string {
	t.Helper()
	if value := envOr("COMBOS_LEG_POSITION_IDS", ""); value != "" {
		var ids []string
		for _, id := range strings.Split(value, ",") {
			if id = strings.TrimSpace(id); id != "" {
				ids = append(ids, id)
			}
		}
		require.GreaterOrEqual(t, len(ids), 2, "COMBOS_LEG_POSITION_IDS must contain at least two comma-separated token ids")
		return ids
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	limit := 10
	markets, err := client.GetComboMarkets(ctx, combostypes.ComboMarketsQuery{Limit: &limit})
	require.NoError(t, err)

	ids := make([]string, 0, 2)
	for _, market := range markets.Markets {
		if len(market.PositionIDs) == 0 {
			continue
		}
		ids = append(ids, market.PositionIDs[0])
		if len(ids) == 2 {
			break
		}
	}
	require.GreaterOrEqual(t, len(ids), 2, "no two combo-eligible YES position ids found; set COMBOS_LEG_POSITION_IDS")
	t.Logf("using discovered taker leg position ids=%s", strings.Join(ids, ","))
	return ids
}

func liveTakerSide(t *testing.T) combostypes.ComboSide {
	t.Helper()
	side := combostypes.ComboSide(envOr("COMBOS_TAKER_SIDE", string(combostypes.ComboSideYes)))
	require.Contains(t, []combostypes.ComboSide{combostypes.ComboSideYes, combostypes.ComboSideNo}, side, "COMBOS_TAKER_SIDE must be YES or NO")
	return side
}

func liveTakerTimeout() time.Duration {
	value := envOr("COMBOS_LIVE_TAKER_TIMEOUT_SECONDS", "120")
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 120 * time.Second
	}
	return time.Duration(n) * time.Second
}

func isRetryableLiveWSError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "websocket: close 1006") ||
		strings.Contains(message, "unexpected EOF") ||
		strings.Contains(message, "websocket: bad handshake")
}

func liveSignatureFunc(t *testing.T, cfg liveComboConfig) func(string, []byte) ([]byte, error) {
	t.Helper()
	require.NotEmpty(t, cfg.PrivateKey, "fill COMBOS_PRIVATE_KEY or liveCombo.PrivateKey")
	privateKeyHex := strings.TrimPrefix(cfg.PrivateKey, "0x")
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(privateKeyHex))
	require.NoError(t, err)

	return func(_ string, digest []byte) ([]byte, error) {
		sig, err := crypto.Sign(digest, privateKey)
		if err != nil {
			return nil, err
		}
		sig[64] += 27
		return sig, nil
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func newQuoteID() string {
	return "quote_go_sdk_" + big.NewInt(time.Now().UnixNano()).String()
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
