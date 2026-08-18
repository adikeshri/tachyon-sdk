import { describe, expect, it } from 'vitest';
import { Tachyon } from '../src/index.js';
import { createMockFetch } from './helpers.js';

describe('analytics', () => {
  it('fetches top queries scoped to a collection', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: {
        queries: [
          {
            query: 'wireless mouse',
            collection: 'products',
            count: 3,
            zero_result_count: 0,
            last_result_count: 12,
            avg_latency_ms: 1.8,
            last_seen: 1786625778150,
          },
        ],
        tracked_queries: 1,
        dropped_queries: 0,
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const result = await client.analytics.top({ collection: 'products', limit: 10 });

    expect(result.queries[0]?.query).toBe('wireless mouse');
    const url = new URL(calls[0]!.url);
    expect(url.pathname).toBe('/analytics/top');
    expect(url.searchParams.get('collection')).toBe('products');
    expect(url.searchParams.get('limit')).toBe('10');
  });

  it('fetches zero-result queries', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { queries: [], tracked_queries: 0, dropped_queries: 0 },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await client.analytics.zeroResults();

    const url = new URL(calls[0]!.url);
    expect(url.pathname).toBe('/analytics/zero-results');
  });

  it('fetches latency percentiles', async () => {
    const { fetchImpl } = createMockFetch(() => ({
      status: 200,
      body: {
        count: 20,
        mean_ms: 1.9,
        p50_ms: 2.0,
        p95_ms: 4.0,
        p99_ms: 4.0,
        max_ms: 3.4,
        total_searches: 20,
        uptime_seconds: 61,
        queries_per_second: 0.33,
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const result = await client.analytics.latency();

    expect(result.p95_ms).toBe(4.0);
  });
});

describe('operations', () => {
  it('checks health', async () => {
    const { fetchImpl } = createMockFetch(() => ({
      status: 200,
      body: { ok: true, version: '0.1.0', uptime_seconds: 61, num_collections: 1 },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const health = await client.health();

    expect(health.ok).toBe(true);
  });

  it('fetches metrics as raw Prometheus text, not JSON', async () => {
    const { fetchImpl } = createMockFetch(() => ({
      status: 200,
      text: '# HELP tachyon_uptime_seconds\ntachyon_uptime_seconds 61\n',
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const metrics = await client.metrics();

    expect(metrics).toContain('tachyon_uptime_seconds 61');
  });
});
