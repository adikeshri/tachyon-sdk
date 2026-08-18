namespace Tachyon.Sdk;

/// <summary><c>/analytics/*</c> — recorded automatically from search traffic, in memory only.</summary>
public sealed class AnalyticsResource
{
    private readonly HttpTransport _transport;

    internal AnalyticsResource(HttpTransport transport) => _transport = transport;

    /// <summary><c>GET /analytics/top</c>.</summary>
    public Task<AnalyticsQueriesResponse> TopAsync(AnalyticsQueryParams? queryParams = null, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<AnalyticsQueriesResponse>(HttpMethod.Get, "/analytics/top", query: ToQuery(queryParams), cancellationToken: cancellationToken);

    /// <summary><c>GET /analytics/zero-results</c>. Ranks by how often a query came back empty.</summary>
    public Task<AnalyticsQueriesResponse> ZeroResultsAsync(AnalyticsQueryParams? queryParams = null, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<AnalyticsQueriesResponse>(HttpMethod.Get, "/analytics/zero-results", query: ToQuery(queryParams), cancellationToken: cancellationToken);

    /// <summary><c>GET /analytics/latency</c>.</summary>
    public Task<AnalyticsLatencyResponse> LatencyAsync(CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<AnalyticsLatencyResponse>(HttpMethod.Get, "/analytics/latency", cancellationToken: cancellationToken);

    private static Dictionary<string, string?> ToQuery(AnalyticsQueryParams? queryParams) => new()
    {
        ["collection"] = queryParams?.Collection,
        ["limit"] = queryParams?.Limit?.ToString(),
    };
}
