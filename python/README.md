# tachyon-sdk (Python)

Official Python client for [Tachyon](https://github.com/adikeshri/tachyon),
the typo-tolerant full-text search engine.

```bash
pip install tachyon-sdk
```

Requires Python 3.8+.

## Quickstart

```python
from tachyon_sdk import Tachyon

client = Tachyon(url="http://localhost:8108", api_key="my-admin-key")

client.collections.create({
    "name": "products",
    "fields": [
        {"name": "title", "type": "text"},
        {"name": "brand", "type": "keyword", "facet": True},
        {"name": "price", "type": "int", "filter": True, "sort": True},
    ],
})

client.collection("products").documents.index([
    {"id": "1", "title": "Wireless Mouse", "brand": "Logitech", "price": 2999},
    {"id": "2", "title": "Mechanical Keyboard", "brand": "Razer", "price": 8999},
])

results = client.collection("products").search(q="wireless mouse")
for hit in results["hits"]:
    print(hit["document"]["title"], hit["text_match"])
```

## Client options

```python
Tachyon(
    url="http://localhost:8108",  # or host="localhost", port=8108, protocol="http"
    api_key="...",                # admin key (read/write) or search key (read-only)
    timeout=15.0,                 # seconds
    headers={"X-Custom": "value"},
    session=None,                 # pass your own requests.Session to share pooling
)
```

## Collections

```python
client.collections.create(schema)   # POST /collections
client.collections.list()           # GET /collections
client.collections.retrieve(name)   # GET /collections/{name}
client.collections.delete(name)     # DELETE /collections/{name}
```

## Documents

```python
collection = client.collection("products")

collection.documents.index(doc_or_list)  # POST   /collections/{name}/documents  (upsert by id)
collection.documents.retrieve(doc_id)    # GET    /collections/{name}/documents/{id}
collection.documents.delete(doc_id)      # DELETE /collections/{name}/documents/{id}
```

`index()` always succeeds at the HTTP level even if individual documents are
rejected — check `num_failed` and `results` on the response.

## Search

```python
collection.search(
    q="wireless mouse",
    query_by=["title", "description"],
    filter="brand:=Logitech && price:<5000",
    sort="_text_match:desc,price:asc",
    facet=["brand", "year"],
    limit=20,
    offset=0,
    prefix=True,
    typo_tolerance=True,
    match_mode="all",  # or "any"
)
```

`query_by` and `facet` accept either a comma-separated string or a list of
field names. Every parameter is optional keyword-only, except `q`.

`found_is_exact` on the response is `False` once block-max WAND pruning has
skipped part of a term's postings for a broad query — at that point `found`
and facet counts are a lower bound, not an exact count. See the [Known
limitations section of the main
README](https://github.com/adikeshri/tachyon#known-limitations).

## Autocomplete

```python
collection.suggest(q="wir", limit=5)
```

## Analytics

```python
client.analytics.top(collection="products", limit=10)
client.analytics.zero_results(collection="products")
client.analytics.latency()
```

Analytics are in-memory only and reset when the server restarts.

## Operations

```python
client.health()   # GET /health — no API key required
client.metrics()  # GET /metrics — Prometheus exposition format, returned as text
```

## Errors

Every non-2xx response raises a `TachyonError` subclass carrying the
server's stable `code` and the HTTP `status`:

```python
from tachyon_sdk import TachyonError, TachyonNotFoundError

try:
    client.collections.retrieve("does-not-exist")
except TachyonNotFoundError as e:
    print(e.code, e.status, e.message)  # collection_not_found 404 ...
except TachyonError as e:
    ...  # any other API error
```

| Class | Status | Codes |
|---|---|---|
| `TachyonRequestError` | 400 | `invalid_schema`, `invalid_document`, `invalid_query`, `invalid_json` |
| `TachyonAuthenticationError` | 401 | `unauthorized` |
| `TachyonAuthorizationError` | 403 | `forbidden` |
| `TachyonNotFoundError` | 404 | `collection_not_found`, `document_not_found` |
| `TachyonConflictError` | 409 | `collection_exists` |
| `TachyonServerError` | 5xx | `corrupt_data`, `io_error`, `internal_error` |

Network failures and timeouts raise `TachyonConnectionError` /
`TachyonTimeoutError` instead, since there's no server response to read a
code from.

## Types

Request and response shapes are `TypedDict`s in `tachyon_sdk.types` — plain
`dict`s at runtime, typed for editor/mypy support:
`CollectionSchema`, `CollectionInfo`, `FieldSchema`, `TachyonDocument`,
`SearchResponse`, `SuggestResponse`, `AnalyticsQueriesResponse`,
`AnalyticsLatencyResponse`, `HealthResponse`, and friends.

## Development

```bash
python -m venv .venv && source .venv/bin/activate
pip install -e ".[dev]"
pytest
mypy tachyon_sdk
```

### Integration tests

Plain `pytest` only runs the mocked unit suite (integration tests are
marker-excluded by default). A second suite in `tests/integration/`
exercises every functional path — collections, documents, search (filters,
sort, facets, pagination, prefix, typo tolerance, match mode, phrases),
suggest, analytics, auth, and error paths — against a real, running Tachyon
server:

```bash
docker run -d -p 8108:8108 \
  -e TACHYON_ADMIN_KEY=admin-key -e TACHYON_SEARCH_KEY=search-key \
  adikeshri/tachyon

pytest -m integration
```

It points at `http://localhost:8108` by default; override with the
`TACHYON_URL`, `TACHYON_ADMIN_KEY`, `TACHYON_SEARCH_KEY` environment
variables. Every test cleans up the collections it creates, even on
failure.

## License

Apache 2.0. See [`LICENSE`](../LICENSE).
