//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

// The search and suggest datasets are shared read-only fixtures across every
// test in this package, set up once here (TestMain is the standard Go
// mechanism for package-level setup/teardown — the equivalent of a
// beforeAll/afterAll) rather than per test.
var (
	searchDatasetName string
	searchDataset     *tachyon.Collection

	suggestDatasetName string
	suggestDataset     *tachyon.Collection
)

func TestMain(m *testing.M) {
	client := adminClient()
	ctx := context.Background()

	searchDatasetName = uniqueName("search")
	if _, err := client.Collections.Create(ctx, tachyon.CollectionSchema{
		Name: searchDatasetName,
		Fields: []tachyon.FieldSchema{
			{Name: "title", Type: tachyon.FieldTypeText},
			{Name: "description", Type: tachyon.FieldTypeText},
			{Name: "brand", Type: tachyon.FieldTypeKeyword, Facet: tachyon.Ptr(true)},
			{Name: "price", Type: tachyon.FieldTypeInt, Filter: tachyon.Ptr(true), Sort: tachyon.Ptr(true)},
			{Name: "in_stock", Type: tachyon.FieldTypeBool, Filter: tachyon.Ptr(true)},
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set up search dataset: %v\n", err)
		os.Exit(1)
	}
	searchDataset = client.Collection(searchDatasetName)
	if _, err := searchDataset.Documents.Index(ctx,
		tachyon.Document{"id": "1", "title": "Wireless Mouse", "description": "A great wireless mouse for everyday use", "brand": "Logitech", "price": 2999, "in_stock": true},
		tachyon.Document{"id": "2", "title": "Mechanical Keyboard", "description": "Clicky keys for typing enthusiasts", "brand": "Razer", "price": 8999, "in_stock": false},
		tachyon.Document{"id": "3", "title": "Wireless Keyboard", "description": "Silent wireless keyboard", "brand": "Logitech", "price": 5999, "in_stock": true},
		tachyon.Document{"id": "4", "title": "Gaming Mouse", "description": "Wired precision gaming mouse", "brand": "Razer", "price": 4999, "in_stock": true},
		tachyon.Document{"id": "5", "title": "USB Cable", "description": "A basic wired cable", "brand": "Anker", "price": 999, "in_stock": true},
	); err != nil {
		fmt.Fprintf(os.Stderr, "failed to index search dataset: %v\n", err)
		os.Exit(1)
	}

	suggestDatasetName = uniqueName("suggest")
	if _, err := client.Collections.Create(ctx, tachyon.CollectionSchema{
		Name:   suggestDatasetName,
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set up suggest dataset: %v\n", err)
		os.Exit(1)
	}
	suggestDataset = client.Collection(suggestDatasetName)
	if _, err := suggestDataset.Documents.Index(ctx,
		tachyon.Document{"id": "1", "title": "wireless mouse"},
		tachyon.Document{"id": "2", "title": "wireless keyboard"},
		tachyon.Document{"id": "3", "title": "wireless mouse"},
		tachyon.Document{"id": "4", "title": "wired cable"},
	); err != nil {
		fmt.Fprintf(os.Stderr, "failed to index suggest dataset: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	_ = client.Collections.Delete(ctx, searchDatasetName)
	_ = client.Collections.Delete(ctx, suggestDatasetName)

	os.Exit(code)
}
