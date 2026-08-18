using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class CollectionsTests
{
    [Fact]
    public async Task Create_FillsInFieldAndCollectionDefaults()
    {
        var client = AdminClient();
        var name = UniqueName("coll");
        try
        {
            var created = await client.Collections.CreateAsync(new CollectionSchema
            {
                Name = name,
                Fields =
                [
                    new FieldSchema { Name = "title", Type = FieldType.Text },
                    new FieldSchema { Name = "brand", Type = FieldType.Keyword, Facet = true },
                    new FieldSchema { Name = "price", Type = FieldType.Int, Filter = true, Sort = true },
                ],
            });

            Assert.Equal(name, created.Name);
            Assert.Equal(0, created.NumDocuments);
            Assert.Equal(0, created.NumSegments);
            Assert.Equal(3, created.Fields.Count);

            var price = created.Fields.Single(f => f.Name == "price");
            Assert.True(price.Filter);
            Assert.True(price.Sort);
            Assert.True(price.Optional);
            Assert.True(price.Index);

            var brand = created.Fields.Single(f => f.Name == "brand");
            Assert.True(brand.Facet);
        }
        finally
        {
            await client.Collections.DeleteAsync(name);
        }
    }

    [Fact]
    public async Task Create_RejectsDuplicateNameWith409()
    {
        var client = AdminClient();
        var name = UniqueName("coll");
        var schema = new CollectionSchema { Name = name, Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] };
        await client.Collections.CreateAsync(schema);
        try
        {
            var exception = await Assert.ThrowsAsync<TachyonConflictException>(() => client.Collections.CreateAsync(schema));
            Assert.Equal("collection_exists", exception.Code);
            Assert.Equal(409, exception.StatusCode);
        }
        finally
        {
            await client.Collections.DeleteAsync(name);
        }
    }

    [Fact]
    public async Task List_IncludesNewlyCreatedCollection()
    {
        var client = AdminClient();
        var name = UniqueName("coll");
        await client.Collections.CreateAsync(new CollectionSchema { Name = name, Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] });
        try
        {
            var list = await client.Collections.ListAsync();
            Assert.Contains(list, c => c.Name == name);
        }
        finally
        {
            await client.Collections.DeleteAsync(name);
        }
    }

    [Fact]
    public async Task Retrieve_UnknownCollection404s()
    {
        var client = AdminClient();
        var exception = await Assert.ThrowsAsync<TachyonNotFoundException>(() => client.Collections.RetrieveAsync(UniqueName("missing")));
        Assert.Equal("collection_not_found", exception.Code);
        Assert.Equal(404, exception.StatusCode);
    }

    [Fact]
    public async Task Delete_RemovesFromRetrieveAndList()
    {
        var client = AdminClient();
        var name = UniqueName("coll");
        await client.Collections.CreateAsync(new CollectionSchema { Name = name, Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] });

        await client.Collections.DeleteAsync(name);

        await Assert.ThrowsAsync<TachyonNotFoundException>(() => client.Collections.RetrieveAsync(name));
        var list = await client.Collections.ListAsync();
        Assert.DoesNotContain(list, c => c.Name == name);
    }

    [Fact]
    public async Task Create_RejectsSchemaWithNoTextField()
    {
        var client = AdminClient();
        var exception = await Assert.ThrowsAsync<TachyonRequestException>(() =>
            client.Collections.CreateAsync(new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "price", Type = FieldType.Int }] }));
        Assert.Equal("invalid_schema", exception.Code);
        Assert.Equal(400, exception.StatusCode);
    }

    [Fact]
    public async Task Create_AppliesTypoToleranceAndDefaultSortingField()
    {
        var client = AdminClient();
        var name = UniqueName("coll");
        try
        {
            var created = await client.Collections.CreateAsync(new CollectionSchema
            {
                Name = name,
                Fields =
                [
                    new FieldSchema { Name = "title", Type = FieldType.Text },
                    new FieldSchema { Name = "popularity", Type = FieldType.Int, Sort = true },
                ],
                TypoTolerance = new TypoToleranceConfig { Enabled = false },
                DefaultSortingField = "popularity",
            });

            Assert.False(created.TypoTolerance?.Enabled);
            Assert.Equal("popularity", created.DefaultSortingField);
        }
        finally
        {
            await client.Collections.DeleteAsync(name);
        }
    }
}
