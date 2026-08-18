package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class SearchIT {

    private static TachyonClient client;
    private static String datasetName;
    private static Collection dataset;

    @BeforeAll
    static void setUp() {
        client = adminClient();
        datasetName = uniqueName("search");
        client.collections.create(CollectionSchema.builder(datasetName)
            .fields(
                FieldSchema.builder("title", FieldType.TEXT).build(),
                FieldSchema.builder("description", FieldType.TEXT).build(),
                FieldSchema.builder("brand", FieldType.KEYWORD).facet(true).build(),
                FieldSchema.builder("price", FieldType.INT).filter(true).sort(true).build(),
                FieldSchema.builder("in_stock", FieldType.BOOL).filter(true).build())
            .build());
        dataset = client.collection(datasetName);
        dataset.documents().index(
            new Document().set("id", "1").set("title", "Wireless Mouse").set("description", "A great wireless mouse for everyday use").set("brand", "Logitech").set("price", 2999).set("in_stock", true),
            new Document().set("id", "2").set("title", "Mechanical Keyboard").set("description", "Clicky keys for typing enthusiasts").set("brand", "Razer").set("price", 8999).set("in_stock", false),
            new Document().set("id", "3").set("title", "Wireless Keyboard").set("description", "Silent wireless keyboard").set("brand", "Logitech").set("price", 5999).set("in_stock", true),
            new Document().set("id", "4").set("title", "Gaming Mouse").set("description", "Wired precision gaming mouse").set("brand", "Razer").set("price", 4999).set("in_stock", true),
            new Document().set("id", "5").set("title", "USB Cable").set("description", "A basic wired cable").set("brand", "Anker").set("price", 999).set("in_stock", true));
    }

    @AfterAll
    static void tearDown() {
        client.collections.delete(datasetName);
    }

    private static List<String> hitIds(SearchResponse results) {
        return results.hits().stream().map(h -> (String) h.document().get("id")).sorted().collect(Collectors.toList());
    }

    @Test
    void matchesTitleAndDescriptionByDefault() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wireless").build());
        assertEquals(2, results.found());
        assertTrue(results.foundIsExact());
        assertEquals(List.of("1", "3"), hitIds(results));
    }

    @Test
    void queryByRestrictsMatchingFields() {
        SearchResponse onlyDescription = dataset.search(SearchParams.builder().q("Clicky").queryBy("description").build());
        assertEquals(1, onlyDescription.found());
        assertEquals("2", onlyDescription.hits().get(0).document().get("id"));

        SearchResponse onlyTitle = dataset.search(SearchParams.builder().q("Clicky").queryBy("title").build());
        assertEquals(0, onlyTitle.found());
    }

    @ParameterizedTest
    @CsvSource({
        "brand:=Logitech, 2",
        "price:<5000, 3",
        "'brand:=[Logitech,Razer]', 4",
        "in_stock:=true, 4",
    })
    void filters_simpleExpressions(String filter, int expectedFound) {
        SearchResponse results = dataset.search(SearchParams.builder().filter(filter).build());
        assertEquals(expectedFound, results.found());
    }

    @Test
    void filters_inclusiveRange() {
        SearchResponse results = dataset.search(SearchParams.builder().filter("price:[1000..5000]").build());
        assertEquals(List.of("1", "4"), hitIds(results));
    }

    @Test
    void filters_andOrWithGrouping() {
        SearchResponse results = dataset.search(SearchParams.builder().filter("(brand:=Logitech || brand:=Razer) && price:<5000").build());
        assertEquals(List.of("1", "4"), hitIds(results));
    }

    @Test
    void filters_negationOnlyMatchesDocumentsThatHaveTheField() {
        SearchResponse results = dataset.search(SearchParams.builder().filter("brand:!=Razer").build());
        assertEquals(List.of("1", "3", "5"), hitIds(results));
    }

    @Test
    void filters_negationExcludesDocumentsMissingTheFieldEntirely() {
        CollectionSchema schema = CollectionSchema.builder(uniqueName("coll"))
            .fields(FieldSchema.builder("title", FieldType.TEXT).build(), FieldSchema.builder("brand", FieldType.KEYWORD).filter(true).build())
            .build();
        withCollection(adminClient(), schema, collection -> {
            collection.documents().index(
                new Document().set("id", "a").set("title", "has brand").set("brand", "Razer"),
                new Document().set("id", "b").set("title", "no brand at all"));

            SearchResponse results = collection.search(SearchParams.builder().filter("brand:!=Razer").build());
            assertEquals(0, results.found());
        });
    }

    @Test
    void sorting_ascending() {
        SearchResponse results = dataset.search(SearchParams.builder().sort("price:asc").limit(10).build());
        List<Integer> prices = results.hits().stream().map(h -> (Integer) h.document().get("price")).collect(Collectors.toList());
        List<Integer> sorted = prices.stream().sorted().collect(Collectors.toList());
        assertEquals(sorted, prices);
    }

    @Test
    void sorting_descending() {
        SearchResponse results = dataset.search(SearchParams.builder().sort("price:desc").limit(10).build());
        List<Integer> prices = results.hits().stream().map(h -> (Integer) h.document().get("price")).collect(Collectors.toList());
        List<Integer> sorted = prices.stream().sorted(Comparator.reverseOrder()).collect(Collectors.toList());
        assertEquals(sorted, prices);
    }

    @Test
    void pagination_movesThroughResultsWithoutOverlap() {
        SearchResponse page1 = dataset.search(SearchParams.builder().sort("price:asc").limit(2).offset(0).build());
        SearchResponse page2 = dataset.search(SearchParams.builder().sort("price:asc").limit(2).offset(2).build());

        assertEquals(5, page1.found());
        assertEquals(5, page2.found());
        assertEquals(2, page1.hits().size());
        assertEquals(2, page2.hits().size());
        List<String> ids1 = hitIds(page1);
        List<String> ids2 = hitIds(page2);
        assertTrue(ids1.stream().noneMatch(ids2::contains));
    }

    @Test
    void prefix_expandsFinalTokenByDefault() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wir").build());
        assertTrue(results.found() >= 2);
    }

    @Test
    void prefix_requiresFullTokenWhenDisabled() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wir").prefix(false).build());
        assertEquals(0, results.found());
    }

    @Test
    void typoTolerance_correctsATypoByDefault() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wirelss").build());
        assertTrue(results.found() >= 1);
    }

    @Test
    void typoTolerance_findsNothingWhenDisabled() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wirelss").typoTolerance(false).build());
        assertEquals(0, results.found());
    }

    @Test
    void matchMode_allRequiresEveryToken() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wireless zzznonexistentterm").build());
        assertEquals(0, results.found());
    }

    @Test
    void matchMode_anyRequiresOneToken() {
        SearchResponse results = dataset.search(SearchParams.builder().q("wireless zzznonexistentterm").matchMode(MatchMode.ANY).build());
        assertTrue(results.found() >= 2);
    }

    @Test
    void phraseQueries_requireAdjacencyWithinAField() {
        SearchResponse results = dataset.search(SearchParams.builder().q("\"wireless mouse\"").build());
        assertEquals(1, results.found());
        assertEquals("1", results.hits().get(0).document().get("id"));
    }

    @Test
    void facets_countEveryMatchingDocumentNotJustThePage() {
        SearchResponse results = dataset.search(SearchParams.builder().facet("brand").limit(1).build());
        assertEquals(1, results.hits().size());
        var brand = results.facets().get("brand");
        assertEquals(2, brand.get("Logitech"));
        assertEquals(2, brand.get("Razer"));
        assertEquals(1, brand.get("Anker"));
    }
}
