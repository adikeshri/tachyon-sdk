package tachyon

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDocumentsIndex(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{
			"num_indexed": 1,
			"num_failed": 1,
			"results": [
				{"success": true, "id": "1"},
				{"success": false, "code": "invalid_document", "error": "field price: expected an integer"}
			]
		}`)
	})

	result, err := client.Collection("products").Documents.Index(context.Background(),
		Document{"id": "1", "title": "Wireless Mouse"},
		Document{"id": "2", "price": "not a number"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NumIndexed != 1 || result.NumFailed != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Results[1].Code != "invalid_document" {
		t.Fatalf("unexpected result[1]: %+v", result.Results[1])
	}

	req := (*requests)[0]
	if req.Method != http.MethodPost || req.Path != "/collections/products/documents" {
		t.Fatalf("unexpected request: %+v", req)
	}
	var body []Document
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	if len(body) != 2 || body[0]["id"] != "1" {
		t.Fatalf("unexpected request body: %+v", body)
	}
}

func TestDocumentsRetrieve(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"id":"1","title":"Wireless Mouse"}`)
	})

	doc, err := client.Collection("products").Documents.Retrieve(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc["title"] != "Wireless Mouse" {
		t.Fatalf("unexpected document: %+v", doc)
	}
	if (*requests)[0].Path != "/collections/products/documents/1" {
		t.Fatalf("unexpected path: %s", (*requests)[0].Path)
	}
}

func TestDocumentsDelete(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	if err := client.Collection("products").Documents.Delete(context.Background(), "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req := (*requests)[0]
	if req.Method != http.MethodDelete || req.Path != "/collections/products/documents/1" {
		t.Fatalf("unexpected request: %+v", req)
	}
}
