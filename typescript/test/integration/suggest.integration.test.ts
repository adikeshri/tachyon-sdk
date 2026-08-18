import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import type { Collection } from '../../src/collection.js';
import { adminClient, uniqueName } from './support.js';

const client = adminClient();
let collectionName: string;
let collection: Collection;

beforeAll(async () => {
  collectionName = uniqueName('suggest');
  await client.collections.create({
    name: collectionName,
    fields: [{ name: 'title', type: 'text' }],
  });
  collection = client.collection(collectionName);
  await collection.documents.index([
    { id: '1', title: 'wireless mouse' },
    { id: '2', title: 'wireless keyboard' },
    { id: '3', title: 'wireless mouse' },
    { id: '4', title: 'wired cable' },
  ]);
});

afterAll(async () => {
  await client.collections.delete(collectionName).catch(() => {});
});

describe('suggest (integration)', () => {
  it('completes a prefix with live-document counts', async () => {
    const result = await collection.suggest({ q: 'wir' });
    const texts = result.suggestions.map((s) => s.text);
    expect(texts).toContain('wireless');
    expect(texts).toContain('wired');

    const wireless = result.suggestions.find((s) => s.text === 'wireless');
    expect(wireless?.count).toBeGreaterThanOrEqual(2);
  });

  it('caps suggestions at the requested limit', async () => {
    const result = await collection.suggest({ q: 'wir', limit: 1 });
    expect(result.suggestions.length).toBeLessThanOrEqual(1);
  });
});
