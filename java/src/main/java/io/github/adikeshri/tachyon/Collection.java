package io.github.adikeshri.tachyon;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * A handle scoped to one collection. Get one via {@link TachyonClient#collection};
 * it does not verify the collection exists until you make a request.
 */
public final class Collection {
    private final HttpTransport transport;
    private final String name;
    private final DocumentsResource documents;

    Collection(HttpTransport transport, String name) {
        this.transport = transport;
        this.name = name;
        this.documents = new DocumentsResource(transport, name);
    }

    public String name() {
        return name;
    }

    public DocumentsResource documents() {
        return documents;
    }

    /** {@code GET /collections/{name}/search}. */
    public SearchResponse search() {
        return search(SearchParams.empty());
    }

    /** {@code GET /collections/{name}/search}. */
    public SearchResponse search(SearchParams params) {
        Map<String, String> query = new LinkedHashMap<>();
        query.put("q", params.q());
        query.put("query_by", join(params.queryBy()));
        query.put("filter", params.filter());
        query.put("sort", params.sort());
        query.put("facet", join(params.facet()));
        query.put("limit", params.limit() == null ? null : params.limit().toString());
        query.put("offset", params.offset() == null ? null : params.offset().toString());
        query.put("prefix", formatBool(params.prefix()));
        query.put("typo_tolerance", formatBool(params.typoTolerance()));
        query.put("match_mode", params.matchMode() == null ? null : params.matchMode().toString());

        String path = "/collections/" + UrlEncoding.pathSegment(name) + "/search";
        return transport.request("GET", path, query, null, SearchResponse.class);
    }

    /** {@code GET /collections/{name}/suggest}. */
    public SuggestResponse suggest(SuggestParams params) {
        Map<String, String> query = new LinkedHashMap<>();
        query.put("q", params.q());
        query.put("query_by", join(params.queryBy()));
        query.put("limit", params.limit() == null ? null : params.limit().toString());
        query.put("typo_tolerance", formatBool(params.typoTolerance()));

        String path = "/collections/" + UrlEncoding.pathSegment(name) + "/suggest";
        return transport.request("GET", path, query, null, SuggestResponse.class);
    }

    private static String join(List<String> values) {
        return values == null || values.isEmpty() ? null : String.join(",", values);
    }

    private static String formatBool(Boolean value) {
        return value == null ? null : value.toString();
    }
}
