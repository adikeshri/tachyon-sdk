using Xunit;
using static Tachyon.Sdk.IntegrationTests.TestSupport;

namespace Tachyon.Sdk.IntegrationTests;

public class ErrorsAndAuthTests
{
    [Fact]
    public async Task NoApiKeyIsRejectedWhenAuthIsEnabled()
    {
        var exception = await Assert.ThrowsAsync<TachyonAuthenticationException>(() => AnonymousClient().Collections.ListAsync());
        Assert.Equal("unauthorized", exception.Code);
        Assert.Equal(401, exception.StatusCode);
    }

    [Fact]
    public async Task SearchOnlyKeyCanRead()
    {
        var list = await SearchOnlyClient().Collections.ListAsync();
        Assert.NotNull(list);
    }

    [Fact]
    public async Task SearchOnlyKeyCannotWrite()
    {
        var exception = await Assert.ThrowsAsync<TachyonAuthorizationException>(() =>
            SearchOnlyClient().Collections.CreateAsync(new CollectionSchema { Name = UniqueName("coll"), Fields = [] }));
        Assert.Equal("forbidden", exception.Code);
        Assert.Equal(403, exception.StatusCode);
    }

    [Fact]
    public async Task RealNetworkFailureRaisesConnectionException()
    {
        // Nothing listens on port 1 (a reserved, unused TCP port), so this
        // exercises a real network failure rather than a mocked one.
        var unreachable = new TachyonClient(new TachyonClientOptions { Url = "http://127.0.0.1:1", Timeout = TimeSpan.FromSeconds(2) });
        await Assert.ThrowsAsync<TachyonConnectionException>(() => unreachable.HealthAsync());
    }

    [Fact]
    public async Task AdminKeyWorksEndToEnd()
    {
        var list = await AdminClient().Collections.ListAsync();
        Assert.NotNull(list);
    }
}
