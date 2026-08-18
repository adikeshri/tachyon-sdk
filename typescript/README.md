# tachyon-sdk (TypeScript / JavaScript)

Official TypeScript/JavaScript client for [Tachyon](https://github.com/adikeshri/tachyon),
the typo-tolerant full-text search engine.

```bash
npm install tachyon-sdk
```

Requires Node 18+ (for global `fetch`), or any runtime that provides one —
browsers, Deno, Bun, Cloudflare Workers. This package ships ESM only.

## Quickstart

```ts
import { Tachyon } from 'tachyon-sdk';

const client = new Tachyon({ url: 'http://localhost:8108', apiKey: 'my-admin-key' });

await client.collections.create({
  name: 'products',
  fields: [
    { name: 'title', type: 'text' },
    { name: 'brand', type: 'keyword', facet: true },
    { name: 'price', type: 'int', filter: true, sort: true },
  ],
});

await client.collection('products').documents.index([
  { id: '1', title: 'Wireless Mouse', brand: 'Logitech', price: 2999 },
  { id: '2', title: 'Mechanical Keyboard', brand: 'Razer', price: 8999 },
]);

const results = await client.collection('products').search({ q: 'wireless mouse' });
for (const hit of results.hits) {
  console.log(hit.document.title, hit.text_match);
}
```

## Client options

```ts
new Tachyon({
  url: 'http://localhost:8108', // or host/port/protocol
  apiKey: '...',                // admin key (read/write) or search key (read-only)
  timeoutMs: 15_000,
  headers: { 'X-Custom': 'value' },
  fetch: myFetch,                // override fetch (testing, or Node < 18)
});
```

## Collections

```ts
client.collections.create(schema); // POST /collections
client.collections.list(); // GET /collections
client.collections.retrieve(name); // GET /collections/{name}
client.collections.delete(name); // DELETE /collections/{name}
```

## Documents

```ts
const collection = client.collection('products');

collection.documents.index(docOrArray); // POST   /collections/{name}/documents  (upsert by id)
collection.documents.retrieve(id); // GET    /collections/{name}/documents/{id}
collection.documents.delete(id); // DELETE /collections/{name}/documents/{id}
```

`index()` always resolves at the HTTP level even if individual documents are
rejected — check `num_failed` and `results` on the response.

## Search

```ts
collection.search({
  q: 'wireless mouse',
  queryBy: ['title', 'description'],
  filter: 'brand:=Logitech && price:<5000',
  sort: '_text_match:desc,price:asc',
  facet: ['brand', 'year'],
  limit: 20,
  offset: 0,
  prefix: true,
  typoTolerance: true,
  matchMode: 'all', // or 'any'
});
```

`queryBy` and `facet` accept either a comma-separated string or a string
array. Pass a generic type parameter to `client.collection<T>(name)` to type
`document` on each hit.

`found_is_exact` on the response is `false` once block-max WAND pruning has
skipped part of a term's postings for a broad query — at that point `found`
and facet counts are a lower bound, not an exact count. See the [Known
limitations section of the main
README](https://github.com/adikeshri/tachyon#known-limitations).

## Autocomplete

```ts
await collection.suggest({ q: 'wir', limit: 5 });
```

## Analytics

```ts
await client.analytics.top({ collection: 'products', limit: 10 });
await client.analytics.zeroResults({ collection: 'products' });
await client.analytics.latency();
```

Analytics are in-memory only and reset when the server restarts.

## Operations

```ts
await client.health(); // GET /health — no API key required
await client.metrics(); // GET /metrics — Prometheus exposition format, returned as text
```

## Errors

Every non-2xx response rejects with a `TachyonError` subclass carrying the
server's stable `code` and the HTTP `status`:

```ts
import { TachyonError, TachyonNotFoundError } from 'tachyon-sdk';

try {
  await client.collections.retrieve('does-not-exist');
} catch (err) {
  if (err instanceof TachyonNotFoundError) {
    console.log(err.code, err.status, err.message); // collection_not_found 404 ...
  } else if (err instanceof TachyonError) {
    // any other API error
  } else {
    throw err;
  }
}
```

| Class | Status | Codes |
|---|---|---|
| `TachyonRequestError` | 400 | `invalid_schema`, `invalid_document`, `invalid_query`, `invalid_json` |
| `TachyonAuthenticationError` | 401 | `unauthorized` |
| `TachyonAuthorizationError` | 403 | `forbidden` |
| `TachyonNotFoundError` | 404 | `collection_not_found`, `document_not_found` |
| `TachyonConflictError` | 409 | `collection_exists` |
| `TachyonServerError` | 5xx | `corrupt_data`, `io_error`, `internal_error` |

Network failures and timeouts reject with `TachyonConnectionError` /
`TachyonTimeoutError` instead, since there's no server response to read a
code from.

## Development

```bash
npm install
npm test
npm run typecheck
npm run build
```

### Integration tests

`npm test` only runs the mocked unit suite. A second suite in
`test/integration/` exercises every functional path — collections,
documents, search (filters, sort, facets, pagination, prefix, typo
tolerance, match mode, phrases), suggest, analytics, auth, and error
paths — against a real, running Tachyon server:

```bash
docker run -d -p 8108:8108 \
  -e TACHYON_ADMIN_KEY=admin-key -e TACHYON_SEARCH_KEY=search-key \
  adikeshri/tachyon

npm run test:integration
```

It points at `http://localhost:8108` by default; override with
`TACHYON_URL`, `TACHYON_ADMIN_KEY`, `TACHYON_SEARCH_KEY`. Every test cleans
up the collections it creates, even on failure.

## License

Apache 2.0. See [`LICENSE`](../LICENSE).
