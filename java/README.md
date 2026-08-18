# tachyon-sdk (Java)

Official Java client for [Tachyon](https://github.com/adikeshri/tachyon), the
typo-tolerant full-text search engine.

```xml
<dependency>
  <groupId>io.github.adikeshri</groupId>
  <artifactId>tachyon-sdk</artifactId>
  <version>1.0.0</version>
</dependency>
```

Requires Java 17+. One dependency: `jackson-databind` (for JSON) — the HTTP
transport is `java.net.http.HttpClient`, built into the JDK since 11.

## Quickstart

```java
import io.github.adikeshri.tachyon.*;

var client = TachyonClient.builder("http://localhost:8108")
    .apiKey("my-admin-key")
    .build();

client.collections.create(CollectionSchema.builder("products")
    .fields(
        FieldSchema.builder("title", FieldType.TEXT).build(),
        FieldSchema.builder("brand", FieldType.KEYWORD).facet(true).build(),
        FieldSchema.builder("price", FieldType.INT).filter(true).sort(true).build())
    .build());

var collection = client.collection("products");
collection.documents().index(
    new Document().set("id", "1").set("title", "Wireless Mouse").set("brand", "Logitech").set("price", 2999),
    new Document().set("id", "2").set("title", "Mechanical Keyboard").set("brand", "Razer").set("price", 8999));

var results = collection.search(SearchParams.builder().q("wireless mouse").build());
for (var hit : results.hits()) {
    System.out.println(hit.document().get("title") + " " + hit.textMatch());
}
```

`Document` is a `LinkedHashMap<String, Object>` with a fluent `.set()` for
building one inline; read it back like any `Map`.

## Client options

```java
var client = TachyonClient.builder("http://localhost:8108")
    .apiKey("...")                              // admin key (read/write) or search key (read-only)
    .timeout(Duration.ofSeconds(15))             // default 15s; ignored if httpClient(...) is supplied
    .headers(Map.of("X-Custom", "value"))
    .httpClient(myHttpClient)                    // override transport (testing, custom pooling)
    .build();
```

Request/response DTOs (`CollectionSchema`, `FieldSchema`, `SearchParams`,
`SuggestParams`, `AnalyticsQueryParams`, `TypoToleranceConfig`) are
immutable Java records; the many-optional-field ones ship a `Builder` for
ergonomic construction. Unset builder fields serialize as absent JSON,
never overriding the server's own default.

## Collections

```java
client.collections.create(schema);   // POST /collections
client.collections.list();           // GET /collections
client.collections.retrieve(name);   // GET /collections/{name}
client.collections.delete(name);     // DELETE /collections/{name}
```

## Documents

```java
var collection = client.collection("products");

collection.documents().index(doc1, doc2, doc3);  // POST   /collections/{name}/documents (upsert by id, varargs)
collection.documents().retrieve(id);             // GET    /collections/{name}/documents/{id}
collection.documents().delete(id);               // DELETE /collections/{name}/documents/{id}
```

`index()` always returns normally even if individual documents are
rejected — check `numFailed()` and `results()` on the response.

## Search

```java
collection.search(SearchParams.builder()
    .q("wireless mouse")
    .queryBy("title", "description")
    .filter("brand:=Logitech && price:<5000")
    .sort("_text_match:desc,price:asc")
    .facet("brand", "year")
    .limit(20)
    .offset(0)
    .prefix(true)
    .typoTolerance(true)
    .matchMode(MatchMode.ALL) // or MatchMode.ANY
    .build());
```

`foundIsExact()` on the response is `false` once block-max WAND pruning has
skipped part of a term's postings for a broad query — at that point
`found()` and facet counts are a lower bound, not an exact count. See the
[Known limitations section of the main
README](https://github.com/adikeshri/tachyon#known-limitations).

## Autocomplete

```java
collection.suggest(SuggestParams.builder("wir").limit(5).build());
```

## Analytics

```java
client.analytics.top(AnalyticsQueryParams.builder().collection("products").limit(10).build());
client.analytics.zeroResults(AnalyticsQueryParams.builder().collection("products").build());
client.analytics.latency();
```

Analytics are in-memory only and reset when the server restarts.

## Operations

```java
client.health();   // GET /health — no API key required
client.metrics();  // GET /metrics — Prometheus exposition format, returned as text
```

## Errors

Every non-2xx response throws a `TachyonException` subclass (unchecked)
carrying the server's stable `code` and the HTTP `statusCode`:

```java
try {
    client.collections.retrieve("does-not-exist");
} catch (TachyonNotFoundException e) {
    System.out.println(e.getCode() + " " + e.getStatusCode() + " " + e.getMessage()); // collection_not_found 404 ...
} catch (TachyonException e) {
    // any other API error
}
```

| Exception | Status | Codes |
|---|---|---|
| `TachyonRequestException` | 400 | `invalid_schema`, `invalid_document`, `invalid_query`, `invalid_json` |
| `TachyonAuthenticationException` | 401 | `unauthorized` |
| `TachyonAuthorizationException` | 403 | `forbidden` |
| `TachyonNotFoundException` | 404 | `collection_not_found`, `document_not_found` |
| `TachyonConflictException` | 409 | `collection_exists` |
| `TachyonServerException` | 5xx | `corrupt_data`, `io_error`, `internal_error` |

Network failures and timeouts throw `TachyonConnectionException` /
`TachyonTimeoutException` instead, since there's no server response to read
a code from.

## Development

```bash
mvn test
```

Unit tests use a real local loopback HTTP server
(`com.sun.net.httpserver.HttpServer`, built into the JDK) rather than a
mocking library — `java.net.http.HttpClient` doesn't have an easy seam to
mock directly.

### Integration tests

`mvn test` only runs the mocked unit suite (Surefire, `src/test/java`,
excluding `*IT.java`). A second suite in `src/it/java`, run by Maven
Failsafe, exercises every functional path — collections, documents, search
(filters, sort, facets, pagination, prefix, typo tolerance, match mode,
phrases), suggest, analytics, auth, and error paths — against a real,
running Tachyon server:

```bash
docker run -d -p 8108:8108 \
  -e TACHYON_ADMIN_KEY=admin-key -e TACHYON_SEARCH_KEY=search-key \
  adikeshri/tachyon

mvn verify
```

It points at `http://localhost:8108` by default; override with the
`TACHYON_URL`, `TACHYON_ADMIN_KEY`, `TACHYON_SEARCH_KEY` environment
variables. Every test cleans up the collections it creates, even on
failure.

## License

Apache 2.0. See [`LICENSE`](../LICENSE).
