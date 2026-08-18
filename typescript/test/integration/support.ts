import { randomUUID } from 'node:crypto';
import { Tachyon } from '../../src/client.js';
import type { CollectionSchema } from '../../src/types.js';

export const BASE_URL = process.env.TACHYON_URL ?? 'http://localhost:8108';
export const ADMIN_KEY = process.env.TACHYON_ADMIN_KEY ?? 'admin-key';
export const SEARCH_KEY = process.env.TACHYON_SEARCH_KEY ?? 'search-key';

export function adminClient(): Tachyon {
  return new Tachyon({ url: BASE_URL, apiKey: ADMIN_KEY });
}

export function searchOnlyClient(): Tachyon {
  return new Tachyon({ url: BASE_URL, apiKey: SEARCH_KEY });
}

export function anonymousClient(): Tachyon {
  return new Tachyon({ url: BASE_URL });
}

export function uniqueName(prefix: string): string {
  return `${prefix}-${randomUUID()}`;
}

/**
 * Creates a collection for the duration of `fn`, then deletes it — even if
 * `fn` throws — so integration tests never leak collections into the
 * shared server between runs.
 */
export async function withCollection<T>(
  client: Tachyon,
  schema: Omit<CollectionSchema, 'name'> & { name?: string },
  fn: (name: string) => Promise<T>,
): Promise<T> {
  const name = schema.name ?? uniqueName('coll');
  await client.collections.create({ ...schema, name });
  try {
    return await fn(name);
  } finally {
    await client.collections.delete(name).catch(() => {});
  }
}
