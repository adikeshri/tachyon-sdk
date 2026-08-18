package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class DocumentsIT {

    @Test
    void index_indexesASingleDocument() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            DocumentsIndexResponse result = collection.documents().index(new Document().set("id", "1").set("title", "Hello World"));
            assertEquals(1, result.numIndexed());
            assertEquals(0, result.numFailed());
        });
    }

    @Test
    void retrieve_returnsDocumentById() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            collection.documents().index(new Document().set("id", "1").set("title", "Hello World"));
            Document doc = collection.documents().retrieve("1");
            assertEquals("Hello World", doc.get("title"));
        });
    }

    @Test
    void index_upsertsById() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            collection.documents().index(new Document().set("id", "1").set("title", "First title"));
            collection.documents().index(new Document().set("id", "1").set("title", "Second title"));

            Document doc = collection.documents().retrieve("1");
            assertEquals("Second title", doc.get("title"));
        });
    }

    @Test
    void delete_thenRetrieve404s() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            collection.documents().index(new Document().set("id", "1").set("title", "Hello World"));
            collection.documents().delete("1");

            TachyonNotFoundException exception = assertThrows(TachyonNotFoundException.class, () -> collection.documents().retrieve("1"));
            assertEquals("document_not_found", exception.getCode());
        });
    }

    @Test
    void retrieveAndDelete_unknownId404s() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            assertThrows(TachyonNotFoundException.class, () -> collection.documents().retrieve("never-existed"));
            assertThrows(TachyonNotFoundException.class, () -> collection.documents().delete("never-existed"));
        });
    }

    @Test
    void index_batchReportsPerDocumentResults() {
        CollectionSchema schema = CollectionSchema.builder(uniqueName("coll"))
            .fields(FieldSchema.builder("title", FieldType.TEXT).build(), FieldSchema.builder("price", FieldType.INT).build())
            .build();
        withCollection(adminClient(), schema, collection -> {
            DocumentsIndexResponse result = collection.documents().index(
                new Document().set("id", "1").set("price", 100),
                new Document().set("id", "2").set("price", "not-a-number"),
                new Document().set("id", "3").set("price", 300));

            assertEquals(2, result.numIndexed());
            assertEquals(1, result.numFailed());
            assertTrue(result.results().get(0).success());
            assertEquals("1", result.results().get(0).id());
            assertFalse(result.results().get(1).success());
            assertEquals("invalid_document", result.results().get(1).code());
            assertTrue(result.results().get(2).success());
            assertEquals("3", result.results().get(2).id());

            assertThrows(TachyonNotFoundException.class, () -> collection.documents().retrieve("2"));
        });
    }

    @Test
    void undeclaredFields_areStoredButNotIndexed() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            collection.documents().index(new Document().set("id", "1").set("title", "Hello").set("undeclared_field", "surprise"));

            Document doc = collection.documents().retrieve("1");
            assertEquals("surprise", doc.get("undeclared_field"));

            SearchResponse results = collection.search(SearchParams.builder().q("surprise").build());
            assertTrue(results.hits().isEmpty());
        });
    }
}
