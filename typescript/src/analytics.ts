import type { HttpClient } from './http.js';
import type { AnalyticsLatencyResponse, AnalyticsQueriesResponse, AnalyticsQueryParams } from './types.js';

/** `/analytics/*` — recorded automatically from search traffic, in memory only. */
export class AnalyticsApi {
  constructor(private readonly http: HttpClient) {}

  /** `GET /analytics/top`. */
  top(params: AnalyticsQueryParams = {}): Promise<AnalyticsQueriesResponse> {
    return this.http.request<AnalyticsQueriesResponse>('GET', '/analytics/top', {
      query: { collection: params.collection, limit: params.limit },
    });
  }

  /** `GET /analytics/zero-results`. Ranks by how often a query came back empty. */
  zeroResults(params: AnalyticsQueryParams = {}): Promise<AnalyticsQueriesResponse> {
    return this.http.request<AnalyticsQueriesResponse>('GET', '/analytics/zero-results', {
      query: { collection: params.collection, limit: params.limit },
    });
  }

  /** `GET /analytics/latency`. */
  latency(): Promise<AnalyticsLatencyResponse> {
    return this.http.request<AnalyticsLatencyResponse>('GET', '/analytics/latency');
  }
}
