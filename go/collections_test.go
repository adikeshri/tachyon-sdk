package tachyon

import (
	"context"
	"net/http"
	"testing"
)

func TestCollectionsCreate(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 201, `{"name":"products","fields":[{"name":"title","type":"text"}],"num_documents":0,"num_segments":0}`)
	})

	info, err := client.Collections.Create(context.Background(), CollectionSchema{
		Name:   "products",
		Fields: []FieldSchema{{Name: "title", Type: FieldTypeText}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "products" || info.NumDocuments != 0 {
		t.Fatalf("unexpected info: %+v", info)
	}

	if len(*requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(*requests))
	}
	req := (*requests)[0]
	if req.Method != http.MethodPost || req.Path != "/collections" {
		t.Fatalf("unexpected request: %+v", req)
	}
	if req.Header.Get("X-TACHYON-API-KEY") != "admin-key" {
		t.Fatalf("missing api key header: %+v", req.Header)
	}
}

func TestCollectionsList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `[{"name":"products","fields":[],"num_documents":5,"num_segments":1}]`)
	})

	list, err := client.Collections.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].Name != "products" || list[0].NumDocuments != 5 {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestCollectionsRetrieve(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"name":"products","fields":[],"num_documents":5,"num_segments":1}`)
	})

	info, err := client.Collections.Retrieve(context.Background(), "products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.NumDocuments != 5 {
		t.Fatalf("unexpected info: %+v", info)
	}
	if (*requests)[0].Path != "/collections/products" {
		t.Fatalf("unexpected path: %s", (*requests)[0].Path)
	}
}

func TestCollectionsRetrieveURLEncodesName(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"name":"my products","fields":[],"num_documents":0,"num_segments":0}`)
	})

	_, err := client.Collections.Retrieve(context.Background(), "my products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*requests)[0].Path != "/collections/my%20products" {
		t.Fatalf("unexpected path: %s", (*requests)[0].Path)
	}
}

func TestCollectionsDelete(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	if err := client.Collections.Delete(context.Background(), "products"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*requests)[0].Method != http.MethodDelete || (*requests)[0].Path != "/collections/products" {
		t.Fatalf("unexpected request: %+v", (*requests)[0])
	}
}
