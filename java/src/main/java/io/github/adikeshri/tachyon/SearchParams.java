package io.github.adikeshri.tachyon;

import java.util.List;

/** Parameters for {@link Collection#search}. Every field is optional. */
public record SearchParams(
    /** Query text. Null/empty matches everything. */
    String q,
    /** Fields to search. Defaults to every `text` field. */
    List<String> queryBy,
    /** Filter expression, e.g. "brand:=Logitech &amp;&amp; price:&lt;5000". */
    String filter,
    /** Sort expression, e.g. "_text_match:desc,price:asc". */
    String sort,
    /** Fields to facet on. */
    List<String> facet,
    /** Page size. Default 10, max 250. */
    Integer limit,
    /** Offset + limit must not exceed 10,000. */
    Integer offset,
    /** Prefix-match the final token. Default true. */
    Boolean prefix,
    /** Allow typo correction. Defaults to the collection's setting. */
    Boolean typoTolerance,
    /** ALL requires every token, ANY requires one. Default ALL. */
    MatchMode matchMode
) {
    private static final SearchParams EMPTY = builder().build();

    public static SearchParams empty() {
        return EMPTY;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static final class Builder {
        private String q;
        private List<String> queryBy;
        private String filter;
        private String sort;
        private List<String> facet;
        private Integer limit;
        private Integer offset;
        private Boolean prefix;
        private Boolean typoTolerance;
        private MatchMode matchMode;

        private Builder() {
        }

        public Builder q(String value) {
            this.q = value;
            return this;
        }

        public Builder queryBy(String... value) {
            this.queryBy = List.of(value);
            return this;
        }

        public Builder filter(String value) {
            this.filter = value;
            return this;
        }

        public Builder sort(String value) {
            this.sort = value;
            return this;
        }

        public Builder facet(String... value) {
            this.facet = List.of(value);
            return this;
        }

        public Builder limit(int value) {
            this.limit = value;
            return this;
        }

        public Builder offset(int value) {
            this.offset = value;
            return this;
        }

        public Builder prefix(boolean value) {
            this.prefix = value;
            return this;
        }

        public Builder typoTolerance(boolean value) {
            this.typoTolerance = value;
            return this;
        }

        public Builder matchMode(MatchMode value) {
            this.matchMode = value;
            return this;
        }

        public SearchParams build() {
            return new SearchParams(q, queryBy, filter, sort, facet, limit, offset, prefix, typoTolerance, matchMode);
        }
    }
}
