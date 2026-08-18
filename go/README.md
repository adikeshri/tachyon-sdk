# tachyon-sdk (Go)

Official Go client for [Tachyon](https://github.com/adikeshri/tachyon), the
typo-tolerant full-text search engine.

```bash
go get github.com/adikeshri/tachyon-sdk/go
```

Zero non-stdlib dependencies. Requires Go 1.21+ (for generics-based `Ptr`).

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func main() {
	ctx := context.Background()
	client := tachyon.NewClient("http://localhost:8108", tachyon.WithAPIKey("my-admin-key"))

	_, err := client.Collections.Create(ctx, tachyon.CollectionSchema{
		Name: "products",
		Fields: []tachyon.FieldSchema{
			{Name: "title", Type: tachyon.FieldTypeText},
			{Name: "brand", Type: tachyon.FieldTypeKeyword, Facet: tachyon.Ptr(true)},
			{Name: "price", Type: tachyon.FieldTypeInt, Filter: tachyon.Ptr(true), Sort: tachyon.Ptr(true)},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Collection("products")
	_, err = collection.Documents.Index(ctx,
		tachyon.Document{"id": "1", "title": "Wireless Mouse", "brand": "Logitech", "price": 2999},
		tachyon.Document{"id": "2", "title": "Mechanical Keyboard", "brand": "Razer", "price": 8999},
	)
	if err != nil {
		log.Fatal(err)
	}

	results, err := collection.Search(ctx, tachyon.SearchParams{Q: "wireless mouse"})
	if err != nil {
		log.Fatal(err)
	}
	for _, hit := range results.Hits {
		fmt.Println(hit.Document["title"], hit.TextMatch)
	}
}
```

## Client options

```go
client := tachyon.NewClient(
	"http://localhost:8108",
	tachyon.WithAPIKey("..."),                  // admin key (read/write) or search key (read-only)
	tachyon.WithTimeout(15*time.Second),        // default 15s
	tachyon.WithHeader("X-Custom", "value"),    // repeatable
	tachyon.WithHTTPClient(customHTTPClient),   // override transport (testing, custom pooling)
)
```

Every method takes a `context.Context` as its first argument, in keeping
with Go convention — cancel it or set a deadline to bound a single call
independently of the client-wide timeout.

## Collections

```go
client.Collections.Create(ctx, schema)  // POST /collections
client.Collections.List(ctx)            // GET /collections
client.Collections.Retrieve(ctx, name)  // GET /collections/{name}
client.Collections.Delete(ctx, name)    // DELETE /collections/{name}
```

`FieldSchema`'s optional attributes (`Facet`, `Filter`, `Sort`, `Index`,
`Optional`, `Boost`) are pointers, so the zero value doesn't accidentally
override the server's own default — set them with `tachyon.Ptr(true)`.

## Documents

```go
collection := client.Collection("products")

collection.Documents.Index(ctx, doc1, doc2, doc3)  // POST   /collections/{name}/documents (upsert by id, variadic)
collection.Documents.Retrieve(ctx, id)             // GET    /collections/{name}/documents/{id}
collection.Documents.Delete(ctx, id)               // DELETE /collections/{name}/documents/{id}
```

`Index` always succeeds at the HTTP level even if individual documents are
rejected — check `NumFailed` and `Results` on the response.

## Search

```go
collection.Search(ctx, tachyon.SearchParams{
	Q:             "wireless mouse",
	QueryBy:       []string{"title", "description"},
	Filter:        "brand:=Logitech && price:<5000",
	Sort:          "_text_match:desc,price:asc",
	Facet:         []string{"brand", "year"},
	Limit:         20,
	Offset:        0,
	Prefix:        tachyon.Ptr(true),
	TypoTolerance: tachyon.Ptr(true),
	MatchMode:     tachyon.MatchModeAll, // or MatchModeAny
})
```

`Prefix` and `TypoTolerance` are `*bool` (their zero value `false` would
otherwise be indistinguishable from "not set, use the default"); every
other field's Go zero value naturally means "omit, use the server's
default."

`FoundIsExact` on the response is `false` once block-max WAND pruning has
skipped part of a term's postings for a broad query — at that point `Found`
and facet counts are a lower bound, not an exact count. See the [Known
limitations section of the main
README](https://github.com/adikeshri/tachyon#known-limitations).

## Autocomplete

```go
collection.Suggest(ctx, tachyon.SuggestParams{Q: "wir", Limit: 5})
```

## Analytics

```go
client.Analytics.Top(ctx, tachyon.AnalyticsQueryParams{Collection: "products", Limit: 10})
client.Analytics.ZeroResults(ctx, tachyon.AnalyticsQueryParams{Collection: "products"})
client.Analytics.Latency(ctx)
```

Analytics are in-memory only and reset when the server restarts.

## Operations

```go
client.Health(ctx)   // GET /health — no API key required
client.Metrics(ctx)  // GET /metrics — Prometheus exposition format, returned as text
```

## Errors

Every non-2xx response returns an `*tachyon.Error` carrying the server's
stable `Code` and the HTTP `Status`, categorized by `Kind`. Use the
`Is*Error` predicates rather than comparing `Kind` directly:

```go
info, err := client.Collections.Retrieve(ctx, "does-not-exist")
if tachyon.IsNotFoundError(err) {
	// err.(*tachyon.Error).Code == "collection_not_found"
}
```

| Predicate | Status | Codes |
|---|---|---|
| `IsRequestError` | 400 | `invalid_schema`, `invalid_document`, `invalid_query`, `invalid_json` |
| `IsAuthenticationError` | 401 | `unauthorized` |
| `IsAuthorizationError` | 403 | `forbidden` |
| `IsNotFoundError` | 404 | `collection_not_found`, `document_not_found` |
| `IsConflictError` | 409 | `collection_exists` |
| `IsServerError` | 5xx | `corrupt_data`, `io_error`, `internal_error` |

Network failures and timeouts return a `*tachyon.ConnectionError` instead,
since there's no server response to read a code from — check with
`IsConnectionError` / `IsTimeoutError` (a timeout is always also a
connection error).

## Development

```bash
go build ./...
go vet ./...
gofmt -l .   # should print nothing
go test ./...
```

### Integration tests

`go test ./...` only runs the mocked unit suite (the `integration` package
is gated behind a build tag so it's not even compiled by default). A
second suite in `integration/` exercises every functional path —
collections, documents, search (filters, sort, facets, pagination, prefix,
typo tolerance, match mode, phrases), suggest, analytics, auth, and error
paths — against a real, running Tachyon server:

```bash
docker run -d -p 8108:8108 \
  -e TACHYON_ADMIN_KEY=admin-key -e TACHYON_SEARCH_KEY=search-key \
  adikeshri/tachyon

go test -tags=integration ./integration/...
```

It points at `http://localhost:8108` by default; override with
`TACHYON_URL`, `TACHYON_ADMIN_KEY`, `TACHYON_SEARCH_KEY`. Every test cleans
up the collections it creates, even on failure (via `t.Cleanup`).

## License

Apache 2.0. See [`LICENSE`](../LICENSE).
