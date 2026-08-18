import { describe, expect, it } from 'vitest';
import { Tachyon } from '../src/index.js';
import {
  TachyonAuthenticationError,
  TachyonAuthorizationError,
  TachyonConflictError,
  TachyonConnectionError,
  TachyonNotFoundError,
  TachyonRequestError,
  TachyonServerError,
  TachyonTimeoutError,
} from '../src/errors.js';
import { createMockFetch } from './helpers.js';

describe('error mapping', () => {
  it.each([
    [400, 'invalid_query', TachyonRequestError],
    [401, 'unauthorized', TachyonAuthenticationError],
    [403, 'forbidden', TachyonAuthorizationError],
    [404, 'collection_not_found', TachyonNotFoundError],
    [409, 'collection_exists', TachyonConflictError],
    [500, 'internal_error', TachyonServerError],
  ] as const)('maps HTTP %i (%s) to %s', async (status, code, ErrorClass) => {
    const { fetchImpl } = createMockFetch(() => ({
      status,
      body: { error: { code, message: `boom: ${code}` } },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const failure = client.collections.retrieve('products');

    await expect(failure).rejects.toBeInstanceOf(ErrorClass);
    await expect(failure).rejects.toMatchObject({ code, status, message: `boom: ${code}` });
  });

  it('wraps network failures in TachyonConnectionError', async () => {
    const fetchImpl = (async () => {
      throw new TypeError('fetch failed');
    }) as typeof fetch;
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await expect(client.health()).rejects.toBeInstanceOf(TachyonConnectionError);
  });

  it('raises TachyonTimeoutError once the configured timeout elapses', async () => {
    const fetchImpl = ((_url: string, init?: RequestInit) =>
      new Promise((_resolve, reject) => {
        init?.signal?.addEventListener('abort', () => {
          const err = new Error('The operation was aborted');
          err.name = 'AbortError';
          reject(err);
        });
      })) as typeof fetch;
    const client = new Tachyon({ url: 'http://localhost:8108', timeoutMs: 5, fetch: fetchImpl });

    await expect(client.health()).rejects.toBeInstanceOf(TachyonTimeoutError);
  });

  it('falls back to a generic error code when the body is not the documented error shape', async () => {
    const { fetchImpl } = createMockFetch(() => ({ status: 500, text: 'upstream exploded' }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await expect(client.health()).rejects.toMatchObject({ code: 'internal_error', status: 500 });
  });
});
