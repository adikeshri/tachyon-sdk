//go:build integration

package integration

import (
	"context"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func TestSuggestCompletesPrefixWithLiveCounts(t *testing.T) {
	result, err := suggestDataset.Suggest(context.Background(), tachyon.SuggestParams{Q: "wir"})
	if err != nil {
		t.Fatalf("suggest failed: %v", err)
	}
	var wireless, wired *tachyon.Suggestion
	for i, s := range result.Suggestions {
		if s.Text == "wireless" {
			wireless = &result.Suggestions[i]
		}
		if s.Text == "wired" {
			wired = &result.Suggestions[i]
		}
	}
	if wireless == nil || wired == nil {
		t.Fatalf("expected both 'wireless' and 'wired' suggestions, got %+v", result.Suggestions)
	}
	if wireless.Count < 2 {
		t.Fatalf("expected wireless count >= 2, got %d", wireless.Count)
	}
}

func TestSuggestCapsAtLimit(t *testing.T) {
	result, err := suggestDataset.Suggest(context.Background(), tachyon.SuggestParams{Q: "wir", Limit: 1})
	if err != nil {
		t.Fatalf("suggest failed: %v", err)
	}
	if len(result.Suggestions) > 1 {
		t.Fatalf("expected at most 1 suggestion, got %d", len(result.Suggestions))
	}
}
