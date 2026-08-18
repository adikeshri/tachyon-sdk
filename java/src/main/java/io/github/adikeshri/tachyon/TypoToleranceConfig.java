package io.github.adikeshri.tachyon;

public record TypoToleranceConfig(Boolean enabled, Integer oneTypoMinLen, Integer twoTypoMinLen, Integer maxTypos) {
    public static Builder builder() {
        return new Builder();
    }

    public static final class Builder {
        private Boolean enabled;
        private Integer oneTypoMinLen;
        private Integer twoTypoMinLen;
        private Integer maxTypos;

        private Builder() {
        }

        public Builder enabled(boolean value) {
            this.enabled = value;
            return this;
        }

        public Builder oneTypoMinLen(int value) {
            this.oneTypoMinLen = value;
            return this;
        }

        public Builder twoTypoMinLen(int value) {
            this.twoTypoMinLen = value;
            return this;
        }

        public Builder maxTypos(int value) {
            this.maxTypos = value;
            return this;
        }

        public TypoToleranceConfig build() {
            return new TypoToleranceConfig(enabled, oneTypoMinLen, twoTypoMinLen, maxTypos);
        }
    }
}
