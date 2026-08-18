using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

/// <summary>
/// A read-only dataset shared across every test in a class, set up once via
/// xUnit's IAsyncLifetime (the equivalent of a beforeAll/afterAll) rather
/// than per test.
/// </summary>
public class SearchDatasetFixture : IAsyncLifetime
{
    private readonly TachyonClient _client = AdminClient();
    private string _name = "";

    public Collection Collection { get; private set; } = null!;

    public async Task InitializeAsync()
    {
        _name = UniqueName("search");
        await _client.Collections.CreateAsync(new CollectionSchema
        {
            Name = _name,
            Fields =
            [
                new FieldSchema { Name = "title", Type = FieldType.Text },
                new FieldSchema { Name = "description", Type = FieldType.Text },
                new FieldSchema { Name = "brand", Type = FieldType.Keyword, Facet = true },
                new FieldSchema { Name = "price", Type = FieldType.Int, Filter = true, Sort = true },
                new FieldSchema { Name = "in_stock", Type = FieldType.Bool, Filter = true },
            ],
        });
        Collection = _client.Collection(_name);
        await Collection.Documents.IndexAsync(
            new Document { ["id"] = "1", ["title"] = "Wireless Mouse", ["description"] = "A great wireless mouse for everyday use", ["brand"] = "Logitech", ["price"] = 2999, ["in_stock"] = true },
            new Document { ["id"] = "2", ["title"] = "Mechanical Keyboard", ["description"] = "Clicky keys for typing enthusiasts", ["brand"] = "Razer", ["price"] = 8999, ["in_stock"] = false },
            new Document { ["id"] = "3", ["title"] = "Wireless Keyboard", ["description"] = "Silent wireless keyboard", ["brand"] = "Logitech", ["price"] = 5999, ["in_stock"] = true },
            new Document { ["id"] = "4", ["title"] = "Gaming Mouse", ["description"] = "Wired precision gaming mouse", ["brand"] = "Razer", ["price"] = 4999, ["in_stock"] = true },
            new Document { ["id"] = "5", ["title"] = "USB Cable", ["description"] = "A basic wired cable", ["brand"] = "Anker", ["price"] = 999, ["in_stock"] = true });
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
