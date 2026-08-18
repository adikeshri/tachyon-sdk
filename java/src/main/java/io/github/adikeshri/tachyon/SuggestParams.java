package io.github.adikeshri.tachyon;

import java.util.List;

/** Parameters for {@link Collection#suggest}. */
public record SuggestParams(
    /** Text being typed; only the final token is completed. */
    String q,
    /** Fields whose terms may be suggested. Defaults to every `text` field. */
    List<String> queryBy,
    /** Suggestions to return. Default 5, max 50. */
    Integer limit,
    /** Also suggest corrections. Defaults to the collection's setting. */
    Boolean typoTolerance
) {
    public static Builder builder(String q) {
        return new Builder(q);
    }

    public static final class Builder {
        private final String q;
        private List<String> queryBy;
        private Integer limit;
        private Boolean typoTolerance;

        private Builder(String q) {
            this.q = q;
        }

        public Builder queryBy(String... value) {
            this.queryBy = List.of(value);
            return this;
        }

        public Builder limit(int value) {
            this.limit = value;
            return this;
        }

        public Builder typoTolerance(boolean value) {
            this.typoTolerance = value;
            return this;
        }

        public SuggestParams build() {
            return new SuggestParams(q, queryBy, limit, typoTolerance);
        }
    }
}
