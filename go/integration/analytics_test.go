//go:build integration

package integration

import (
	"context"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func TestAnalyticsRecordsCountsScopedByCollection(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	ctx := context.Background()
	if _, err := collection.Documents.Index(ctx, tachyon.Document{"id": "1", "title": "gizmo"}); err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if _, err := collection.Search(ctx, tachyon.SearchParams{Q: "gizmo"}); err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if _, err := collection.Search(ctx, tachyon.SearchParams{Q: "gizmo"}); err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if _, err := collection.Search(ctx, tachyon.SearchParams{Q: "zzz-totally-absent-term"}); err != nil {
		t.Fatalf("search failed: %v", err)
	}

	top, err := client.Analytics.Top(ctx, tachyon.AnalyticsQueryParams{Collection: collection.Name()})
	if err != nil {
		t.Fatalf("analytics top failed: %v", err)
	}
	var gizmo *tachyon.AnalyticsQuery
	for i, q := range top.Queries {
		if q.Query == "gizmo" {
			gizmo = &top.Queries[i]
		}
	}
	if gizmo == nil {
		t.Fatalf("expected a 'gizmo' query in top: %+v", top.Queries)
	}
	if gizmo.Count != 2 || gizmo.Collection != collection.Name() {
		t.Fatalf("unexpected gizmo query: %+v", gizmo)
	}

	zeroResults, err := client.Analytics.ZeroResults(ctx, tachyon.AnalyticsQueryParams{Collection: collection.Name()})
	if err != nil {
		t.Fatalf("analytics zero-results failed: %v", err)
	}
	found := false
	for _, q := range zeroResults.Queries {
		if q.Query == "zzz-totally-absent-term" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected zzz-totally-absent-term in zero-results: %+v", zeroResults.Queries)
	}
}

func TestAnalyticsTopRespectsLimit(t *testing.T) {
	client := adminClient()
	top, err := client.Analytics.Top(context.Background(), tachyon.AnalyticsQueryParams{Limit: 1})
	if err != nil {
		t.Fatalf("analytics top failed: %v", err)
	}
	if len(top.Queries) > 1 {
		t.Fatalf("expected at most 1 query, got %d", len(top.Queries))
	}
}

func TestAnalyticsLatencyReportsPercentiles(t *testing.T) {
	client := adminClient()
	latency, err := client.Analytics.Latency(context.Background())
	if err != nil {
		t.Fatalf("analytics latency failed: %v", err)
	}
	if latency.TotalSearches <= 0 {
		t.Fatalf("expected total_searches > 0, got %d", latency.TotalSearches)
	}
	if latency.P95Ms < latency.P50Ms || latency.P99Ms < latency.P95Ms {
		t.Fatalf("expected percentiles to be non-decreasing: p50=%v p95=%v p99=%v", latency.P50Ms, latency.P95Ms, latency.P99Ms)
	}
}
