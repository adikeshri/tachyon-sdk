import { describe, expect, it } from 'vitest';
import { TachyonConflictError, TachyonError, TachyonNotFoundError, TachyonRequestError } from '../../src/errors.js';
import { adminClient, uniqueName } from './support.js';

const client = adminClient();

describe('collections (integration)', () => {
  it('creates a collection and fills in field/collection defaults', async () => {
    const name = uniqueName('coll');
    try {
      const created = await client.collections.create({
        name,
        fields: [
          { name: 'title', type: 'text' },
          { name: 'brand', type: 'keyword', facet: true },
          { name: 'price', type: 'int', filter: true, sort: true },
        ],
      });
      expect(created.name).toBe(name);
      expect(created.num_documents).toBe(0);
      expect(created.num_segments).toBe(0);
      expect(created.fields).toHaveLength(3);

      const price = created.fields.find((f) => f.name === 'price');
      expect(price).toMatchObject({ filter: true, sort: true, optional: true, index: true });

      const brand = created.fields.find((f) => f.name === 'brand');
      expect(brand?.facet).toBe(true);
    } finally {
      await client.collections.delete(name).catch(() => {});
    }
  });

  it('rejects a duplicate collection name with collection_exists / 409', async () => {
    const name = uniqueName('coll');
    await client.collections.create({ name, fields: [{ name: 'title', type: 'text' }] });
    try {
      const failure = client.collections.create({ name, fields: [{ name: 'title', type: 'text' }] });
      await expect(failure).rejects.toBeInstanceOf(TachyonConflictError);
      await expect(failure).rejects.toMatchObject({ code: 'collection_exists', status: 409 });
    } finally {
      await client.collections.delete(name).catch(() => {});
    }
  });

  it('lists collections, including one just created', async () => {
    const name = uniqueName('coll');
    await client.collections.create({ name, fields: [{ name: 'title', type: 'text' }] });
    try {
      const list = await client.collections.list();
      expect(list.some((c) => c.name === name)).toBe(true);
    } finally {
      await client.collections.delete(name).catch(() => {});
    }
  });

  it('retrieves a single collection by name', async () => {
    const name = uniqueName('coll');
    await client.collections.create({ name, fields: [{ name: 'title', type: 'text' }] });
    try {
      const retrieved = await client.collections.retrieve(name);
      expect(retrieved.name).toBe(name);
    } finally {
      await client.collections.delete(name).catch(() => {});
    }
  });

  it('404s retrieving an unknown collection', async () => {
    const failure = client.collections.retrieve(uniqueName('missing'));
    await expect(failure).rejects.toBeInstanceOf(TachyonNotFoundError);
    await expect(failure).rejects.toMatchObject({ code: 'collection_not_found', status: 404 });
  });

  it('deletes a collection so it stops appearing in retrieve/list', async () => {
    const name = uniqueName('coll');
    await client.collections.create({ name, fields: [{ name: 'title', type: 'text' }] });

    await client.collections.delete(name);

    await expect(client.collections.retrieve(name)).rejects.toBeInstanceOf(TachyonNotFoundError);
    const list = await client.collections.list();
    expect(list.some((c) => c.name === name)).toBe(false);
  });

  it('rejects a semantically invalid schema with invalid_schema / 400', async () => {
    // Tachyon requires at least one indexed `text` field so the collection is searchable.
    const failure = client.collections.create({
      name: uniqueName('coll'),
      fields: [{ name: 'price', type: 'int' }],
    });
    await expect(failure).rejects.toBeInstanceOf(TachyonRequestError);
    await expect(failure).rejects.toMatchObject({ code: 'invalid_schema', status: 400 });
  });

  it('rejects a field with an unrecognized type at the JSON layer (422, no error envelope)', async () => {
    // An unknown `type` enum value fails Axum's JSON extraction before the request reaches
    // the handler, so the server never gets to wrap it in the documented {error:{code,...}}
    // shape — it comes back as plain text. The SDK still surfaces status + message for it,
    // just without a stable `code`, via the base TachyonError rather than a subclass.
    const failure = client.collections.create({
      name: uniqueName('coll'),
      fields: [{ name: 'x', type: 'not-a-real-type' as never }],
    });
    await expect(failure).rejects.toBeInstanceOf(TachyonError);
    await expect(failure).rejects.toMatchObject({ status: 422 });
  });

  it('applies collection-level typo_tolerance and default_sorting_field options', async () => {
    const name = uniqueName('coll');
    try {
      const created = await client.collections.create({
        name,
        fields: [
          { name: 'title', type: 'text' },
          { name: 'popularity', type: 'int', sort: true },
        ],
        typo_tolerance: { enabled: false },
        default_sorting_field: 'popularity',
      });
      expect(created.typo_tolerance?.enabled).toBe(false);
      expect(created.default_sorting_field).toBe('popularity');
    } finally {
      await client.collections.delete(name).catch(() => {});
    }
  });
});
