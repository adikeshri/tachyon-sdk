# tachyon-sdk (.NET / C#)

Official .NET client for [Tachyon](https://github.com/adikeshri/tachyon), the
typo-tolerant full-text search engine.

```bash
dotnet add package Tachyon.Sdk
```

Targets .NET 8 (LTS). Zero third-party dependencies — just
`System.Net.Http` and `System.Text.Json`, both in the BCL.

## Quickstart

```csharp
using Tachyon.Sdk;

var client = new TachyonClient(new TachyonClientOptions
{
    Url = "http://localhost:8108",
    ApiKey = "my-admin-key",
});

await client.Collections.CreateAsync(new CollectionSchema
{
    Name = "products",
    Fields =
    [
        new FieldSchema { Name = "title", Type = FieldType.Text },
        new FieldSchema { Name = "brand", Type = FieldType.Keyword, Facet = true },
        new FieldSchema { Name = "price", Type = FieldType.Int, Filter = true, Sort = true },
    ],
});

var collection = client.Collection("products");
await collection.Documents.IndexAsync(
    new Document { ["id"] = "1", ["title"] = "Wireless Mouse", ["brand"] = "Logitech", ["price"] = 2999 },
    new Document { ["id"] = "2", ["title"] = "Mechanical Keyboard", ["brand"] = "Razer", ["price"] = 8999 });

var results = await collection.SearchAsync(new SearchParams { Q = "wireless mouse" });
foreach (var hit in results.Hits)
{
    Console.WriteLine($"{hit.Document["title"]} {hit.TextMatch}");
}
```

`Document` is a type alias for `System.Text.Json.Nodes.JsonObject` declared
via a `global using` inside this package — aliases don't cross assembly
boundaries in C#, so add the same line to your own project if you want to
write `Document` instead of `JsonObject`:

```csharp
global using Document = System.Text.Json.Nodes.JsonObject;
```

Or just use `JsonObject` directly; it's the exact same type either way.

## Client options

```csharp
var client = new TachyonClient(new TachyonClientOptions
{
    Url = "http://localhost:8108",
    ApiKey = "...",                          // admin key (read/write) or search key (read-only)
    Timeout = TimeSpan.FromSeconds(15),       // default 15s; ignored if HttpClient is supplied
    Headers = new Dictionary<string, string> { ["X-Custom"] = "value" },
    HttpClient = myHttpClient,                // override transport (testing, custom pooling)
});
```

Every method takes an optional trailing `CancellationToken`, in keeping with
.NET convention.

## Collections

```csharp
client.Collections.CreateAsync(schema);   // POST /collections
client.Collections.ListAsync();           // GET /collections
client.Collections.RetrieveAsync(name);   // GET /collections/{name}
client.Collections.DeleteAsync(name);     // DELETE /collections/{name}
```

`FieldSchema`'s optional attributes (`Facet`, `Filter`, `Sort`, `Index`,
`Optional`, `Boost`) are nullable, so leaving them unset doesn't override
the server's own default.

## Documents

```csharp
var collection = client.Collection("products");

collection.Documents.IndexAsync(doc1, doc2, doc3);  // POST   /collections/{name}/documents (upsert by id, params array)
collection.Documents.RetrieveAsync(id);             // GET    /collections/{name}/documents/{id}
collection.Documents.DeleteAsync(id);                // DELETE /collections/{name}/documents/{id}
```

`IndexAsync` always resolves at the HTTP level even if individual documents
are rejected — check `NumFailed` and `Results` on the response.

## Search

```csharp
await collection.SearchAsync(new SearchParams
{
    Q = "wireless mouse",
    QueryBy = ["title", "description"],
    Filter = "brand:=Logitech && price:<5000",
    Sort = "_text_match:desc,price:asc",
    Facet = ["brand", "year"],
    Limit = 20,
    Offset = 0,
    Prefix = true,
    TypoTolerance = true,
    MatchMode = MatchMode.All, // or MatchMode.Any
});
```

`FoundIsExact` on the response is `false` once block-max WAND pruning has
skipped part of a term's postings for a broad query — at that point `Found`
and facet counts are a lower bound, not an exact count. See the [Known
limitations section of the main
README](https://github.com/adikeshri/tachyon#known-limitations).

## Autocomplete

```csharp
await collection.SuggestAsync(new SuggestParams { Q = "wir", Limit = 5 });
```

## Analytics

```csharp
await client.Analytics.TopAsync(new AnalyticsQueryParams { Collection = "products", Limit = 10 });
await client.Analytics.ZeroResultsAsync(new AnalyticsQueryParams { Collection = "products" });
await client.Analytics.LatencyAsync();
```

Analytics are in-memory only and reset when the server restarts.

## Operations

```csharp
await client.HealthAsync();   // GET /health — no API key required
await client.MetricsAsync();  // GET /metrics — Prometheus exposition format, returned as text
```

## Errors

Every non-2xx response throws a `TachyonException` subclass carrying the
server's stable `Code` and the HTTP `StatusCode`:

```csharp
try
{
    await client.Collections.RetrieveAsync("does-not-exist");
}
catch (TachyonNotFoundException e)
{
    Console.WriteLine($"{e.Code} {e.StatusCode} {e.Message}"); // collection_not_found 404 ...
}
catch (TachyonException e)
{
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
dotnet build
dotnet test
```

`dotnet test` only runs the mocked unit suite (`Tachyon.Sdk.Tests`) — it's
the only test project in the `.sln`. The library itself targets `net8.0`
for maximum consumer reach; the test projects target `net9.0` so they run
on whatever's installed, which works fine since a net9.0 host can consume a
net8.0 library directly.

### Integration tests

A separate project, `Tachyon.Sdk.IntegrationTests`, isn't part of the
solution's default test run — it exercises every functional path
(collections, documents, search, suggest, analytics, auth, error paths)
against a real, running Tachyon server:

```bash
docker run -d -p 8108:8108 \
  -e TACHYON_ADMIN_KEY=admin-key -e TACHYON_SEARCH_KEY=search-key \
  adikeshri/tachyon

dotnet test Tachyon.Sdk.IntegrationTests/Tachyon.Sdk.IntegrationTests.csproj
```

It points at `http://localhost:8108` by default; override with the
`TACHYON_URL`, `TACHYON_ADMIN_KEY`, `TACHYON_SEARCH_KEY` environment
variables. Every test cleans up the collections it creates, even on
failure.

## License

Apache 2.0. See [`LICENSE`](../LICENSE).
