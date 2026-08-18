namespace Tachyon.Sdk.IntegrationTests;

public static class TestSupport
{
    public static string BaseUrl => Environment.GetEnvironmentVariable("TACHYON_URL") ?? "http://localhost:8108";
    public static string AdminKey => Environment.GetEnvironmentVariable("TACHYON_ADMIN_KEY") ?? "admin-key";
    public static string SearchKey => Environment.GetEnvironmentVariable("TACHYON_SEARCH_KEY") ?? "search-key";

    public static TachyonClient AdminClient() => new(new TachyonClientOptions { Url = BaseUrl, ApiKey = AdminKey });
    public static TachyonClient SearchOnlyClient() => new(new TachyonClientOptions { Url = BaseUrl, ApiKey = SearchKey });
    public static TachyonClient AnonymousClient() => new(new TachyonClientOptions { Url = BaseUrl });

    public static string UniqueName(string prefix) => $"{prefix}-{Guid.NewGuid():N}";

    /// <summary>
    /// Creates a collection, runs <paramref name="body"/>, and deletes the
    /// collection afterwards — even if <paramref name="body"/> throws — so
    /// integration tests never leak collections into the shared server.
    /// </summary>
    public static async Task WithCollectionAsync(TachyonClient client, CollectionSchema schema, Func<Collection, Task> body)
    {
        await client.Collections.CreateAsync(schema);
        try
        {
            await body(client.Collection(schema.Name));
        }
        finally
        {
            try
            {
                await client.Collections.DeleteAsync(schema.Name);
            }
            catch
            {
                // best-effort cleanup
            }
        }
    }
}
