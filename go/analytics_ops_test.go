package tachyon

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestAnalyticsTop(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{
			"queries": [{"query":"wireless mouse","collection":"products","count":3,"zero_result_count":0,"last_result_count":12,"avg_latency_ms":1.8,"last_seen":1786625778150}],
			"tracked_queries": 1,
			"dropped_queries": 0
		}`)
	})

	result, err := client.Analytics.Top(context.Background(), AnalyticsQueryParams{Collection: "products", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Queries[0].Query != "wireless mouse" {
		t.Fatalf("unexpected result: %+v", result)
	}

	req := (*requests)[0]
	if req.Path != "/analytics/top" {
		t.Fatalf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("collection") != "products" || req.Query.Get("limit") != "10" {
		t.Fatalf("unexpected query: %v", req.Query)
	}
}

func TestAnalyticsZeroResults(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"queries":[],"tracked_queries":0,"dropped_queries":0}`)
	})

	_, err := client.Analytics.ZeroResults(context.Background(), AnalyticsQueryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*requests)[0].Path != "/analytics/zero-results" {
		t.Fatalf("unexpected path: %s", (*requests)[0].Path)
	}
}

func TestAnalyticsLatency(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"count":20,"mean_ms":1.9,"p50_ms":2.0,"p95_ms":4.0,"p99_ms":4.0,"max_ms":3.4,"total_searches":20,"uptime_seconds":61,"queries_per_second":0.33}`)
	})

	result, err := client.Analytics.Latency(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.P95Ms != 4.0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestHealth(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"ok":true,"version":"0.1.0","uptime_seconds":61,"num_collections":1}`)
	})

	health, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !health.OK {
		t.Fatalf("unexpected health: %+v", health)
	}
}

func TestMetricsReturnsRawTextNotJSON(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("# HELP tachyon_uptime_seconds\ntachyon_uptime_seconds 61\n"))
	})

	metrics, err := client.Metrics(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(metrics, "tachyon_uptime_seconds 61") {
		t.Fatalf("unexpected metrics: %s", metrics)
	}
}
