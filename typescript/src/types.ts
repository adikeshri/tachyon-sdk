/** The five scalar field types Tachyon supports in a collection schema. */
export type FieldType = 'text' | 'keyword' | 'int' | 'float' | 'bool' | 'date';

export interface FieldSchema {
  name: string;
  type: FieldType;
  /** Build a facet column. Implies filterable. Default: false. */
  facet?: boolean;
  /** Allow filter expressions on this field. Default: false. */
  filter?: boolean;
  /** Allow sorting on this field. Default: false. */
  sort?: boolean;
  /** Include `text` content in the inverted index. Default: true. */
  index?: boolean;
  /** Allow documents that omit the field. Default: true. */
  optional?: boolean;
  /**
   * Per-field relevance multiplier. Defaults to 10 for a field named
   * `title`, 6 for `brand`, 2 for `description`, and 1 otherwise.
   */
  boost?: number;
}

export interface TypoToleranceConfig {
  enabled?: boolean;
  one_typo_min_len?: number;
  two_typo_min_len?: number;
  max_typos?: number;
}

export interface CollectionSchema {
  name: string;
  fields: FieldSchema[];
  typo_tolerance?: TypoToleranceConfig;
  default_sorting_field?: string;
}

/** A collection schema as returned by the server, with live counters. */
export interface CollectionInfo extends CollectionSchema {
  num_documents: number;
  num_segments: number;
}

/** A document is an arbitrary JSON object; `id` is always a string. */
export type TachyonDocument = Record<string, unknown>;

export interface DocumentIndexResult {
  success: boolean;
  id?: string;
  code?: string;
  error?: string;
}

export interface DocumentsIndexResponse {
  num_indexed: number;
  num_failed: number;
  results: DocumentIndexResult[];
}

export type MatchMode = 'all' | 'any';

export interface SearchParams {
  /** Query text. Empty (or omitted) matches everything. */
  q?: string;
  /** Fields to search. Defaults to every `text` field. */
  queryBy?: string | string[];
  /** Filter expression, e.g. `brand:=Logitech && price:<5000`. */
  filter?: string;
  /** Sort expression, e.g. `_text_match:desc,price:asc`. */
  sort?: string;
  /** Fields to facet on. */
  facet?: string | string[];
  /** Page size. Default 10, max 250. */
  limit?: number;
  /** Offset. `offset + limit` must not exceed 10,000. */
  offset?: number;
  /** Prefix-match the final token. Default true. */
  prefix?: boolean;
  /** Allow typo correction. Defaults to the collection's setting. */
  typoTolerance?: boolean;
  /** `all` requires every token, `any` requires one. Default `all`. */
  matchMode?: MatchMode;
}

export interface SearchHit<T extends TachyonDocument = TachyonDocument> {
  document: T;
  text_match: number;
}

export interface SearchResponse<T extends TachyonDocument = TachyonDocument> {
  found: number;
  /**
   * `false` once block-max WAND pruning has skipped part of a term's
   * postings for this query, at which point `found` (and facet counts)
   * become a lower bound rather than an exact count.
   */
  found_is_exact: boolean;
  search_time_ms: number;
  hits: SearchHit<T>[];
  facets?: Record<string, Record<string, number>>;
}

export interface SuggestParams {
  /** Text being typed; only the final token is completed. */
  q: string;
  /** Fields whose terms may be suggested. Defaults to every `text` field. */
  queryBy?: string | string[];
  /** Suggestions to return. Default 5, max 50. */
  limit?: number;
  /** Also suggest corrections. Defaults to the collection's setting. */
  typoTolerance?: boolean;
}

export interface Suggestion {
  text: string;
  count: number;
  typos: number;
}

export interface SuggestResponse {
  suggestions: Suggestion[];
  search_time_ms: number;
}

export interface AnalyticsQueryParams {
  collection?: string;
  /** Default 20, max 500. */
  limit?: number;
}

export interface AnalyticsQuery {
  query: string;
  collection: string;
  count: number;
  zero_result_count: number;
  last_result_count: number;
  avg_latency_ms: number;
  last_seen: number;
}

export interface AnalyticsQueriesResponse {
  queries: AnalyticsQuery[];
  tracked_queries: number;
  dropped_queries: number;
}

export interface AnalyticsLatencyResponse {
  count: number;
  mean_ms: number;
  p50_ms: number;
  p95_ms: number;
  p99_ms: number;
  max_ms: number;
  total_searches: number;
  uptime_seconds: number;
  queries_per_second: number;
}

export interface HealthResponse {
  ok: boolean;
  version: string;
  uptime_seconds: number;
  num_collections: number;
}
