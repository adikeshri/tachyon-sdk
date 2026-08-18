package io.github.adikeshri.tachyon;

/** Parameters for {@link AnalyticsResource#top} and {@link AnalyticsResource#zeroResults}. */
public record AnalyticsQueryParams(String collection, /** Default 20, max 500. */ Integer limit) {
    private static final AnalyticsQueryParams EMPTY = new AnalyticsQueryParams(null, null);

    public static AnalyticsQueryParams empty() {
        return EMPTY;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static final class Builder {
        private String collection;
        private Integer limit;

        private Builder() {
        }

        public Builder collection(String value) {
            this.collection = value;
            return this;
        }

        public Builder limit(int value) {
            this.limit = value;
            return this;
        }

        public AnalyticsQueryParams build() {
            return new AnalyticsQueryParams(collection, limit);
        }
    }
}
