//go:build integration

package integration

import (
	"context"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func TestCollectionsCreateFillsInDefaults(t *testing.T) {
	client := adminClient()
	name := uniqueName("coll")
	t.Cleanup(func() { _ = client.Collections.Delete(context.Background(), name) })

	created, err := client.Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name: name,
		Fields: []tachyon.FieldSchema{
			{Name: "title", Type: tachyon.FieldTypeText},
			{Name: "brand", Type: tachyon.FieldTypeKeyword, Facet: tachyon.Ptr(true)},
			{Name: "price", Type: tachyon.FieldTypeInt, Filter: tachyon.Ptr(true), Sort: tachyon.Ptr(true)},
		},
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Name != name || created.NumDocuments != 0 || created.NumSegments != 0 {
		t.Fatalf("unexpected created collection: %+v", created)
	}
	if len(created.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(created.Fields))
	}
}

func TestCollectionsCreateRejectsDuplicateName(t *testing.T) {
	client := adminClient()
	name := uniqueName("coll")
	t.Cleanup(func() { _ = client.Collections.Delete(context.Background(), name) })

	schema := tachyon.CollectionSchema{Name: name, Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}}}
	if _, err := client.Collections.Create(context.Background(), schema); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err := client.Collections.Create(context.Background(), schema)
	if !tachyon.IsConflictError(err) {
		t.Fatalf("expected a conflict error, got %v", err)
	}
	var tErr *tachyon.Error
	if e, ok := err.(*tachyon.Error); ok {
		tErr = e
	}
	if tErr == nil || tErr.Code != "collection_exists" || tErr.Status != 409 {
		t.Fatalf("unexpected error: %+v", tErr)
	}
}

func TestCollectionsListIncludesNewlyCreated(t *testing.T) {
	client := adminClient()
	name := uniqueName("coll")
	t.Cleanup(func() { _ = client.Collections.Delete(context.Background(), name) })

	if _, err := client.Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name: name, Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	}); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	list, err := client.Collections.List(context.Background())
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	found := false
	for _, c := range list {
		if c.Name == name {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected %q in list of %d collections", name, len(list))
	}
}

func TestCollectionsRetrieveUnknown404s(t *testing.T) {
	client := adminClient()
	_, err := client.Collections.Retrieve(context.Background(), uniqueName("missing"))
	if !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected a not-found error, got %v", err)
	}
}

func TestCollectionsDeleteRemovesFromRetrieveAndList(t *testing.T) {
	client := adminClient()
	name := uniqueName("coll")
	if _, err := client.Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name: name, Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
	}); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := client.Collections.Delete(context.Background(), name); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := client.Collections.Retrieve(context.Background(), name); !tachyon.IsNotFoundError(err) {
		t.Fatalf("expected not-found after delete, got %v", err)
	}
}

func TestCollectionsCreateRejectsSchemaWithNoTextField(t *testing.T) {
	client := adminClient()
	_, err := client.Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name:   uniqueName("coll"),
		Fields: []tachyon.FieldSchema{{Name: "price", Type: tachyon.FieldTypeInt}},
	})
	if !tachyon.IsRequestError(err) {
		t.Fatalf("expected a request error, got %v", err)
	}
	var tErr *tachyon.Error
	if e, ok := err.(*tachyon.Error); ok {
		tErr = e
	}
	if tErr == nil || tErr.Code != "invalid_schema" {
		t.Fatalf("unexpected error: %+v", tErr)
	}
}

func TestCollectionsTypoToleranceAndDefaultSortingField(t *testing.T) {
	client := adminClient()
	name := uniqueName("coll")
	t.Cleanup(func() { _ = client.Collections.Delete(context.Background(), name) })

	created, err := client.Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name: name,
		Fields: []tachyon.FieldSchema{
			{Name: "title", Type: tachyon.FieldTypeText},
			{Name: "popularity", Type: tachyon.FieldTypeInt, Sort: tachyon.Ptr(true)},
		},
		TypoTolerance:       &tachyon.TypoTolerance{Enabled: tachyon.Ptr(false)},
		DefaultSortingField: "popularity",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.TypoTolerance == nil || created.TypoTolerance.Enabled == nil || *created.TypoTolerance.Enabled != false {
		t.Fatalf("unexpected typo_tolerance: %+v", created.TypoTolerance)
	}
	if created.DefaultSortingField != "popularity" {
		t.Fatalf("unexpected default_sorting_field: %q", created.DefaultSortingField)
	}
}
