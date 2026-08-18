import { describe, expect, it } from 'vitest';
import { Tachyon } from '../../src/client.js';
import { TachyonAuthenticationError, TachyonAuthorizationError, TachyonConnectionError } from '../../src/errors.js';
import { adminClient, anonymousClient, searchOnlyClient, uniqueName } from './support.js';

describe('auth and error paths (integration)', () => {
  it('rejects requests with no API key when auth is enabled (401)', async () => {
    const failure = anonymousClient().collections.list();
    await expect(failure).rejects.toBeInstanceOf(TachyonAuthenticationError);
    await expect(failure).rejects.toMatchObject({ code: 'unauthorized', status: 401 });
  });

  it('lets a search-only key read', async () => {
    const list = await searchOnlyClient().collections.list();
    expect(Array.isArray(list)).toBe(true);
  });

  it('rejects a search-only key attempting a write (403)', async () => {
    const failure = searchOnlyClient().collections.create({ name: uniqueName('coll'), fields: [] });
    await expect(failure).rejects.toBeInstanceOf(TachyonAuthorizationError);
    await expect(failure).rejects.toMatchObject({ code: 'forbidden', status: 403 });
  });

  it('wraps a real network failure in TachyonConnectionError', async () => {
    // Nothing listens on port 1 (a reserved, unused TCP port), so this exercises
    // the real fetch failure path rather than a mocked one.
    const unreachable = new Tachyon({ url: 'http://127.0.0.1:1', timeoutMs: 2000 });
    await expect(unreachable.health()).rejects.toBeInstanceOf(TachyonConnectionError);
  });

  it('accepts a valid admin key end to end', async () => {
    const list = await adminClient().collections.list();
    expect(Array.isArray(list)).toBe(true);
  });
});
