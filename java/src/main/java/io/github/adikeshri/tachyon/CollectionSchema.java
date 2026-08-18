package io.github.adikeshri.tachyon;

import java.util.List;

public record CollectionSchema(
    String name,
    List<FieldSchema> fields,
    TypoToleranceConfig typoTolerance,
    String defaultSortingField
) {
    public static Builder builder(String name) {
        return new Builder(name);
    }

    public static final class Builder {
        private final String name;
        private List<FieldSchema> fields = List.of();
        private TypoToleranceConfig typoTolerance;
        private String defaultSortingField;

        private Builder(String name) {
            this.name = name;
        }

        public Builder fields(List<FieldSchema> value) {
            this.fields = value;
            return this;
        }

        public Builder fields(FieldSchema... value) {
            this.fields = List.of(value);
            return this;
        }

        public Builder typoTolerance(TypoToleranceConfig value) {
            this.typoTolerance = value;
            return this;
        }

        public Builder defaultSortingField(String value) {
            this.defaultSortingField = value;
            return this;
        }

        public CollectionSchema build() {
            return new CollectionSchema(name, fields, typoTolerance, defaultSortingField);
        }
    }
}
