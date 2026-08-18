using System.Net;
using Xunit;

namespace Tachyon.Sdk.Tests;

public class CollectionsTests
{
    [Fact]
    public async Task Create_SendsSchemaAndReturnsInfo()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.Created, """{"name":"products","fields":[{"name":"title","type":"text"}],"num_documents":0,"num_segments":0}"""));

        var info = await client.Collections.CreateAsync(new CollectionSchema
        {
            Name = "products",
            Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }],
        });

        Assert.Equal("products", info.Name);
        Assert.Equal(0, info.NumDocuments);
        Assert.Single(handler.Requests);
        var request = handler.Requests[0];
        Assert.Equal(HttpMethod.Post, request.Method);
        Assert.Equal("/collections", request.RequestUri.AbsolutePath);
        Assert.Equal("admin-key", request.Headers.Get("X-TACHYON-API-KEY"));
        Assert.Contains("\"title\"", request.Body);
    }

    [Fact]
    public async Task List_ReturnsCollections()
    {
        var (client, _) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """[{"name":"products","fields":[],"num_documents":5,"num_segments":1}]"""));

        var list = await client.Collections.ListAsync();

        Assert.Single(list);
        Assert.Equal("products", list[0].Name);
        Assert.Equal(5, list[0].NumDocuments);
    }

    [Fact]
    public async Task Retrieve_ReturnsCollection()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"name":"products","fields":[],"num_documents":5,"num_segments":1}"""));

        var info = await client.Collections.RetrieveAsync("products");

        Assert.Equal(5, info.NumDocuments);
        Assert.Equal("/collections/products", handler.Requests[0].RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Retrieve_UrlEncodesCollectionName()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"name":"my products","fields":[],"num_documents":0,"num_segments":0}"""));

        await client.Collections.RetrieveAsync("my products");

        Assert.Equal("/collections/my%20products", handler.Requests[0].RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Delete_SendsDeleteAndCompletesOn204()
    {
        var (client, handler) = TestClientFactory.Create(_ => FakeResponse.NoContent());

        await client.Collections.DeleteAsync("products");

        Assert.Equal(HttpMethod.Delete, handler.Requests[0].Method);
        Assert.Equal("/collections/products", handler.Requests[0].RequestUri.AbsolutePath);
    }
}
