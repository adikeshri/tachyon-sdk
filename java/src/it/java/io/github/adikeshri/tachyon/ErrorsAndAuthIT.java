package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.util.List;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class ErrorsAndAuthIT {

    @Test
    void noApiKeyIsRejectedWhenAuthIsEnabled() {
        TachyonAuthenticationException exception = assertThrows(TachyonAuthenticationException.class, () -> anonymousClient().collections.list());
        assertEquals("unauthorized", exception.getCode());
        assertEquals(401, exception.getStatusCode());
    }

    @Test
    void searchOnlyKeyCanRead() {
        List<CollectionInfo> list = searchOnlyClient().collections.list();
        assertNotNull(list);
    }

    @Test
    void searchOnlyKeyCannotWrite() {
        TachyonAuthorizationException exception = assertThrows(TachyonAuthorizationException.class, () ->
            searchOnlyClient().collections.create(CollectionSchema.builder(uniqueName("coll")).fields(List.of()).build()));
        assertEquals("forbidden", exception.getCode());
        assertEquals(403, exception.getStatusCode());
    }

    @Test
    void realNetworkFailureRaisesConnectionException() {
        // Nothing listens on port 1 (a reserved, unused TCP port), so this
        // exercises a real network failure rather than a mocked one.
        TachyonClient unreachable = TachyonClient.builder("http://127.0.0.1:1").timeout(Duration.ofSeconds(2)).build();
        assertThrows(TachyonConnectionException.class, unreachable::health);
    }

    @Test
    void adminKeyWorksEndToEnd() {
        List<CollectionInfo> list = adminClient().collections.list();
        assertNotNull(list);
    }
}
