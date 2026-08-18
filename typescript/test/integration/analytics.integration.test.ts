import { describe, expect, it } from 'vitest';
import { adminClient, uniqueName, withCollection } from './support.js';

const client = adminClient();

describe('analytics (integration)', () => {
  it('records search counts and zero-result queries, scoped by collection', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await collection.documents.index({ id: '1', title: 'gizmo' });

      await collection.search({ q: 'gizmo' });
      await collection.search({ q: 'gizmo' });
      await collection.search({ q: 'zzz-totally-absent-term' });

      const top = await client.analytics.top({ collection: name });
      const gizmoQuery = top.queries.find((q) => q.query === 'gizmo');
      expect(gizmoQuery).toBeDefined();
      expect(gizmoQuery?.count).toBe(2);
      expect(gizmoQuery?.collection).toBe(name);

      const zeroResults = await client.analytics.zeroResults({ collection: name });
      expect(zeroResults.queries.some((q) => q.query === 'zzz-totally-absent-term')).toBe(true);
    });
  });

  it('respects the limit parameter', async () => {
    const top = await client.analytics.top({ limit: 1 });
    expect(top.queries.length).toBeLessThanOrEqual(1);
  });

  it('reports latency percentiles across all searches so far', async () => {
    const latency = await client.analytics.latency();
    expect(latency.total_searches).toBeGreaterThan(0);
    expect(latency.p50_ms).toBeGreaterThanOrEqual(0);
    expect(latency.p95_ms).toBeGreaterThanOrEqual(latency.p50_ms);
    expect(latency.p99_ms).toBeGreaterThanOrEqual(latency.p95_ms);
  });
});
