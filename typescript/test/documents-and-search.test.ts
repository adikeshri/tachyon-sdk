import { describe, expect, it } from 'vitest';
import { Tachyon } from '../src/index.js';
import { createMockFetch } from './helpers.js';

describe('documents', () => {
  it('indexes a batch of documents and reports per-document results', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: {
        num_indexed: 1,
        num_failed: 1,
        results: [
          { success: true, id: '1' },
          { success: false, code: 'invalid_document', error: 'field `price`: expected an integer, got a string' },
        ],
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const result = await client.collection('products').documents.index([
      { id: '1', title: 'Wireless Mouse' },
      { id: '2', price: 'not a number' },
    ]);

    expect(result.num_indexed).toBe(1);
    expect(result.num_failed).toBe(1);
    expect(result.results[1]?.code).toBe('invalid_document');
    expect(calls[0]?.url).toBe('http://localhost:8108/collections/products/documents');
    expect(calls[0]?.method).toBe('POST');
  });

  it('retrieves a document by id', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { id: '1', title: 'Wireless Mouse' },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const doc = await client.collection('products').documents.retrieve('1');

    expect(doc.title).toBe('Wireless Mouse');
    expect(calls[0]?.url).toBe('http://localhost:8108/collections/products/documents/1');
  });

  it('deletes a document', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({ status: 204, text: '' }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await client.collection('products').documents.delete('1');

    expect(calls[0]?.method).toBe('DELETE');
    expect(calls[0]?.url).toBe('http://localhost:8108/collections/products/documents/1');
  });
});

describe('search', () => {
  it('serializes search params into the expected query string', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { found: 1, found_is_exact: true, search_time_ms: 1, hits: [] },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await client.collection('products').search({
      q: 'wireless mouse',
      queryBy: ['title', 'description'],
      filter: 'brand:=Logitech && price:<5000',
      sort: '_text_match:desc,price:asc',
      facet: ['brand', 'year'],
      limit: 20,
      offset: 40,
      prefix: false,
      typoTolerance: true,
      matchMode: 'any',
    });

    const url = new URL(calls[0]!.url);
    expect(url.pathname).toBe('/collections/products/search');
    expect(url.searchParams.get('q')).toBe('wireless mouse');
    expect(url.searchParams.get('query_by')).toBe('title,description');
    expect(url.searchParams.get('filter')).toBe('brand:=Logitech && price:<5000');
    expect(url.searchParams.get('sort')).toBe('_text_match:desc,price:asc');
    expect(url.searchParams.get('facet')).toBe('brand,year');
    expect(url.searchParams.get('limit')).toBe('20');
    expect(url.searchParams.get('offset')).toBe('40');
    expect(url.searchParams.get('prefix')).toBe('false');
    expect(url.searchParams.get('typo_tolerance')).toBe('true');
    expect(url.searchParams.get('match_mode')).toBe('any');
  });

  it('omits unset params rather than sending empty strings', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { found: 0, found_is_exact: true, search_time_ms: 0, hits: [] },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await client.collection('products').search();

    const url = new URL(calls[0]!.url);
    expect(url.searchParams.has('filter')).toBe(false);
    expect(url.searchParams.has('limit')).toBe(false);
  });

  it('returns hits and facets, and surfaces found_is_exact', async () => {
    const { fetchImpl } = createMockFetch(() => ({
      status: 200,
      body: {
        found: 1240,
        found_is_exact: false,
        search_time_ms: 12,
        hits: [{ document: { id: '1', title: 'Wireless Mouse' }, text_match: 554.788 }],
        facets: { brand: { Logitech: 1240, Razer: 830 } },
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const results = await client.collection('products').search({ q: 'wireless mouse' });

    expect(results.found).toBe(1240);
    expect(results.found_is_exact).toBe(false);
    expect(results.hits[0]?.document.title).toBe('Wireless Mouse');
    expect(results.facets?.brand?.Logitech).toBe(1240);
  });

  it('suggests completions', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: {
        suggestions: [
          { text: 'wireless', count: 3, typos: 0 },
          { text: 'wired', count: 2, typos: 0 },
        ],
        search_time_ms: 0,
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const result = await client.collection('products').suggest({ q: 'wir', limit: 5 });

    expect(result.suggestions).toHaveLength(2);
    const url = new URL(calls[0]!.url);
    expect(url.pathname).toBe('/collections/products/suggest');
    expect(url.searchParams.get('q')).toBe('wir');
    expect(url.searchParams.get('limit')).toBe('5');
  });
});
