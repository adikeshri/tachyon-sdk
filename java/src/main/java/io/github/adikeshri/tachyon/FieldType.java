package io.github.adikeshri.tachyon;

/** One of the scalar field types Tachyon supports in a collection schema. */
public enum FieldType {
    TEXT("text"),
    KEYWORD("keyword"),
    INT("int"),
    FLOAT("float"),
    BOOL("bool"),
    DATE("date");

    private final String wireValue;

    FieldType(String wireValue) {
        this.wireValue = wireValue;
    }

    @Override
    public String toString() {
        return wireValue;
    }
}
