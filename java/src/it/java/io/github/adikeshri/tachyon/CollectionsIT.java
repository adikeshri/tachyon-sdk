package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import java.util.List;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class CollectionsIT {

    @Test
    void create_fillsInFieldAndCollectionDefaults() {
        TachyonClient client = adminClient();
        String name = uniqueName("coll");
        try {
            CollectionInfo created = client.collections.create(CollectionSchema.builder(name)
                .fields(
                    FieldSchema.builder("title", FieldType.TEXT).build(),
                    FieldSchema.builder("brand", FieldType.KEYWORD).facet(true).build(),
                    FieldSchema.builder("price", FieldType.INT).filter(true).sort(true).build())
                .build());

            assertEquals(name, created.name());
            assertEquals(0, created.numDocuments());
            assertEquals(0, created.numSegments());
            assertEquals(3, created.fields().size());

            FieldSchema price = created.fields().stream().filter(f -> f.name().equals("price")).findFirst().orElseThrow();
            assertTrue(price.filter());
            assertTrue(price.sort());
            assertTrue(price.optional());
            assertTrue(price.index());

            FieldSchema brand = created.fields().stream().filter(f -> f.name().equals("brand")).findFirst().orElseThrow();
            assertTrue(brand.facet());
        } finally {
            client.collections.delete(name);
        }
    }

    @Test
    void create_rejectsDuplicateNameWith409() {
        TachyonClient client = adminClient();
        String name = uniqueName("coll");
        CollectionSchema schema = CollectionSchema.builder(name).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build();
        client.collections.create(schema);
        try {
            TachyonConflictException exception = assertThrows(TachyonConflictException.class, () -> client.collections.create(schema));
            assertEquals("collection_exists", exception.getCode());
            assertEquals(409, exception.getStatusCode());
        } finally {
            client.collections.delete(name);
        }
    }

    @Test
    void list_includesNewlyCreatedCollection() {
        TachyonClient client = adminClient();
        String name = uniqueName("coll");
        client.collections.create(CollectionSchema.builder(name).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build());
        try {
            List<CollectionInfo> list = client.collections.list();
            assertTrue(list.stream().anyMatch(c -> c.name().equals(name)));
        } finally {
            client.collections.delete(name);
        }
    }

    @Test
    void retrieve_unknownCollection404s() {
        TachyonClient client = adminClient();
        TachyonNotFoundException exception = assertThrows(TachyonNotFoundException.class, () -> client.collections.retrieve(uniqueName("missing")));
        assertEquals("collection_not_found", exception.getCode());
        assertEquals(404, exception.getStatusCode());
    }

    @Test
    void delete_removesFromRetrieveAndList() {
        TachyonClient client = adminClient();
        String name = uniqueName("coll");
        client.collections.create(CollectionSchema.builder(name).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build());

        client.collections.delete(name);

        assertThrows(TachyonNotFoundException.class, () -> client.collections.retrieve(name));
        assertTrue(client.collections.list().stream().noneMatch(c -> c.name().equals(name)));
    }

    @Test
    void create_rejectsSchemaWithNoTextField() {
        TachyonClient client = adminClient();
        TachyonRequestException exception = assertThrows(TachyonRequestException.class, () ->
            client.collections.create(CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("price", FieldType.INT).build()).build()));
        assertEquals("invalid_schema", exception.getCode());
        assertEquals(400, exception.getStatusCode());
    }

    @Test
    void create_appliesTypoToleranceAndDefaultSortingField() {
        TachyonClient client = adminClient();
        String name = uniqueName("coll");
        try {
            CollectionInfo created = client.collections.create(CollectionSchema.builder(name)
                .fields(
                    FieldSchema.builder("title", FieldType.TEXT).build(),
                    FieldSchema.builder("popularity", FieldType.INT).sort(true).build())
                .typoTolerance(TypoToleranceConfig.builder().enabled(false).build())
                .defaultSortingField("popularity")
                .build());

            assertFalse(created.typoTolerance().enabled());
            assertEquals("popularity", created.defaultSortingField());
        } finally {
            client.collections.delete(name);
        }
    }
}
