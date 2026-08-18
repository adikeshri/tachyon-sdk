import { DocumentsApi } from './documents.js';
import type { HttpClient } from './http.js';
import type {
  SearchParams,
  SearchResponse,
  SuggestParams,
  SuggestResponse,
  TachyonDocument,
} from './types.js';

function joinIfArray(value: string | string[] | undefined): string | undefined {
  return Array.isArray(value) ? value.join(',') : value;
}

function serializeSearchParams(params: SearchParams): Record<string, unknown> {
  return {
    q: params.q,
    query_by: joinIfArray(params.queryBy),
    filter: params.filter,
    sort: params.sort,
    facet: joinIfArray(params.facet),
    limit: params.limit,
    offset: params.offset,
    prefix: params.prefix,
    typo_tolerance: params.typoTolerance,
    match_mode: params.matchMode,
  };
}

function serializeSuggestParams(params: SuggestParams): Record<string, unknown> {
  return {
    q: params.q,
    query_by: joinIfArray(params.queryBy),
    limit: params.limit,
    typo_tolerance: params.typoTolerance,
  };
}

/**
 * A handle scoped to one collection. Get one via `client.collection(name)`;
 * it does not verify the collection exists until you make a request.
 */
export class Collection<T extends TachyonDocument = TachyonDocument> {
  readonly documents: DocumentsApi<T>;

  constructor(
    private readonly http: HttpClient,
    readonly name: string,
  ) {
    this.documents = new DocumentsApi<T>(http, name);
  }

  /** `GET /collections/{name}/search`. */
  search(params: SearchParams = {}): Promise<SearchResponse<T>> {
    return this.http.request<SearchResponse<T>>(
      'GET',
      `/collections/${encodeURIComponent(this.name)}/search`,
      { query: serializeSearchParams(params) },
    );
  }

  /** `GET /collections/{name}/suggest`. */
  suggest(params: SuggestParams): Promise<SuggestResponse> {
    return this.http.request<SuggestResponse>(
      'GET',
      `/collections/${encodeURIComponent(this.name)}/suggest`,
      { query: serializeSuggestParams(params) },
    );
  }
}
