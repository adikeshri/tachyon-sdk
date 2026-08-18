package tachyon

import (
	"context"
	"net/http"
	"testing"
)

func TestSearchSerializesParamsIntoTheExpectedQueryString(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"found":1,"found_is_exact":true,"search_time_ms":1,"hits":[]}`)
	})

	_, err := client.Collection("products").Search(context.Background(), SearchParams{
		Q:             "wireless mouse",
		QueryBy:       []string{"title", "description"},
		Filter:        "brand:=Logitech && price:<5000",
		Sort:          "_text_match:desc,price:asc",
		Facet:         []string{"brand", "year"},
		Limit:         20,
		Offset:        40,
		Prefix:        Ptr(false),
		TypoTolerance: Ptr(true),
		MatchMode:     MatchModeAny,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := (*requests)[0]
	if req.Path != "/collections/products/search" {
		t.Fatalf("unexpected path: %s", req.Path)
	}
	want := map[string]string{
		"q":              "wireless mouse",
		"query_by":       "title,description",
		"filter":         "brand:=Logitech && price:<5000",
		"sort":           "_text_match:desc,price:asc",
		"facet":          "brand,year",
		"limit":          "20",
		"offset":         "40",
		"prefix":         "false",
		"typo_tolerance": "true",
		"match_mode":     "any",
	}
	for key, expected := range want {
		if got := req.Query.Get(key); got != expected {
			t.Errorf("query[%q] = %q, want %q", key, got, expected)
		}
	}
}

func TestSearchOmitsUnsetParams(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"found":0,"found_is_exact":true,"search_time_ms":0,"hits":[]}`)
	})

	_, err := client.Collection("products").Search(context.Background(), SearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := (*requests)[0]
	if req.Query.Has("filter") || req.Query.Has("limit") || req.Query.Has("prefix") {
		t.Fatalf("expected unset params to be omitted, got: %v", req.Query)
	}
}

func TestSearchReturnsHitsFacetsAndFoundIsExact(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{
			"found": 1240,
			"found_is_exact": false,
			"search_time_ms": 12,
			"hits": [{"document": {"id": "1", "title": "Wireless Mouse"}, "text_match": 554.788}],
			"facets": {"brand": {"Logitech": 1240, "Razer": 830}}
		}`)
	})

	results, err := client.Collection("products").Search(context.Background(), SearchParams{Q: "wireless mouse"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results.Found != 1240 || results.FoundIsExact != false {
		t.Fatalf("unexpected results: %+v", results)
	}
	if results.Hits[0].Document["title"] != "Wireless Mouse" {
		t.Fatalf("unexpected hit: %+v", results.Hits[0])
	}
	if results.Facets["brand"]["Logitech"] != 1240 {
		t.Fatalf("unexpected facets: %+v", results.Facets)
	}
}

func TestSuggest(t *testing.T) {
	client, requests := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{
			"suggestions": [
				{"text": "wireless", "count": 3, "typos": 0},
				{"text": "wired", "count": 2, "typos": 0}
			],
			"search_time_ms": 0
		}`)
	})

	result, err := client.Collection("products").Suggest(context.Background(), SuggestParams{Q: "wir", Limit: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Suggestions) != 2 {
		t.Fatalf("unexpected suggestions: %+v", result.Suggestions)
	}

	req := (*requests)[0]
	if req.Path != "/collections/products/suggest" {
		t.Fatalf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("q") != "wir" || req.Query.Get("limit") != "5" {
		t.Fatalf("unexpected query: %v", req.Query)
	}
}
