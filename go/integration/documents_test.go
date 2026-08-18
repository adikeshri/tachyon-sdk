//go:build integration

package integration

import (
	"context"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func TestDocumentsIndexSingle(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})

	result, err := collection.Documents.Index(context.Background(), tachyon.Document{"id": "1", "title": "Hello World"})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}
	if result.NumIndexed != 1 || result.NumFailed != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestDocumentsRetrieveByID(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	if _, err := collection.Documents.Index(context.Background(), tachyon.Document{"id": "1", "title": "Hello World"}); err != nil {
		t.Fatalf("index failed: %v", err)
	}

	doc, err := collection.Documents.Retrieve(context.Background(), "1")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if doc["title"] != "Hello World" {
		t.Fatalf("unexpected document: %+v", doc)
	}
}

func TestDocumentsIndexUpsertsByID(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	ctx := context.Background()
	if _, err := collection.Documents.Index(ctx, tachyon.Document{"id": "1", "title": "First title"}); err != nil {
		t.Fatalf("first index failed: %v", err)
	}
	if _, err := collection.Documents.Index(ctx, tachyon.Document{"id": "1", "title": "Second title"}); err != nil {
		t.Fatalf("second index failed: %v", err)
	}

	doc, err := collection.Documents.Retrieve(ctx, "1")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if doc["title"] != "Second title" {
		t.Fatalf("expected upsert to overwrite title, got: %+v", doc)
	}
}

func TestDocumentsDeleteThenRetrieve404s(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	ctx := context.Background()
	if _, err := collection.Documents.Index(ctx, tachyon.Document{"id": "1", "title": "Hello World"}); err != nil {
		t.Fatalf("index failed: %v", err)
	}
	if err := collection.Documents.Delete(ctx, "1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := collection.Documents.Retrieve(ctx, "1"); !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected not-found after delete, got %v", err)
	}
}

func TestDocumentsRetrieveAndDeleteUnknown404s(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	ctx := context.Background()
	if _, err := collection.Documents.Retrieve(ctx, "never-existed"); !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected not-found on retrieve, got %v", err)
	}
	if err := collection.Documents.Delete(ctx, "never-existed"); !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected not-found on delete, got %v", err)
	}
}

func TestDocumentsIndexBatchReportsPerDocumentResults(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{
			{Name: "title", Type: tachyon.FieldTypeText},
			{Name: "price", Type: tachyon.FieldTypeInt},
		},
	})
	ctx := context.Background()

	result, err := collection.Documents.Index(ctx,
		tachyon.Document{"id": "1", "price": 100},
		tachyon.Document{"id": "2", "price": "not-a-number"},
		tachyon.Document{"id": "3", "price": 300},
	)
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}
	if result.NumIndexed != 2 || result.NumFailed != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if !result.Results[0].Success || result.Results[0].ID != "1" {
		t.Fatalf("unexpected results[0]: %+v", result.Results[0])
	}
	if result.Results[1].Success || result.Results[1].Code != "invalid_document" {
		t.Fatalf("unexpected results[1]: %+v", result.Results[1])
	}
	if !result.Results[2].Success || result.Results[2].ID != "3" {
		t.Fatalf("unexpected results[2]: %+v", result.Results[2])
	}

	if _, err := collection.Documents.Retrieve(ctx, "2"); !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected the failed document to not exist, got %v", err)
	}
}

func TestDocumentsUndeclaredFieldsStoredNotIndexed(t *testing.T) {
	client := adminClient()
	collection := withCollection(t, client, tachyon.CollectionSchema{
		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	})
	ctx := context.Background()
	if _, err := collection.Documents.Index(ctx, tachyon.Document{"id": "1", "title": "Hello", "undeclared_field": "surprise"}); err != nil {
		t.Fatalf("index failed: %v", err)
	}

	doc, err := collection.Documents.Retrieve(ctx, "1")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if doc["undeclared_field"] != "surprise" {
		t.Fatalf("expected undeclared field to be stored, got: %+v", doc)
	}

	results, err := collection.Search(ctx, tachyon.SearchParams{Q: "surprise"})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results.Hits) != 0 {
		t.Fatalf("expected undeclared field to not be searchable, got %d hits", len(results.Hits))
	}
}
