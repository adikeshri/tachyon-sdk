package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.io.IOException;
import java.time.Duration;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.*;

class ErrorsTest {

    static Stream<Arguments> statusToException() {
        return Stream.of(
            Arguments.of(400, "invalid_query", TachyonRequestException.class),
            Arguments.of(401, "unauthorized", TachyonAuthenticationException.class),
            Arguments.of(403, "forbidden", TachyonAuthorizationException.class),
            Arguments.of(404, "collection_not_found", TachyonNotFoundException.class),
            Arguments.of(409, "collection_exists", TachyonConflictException.class),
            Arguments.of(500, "internal_error", TachyonServerException.class));
    }

    @ParameterizedTest
    @MethodSource("statusToException")
    void mapsHttpStatusToExceptionType(int status, String code, Class<? extends TachyonException> expectedType) throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(status,
            "{\"error\":{\"code\":\"" + code + "\",\"message\":\"boom: " + code + "\"}}"))) {
            TachyonClient client = TestClients.forServer(server);

            TachyonException exception = assertThrows(expectedType, () -> client.collections.retrieve("products"));

            assertEquals(code, exception.getCode());
            assertEquals(status, exception.getStatusCode());
            assertEquals("boom: " + code, exception.getMessage());
        }
    }

    @Test
    void wrapsNetworkFailuresInConnectionException() {
        // Nothing listens on port 1 (a reserved, unused TCP port), so this
        // exercises a real network failure rather than a mocked one.
        TachyonClient client = TachyonClient.builder("http://127.0.0.1:1").timeout(Duration.ofSeconds(2)).build();

        assertThrows(TachyonConnectionException.class, client::health);
    }

    @Test
    void raisesTimeoutExceptionWhenServerIsSlow() throws IOException {
        try (TestServer server = new TestServer(req -> {
            try {
                Thread.sleep(200);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            return TestServer.FakeResponse.json(200, "{\"ok\":true,\"version\":\"x\",\"uptime_seconds\":0,\"num_collections\":0}");
        })) {
            TachyonClient client = TachyonClient.builder(server.baseUrl).timeout(Duration.ofMillis(10)).build();

            assertThrows(TachyonTimeoutException.class, client::health);
        }
    }

    @Test
    void fallsBackToGenericCodeWhenBodyIsNotTheDocumentedShape() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.text(500, "upstream exploded"))) {
            TachyonClient client = TestClients.forServer(server);

            TachyonException exception = assertThrows(TachyonServerException.class, client::health);

            assertEquals("internal_error", exception.getCode());
            assertEquals(500, exception.getStatusCode());
        }
    }
}
