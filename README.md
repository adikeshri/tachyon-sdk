<p align="center">
  <img src="assets/logo.png" alt="Tachyon logo" width="160">
</p>

# tachyon-sdk

Official client libraries for [Tachyon](https://github.com/adikeshri/tachyon),
the open-source, typo-tolerant full-text search engine.

Both SDKs are thin, fully-typed wrappers over Tachyon's JSON/HTTP API — same
resources, same error codes, same defaults, in idiomatic TypeScript and
Python. See [`docs/api.md`](https://github.com/adikeshri/tachyon/blob/main/docs/api.md)
in the main repo for the underlying API reference either one is built on.

| | |
|---|---|
| [**TypeScript / JavaScript**](typescript) | `npm install tachyon-sdk` — Node 18+, ESM |
| [**Python**](python) | `pip install tachyon-sdk` — Python 3.8+ |

## Quickstart

Start a Tachyon server ([Docker Hub](https://hub.docker.com/r/adikeshri/tachyon)):

```bash
docker run -p 8108:8108 adikeshri/tachyon
```

TypeScript:

```ts
import { Tachyon } from 'tachyon-sdk';

const client = new Tachyon({ url: 'http://localhost:8108' });

await client.collections.create({
  name: 'products',
  fields: [{ name: 'title', type: 'text' }],
});
await client.collection('products').documents.index({ id: '1', title: 'Wireless Mouse' });
const results = await client.collection('products').search({ q: 'wireless mouse' });
```

Python:

```python
from tachyon_sdk import Tachyon

client = Tachyon(url="http://localhost:8108")

client.collections.create({"name": "products", "fields": [{"name": "title", "type": "text"}]})
client.collection("products").documents.index({"id": "1", "title": "Wireless Mouse"})
results = client.collection("products").search(q="wireless mouse")
```

Full usage — client options, every resource, error handling — is documented
in each SDK's own README: [`typescript/README.md`](typescript/README.md),
[`python/README.md`](python/README.md).

## What's covered

Both clients cover the whole of Tachyon's HTTP API:

- **Collections** — create, list, retrieve, delete
- **Documents** — index (upsert), retrieve, delete
- **Search** — full parameter set (`filter`, `sort`, `facet`, pagination,
  prefix, typo tolerance, match mode), including `found_is_exact` for
  block-max WAND pruning
- **Autocomplete** — `suggest`
- **Analytics** — top queries, zero-result queries, latency percentiles
- **Operations** — `health`, `metrics`

Both map every documented error code (`invalid_schema`, `unauthorized`,
`collection_not_found`, ...) onto a typed exception carrying the server's
`code` and HTTP `status`, and distinguish API errors from connection/timeout
failures.

## Repository layout

```text
typescript/   npm package "tachyon-sdk"
python/       PyPI package "tachyon-sdk" (import tachyon_sdk)
```

Each is independently versioned, tested, and published; see the README in
each directory for development setup.

## Contributing

Issues and pull requests are welcome. For the server itself, see
[adikeshri/tachyon](https://github.com/adikeshri/tachyon).

## License

Apache 2.0. See [`LICENSE`](LICENSE).
