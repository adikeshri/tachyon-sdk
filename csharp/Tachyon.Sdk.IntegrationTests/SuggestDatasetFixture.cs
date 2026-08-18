using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class SuggestDatasetFixture : IAsyncLifetime
{
    private readonly TachyonClient _client = AdminClient();
    private string _name = "";

    public Collection Collection { get; private set; } = null!;

    public async Task InitializeAsync()
    {
        _name = UniqueName("suggest");
        await _client.Collections.CreateAsync(new CollectionSchema
        {
            Name = _name,
            Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }],
        });
        Collection = _client.Collection(_name);
        await Collection.Documents.IndexAsync(
            new Document { ["id"] = "1", ["title"] = "wireless mouse" },
            new Document { ["id"] = "2", ["title"] = "wireless keyboard" },
            new Document { ["id"] = "3", ["title"] = "wireless mouse" },
            new Document { ["id"] = "4", ["title"] = "wired cable" });
    }

    public async Task DisposeAsync()
    {
        try
        {
            await _client.Collections.DeleteAsync(_name);
        }
        catch
        {
            // best-effort cleanup
        }
    }
}

public class SuggestTests(SuggestDatasetFixture fixture) : IClassFixture<SuggestDatasetFixture>
{
    [Fact]
    public async Task CompletesPrefixWithLiveDocumentCounts()
    {
        var result = await fixture.Collection.SuggestAsync(new SuggestParams { Q = "wir" });
        var wireless = result.Suggestions.SingleOrDefault(s => s.Text == "wireless");
        var wired = result.Suggestions.SingleOrDefault(s => s.Text == "wired");
        Assert.NotNull(wireless);
        Assert.NotNull(wired);
        Assert.True(wireless!.Count >= 2);
    }

    [Fact]
    public async Task CapsSuggestionsAtTheRequestedLimit()
    {
        var result = await fixture.Collection.SuggestAsync(new SuggestParams { Q = "wir", Limit = 1 });
        Assert.True(result.Suggestions.Count <= 1);
    }
}
