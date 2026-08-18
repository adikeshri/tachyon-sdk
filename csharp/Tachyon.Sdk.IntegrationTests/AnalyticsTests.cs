using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class AnalyticsTests
{
    [Fact]
    public async Task RecordsSearchCountsAndZeroResultsScopedByCollection()
    {
        var client = AdminClient();
        await WithCollectionAsync(client, new CollectionSchema { Name = UniqueName("coll"), Fields = [new FieldSchema { Name = "title", Type = FieldType.Text }] }, async collection =>
        {
            await collection.Documents.IndexAsync(new Document { ["id"] = "1", ["title"] = "gizmo" });

            await collection.SearchAsync(new SearchParams { Q = "gizmo" });
            await collection.SearchAsync(new SearchParams { Q = "gizmo" });
            await collection.SearchAsync(new SearchParams { Q = "zzz-totally-absent-term" });

            var top = await client.Analytics.TopAsync(new AnalyticsQueryParams { Collection = collection.Name });
            var gizmo = top.Queries.SingleOrDefault(q => q.Query == "gizmo");
            Assert.NotNull(gizmo);
            Assert.Equal(2, gizmo!.Count);
            Assert.Equal(collection.Name, gizmo.Collection);

            var zeroResults = await client.Analytics.ZeroResultsAsync(new AnalyticsQueryParams { Collection = collection.Name });
            Assert.Contains(zeroResults.Queries, q => q.Query == "zzz-totally-absent-term");
        });
    }

    [Fact]
    public async Task Top_RespectsLimit()
    {
        var client = AdminClient();
        var top = await client.Analytics.TopAsync(new AnalyticsQueryParams { Limit = 1 });
        Assert.True(top.Queries.Count <= 1);
    }

    [Fact]
    public async Task Latency_ReportsPercentilesAcrossAllSearchesSoFar()
    {
        var client = AdminClient();
        var latency = await client.Analytics.LatencyAsync();
        Assert.True(latency.TotalSearches > 0);
        Assert.True(latency.P95Ms >= latency.P50Ms);
        Assert.True(latency.P99Ms >= latency.P95Ms);
    }
}
