//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
)

func TestHealthDoesNotRequireAPIKey(t *testing.T) {
	health, err := anonymousClient().Health(context.Background())
	if err != nil {
		t.Fatalf("health failed: %v", err)
	}
	if !health.OK {
		t.Fatalf("unexpected health: %+v", health)
	}
	if health.NumCollections < 0 {
		t.Fatalf("unexpected num_collections: %d", health.NumCollections)
	}
}

func TestMetricsExposesPrometheusText(t *testing.T) {
	metrics, err := adminClient().Metrics(context.Background())
	if err != nil {
		t.Fatalf("metrics failed: %v", err)
	}
	if !strings.Contains(metrics, "tachyon_uptime_seconds") {
		t.Fatalf("expected tachyon_uptime_seconds in metrics output")
	}
	if !strings.Contains(metrics, "tachyon_collections") {
		t.Fatalf("expected tachyon_collections in metrics output")
	}
}
