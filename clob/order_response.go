package clob

import (
	"context"
	"strings"
	"time"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

const (
	defaultResolveTradesTimeout      = 30 * time.Second
	defaultResolveTradesPollInterval = 250 * time.Millisecond
	failedTradeStatus                = "FAILED"
)

// resolveTransactionsHashes fills TransactionsHashes for asynchronously
// committed fills. It returns the original response unchanged when the polling
// window expires or a trade cannot be fetched.
func (c *Client) resolveTransactionsHashes(ctx context.Context, response *clobtypes.OrderResponse, auth *sdktypes.AuthOption) {
	if response == nil || len(response.TransactionsHashes) > 0 {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	tradeIDs := uniqueTradeIDs(response.TradeIDs)
	if len(tradeIDs) == 0 {
		return
	}

	resolved := make(map[string]clobtypes.Trade, len(tradeIDs))
	pollCtx, cancel := context.WithTimeout(ctx, c.resolveTimeout())
	defer cancel()

poll:
	for {
		for _, tradeID := range tradeIDs {
			if _, ok := resolved[tradeID]; ok {
				continue
			}

			trades, err := c.GetTrades(pollCtx, clobtypes.GetTradesRequest{ID: &tradeID}, auth)
			if err != nil || trades == nil {
				continue
			}
			for _, trade := range trades.Data {
				if trade.ID == tradeID && isResolvedTrade(trade) {
					resolved[tradeID] = trade
				}
			}
		}

		if len(resolved) == len(tradeIDs) || pollCtx.Err() != nil {
			break
		}

		timer := time.NewTimer(c.resolvePollInterval())
		select {
		case <-pollCtx.Done():
			timer.Stop()
			break poll
		case <-timer.C:
		}
	}

	transactionHashes := make([]string, 0, len(tradeIDs))
	for _, tradeID := range tradeIDs {
		trade, ok := resolved[tradeID]
		if !ok || strings.EqualFold(trade.Status, failedTradeStatus) || trade.TransactionHash == "" {
			continue
		}
		transactionHashes = append(transactionHashes, trade.TransactionHash)
	}
	if len(transactionHashes) > 0 {
		response.TransactionsHashes = transactionHashes
	}
}

func (c *Client) resolveTimeout() time.Duration {
	if c.resolveTradesTimeout > 0 {
		return c.resolveTradesTimeout
	}
	return defaultResolveTradesTimeout
}

func (c *Client) resolvePollInterval() time.Duration {
	if c.resolveTradesPollInterval > 0 {
		return c.resolveTradesPollInterval
	}
	return defaultResolveTradesPollInterval
}

func uniqueTradeIDs(tradeIDs []string) []string {
	seen := make(map[string]struct{}, len(tradeIDs))
	unique := make([]string, 0, len(tradeIDs))
	for _, tradeID := range tradeIDs {
		tradeID = strings.TrimSpace(tradeID)
		if tradeID == "" {
			continue
		}
		if _, ok := seen[tradeID]; ok {
			continue
		}
		seen[tradeID] = struct{}{}
		unique = append(unique, tradeID)
	}
	return unique
}

func isResolvedTrade(trade clobtypes.Trade) bool {
	return strings.EqualFold(trade.Status, failedTradeStatus) || strings.TrimSpace(trade.TransactionHash) != ""
}
