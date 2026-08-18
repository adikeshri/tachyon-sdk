using System.Net;
using System.Web;
using Xunit;

namespace Tachyon.Sdk.Tests;

public class DocumentsAndSearchTests
{
    [Fact]
    public async Task Documents_Index_ReportsPerDocumentResults()
    {
        var (client, handler) = TestClientFactory.Create(_ => FakeResponse.Json(HttpStatusCode.OK, """
            {
              "num_indexed": 1,
              "num_failed": 1,
              "results": [
                {"success": true, "id": "1"},
                {"success": false, "code": "invalid_document", "error": "field price: expected an integer"}
              ]
            }
            """));

        var result = await client.Collection("products").Documents.IndexAsync(
            new Document { ["id"] = "1", ["title"] = "Wireless Mouse" },
            new Document { ["id"] = "2", ["price"] = "not a number" });

        Assert.Equal(1, result.NumIndexed);
        Assert.Equal(1, result.NumFailed);
        Assert.Equal("invalid_document", result.Results[1].Code);
        var request = handler.Requests[0];
        Assert.Equal(HttpMethod.Post, request.Method);
        Assert.Equal("/collections/products/documents", request.RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Documents_Retrieve_ReturnsDocument()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"id":"1","title":"Wireless Mouse"}"""));

        var doc = await client.Collection("products").Documents.RetrieveAsync("1");

        Assert.Equal("Wireless Mouse", (string?)doc["title"]);
        Assert.Equal("/collections/products/documents/1", handler.Requests[0].RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Documents_Delete_SendsDelete()
    {
        var (client, handler) = TestClientFactory.Create(_ => FakeResponse.NoContent());

        await client.Collection("products").Documents.DeleteAsync("1");

        Assert.Equal(HttpMethod.Delete, handler.Requests[0].Method);
        Assert.Equal("/collections/products/documents/1", handler.Requests[0].RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Search_SerializesParamsIntoExpectedQueryString()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"found":1,"found_is_exact":true,"search_time_ms":1,"hits":[]}"""));

        await client.Collection("products").SearchAsync(new SearchParams
        {
            Q = "wireless mouse",
            QueryBy = ["title", "description"],
            Filter = "brand:=Logitech && price:<5000",
            Sort = "_text_match:desc,price:asc",
            Facet = ["brand", "year"],
            Limit = 20,
            Offset = 40,
            Prefix = false,
            TypoTolerance = true,
            MatchMode = MatchMode.Any,
        });

        var query = HttpUtility.ParseQueryString(handler.Requests[0].RequestUri.Query);
        Assert.Equal("/collections/products/search", handler.Requests[0].RequestUri.AbsolutePath);
        Assert.Equal("wireless mouse", query["q"]);
        Assert.Equal("title,description", query["query_by"]);
        Assert.Equal("brand:=Logitech && price:<5000", query["filter"]);
        Assert.Equal("_text_match:desc,price:asc", query["sort"]);
        Assert.Equal("brand,year", query["facet"]);
        Assert.Equal("20", query["limit"]);
        Assert.Equal("40", query["offset"]);
        Assert.Equal("false", query["prefix"]);
        Assert.Equal("true", query["typo_tolerance"]);
        Assert.Equal("any", query["match_mode"]);
    }

    [Fact]
    public async Task Search_OmitsUnsetParams()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"found":0,"found_is_exact":true,"search_time_ms":0,"hits":[]}"""));

        await client.Collection("products").SearchAsync();

        var query = HttpUtility.ParseQueryString(handler.Requests[0].RequestUri.Query);
        Assert.Null(query["filter"]);
        Assert.Null(query["limit"]);
    }

    [Fact]
    public async Task Search_ReturnsHitsFacetsAndFoundIsExact()
    {
        var (client, _) = TestClientFactory.Create(_ => FakeResponse.Json(HttpStatusCode.OK, """
            {
              "found": 1240,
              "found_is_exact": false,
              "search_time_ms": 12,
              "hits": [{"document": {"id": "1", "title": "Wireless Mouse"}, "text_match": 554.788}],
              "facets": {"brand": {"Logitech": 1240, "Razer": 830}}
            }
            """));

        var results = await client.Collection("products").SearchAsync(new SearchParams { Q = "wireless mouse" });

        Assert.Equal(1240, results.Found);
        Assert.False(results.FoundIsExact);
        Assert.Equal("Wireless Mouse", (string?)results.Hits[0].Document["title"]);
        Assert.Equal(1240, results.Facets!["brand"]["Logitech"]);
    }

    [Fact]
    public async Task Suggest_ReturnsSuggestions()
    {
        var (client, handler) = TestClientFactory.Create(_ => FakeResponse.Json(HttpStatusCode.OK, """
            {
              "suggestions": [
                {"text": "wireless", "count": 3, "typos": 0},
                {"text": "wired", "count": 2, "typos": 0}
              ],
              "search_time_ms": 0
            }
            """));

        var result = await client.Collection("products").SuggestAsync(new SuggestParams { Q = "wir", Limit = 5 });

        Assert.Equal(2, result.Suggestions.Count);
        var query = HttpUtility.ParseQueryString(handler.Requests[0].RequestUri.Query);
        Assert.Equal("/collections/products/suggest", handler.Requests[0].RequestUri.AbsolutePath);
        Assert.Equal("wir", query["q"]);
        Assert.Equal("5", query["limit"]);
    }
}
