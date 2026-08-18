import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import type { Collection } from '../../src/collection.js';
import { adminClient, uniqueName, withCollection } from './support.js';

const client = adminClient();
let collectionName: string;
let collection: Collection;

beforeAll(async () => {
  collectionName = uniqueName('search');
  await client.collections.create({
    name: collectionName,
    fields: [
      { name: 'title', type: 'text' },
      { name: 'description', type: 'text' },
      { name: 'brand', type: 'keyword', facet: true },
      { name: 'price', type: 'int', filter: true, sort: true },
      { name: 'in_stock', type: 'bool', filter: true },
    ],
  });
  collection = client.collection(collectionName);
  await collection.documents.index([
    {
      id: '1',
      title: 'Wireless Mouse',
      description: 'A great wireless mouse for everyday use',
      brand: 'Logitech',
      price: 2999,
      in_stock: true,
    },
    {
      id: '2',
      title: 'Mechanical Keyboard',
      description: 'Clicky keys for typing enthusiasts',
      brand: 'Razer',
      price: 8999,
      in_stock: false,
    },
    {
      id: '3',
      title: 'Wireless Keyboard',
      description: 'Silent wireless keyboard',
      brand: 'Logitech',
      price: 5999,
      in_stock: true,
    },
    {
      id: '4',
      title: 'Gaming Mouse',
      description: 'Wired precision gaming mouse',
      brand: 'Razer',
      price: 4999,
      in_stock: true,
    },
    { id: '5', title: 'USB Cable', description: 'A basic wired cable', brand: 'Anker', price: 999, in_stock: true },
  ]);
});

afterAll(async () => {
  await client.collections.delete(collectionName).catch(() => {});
});

describe('search (integration)', () => {
  it('matches on title and description by default, with an exact count', async () => {
    const results = await collection.search({ q: 'wireless' });
    expect(results.found).toBe(2);
    expect(results.found_is_exact).toBe(true);
    const ids = results.hits.map((h) => h.document.id).sort();
    expect(ids).toEqual(['1', '3']);
  });

  it('restricts matching to the fields named in query_by', async () => {
    const onlyDescription = await collection.search({ q: 'Clicky', queryBy: 'description' });
    expect(onlyDescription.found).toBe(1);
    expect(onlyDescription.hits[0]?.document.id).toBe('2');

    const onlyTitle = await collection.search({ q: 'Clicky', queryBy: 'title' });
    expect(onlyTitle.found).toBe(0);
  });

  it('accepts query_by as an array and joins it with commas', async () => {
    const results = await collection.search({ q: 'Clicky', queryBy: ['title', 'description'] });
    expect(results.found).toBe(1);
  });

  describe('filters', () => {
    it('equality', async () => {
      const results = await collection.search({ filter: 'brand:=Logitech' });
      expect(results.found).toBe(2);
    });

    it('numeric comparison', async () => {
      const results = await collection.search({ filter: 'price:<5000' });
      expect(results.found).toBe(3); // 999, 2999, 4999
    });

    it('inclusive range', async () => {
      const results = await collection.search({ filter: 'price:[1000..5000]' });
      const ids = results.hits.map((h) => h.document.id).sort();
      expect(ids).toEqual(['1', '4']); // 2999 and 4999; 999 and 5999 fall outside
    });

    it('set membership', async () => {
      const results = await collection.search({ filter: 'brand:=[Logitech,Razer]' });
      expect(results.found).toBe(4);
    });

    it('boolean equality', async () => {
      const results = await collection.search({ filter: 'in_stock:=true' });
      expect(results.found).toBe(4);
    });

    it('&& and || combined with grouping', async () => {
      const results = await collection.search({
        filter: '(brand:=Logitech || brand:=Razer) && price:<5000',
      });
      const ids = results.hits.map((h) => h.document.id).sort();
      expect(ids).toEqual(['1', '4']);
    });

    it('negation only matches documents that have the field', async () => {
      const results = await collection.search({ filter: 'brand:!=Razer' });
      const ids = results.hits.map((h) => h.document.id).sort();
      expect(ids).toEqual(['1', '3', '5']);
    });

    it('negation excludes documents missing the field entirely, not just non-matches', async () => {
      await withCollection(
        client,
        { fields: [{ name: 'title', type: 'text' }, { name: 'brand', type: 'keyword', filter: true }] },
        async (name) => {
          const scoped = client.collection(name);
          await scoped.documents.index([
            { id: 'a', title: 'has brand', brand: 'Razer' },
            { id: 'b', title: 'no brand at all' },
          ]);
          const results = await scoped.search({ filter: 'brand:!=Razer' });
          // 'b' has no brand field at all, so "not Razer" must not match it.
          expect(results.found).toBe(0);
        },
      );
    });
  });

  describe('sorting', () => {
    it('sorts ascending', async () => {
      const results = await collection.search({ q: '', sort: 'price:asc', limit: 10 });
      const prices = results.hits.map((h) => h.document.price as number);
      expect(prices).toEqual([...prices].sort((a, b) => a - b));
    });

    it('sorts descending', async () => {
      const results = await collection.search({ q: '', sort: 'price:desc', limit: 10 });
      const prices = results.hits.map((h) => h.document.price as number);
      expect(prices).toEqual([...prices].sort((a, b) => b - a));
    });
  });

  describe('pagination', () => {
    it('limit/offset moves through the result set without overlap', async () => {
      const page1 = await collection.search({ q: '', sort: 'price:asc', limit: 2, offset: 0 });
      const page2 = await collection.search({ q: '', sort: 'price:asc', limit: 2, offset: 2 });

      expect(page1.found).toBe(5);
      expect(page2.found).toBe(5);
      expect(page1.hits).toHaveLength(2);
      expect(page2.hits).toHaveLength(2);
      const page1Ids = page1.hits.map((h) => h.document.id);
      const page2Ids = page2.hits.map((h) => h.document.id);
      expect(page1Ids.some((id) => page2Ids.includes(id))).toBe(false);
    });
  });

  describe('prefix matching', () => {
    it('prefix-expands the final token by default', async () => {
      const results = await collection.search({ q: 'wir' });
      expect(results.found).toBeGreaterThanOrEqual(2);
    });

    it('requires a full token match when prefix is disabled', async () => {
      const results = await collection.search({ q: 'wir', prefix: false });
      expect(results.found).toBe(0);
    });
  });

  describe('typo tolerance', () => {
    it('corrects a typo by default', async () => {
      const results = await collection.search({ q: 'wirelss' });
      expect(results.found).toBeGreaterThanOrEqual(1);
    });

    it('finds nothing when typo tolerance is explicitly disabled', async () => {
      const results = await collection.search({ q: 'wirelss', typoTolerance: false });
      expect(results.found).toBe(0);
    });
  });

  describe('match_mode', () => {
    it('all (default) requires every token to be present', async () => {
      const results = await collection.search({ q: 'wireless zzznonexistentterm' });
      expect(results.found).toBe(0);
    });

    it('any requires only one token to be present', async () => {
      const results = await collection.search({ q: 'wireless zzznonexistentterm', matchMode: 'any' });
      expect(results.found).toBeGreaterThanOrEqual(2);
    });
  });

  it('phrase queries require adjacency within a single field', async () => {
    const results = await collection.search({ q: '"wireless mouse"' });
    expect(results.found).toBe(1);
    expect(results.hits[0]?.document.id).toBe('1');
  });

  it('facets count every matching document, not just the page', async () => {
    const results = await collection.search({ q: '', facet: 'brand', limit: 1 });
    expect(results.hits).toHaveLength(1); // page size respected
    expect(results.facets?.brand).toEqual({ Logitech: 2, Razer: 2, Anker: 1 });
  });
});
