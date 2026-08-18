namespace Tachyon.Sdk.Tests;

public static class TestClientFactory
{
    public static (TachyonClient Client, FakeHttpMessageHandler Handler) Create(
        Func<HttpRequestMessage, FakeResponse> handler,
        string? apiKey = "admin-key")
    {
        var fakeHandler = new FakeHttpMessageHandler(handler);
        var httpClient = new HttpClient(fakeHandler);
        var client = new TachyonClient(new TachyonClientOptions
        {
            Url = "http://localhost:8108",
            ApiKey = apiKey,
            HttpClient = httpClient,
        });
        return (client, fakeHandler);
    }
}
