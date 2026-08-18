package io.github.adikeshri.tachyon;

import java.util.LinkedHashMap;
import java.util.Map;

/** {@code /analytics/*} — recorded automatically from search traffic, in memory only. */
public final class AnalyticsResource {
    private final HttpTransport transport;

    AnalyticsResource(HttpTransport transport) {
        this.transport = transport;
    }

    /** {@code GET /analytics/top}. */
    public AnalyticsQueriesResponse top() {
        return top(AnalyticsQueryParams.empty());
    }

    /** {@code GET /analytics/top}. */
    public AnalyticsQueriesResponse top(AnalyticsQueryParams params) {
        return transport.request("GET", "/analytics/top", toQuery(params), null, AnalyticsQueriesResponse.class);
    }

    /** {@code GET /analytics/zero-results}. Ranks by how often a query came back empty. */
    public AnalyticsQueriesResponse zeroResults() {
        return zeroResults(AnalyticsQueryParams.empty());
    }

    /** {@code GET /analytics/zero-results}. Ranks by how often a query came back empty. */
    public AnalyticsQueriesResponse zeroResults(AnalyticsQueryParams params) {
        return transport.request("GET", "/analytics/zero-results", toQuery(params), null, AnalyticsQueriesResponse.class);
    }

    /** {@code GET /analytics/latency}. */
    public AnalyticsLatencyResponse latency() {
        return transport.request("GET", "/analytics/latency", null, null, AnalyticsLatencyResponse.class);
    }

    private static Map<String, String> toQuery(AnalyticsQueryParams params) {
        Map<String, String> query = new LinkedHashMap<>();
        query.put("collection", params.collection());
        query.put("limit", params.limit() == null ? null : params.limit().toString());
        return query;
    }
}
