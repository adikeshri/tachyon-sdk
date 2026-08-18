package io.github.adikeshri.tachyon;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.SerializationFeature;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.Map;

/**
 * Client for a single Tachyon server.
 *
 * <pre>{@code
 * var client = TachyonClient.builder("http://localhost:8108").apiKey("my-admin-key").build();
 * client.collections.create(CollectionSchema.builder("products")
 *     .fields(FieldSchema.builder("title", FieldType.TEXT).build())
 *     .build());
 * client.collection("products").documents().index(new Document().set("id", "1").set("title", "Wireless Mouse"));
 * var results = client.collection("products").search(SearchParams.builder().q("wireless mouse").build());
 * }</pre>
 */
public final class TachyonClient {
    public final CollectionsResource collections;
    public final AnalyticsResource analytics;

    private final HttpTransport transport;

    private TachyonClient(Builder builder) {
        ObjectMapper mapper = new ObjectMapper();
        mapper.setPropertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE);
        mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);
        mapper.configure(SerializationFeature.WRITE_ENUMS_USING_TO_STRING, true);
        mapper.configure(DeserializationFeature.READ_ENUMS_USING_TO_STRING, true);
        // Tachyon's responses carry fields this SDK doesn't model yet (e.g.
        // created_at on a collection) — ignore rather than fail on those.
        mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);

        HttpClient httpClient = builder.httpClient != null ? builder.httpClient : HttpClient.newHttpClient();

        this.transport = new HttpTransport(httpClient, builder.url, builder.apiKey, builder.headers, builder.timeout, mapper);
        this.collections = new CollectionsResource(transport);
        this.analytics = new AnalyticsResource(transport);
    }

    /** @param url Base URL of the Tachyon server, e.g. "http://localhost:8108". */
    public static Builder builder(String url) {
        return new Builder(url);
    }

    /** Get a handle scoped to one collection, for documents/search/suggest. */
    public Collection collection(String name) {
        return new Collection(transport, name);
    }

    /** {@code GET /health}. Always reachable without an API key. */
    public HealthResponse health() {
        return transport.request("GET", "/health", null, null, HealthResponse.class);
    }

    /** {@code GET /metrics}. Prometheus exposition format, returned as plain text. */
    public String metrics() {
        return transport.requestText("GET", "/metrics");
    }

    public static final class Builder {
        private final String url;
        private String apiKey;
        private Duration timeout = Duration.ofSeconds(15);
        private Map<String, String> headers;
        private HttpClient httpClient;

        private Builder(String url) {
            this.url = url;
        }

        /** Sent as X-TACHYON-API-KEY. Use an admin key for writes, a search key for read-only access. */
        public Builder apiKey(String value) {
            this.apiKey = value;
            return this;
        }

        /** Per-request timeout. Default 15s. Ignored if {@link #httpClient} is supplied. */
        public Builder timeout(Duration value) {
            this.timeout = value;
            return this;
        }

        /** Extra headers merged into every request. */
        public Builder headers(Map<String, String> value) {
            this.headers = value;
            return this;
        }

        /** Override the {@link HttpClient} (mainly for testing, or to share connection pooling). */
        public Builder httpClient(HttpClient value) {
            this.httpClient = value;
            return this;
        }

        public TachyonClient build() {
            return new TachyonClient(this);
        }
    }
}
