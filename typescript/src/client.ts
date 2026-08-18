import { AnalyticsApi } from './analytics.js';
import { Collection } from './collection.js';
import { CollectionsApi } from './collections.js';
import { HttpClient } from './http.js';
import type { HealthResponse, TachyonDocument } from './types.js';

export interface TachyonClientOptions {
  /** Full base URL, e.g. `http://localhost:8108`. Takes precedence over host/port/protocol. */
  url?: string;
  /** Default `localhost`. Ignored if `url` is set. */
  host?: string;
  /** Default `8108`. Ignored if `url` is set. */
  port?: number;
  /** Default `http`. Ignored if `url` is set. */
  protocol?: string;
  /** Sent as `X-TACHYON-API-KEY`. Use an admin key for writes, a search key for read-only access. */
  apiKey?: string;
  /** Per-request timeout in milliseconds. Default 15000. */
  timeoutMs?: number;
  /** Extra headers merged into every request. */
  headers?: Record<string, string>;
  /** Override the `fetch` implementation (mainly for testing, or Node < 18). */
  fetch?: typeof fetch;
}

/**
 * Client for a single Tachyon server.
 *
 * ```ts
 * const client = new Tachyon({ url: 'http://localhost:8108', apiKey: 'my-admin-key' });
 * await client.collections.create({ name: 'products', fields: [{ name: 'title', type: 'text' }] });
 * await client.collection('products').documents.index({ id: '1', title: 'Wireless Mouse' });
 * const results = await client.collection('products').search({ q: 'wireless mouse' });
 * ```
 */
export class Tachyon {
  readonly collections: CollectionsApi;
  readonly analytics: AnalyticsApi;
  private readonly http: HttpClient;

  constructor(options: TachyonClientOptions = {}) {
    const url = options.url ?? `${options.protocol ?? 'http'}://${options.host ?? 'localhost'}:${options.port ?? 8108}`;
    this.http = new HttpClient({
      url,
      apiKey: options.apiKey,
      timeoutMs: options.timeoutMs,
      headers: options.headers,
      fetch: options.fetch,
    });
    this.collections = new CollectionsApi(this.http);
    this.analytics = new AnalyticsApi(this.http);
  }

  /** Get a handle scoped to one collection, for documents/search/suggest. */
  collection<T extends TachyonDocument = TachyonDocument>(name: string): Collection<T> {
    return new Collection<T>(this.http, name);
  }

  /** `GET /health`. Always reachable without an API key. */
  health(): Promise<HealthResponse> {
    return this.http.request<HealthResponse>('GET', '/health');
  }

  /** `GET /metrics`. Prometheus exposition format, returned as plain text. */
  metrics(): Promise<string> {
    return this.http.requestText('GET', '/metrics');
  }
}
