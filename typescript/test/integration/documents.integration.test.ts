import { describe, expect, it } from 'vitest';
import { TachyonNotFoundError } from '../../src/errors.js';
import { adminClient, withCollection } from './support.js';

const client = adminClient();

describe('documents (integration)', () => {
  it('indexes a single document (not wrapped in an array)', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      const result = await collection.documents.index({ id: '1', title: 'Hello World' });
      expect(result.num_indexed).toBe(1);
      expect(result.num_failed).toBe(0);
      expect(result.results).toEqual([{ success: true, id: '1' }]);
    });
  });

  it('retrieves a document by id', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await collection.documents.index({ id: '1', title: 'Hello World' });

      const doc = await collection.documents.retrieve('1');
      expect(doc.id).toBe('1');
      expect(doc.title).toBe('Hello World');
    });
  });

  it('upserts by id: re-indexing overwrites the previous document', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await collection.documents.index({ id: '1', title: 'First title' });
      await collection.documents.index({ id: '1', title: 'Second title' });

      const doc = await collection.documents.retrieve('1');
      expect(doc.title).toBe('Second title');
    });
  });

  it('deletes a document, after which retrieval 404s', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await collection.documents.index({ id: '1', title: 'Hello World' });

      await collection.documents.delete('1');

      const failure = collection.documents.retrieve('1');
      await expect(failure).rejects.toBeInstanceOf(TachyonNotFoundError);
      await expect(failure).rejects.toMatchObject({ code: 'document_not_found', status: 404 });
    });
  });

  it('404s retrieving and deleting an id that never existed', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await expect(collection.documents.retrieve('never-existed')).rejects.toBeInstanceOf(TachyonNotFoundError);
      await expect(collection.documents.delete('never-existed')).rejects.toBeInstanceOf(TachyonNotFoundError);
    });
  });

  it('indexes a batch, reporting each success/failure without failing the whole request', async () => {
    const schema = {
      fields: [
        { name: 'title', type: 'text' as const },
        { name: 'price', type: 'int' as const },
      ],
    };
    await withCollection(client, schema, async (name) => {
      const collection = client.collection(name);
      const result = await collection.documents.index([
        { id: '1', price: 100 },
        { id: '2', price: 'not-a-number' },
        { id: '3', price: 300 },
      ]);

      expect(result.num_indexed).toBe(2);
      expect(result.num_failed).toBe(1);
      expect(result.results[0]).toMatchObject({ success: true, id: '1' });
      expect(result.results[1]).toMatchObject({ success: false, code: 'invalid_document' });
      expect(result.results[2]).toMatchObject({ success: true, id: '3' });

      await expect(collection.documents.retrieve('2')).rejects.toBeInstanceOf(TachyonNotFoundError);
    });
  });

  it('stores fields not declared in the schema, returned but not indexed', async () => {
    await withCollection(client, { fields: [{ name: 'title', type: 'text' }] }, async (name) => {
      const collection = client.collection(name);
      await collection.documents.index({ id: '1', title: 'Hello', undeclared_field: 'surprise' });

      const doc = await collection.documents.retrieve('1');
      expect(doc.undeclared_field).toBe('surprise');

      // Not indexed/filterable: searching for it by content should not match.
      const results = await collection.search({ q: 'surprise' });
      expect(results.hits).toHaveLength(0);
    });
  });
});
