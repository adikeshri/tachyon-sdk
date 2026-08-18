using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class OpsTests
{
    [Fact]
    public async Task Health_DoesNotRequireAnApiKey()
    {
        var health = await AnonymousClient().HealthAsync();
        Assert.True(health.Ok);
        Assert.True(health.NumCollections >= 0);
    }

    [Fact]
    public async Task Metrics_ExposesPrometheusText()
    {
        var metrics = await AdminClient().MetricsAsync();
        Assert.Contains("tachyon_uptime_seconds", metrics);
        Assert.Contains("tachyon_collections", metrics);
    }
}
