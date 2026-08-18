package io.github.adikeshri.tachyon;

/** One field in a collection schema. See {@link CollectionSchema}. */
public record FieldSchema(
    String name,
    FieldType type,
    /** Build a facet column. Implies filterable. Default: false. */
    Boolean facet,
    /** Allow filter expressions on this field. Default: false. */
    Boolean filter,
    /** Allow sorting on this field. Default: false. */
    Boolean sort,
    /** Include `text` content in the inverted index. Default: true. */
    Boolean index,
    /** Allow documents that omit the field. Default: true. */
    Boolean optional,
    /** Per-field relevance multiplier. Defaults to 10/6/2/1 by field name. */
    Double boost
) {
    public static Builder builder(String name, FieldType type) {
        return new Builder(name, type);
    }

    public static final class Builder {
        private final String name;
        private final FieldType type;
        private Boolean facet;
        private Boolean filter;
        private Boolean sort;
        private Boolean index;
        private Boolean optional;
        private Double boost;

        private Builder(String name, FieldType type) {
            this.name = name;
            this.type = type;
        }

        public Builder facet(boolean value) {
            this.facet = value;
            return this;
        }

        public Builder filter(boolean value) {
            this.filter = value;
            return this;
        }

        public Builder sort(boolean value) {
            this.sort = value;
            return this;
        }

        public Builder index(boolean value) {
            this.index = value;
            return this;
        }

        public Builder optional(boolean value) {
            this.optional = value;
            return this;
        }

        public Builder boost(double value) {
            this.boost = value;
            return this;
        }

        public FieldSchema build() {
            return new FieldSchema(name, type, facet, filter, sort, index, optional, boost);
        }
    }
}
