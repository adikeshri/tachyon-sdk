package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class CollectionsTest {

    @Test
    void create_sendsSchemaAndReturnsInfo() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(201,
            """
            {"name":"products","fields":[{"name":"title","type":"text"}],"num_documents":0,"num_segments":0}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            CollectionInfo info = client.collections.create(CollectionSchema.builder("products")
                .fields(FieldSchema.builder("title", FieldType.TEXT).build())
                .build());

            assertEquals("products", info.name());
            assertEquals(0, info.numDocuments());
            assertEquals(1, server.requests.size());
            TestServer.CapturedRequest request = server.requests.get(0);
            assertEquals("POST", request.method());
            assertEquals("/collections", request.path());
            assertEquals("admin-key", request.header("X-TACHYON-API-KEY"));
            assertTrue(request.body().contains("\"title\""));
        }
    }

    @Test
    void list_returnsCollections() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            [{"name":"products","fields":[],"num_documents":5,"num_segments":1}]
            """))) {
            TachyonClient client = TestClients.forServer(server);

            List<CollectionInfo> list = client.collections.list();

            assertEquals(1, list.size());
            assertEquals("products", list.get(0).name());
            assertEquals(5, list.get(0).numDocuments());
        }
    }

    @Test
    void retrieve_returnsCollection() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"name":"products","fields":[],"num_documents":5,"num_segments":1}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            CollectionInfo info = client.collections.retrieve("products");

            assertEquals(5, info.numDocuments());
            assertEquals("/collections/products", server.requests.get(0).path());
        }
    }

    @Test
    void retrieve_urlEncodesCollectionName() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"name":"my products","fields":[],"num_documents":0,"num_segments":0}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            client.collections.retrieve("my products");

            assertEquals("/collections/my%20products", server.requests.get(0).path());
        }
    }

    @Test
    void delete_sendsDeleteAndCompletesOn204() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.noContent())) {
            TachyonClient client = TestClients.forServer(server);

            assertDoesNotThrow(() -> client.collections.delete("products"));

            assertEquals("DELETE", server.requests.get(0).method());
            assertEquals("/collections/products", server.requests.get(0).path());
        }
    }
}
