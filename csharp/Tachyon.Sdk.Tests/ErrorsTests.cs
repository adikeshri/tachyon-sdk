using System.Net;
using Xunit;

namespace Tachyon.Sdk.Tests;

public class ErrorsTests
{
    public static IEnumerable<object[]> StatusToException()
    {
        yield return [HttpStatusCode.BadRequest, "invalid_query", typeof(TachyonRequestException)];
        yield return [HttpStatusCode.Unauthorized, "unauthorized", typeof(TachyonAuthenticationException)];
        yield return [HttpStatusCode.Forbidden, "forbidden", typeof(TachyonAuthorizationException)];
        yield return [HttpStatusCode.NotFound, "collection_not_found", typeof(TachyonNotFoundException)];
        yield return [HttpStatusCode.Conflict, "collection_exists", typeof(TachyonConflictException)];
        yield return [HttpStatusCode.InternalServerError, "internal_error", typeof(TachyonServerException)];
    }

    [Theory]
    [MemberData(nameof(StatusToException))]
    public async Task MapsHttpStatusToExceptionType(HttpStatusCode status, string code, Type expectedType)
    {
        var (client, _) = TestClientFactory.Create(_ =>
            FakeResponse.Json(status, $$$"""{"error":{"code":"{{{code}}}","message":"boom: {{{code}}}"}}"""));

        var exception = await Assert.ThrowsAsync(expectedType, () => client.Collections.RetrieveAsync("products"));

        var tachyonException = Assert.IsAssignableFrom<TachyonException>(exception);
        Assert.Equal(code, tachyonException.Code);
        Assert.Equal((int)status, tachyonException.StatusCode);
        Assert.Equal($"boom: {code}", tachyonException.Message);
    }

    [Fact]
    public async Task WrapsNetworkFailuresInConnectionException()
    {
        // Nothing listens on port 1 (a reserved, unused TCP port), so this
        // exercises a real network failure rather than a mocked one.
        var client = new TachyonClient(new TachyonClientOptions { Url = "http://127.0.0.1:1", Timeout = TimeSpan.FromSeconds(2) });

        await Assert.ThrowsAsync<TachyonConnectionException>(() => client.HealthAsync());
    }

    [Fact]
    public async Task RaisesTimeoutExceptionWhenServerIsSlow()
    {
        var (client, _) = TestClientFactory.Create(_ =>
        {
            Thread.Sleep(200);
            return FakeResponse.Json(HttpStatusCode.OK, """{"ok":true,"version":"x","uptime_seconds":0,"num_collections":0}""");
        });
        // The client-wide Timeout above is 15s (default) and irrelevant here — we
        // pass a short-lived CancellationToken instead, so this also exercises
        // that HttpClient.Timeout firing is distinguished from the caller
        // cancelling their own token (see HttpTransport.SendAsync).
        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(10));

        var exception = await Assert.ThrowsAsync<TaskCanceledException>(() => client.HealthAsync(cts.Token));
        Assert.IsNotType<TachyonTimeoutException>(exception);
    }

    [Fact]
    public async Task HttpClientTimeoutRaisesTachyonTimeoutException()
    {
        var handler = new FakeHttpMessageHandler(_ =>
        {
            Thread.Sleep(200);
            return FakeResponse.Json(HttpStatusCode.OK, """{"ok":true,"version":"x","uptime_seconds":0,"num_collections":0}""");
        });
        var httpClient = new HttpClient(handler) { Timeout = TimeSpan.FromMilliseconds(10) };
        var client = new TachyonClient(new TachyonClientOptions { Url = "http://localhost:8108", HttpClient = httpClient });

        await Assert.ThrowsAsync<TachyonTimeoutException>(() => client.HealthAsync());
    }

    [Fact]
    public async Task FallsBackToGenericCodeWhenBodyIsNotTheDocumentedShape()
    {
        var (client, _) = TestClientFactory.Create(_ => new FakeResponse(HttpStatusCode.InternalServerError, "upstream exploded", "text/plain"));

        var exception = await Assert.ThrowsAsync<TachyonServerException>(() => client.HealthAsync());

        Assert.Equal("internal_error", exception.Code);
        Assert.Equal(500, exception.StatusCode);
    }
}
