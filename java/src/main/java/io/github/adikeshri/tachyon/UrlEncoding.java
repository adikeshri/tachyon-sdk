package io.github.adikeshri.tachyon;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

final class UrlEncoding {
    private UrlEncoding() {
    }

    /**
     * Percent-encodes a single URL path segment. {@link URLEncoder} encodes
     * for {@code application/x-www-form-urlencoded} (query strings), which
     * represents a space as {@code +} — literal in a path, not a decoded
     * space — so that gets converted to {@code %20} afterward.
     */
    static String pathSegment(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8).replace("+", "%20");
    }
}
