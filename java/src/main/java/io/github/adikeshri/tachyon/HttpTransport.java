package io.github.adikeshri.tachyon;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.http.HttpTimeoutException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Map;

/** Thin JSON-over-HTTP client shared by every resource in the SDK. */
final class HttpTransport {
    private final HttpClient httpClient;
    private final String baseUrl;
    private final String apiKey;
    private final Map<String, String> headers;
    private final Duration timeout;
    private final ObjectMapper mapper;

    HttpTransport(HttpClient httpClient, String baseUrl, String apiKey, Map<String, String> headers, Duration timeout, ObjectMapper mapper) {
        this.httpClient = httpClient;
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.apiKey = apiKey;
        this.headers = headers;
        this.timeout = timeout;
        this.mapper = mapper;
    }

    <T> T request(String method, String path, Map<String, String> query, Object body, Class<T> type) {
        RawResponse response = send(method, path, query, body);
        if (response.status() >= 400) {
            throw errorFromBody(response.status(), response.body());
        }
        try {
            return mapper.readValue(response.body(), type);
        } catch (JsonProcessingException e) {
            throw new TachyonConnectionException("Tachyon returned a response that could not be decoded: " + e.getMessage(), e);
        }
    }

    <T> T request(String method, String path, Map<String, String> query, Object body, TypeReference<T> typeRef) {
        RawResponse response = send(method, path, query, body);
        if (response.status() >= 400) {
            throw errorFromBody(response.status(), response.body());
        }
        try {
            return mapper.readValue(response.body(), typeRef);
        } catch (JsonProcessingException e) {
            throw new TachyonConnectionException("Tachyon returned a response that could not be decoded: " + e.getMessage(), e);
        }
    }

    void requestVoid(String method, String path) {
        RawResponse response = send(method, path, null, null);
        if (response.status() >= 400) {
            throw errorFromBody(response.status(), response.body());
        }
    }

    String requestText(String method, String path) {
        RawResponse response = send(method, path, null, null);
        if (response.status() >= 400) {
            throw errorFromBody(response.status(), response.body());
        }
        return response.body();
    }

    private record RawResponse(int status, String body) {
    }

    private RawResponse send(String method, String path, Map<String, String> query, Object body) {
        String url = buildUrl(path, query);
        HttpRequest.Builder builder = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .timeout(timeout)
            .header("Accept", "application/json");
        if (apiKey != null && !apiKey.isEmpty()) {
            builder.header("X-TACHYON-API-KEY", apiKey);
        }
        if (headers != null) {
            for (Map.Entry<String, String> entry : headers.entrySet()) {
                builder.header(entry.getKey(), entry.getValue());
            }
        }
        if (body != null) {
            String json;
            try {
                json = mapper.writeValueAsString(body);
            } catch (JsonProcessingException e) {
                throw new IllegalArgumentException("Failed to encode request body: " + e.getMessage(), e);
            }
            builder.header("Content-Type", "application/json")
                .method(method, HttpRequest.BodyPublishers.ofString(json, StandardCharsets.UTF_8));
        } else {
            builder.method(method, HttpRequest.BodyPublishers.noBody());
        }

        HttpResponse<String> response;
        try {
            response = httpClient.send(builder.build(), HttpResponse.BodyHandlers.ofString());
        } catch (HttpTimeoutException e) {
            throw new TachyonTimeoutException("Request to " + url + " timed out: " + e.getMessage(), e);
        } catch (IOException e) {
            throw new TachyonConnectionException("Failed to reach Tachyon at " + url + ": " + e.getMessage(), e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new TachyonConnectionException("Request to " + url + " was interrupted: " + e.getMessage(), e);
        }

        return new RawResponse(response.statusCode(), response.body());
    }

    private String buildUrl(String path, Map<String, String> query) {
        StringBuilder url = new StringBuilder(baseUrl).append(path);
        if (query != null && !query.isEmpty()) {
            StringBuilder qs = new StringBuilder();
            for (Map.Entry<String, String> entry : query.entrySet()) {
                if (entry.getValue() == null) {
                    continue;
                }
                if (!qs.isEmpty()) {
                    qs.append('&');
                }
                qs.append(URLEncoder.encode(entry.getKey(), StandardCharsets.UTF_8))
                    .append('=')
                    .append(URLEncoder.encode(entry.getValue(), StandardCharsets.UTF_8));
            }
            if (!qs.isEmpty()) {
                url.append('?').append(qs);
            }
        }
        return url.toString();
    }

    private TachyonException errorFromBody(int status, String body) {
        try {
            JsonNode root = mapper.readTree(body);
            JsonNode error = root == null ? null : root.get("error");
            if (error != null) {
                String code = error.has("code") ? error.get("code").asText() : null;
                String message = error.has("message") ? error.get("message").asText() : body;
                if (code != null && !code.isEmpty()) {
                    return TachyonExceptionFactory.fromResponse(status, message, code);
                }
            }
        } catch (JsonProcessingException e) {
            // Falls through to the generic mapping below.
        }
        return TachyonExceptionFactory.fromResponse(status, body, "internal_error");
    }
}
