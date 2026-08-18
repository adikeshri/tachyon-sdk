package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class SuggestIT {

    private static TachyonClient client;
    private static String datasetName;
    private static Collection dataset;

    @BeforeAll
    static void setUp() {
        client = adminClient();
        datasetName = uniqueName("suggest");
        client.collections.create(CollectionSchema.builder(datasetName).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build());
        dataset = client.collection(datasetName);
        dataset.documents().index(
            new Document().set("id", "1").set("title", "wireless mouse"),
            new Document().set("id", "2").set("title", "wireless keyboard"),
            new Document().set("id", "3").set("title", "wireless mouse"),
            new Document().set("id", "4").set("title", "wired cable"));
    }

    @AfterAll
    static void tearDown() {
        client.collections.delete(datasetName);
    }

    @Test
    void completesPrefixWithLiveDocumentCounts() {
        SuggestResponse result = dataset.suggest(SuggestParams.builder("wir").build());
        Suggestion wireless = result.suggestions().stream().filter(s -> s.text().equals("wireless")).findFirst().orElse(null);
        Suggestion wired = result.suggestions().stream().filter(s -> s.text().equals("wired")).findFirst().orElse(null);
        assertNotNull(wireless);
        assertNotNull(wired);
        assertTrue(wireless.count() >= 2);
    }

    @Test
    void capsSuggestionsAtTheRequestedLimit() {
        SuggestResponse result = dataset.suggest(SuggestParams.builder("wir").limit(1).build());
        assertTrue(result.suggestions().size() <= 1);
    }
}
