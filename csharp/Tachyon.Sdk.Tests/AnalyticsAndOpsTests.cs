using System.Net;
using System.Web;
using Xunit;

namespace Tachyon.Sdk.Tests;

public class AnalyticsAndOpsTests
{
    [Fact]
    public async Task Analytics_Top_ScopedToCollection()
    {
        var (client, handler) = TestClientFactory.Create(_ => FakeResponse.Json(HttpStatusCode.OK, """
            {
              "queries": [{"query":"wireless mouse","collection":"products","count":3,"zero_result_count":0,"last_result_count":12,"avg_latency_ms":1.8,"last_seen":1786625778150}],
              "tracked_queries": 1,
              "dropped_queries": 0
            }
            """));

        var result = await client.Analytics.TopAsync(new AnalyticsQueryParams { Collection = "products", Limit = 10 });

        Assert.Equal("wireless mouse", result.Queries[0].Query);
        var query = HttpUtility.ParseQueryString(handler.Requests[0].RequestUri.Query);
        Assert.Equal("/analytics/top", handler.Requests[0].RequestUri.AbsolutePath);
        Assert.Equal("products", query["collection"]);
        Assert.Equal("10", query["limit"]);
    }

    [Fact]
    public async Task Analytics_ZeroResults_HitsRightPath()
    {
        var (client, handler) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"queries":[],"tracked_queries":0,"dropped_queries":0}"""));

        await client.Analytics.ZeroResultsAsync();

        Assert.Equal("/analytics/zero-results", handler.Requests[0].RequestUri.AbsolutePath);
    }

    [Fact]
    public async Task Analytics_Latency_ReturnsPercentiles()
    {
        var (client, _) = TestClientFactory.Create(_ => FakeResponse.Json(HttpStatusCode.OK,
            """{"count":20,"mean_ms":1.9,"p50_ms":2.0,"p95_ms":4.0,"p99_ms":4.0,"max_ms":3.4,"total_searches":20,"uptime_seconds":61,"queries_per_second":0.33}"""));

        var result = await client.Analytics.LatencyAsync();

        Assert.Equal(4.0, result.P95Ms);
    }

    [Fact]
    public async Task Health_ReturnsStatus()
    {
        var (client, _) = TestClientFactory.Create(_ =>
            FakeResponse.Json(HttpStatusCode.OK, """{"ok":true,"version":"0.1.0","uptime_seconds":61,"num_collections":1}"""));

        var health = await client.HealthAsync();

        Assert.True(health.Ok);
    }

    [Fact]
    public async Task Metrics_ReturnsRawTextNotJson()
    {
        var (client, _) = TestClientFactory.Create(_ =>
            FakeResponse.Text(HttpStatusCode.OK, "# HELP tachyon_uptime_seconds\ntachyon_uptime_seconds 61\n"));

        var metrics = await client.MetricsAsync();

        Assert.Contains("tachyon_uptime_seconds 61", metrics);
    }
}
