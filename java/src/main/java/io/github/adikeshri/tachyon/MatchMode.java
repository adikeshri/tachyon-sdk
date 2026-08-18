package io.github.adikeshri.tachyon;

/** Whether a search requires every token (default) or just one. */
public enum MatchMode {
    ALL("all"),
    ANY("any");

    private final String wireValue;

    MatchMode(String wireValue) {
        this.wireValue = wireValue;
    }

    @Override
    public String toString() {
        return wireValue;
    }
}
