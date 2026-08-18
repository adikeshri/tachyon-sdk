import { describe, expect, it } from 'vitest';
import { adminClient, anonymousClient } from './support.js';

describe('operations (integration)', () => {
  it('reports health without requiring an API key', async () => {
    const health = await anonymousClient().health();
    expect(health.ok).toBe(true);
    expect(typeof health.version).toBe('string');
    expect(health.num_collections).toBeGreaterThanOrEqual(0);
  });

  it('exposes Prometheus metrics as plain text', async () => {
    const metrics = await adminClient().metrics();
    expect(metrics).toContain('tachyon_uptime_seconds');
    expect(metrics).toContain('tachyon_collections');
  });
});
