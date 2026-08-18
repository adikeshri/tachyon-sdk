using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class SearchTests(SearchDatasetFixture fixture) : IClassFixture<SearchDatasetFixture>
{
    private Collection Dataset => fixture.Collection;

    private static List<string> HitIds(SearchResponse results) =>
        results.Hits.Select(h => (string)h.Document["id"]!).OrderBy(id => id, StringComparer.Ordinal).ToList();

    [Fact]
    public async Task MatchesTitleAndDescriptionByDefault()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wireless" });
        Assert.Equal(2, results.Found);
        Assert.True(results.FoundIsExact);
        Assert.Equal(["1", "3"], HitIds(results));
    }

    [Fact]
    public async Task QueryByRestrictsMatchingFields()
    {
        var onlyDescription = await Dataset.SearchAsync(new SearchParams { Q = "Clicky", QueryBy = ["description"] });
        Assert.Equal(1, onlyDescription.Found);
        Assert.Equal("2", (string?)onlyDescription.Hits[0].Document["id"]);

        var onlyTitle = await Dataset.SearchAsync(new SearchParams { Q = "Clicky", QueryBy = ["title"] });
        Assert.Equal(0, onlyTitle.Found);
    }

    [Theory]
    [InlineData("brand:=Logitech", 2)]
    [InlineData("price:<5000", 3)]
    [InlineData("brand:=[Logitech,Razer]", 4)]
    [InlineData("in_stock:=true", 4)]
    public async Task Filters_SimpleExpressions(string filter, int expectedFound)
    {
        var results = await Dataset.SearchAsync(new SearchParams { Filter = filter });
        Assert.Equal(expectedFound, results.Found);
    }

    [Fact]
    public async Task Filters_InclusiveRange()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Filter = "price:[1000..5000]" });
        Assert.Equal(["1", "4"], HitIds(results));
    }

    [Fact]
    public async Task Filters_AndOrWithGrouping()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Filter = "(brand:=Logitech || brand:=Razer) && price:<5000" });
        Assert.Equal(["1", "4"], HitIds(results));
    }

    [Fact]
    public async Task Filters_NegationOnlyMatchesDocumentsThatHaveTheField()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Filter = "brand:!=Razer" });
        Assert.Equal(["1", "3", "5"], HitIds(results));
    }

    [Fact]
    public async Task Filters_NegationExcludesDocumentsMissingTheFieldEntirely()
    {
        var client = AdminClient();
        var schema = new CollectionSchema
        {
            Name = UniqueName("coll"),
            Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }, new FieldSchema { Name = "brand", Type = FieldType.Keyword, Filter = true }],
        };
        await WithCollectionAsync(client, schema, async collection =>
        {
            await collection.Documents.IndexAsync(
                new Document { ["id"] = "a", ["title"] = "has brand", ["brand"] = "Razer" },
                new Document { ["id"] = "b", ["title"] = "no brand at all" });

            var results = await collection.SearchAsync(new SearchParams { Filter = "brand:!=Razer" });
            Assert.Equal(0, results.Found);
        });
    }

    [Fact]
    public async Task Sorting_Ascending()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Sort = "price:asc", Limit = 10 });
        var prices = results.Hits.Select(h => (double)h.Document["price"]!).ToList();
        Assert.Equal(prices.OrderBy(p => p), prices);
    }

    [Fact]
    public async Task Sorting_Descending()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Sort = "price:desc", Limit = 10 });
        var prices = results.Hits.Select(h => (double)h.Document["price"]!).ToList();
        Assert.Equal(prices.OrderByDescending(p => p), prices);
    }

    [Fact]
    public async Task Pagination_MovesThroughResultsWithoutOverlap()
    {
        var page1 = await Dataset.SearchAsync(new SearchParams { Sort = "price:asc", Limit = 2, Offset = 0 });
        var page2 = await Dataset.SearchAsync(new SearchParams { Sort = "price:asc", Limit = 2, Offset = 2 });

        Assert.Equal(5, page1.Found);
        Assert.Equal(5, page2.Found);
        Assert.Equal(2, page1.Hits.Count);
        Assert.Equal(2, page2.Hits.Count);
        Assert.Empty(HitIds(page1).Intersect(HitIds(page2)));
    }

    [Fact]
    public async Task Prefix_ExpandsFinalTokenByDefault()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wir" });
        Assert.True(results.Found >= 2);
    }

    [Fact]
    public async Task Prefix_RequiresFullTokenWhenDisabled()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wir", Prefix = false });
        Assert.Equal(0, results.Found);
    }

    [Fact]
    public async Task TypoTolerance_CorrectsATypoByDefault()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wirelss" });
        Assert.True(results.Found >= 1);
    }

    [Fact]
    public async Task TypoTolerance_FindsNothingWhenDisabled()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wirelss", TypoTolerance = false });
        Assert.Equal(0, results.Found);
    }

    [Fact]
    public async Task MatchMode_AllRequiresEveryToken()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wireless zzznonexistentterm" });
        Assert.Equal(0, results.Found);
    }

    [Fact]
    public async Task MatchMode_AnyRequiresOneToken()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "wireless zzznonexistentterm", MatchMode = MatchMode.Any });
        Assert.True(results.Found >= 2);
    }

    [Fact]
    public async Task PhraseQueries_RequireAdjacencyWithinAField()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Q = "\"wireless mouse\"" });
        Assert.Equal(1, results.Found);
        Assert.Equal("1", (string?)results.Hits[0].Document["id"]);
    }

    [Fact]
    public async Task Facets_CountEveryMatchingDocumentNotJustThePage()
    {
        var results = await Dataset.SearchAsync(new SearchParams { Facet = ["brand"], Limit = 1 });
        Assert.Single(results.Hits);
        var brand = results.Facets!["brand"];
        Assert.Equal(2, brand["Logitech"]);
        Assert.Equal(2, brand["Razer"]);
        Assert.Equal(1, brand["Anker"]);
    }
}
