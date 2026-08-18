package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;

class DocumentsAndSearchTest {

    @Test
    void documents_index_reportsPerDocumentResults() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {
              "num_indexed": 1,
              "num_failed": 1,
              "results": [
                {"success": true, "id": "1"},
                {"success": false, "code": "invalid_document", "error": "field price: expected an integer"}
              ]
            }
            """))) {
            TachyonClient client = TestClients.forServer(server);

            DocumentsIndexResponse result = client.collection("products").documents().index(
                new Document().set("id", "1").set("title", "Wireless Mouse"),
                new Document().set("id", "2").set("price", "not a number"));

            assertEquals(1, result.numIndexed());
            assertEquals(1, result.numFailed());
            assertEquals("invalid_document", result.results().get(1).code());
            assertEquals("POST", server.requests.get(0).method());
            assertEquals("/collections/products/documents", server.requests.get(0).path());
        }
    }

    @Test
    void documents_retrieve_returnsDocument() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"id":"1","title":"Wireless Mouse"}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            Document doc = client.collection("products").documents().retrieve("1");

            assertEquals("Wireless Mouse", doc.get("title"));
            assertEquals("/collections/products/documents/1", server.requests.get(0).path());
        }
    }

    @Test
    void documents_delete_sendsDelete() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.noContent())) {
            TachyonClient client = TestClients.forServer(server);

            client.collection("products").documents().delete("1");

            assertEquals("DELETE", server.requests.get(0).method());
            assertEquals("/collections/products/documents/1", server.requests.get(0).path());
        }
    }

    @Test
    void search_serializesParamsIntoExpectedQueryString() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"found":1,"found_is_exact":true,"search_time_ms":1,"hits":[]}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            client.collection("products").search(SearchParams.builder()
                .q("wireless mouse")
                .queryBy("title", "description")
                .filter("brand:=Logitech && price:<5000")
                .sort("_text_match:desc,price:asc")
                .facet("brand", "year")
                .limit(20)
                .offset(40)
                .prefix(false)
                .typoTolerance(true)
                .matchMode(MatchMode.ANY)
                .build());

            var request = server.requests.get(0);
            assertEquals("/collections/products/search", request.path());
            assertEquals("wireless mouse", request.query().get("q"));
            assertEquals("title,description", request.query().get("query_by"));
            assertEquals("brand:=Logitech && price:<5000", request.query().get("filter"));
            assertEquals("_text_match:desc,price:asc", request.query().get("sort"));
            assertEquals("brand,year", request.query().get("facet"));
            assertEquals("20", request.query().get("limit"));
            assertEquals("40", request.query().get("offset"));
            assertEquals("false", request.query().get("prefix"));
            assertEquals("true", request.query().get("typo_tolerance"));
            assertEquals("any", request.query().get("match_mode"));
        }
    }

    @Test
    void search_omitsUnsetParams() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"found":0,"found_is_exact":true,"search_time_ms":0,"hits":[]}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            client.collection("products").search();

            var query = server.requests.get(0).query();
            assertFalse(query.containsKey("filter"));
            assertFalse(query.containsKey("limit"));
        }
    }

    @Test
    void search_returnsHitsFacetsAndFoundIsExact() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {
              "found": 1240,
              "found_is_exact": false,
              "search_time_ms": 12,
              "hits": [{"document": {"id": "1", "title": "Wireless Mouse"}, "text_match": 554.788}],
              "facets": {"brand": {"Logitech": 1240, "Razer": 830}}
            }
            """))) {
            TachyonClient client = TestClients.forServer(server);

            SearchResponse results = client.collection("products").search(SearchParams.builder().q("wireless mouse").build());

            assertEquals(1240, results.found());
            assertFalse(results.foundIsExact());
            assertEquals("Wireless Mouse", results.hits().get(0).document().get("title"));
            assertEquals(1240, results.facets().get("brand").get("Logitech"));
        }
    }

    @Test
    void suggest_returnsSuggestions() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {
              "suggestions": [
                {"text": "wireless", "count": 3, "typos": 0},
                {"text": "wired", "count": 2, "typos": 0}
              ],
              "search_time_ms": 0
            }
            """))) {
            TachyonClient client = TestClients.forServer(server);

            SuggestResponse result = client.collection("products").suggest(SuggestParams.builder("wir").limit(5).build());

            assertEquals(2, result.suggestions().size());
            var request = server.requests.get(0);
            assertEquals("/collections/products/suggest", request.path());
            assertEquals("wir", request.query().get("q"));
            assertEquals("5", request.query().get("limit"));
        }
    }
}
