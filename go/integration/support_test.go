//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func baseURL() string {
	if v := os.Getenv("TACHYON_URL"); v != "" {
		return v
	}
	return "http://localhost:8108"
}

func adminKey() string {
	if v := os.Getenv("TACHYON_ADMIN_KEY"); v != "" {
		return v
	}
	return "admin-key"
}

func searchKey() string {
	if v := os.Getenv("TACHYON_SEARCH_KEY"); v != "" {
		return v
	}
	return "search-key"
}

func adminClient() *tachyon.Client {
	return tachyon.NewClient(baseURL(), tachyon.WithAPIKey(adminKey()))
}

func searchOnlyClient() *tachyon.Client {
	return tachyon.NewClient(baseURL(), tachyon.WithAPIKey(searchKey()))
}

func anonymousClient() *tachyon.Client {
	return tachyon.NewClient(baseURL())
}

func uniqueName(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf))
}

// withCollection creates a collection for the duration of the test (deleted
// via t.Cleanup, which runs even if the test fails), and returns a handle
// scoped to it.
func withCollection(t *testing.T, client *tachyon.Client, schema tachyon.CollectionSchema) *tachyon.Collection {
	t.Helper()
	if schema.Name == "" {
		schema.Name = uniqueName("coll")
	}
	if _, err := client.Collections.Create(context.Background(), schema); err != nil {
		t.Fatalf("failed to create collection %q: %v", schema.Name, err)
	}
	t.Cleanup(func() {
		_ = client.Collections.Delete(context.Background(), schema.Name)
	})
	return client.Collection(schema.Name)
}
