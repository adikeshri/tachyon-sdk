import type { HttpClient } from './http.js';
import type { CollectionInfo, CollectionSchema } from './types.js';

/** `/collections` — create, list, and remove collections. */
export class CollectionsApi {
  constructor(private readonly http: HttpClient) {}

  /** `POST /collections`. Field types are immutable after creation. */
  create(schema: CollectionSchema): Promise<CollectionInfo> {
    return this.http.request<CollectionInfo>('POST', '/collections', { body: schema });
  }

  /** `GET /collections`. */
  list(): Promise<CollectionInfo[]> {
    return this.http.request<CollectionInfo[]>('GET', '/collections');
  }

  /** `GET /collections/{name}`. */
  retrieve(name: string): Promise<CollectionInfo> {
    return this.http.request<CollectionInfo>('GET', `/collections/${encodeURIComponent(name)}`);
  }

  /** `DELETE /collections/{name}`. Removes the collection and all its data. */
  async delete(name: string): Promise<void> {
    await this.http.request<void>('DELETE', `/collections/${encodeURIComponent(name)}`);
  }
}
