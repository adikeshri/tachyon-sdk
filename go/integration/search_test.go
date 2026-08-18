//go:build integration

package integration

import (
	"context"
	"sort"
	"testing"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func hitIDs(results tachyon.SearchResponse) []string {
	ids := make([]string, len(results.Hits))
	for i, h := range results.Hits {
		ids[i] = h.Document["id"].(string)
	}
	sort.Strings(ids)
	return ids
}

func TestSearchMatchesTitleAndDescriptionByDefault(t *testing.T) {
	results, err := searchDataset.Search(context.Background(), tachyon.SearchParams{Q: "wireless"})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if results.Found != 2 || !results.FoundIsExact {
		t.Fatalf("unexpected results: found=%d exact=%v", results.Found, results.FoundIsExact)
	}
	ids := hitIDs(results)
	if ids[0] != "1" || ids[1] != "3" {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

func TestSearchQueryByRestrictsFields(t *testing.T) {
	ctx := context.Background()
	onlyDescription, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "Clicky", QueryBy: []string{"description"}})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if onlyDescription.Found != 1 || onlyDescription.Hits[0].Document["id"] != "2" {
		t.Fatalf("unexpected results: %+v", onlyDescription)
	}

	onlyTitle, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "Clicky", QueryBy: []string{"title"}})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if onlyTitle.Found != 0 {
		t.Fatalf("expected 0 results, got %d", onlyTitle.Found)
	}
}

func TestSearchFilters(t *testing.T) {
	ctx := context.Background()

	t.Run("equality", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "brand:=Logitech"})
		if err != nil || r.Found != 2 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})

	t.Run("numeric comparison", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "price:<5000"})
		if err != nil || r.Found != 3 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})

	t.Run("inclusive range", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "price:[1000..5000]"})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		ids := hitIDs(r)
		if len(ids) != 2 || ids[0] != "1" || ids[1] != "4" {
			t.Fatalf("unexpected ids: %v", ids)
		}
	})

	t.Run("set membership", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "brand:=[Logitech,Razer]"})
		if err != nil || r.Found != 4 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})

	t.Run("boolean equality", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "in_stock:=true"})
		if err != nil || r.Found != 4 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})

	t.Run("and or with grouping", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "(brand:=Logitech || brand:=Razer) && price:<5000"})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		ids := hitIDs(r)
		if len(ids) != 2 || ids[0] != "1" || ids[1] != "4" {
			t.Fatalf("unexpected ids: %v", ids)
		}
	})

	t.Run("negation only matches documents that have the field", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Filter: "brand:!=Razer"})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		ids := hitIDs(r)
		if len(ids) != 3 || ids[0] != "1" || ids[1] != "3" || ids[2] != "5" {
			t.Fatalf("unexpected ids: %v", ids)
		}
	})

	t.Run("negation excludes documents missing the field entirely", func(t *testing.T) {
		client := adminClient()
		collection := withCollection(t, client, tachyon.CollectionSchema{
			Fields: []tachyon.FieldSchema{
				{Name: "title", Type: tachyon.FieldTypeText},
				{Name: "brand", Type: tachyon.FieldTypeKeyword, Filter: tachyon.Ptr(true)},
			},
		})
		if _, err := collection.Documents.Index(ctx,
			tachyon.Document{"id": "a", "title": "has brand", "brand": "Razer"},
			tachyon.Document{"id": "b", "title": "no brand at all"},
		); err != nil {
			t.Fatalf("index failed: %v", err)
		}
		r, err := collection.Search(ctx, tachyon.SearchParams{Filter: "brand:!=Razer"})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if r.Found != 0 {
			t.Fatalf("expected 0 results (doc with no brand must not match), got %d", r.Found)
		}
	})
}

func TestSearchSorting(t *testing.T) {
	ctx := context.Background()

	t.Run("ascending", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Sort: "price:asc", Limit: 10})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		for i := 1; i < len(r.Hits); i++ {
			prev := r.Hits[i-1].Document["price"].(float64)
			cur := r.Hits[i].Document["price"].(float64)
			if prev > cur {
				t.Fatalf("not ascending at index %d: %v", i, r.Hits)
			}
		}
	})

	t.Run("descending", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Sort: "price:desc", Limit: 10})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		for i := 1; i < len(r.Hits); i++ {
			prev := r.Hits[i-1].Document["price"].(float64)
			cur := r.Hits[i].Document["price"].(float64)
			if prev < cur {
				t.Fatalf("not descending at index %d: %v", i, r.Hits)
			}
		}
	})
}

func TestSearchPagination(t *testing.T) {
	ctx := context.Background()
	page1, err := searchDataset.Search(ctx, tachyon.SearchParams{Sort: "price:asc", Limit: 2, Offset: 0})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	page2, err := searchDataset.Search(ctx, tachyon.SearchParams{Sort: "price:asc", Limit: 2, Offset: 2})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if page1.Found != 5 || page2.Found != 5 {
		t.Fatalf("unexpected found counts: %d, %d", page1.Found, page2.Found)
	}
	if len(page1.Hits) != 2 || len(page2.Hits) != 2 {
		t.Fatalf("unexpected page sizes: %d, %d", len(page1.Hits), len(page2.Hits))
	}
	p1 := hitIDs(page1)
	p2 := hitIDs(page2)
	for _, id := range p1 {
		for _, other := range p2 {
			if id == other {
				t.Fatalf("page overlap: %s appears in both pages", id)
			}
		}
	}
}

func TestSearchPrefixMatching(t *testing.T) {
	ctx := context.Background()
	t.Run("prefix-expands by default", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wir"})
		if err != nil || r.Found < 2 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
	t.Run("requires full token when disabled", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wir", Prefix: tachyon.Ptr(false)})
		if err != nil || r.Found != 0 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
}

func TestSearchTypoTolerance(t *testing.T) {
	ctx := context.Background()
	t.Run("corrects a typo by default", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wirelss"})
		if err != nil || r.Found < 1 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
	t.Run("finds nothing when explicitly disabled", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wirelss", TypoTolerance: tachyon.Ptr(false)})
		if err != nil || r.Found != 0 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
}

func TestSearchMatchMode(t *testing.T) {
	ctx := context.Background()
	t.Run("all requires every token", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wireless zzznonexistentterm"})
		if err != nil || r.Found != 0 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
	t.Run("any requires one token", func(t *testing.T) {
		r, err := searchDataset.Search(ctx, tachyon.SearchParams{Q: "wireless zzznonexistentterm", MatchMode: tachyon.MatchModeAny})
		if err != nil || r.Found < 2 {
			t.Fatalf("found=%d err=%v", r.Found, err)
		}
	})
}

func TestSearchPhraseQueriesRequireAdjacency(t *testing.T) {
	r, err := searchDataset.Search(context.Background(), tachyon.SearchParams{Q: `"wireless mouse"`})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if r.Found != 1 || r.Hits[0].Document["id"] != "1" {
		t.Fatalf("unexpected results: %+v", r)
	}
}

func TestSearchFacetsCountEveryMatch(t *testing.T) {
	r, err := searchDataset.Search(context.Background(), tachyon.SearchParams{Facet: []string{"brand"}, Limit: 1})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(r.Hits) != 1 {
		t.Fatalf("expected page size 1, got %d", len(r.Hits))
	}
	want := map[string]int{"Logitech": 2, "Razer": 2, "Anker": 1}
	got := r.Facets["brand"]
	if len(got) != len(want) {
		t.Fatalf("unexpected facets: %+v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("facets[%q] = %d, want %d (all: %+v)", k, got[k], v, got)
		}
	}
}
