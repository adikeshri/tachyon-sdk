using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class DocumentsTests
{
    [Fact]
    public async Task Index_IndexesASingleDocument()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            var result = await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Hello World" });
            Assert.Equal(1, result.NumIndexed);
            Assert.Equal(0, result.NumFailed);
        });
    }

    [Fact]
    public async Task Retrieve_ReturnsDocumentById()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Hello World" });
            var doc = await collection.Documents.RetrieveAsync("1");
            Assert.Equal("Hello World", (string?)doc["title"]);
        });
    }

    [Fact]
    public async Task Index_UpsertsById()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "First title" });
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Second title" });

            var doc = await collection.Documents.RetrieveAsync("1");
            Assert.Equal("Second title", (string?)doc["title"]);
        });
    }

    [Fact]
    public async Task Delete_ThenRetrieve404s()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Hello World" });
            await collection.Documents.DeleteAsync("1");

            var exception = await Assert.ThrowsAsync<TachyonNotFoundException>(() => collection.Documents.RetrieveAsync("1"));
            Assert.Equal("document_not_found", exception.Code);
        });
    }

    [Fact]
    public async Task RetrieveAndDelete_UnknownId404s()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await Assert.ThrowsAsync<TachyonNotFoundException>(() => collection.Documents.RetrieveAsync("never-existed"));
            await Assert.ThrowsAsync<TachyonNotFoundException>(() => collection.Documents.DeleteAsync("never-existed"));
        });
    }

    [Fact]
    public async Task Index_BatchReportsPerDocumentResults()
    {
        var client = AdminClient();
        var schema = new CollectionSchema
        {
            Name = UniqueName("coll"),
            Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }, new FieldSchema { Name = "price", Type = FieldType.Int }],
        };
        await WithCollectionAsync(client, schema, async collection =>
        {
            var result = await collection.Documents.IndexAsync(
                new Document { ["id"] = "1", ["price"] = 100 },
                new Document { ["id"] = "2", ["price"] = "not-a-number" },
                new Document { ["id"] = "3", ["price"] = 300 });

            Assert.Equal(2, result.NumIndexed);
            Assert.Equal(1, result.NumFailed);
            Assert.True(result.Results[0].Success);
            Assert.Equal("1", result.Results[0].Id);
            Assert.False(result.Results[1].Success);
            Assert.Equal("invalid_document", result.Results[1].Code);
            Assert.True(result.Results[2].Success);
            Assert.Equal("3", result.Results[2].Id);

            await Assert.ThrowsAsync<TachyonNotFoundException>(() => collection.Documents.RetrieveAsync("2"));
        });
    }

    [Fact]
    public async Task UndeclaredFields_AreStoredButNotIndexed()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Hello", ["undeclared_field"] = "surprise" });

            var doc = await collection.Documents.RetrieveAsync("1");
            Assert.Equal("surprise", (string?)doc["undeclared_field"]);

            var results = await collection.SearchAsync(new SearchParams { Q = "surprise" });
            Assert.Empty(results.Hits);
        });
    }
}
