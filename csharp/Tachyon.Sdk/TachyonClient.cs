namespace Tachyon.Sdk;

/// <summary>Options for constructing a <see cref="TachyonClient"/>.</summary>
public sealed class TachyonClientOptions
{
    /// <summary>Base URL of the Tachyon server, e.g. "http://localhost:8108".</summary>
    public required string Url { get; init; }

    /// <summary>Sent as X-TACHYON-API-KEY. Use an admin key for writes, a search key for read-only access.</summary>
    public string? ApiKey { get; init; }

    /// <summary>Per-request timeout. Default 15s. Ignored if <see cref="HttpClient"/> is supplied.</summary>
    public TimeSpan Timeout { get; init; } = TimeSpan.FromSeconds(15);

    /// <summary>Extra headers merged into every request.</summary>
    public IReadOnlyDictionary<string, string>? Headers { get; init; }

    /// <summary>Override the <see cref="System.Net.Http.HttpClient"/> (mainly for testing, or to share connection pooling).</summary>
    public HttpClient? HttpClient { get; init; }
}

/// <summary>
/// Client for a single Tachyon server.
///
/// <code>
/// var client = new TachyonClient(new TachyonClientOptions { Url = "http://localhost:8108", ApiKey = "my-admin-key" });
/// await client.Collections.CreateAsync(new CollectionSchema { Name = "products", Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] });
/// await client.Collection("products").Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "Wireless Mouse" });
/// var results = await client.Collection("products").SearchAsync(new SearchParams { Q = "wireless mouse" });
/// </code>
/// </summary>
public sealed class TachyonClient
{
    private readonly HttpTransport _transport;

    public CollectionsResource Collections { get; }
    public AnalyticsResource Analytics { get; }

    public TachyonClient(TachyonClientOptions options)
    {
        var httpClient = options.HttpClient ?? new HttpClient { Timeout = options.Timeout };
        _transport = new HttpTransport(httpClient, options.Url, options.ApiKey, options.Headers);
        Collections = new CollectionsResource(_transport);
        Analytics = new AnalyticsResource(_transport);
    }

    /// <summary>Get a handle scoped to one collection, for documents/search/suggest.</summary>
    public Collection Collection(string name) => new(_transport, name);

    /// <summary><c>GET /health</c>. Always reachable without an API key.</summary>
    public Task<HealthResponse> HealthAsync(CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<HealthResponse>(HttpMethod.Get, "/health", cancellationToken: cancellationToken);

    /// <summary><c>GET /metrics</c>. Prometheus exposition format, returned as plain text.</summary>
    public Task<string> MetricsAsync(CancellationToken cancellationToken = default) =>
        _transport.RequestTextAsync(HttpMethod.Get, "/metrics", cancellationToken);
}
