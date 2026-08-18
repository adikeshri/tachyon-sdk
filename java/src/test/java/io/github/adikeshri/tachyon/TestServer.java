package io.github.adikeshri.tachyon;

import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Function;

/**
 * A real local loopback HTTP server (java.net.http.HttpClient isn't easily
 * mockable directly) — the same idea as Go's httptest.NewServer. Records
 * every request it receives and replies according to a handler function.
 */
final class TestServer implements AutoCloseable {
    private final HttpServer server;
    final List<CapturedRequest> requests = new ArrayList<>();
    final String baseUrl;

    record CapturedRequest(String method, String path, Map<String, String> query, Map<String, List<String>> headers, String body) {
        String header(String name) {
            List<String> values = headers.get(name);
            return values == null || values.isEmpty() ? null : values.get(0);
        }
    }

    record FakeResponse(int status, String body, String contentType) {
        static FakeResponse json(int status, String body) {
            return new FakeResponse(status, body, "application/json");
        }

        static FakeResponse text(int status, String body) {
            return new FakeResponse(status, body, "text/plain");
        }

        static FakeResponse noContent() {
            return new FakeResponse(204, "", "application/json");
        }
    }

    TestServer(Function<CapturedRequest, FakeResponse> handler) throws IOException {
        server = HttpServer.create(new InetSocketAddress("localhost", 0), 0);
        server.createContext("/", exchange -> {
            try {
                String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
                CapturedRequest captured = new CapturedRequest(
                    exchange.getRequestMethod(),
                    // getRawPath(), not getPath(): URI.getPath() is auto-decoded
                    // (so "my%20products" would already read back as "my
                    // products"), which would hide an encoding bug in the
                    // client. getRawPath() reflects what was actually sent.
                    exchange.getRequestURI().getRawPath(),
                    parseQuery(exchange.getRequestURI().getRawQuery()),
                    exchange.getRequestHeaders(),
                    body);
                requests.add(captured);

                FakeResponse response = handler.apply(captured);
                byte[] responseBytes = response.body().getBytes(StandardCharsets.UTF_8);
                exchange.getResponseHeaders().add("Content-Type", response.contentType());
                if (response.status() == 204) {
                    exchange.sendResponseHeaders(204, -1);
                } else {
                    exchange.sendResponseHeaders(response.status(), responseBytes.length);
                    try (OutputStream os = exchange.getResponseBody()) {
                        os.write(responseBytes);
                    }
                }
            } finally {
                exchange.close();
            }
        });
        server.start();
        baseUrl = "http://localhost:" + server.getAddress().getPort();
    }

    private static Map<String, String> parseQuery(String raw) {
        Map<String, String> result = new LinkedHashMap<>();
        if (raw == null || raw.isEmpty()) {
            return result;
        }
        for (String pair : raw.split("&")) {
            int eq = pair.indexOf('=');
            String key = eq >= 0 ? pair.substring(0, eq) : pair;
            String value = eq >= 0 ? pair.substring(eq + 1) : "";
            result.put(URLDecoder.decode(key, StandardCharsets.UTF_8), URLDecoder.decode(value, StandardCharsets.UTF_8));
        }
        return result;
    }

    @Override
    public void close() {
        server.stop(0);
    }
}
