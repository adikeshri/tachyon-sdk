import type { HttpClient } from './http.js';
import type { DocumentsIndexResponse, TachyonDocument } from './types.js';

/** `/collections/{name}/documents` — index, fetch, and delete documents. */
export class DocumentsApi<T extends TachyonDocument = TachyonDocument> {
  constructor(
    private readonly http: HttpClient,
    private readonly collectionName: string,
  ) {}

  /**
   * `POST /collections/{name}/documents`. Accepts one document or an array;
   * documents are upserted by `id`. Individual documents can fail without
   * failing their neighbours — check `num_failed` and `results`.
   */
  index(documents: T | T[]): Promise<DocumentsIndexResponse> {
    return this.http.request<DocumentsIndexResponse>(
      'POST',
      `/collections/${encodeURIComponent(this.collectionName)}/documents`,
      { body: documents },
    );
  }

  /** `GET /collections/{name}/documents/{id}`. */
  retrieve(id: string): Promise<T> {
    return this.http.request<T>(
      'GET',
      `/collections/${encodeURIComponent(this.collectionName)}/documents/${encodeURIComponent(id)}`,
    );
  }

  /** `DELETE /collections/{name}/documents/{id}`. */
  async delete(id: string): Promise<void> {
    await this.http.request<void>(
      'DELETE',
      `/collections/${encodeURIComponent(this.collectionName)}/documents/${encodeURIComponent(id)}`,
    );
  }
}
