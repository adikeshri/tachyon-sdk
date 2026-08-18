package io.github.adikeshri.tachyon;

import java.util.UUID;
import java.util.function.Consumer;

final class TestSupport {
    private TestSupport() {
    }

    static String baseUrl() {
        String value = System.getenv("TACHYON_URL");
        return value == null || value.isEmpty() ? "http://localhost:8108" : value;
    }

    static String adminKey() {
        String value = System.getenv("TACHYON_ADMIN_KEY");
        return value == null || value.isEmpty() ? "admin-key" : value;
    }

    static String searchKey() {
        String value = System.getenv("TACHYON_SEARCH_KEY");
        return value == null || value.isEmpty() ? "search-key" : value;
    }

    static TachyonClient adminClient() {
        return TachyonClient.builder(baseUrl()).apiKey(adminKey()).build();
    }

    static TachyonClient searchOnlyClient() {
        return TachyonClient.builder(baseUrl()).apiKey(searchKey()).build();
    }

    static TachyonClient anonymousClient() {
        return TachyonClient.builder(baseUrl()).build();
    }

    static String uniqueName(String prefix) {
        return prefix + "-" + UUID.randomUUID();
    }

    /**
     * Creates a collection, runs {@code body}, and deletes the collection
     * afterwards — even if {@code body} throws — so integration tests never
     * leak collections into the shared server.
     */
    static void withCollection(TachyonClient client, CollectionSchema schema, Consumer<Collection> body) {
        client.collections.create(schema);
        try {
            body.accept(client.collection(schema.name()));
        } finally {
            try {
                client.collections.delete(schema.name());
            } catch (RuntimeException e) {
                // best-effort cleanup
            }
        }
    }
}
