package dataapi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var conditionPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)

func addParam[T any](params map[string]any, name string, value *T) {
	if value != nil {
		params[name] = *value
	}
}

func addPageParams(params map[string]any, limit *int, cursor *string) error {
	if cursor != nil && *cursor != "" {
		params["cursor"] = *cursor
		return nil // The cursor's page size wins.
	}
	if limit != nil && (*limit < 0 || *limit > 1000) {
		return fmt.Errorf("limit out of range (0..1000)")
	}
	addParam(params, "limit", limit)
	return nil
}

func distinctCount[T comparable](values []T) int {
	seen := make(map[T]struct{}, len(values))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	return len(seen)
}

func validateConditions(values []string) error {
	if distinctCount(values) > 20 {
		return fmt.Errorf("condition supports at most 20 distinct IDs")
	}
	for _, value := range values {
		if !conditionPattern.MatchString(value) {
			return fmt.Errorf("invalid condition ID: %s (must be 0x + 64 hex chars)", value)
		}
	}
	return nil
}

func validateEventIDs(values []int64) error {
	if distinctCount(values) > 20 {
		return fmt.Errorf("event_id supports at most 20 distinct IDs")
	}
	return nil
}

func joinEventIDs(values []int64) string {
	ids := make([]string, len(values))
	for i, value := range values {
		ids[i] = strconv.FormatInt(value, 10)
	}
	return strings.Join(ids, ",")
}
