package io.github.adikeshri.tachyon;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * A document is an arbitrary JSON object; {@code id} is always a string.
 * A plain {@link Map} in every way that matters — {@link #set} just adds a
 * fluent chain for building one inline: {@code new Document().set("id",
 * "1").set("title", "Wireless Mouse")}.
 */
public final class Document extends LinkedHashMap<String, Object> {
    public Document() {
        super();
    }

    public Document(Map<String, ?> initial) {
        super(initial);
    }

    public Document set(String key, Object value) {
        put(key, value);
        return this;
    }
}
