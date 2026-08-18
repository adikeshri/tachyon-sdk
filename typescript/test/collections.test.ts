import { describe, expect, it } from 'vitest';
import { Tachyon } from '../src/index.js';
import { createMockFetch } from './helpers.js';

describe('collections', () => {
  it('creates a collection with the given schema', async () => {
    const { fetchImpl, calls } = createMockFetch((call) => ({
      status: 201,
      body: {
        name: 'products',
        fields: [{ name: 'title', type: 'text' }],
        num_documents: 0,
        num_segments: 0,
      },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', apiKey: 'admin-key', fetch: fetchImpl });

    const info = await client.collections.create({
      name: 'products',
      fields: [{ name: 'title', type: 'text' }],
    });

    expect(info.name).toBe('products');
    expect(info.num_documents).toBe(0);
    expect(calls).toHaveLength(1);
    expect(calls[0]?.method).toBe('POST');
    expect(calls[0]?.url).toBe('http://localhost:8108/collections');
    expect(calls[0]?.headers['x-tachyon-api-key']).toBe('admin-key');
    expect(JSON.parse(calls[0]?.body ?? '{}')).toEqual({
      name: 'products',
      fields: [{ name: 'title', type: 'text' }],
    });
  });

  it('lists collections', async () => {
    const { fetchImpl } = createMockFetch(() => ({
      status: 200,
      body: [{ name: 'products', fields: [], num_documents: 5, num_segments: 1 }],
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const collections = await client.collections.list();

    expect(collections).toHaveLength(1);
    expect(collections[0]?.name).toBe('products');
  });

  it('retrieves a single collection', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { name: 'products', fields: [], num_documents: 5, num_segments: 1 },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    const info = await client.collections.retrieve('products');

    expect(info.num_documents).toBe(5);
    expect(calls[0]?.url).toBe('http://localhost:8108/collections/products');
  });

  it('URL-encodes collection names', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({
      status: 200,
      body: { name: 'my products', fields: [], num_documents: 0, num_segments: 0 },
    }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await client.collections.retrieve('my products');

    expect(calls[0]?.url).toBe('http://localhost:8108/collections/my%20products');
  });

  it('deletes a collection and resolves with no body on 204', async () => {
    const { fetchImpl, calls } = createMockFetch(() => ({ status: 204, text: '' }));
    const client = new Tachyon({ url: 'http://localhost:8108', fetch: fetchImpl });

    await expect(client.collections.delete('products')).resolves.toBeUndefined();
    expect(calls[0]?.method).toBe('DELETE');
    expect(calls[0]?.url).toBe('http://localhost:8108/collections/products');
  });
});
