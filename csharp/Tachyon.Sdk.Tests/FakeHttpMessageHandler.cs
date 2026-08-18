using System.Net;
using System.Net.Http.Headers;
using System.Text;

namespace Tachyon.Sdk.Tests;

public sealed record RecordedRequest
{
    public required HttpMethod Method { get; init; }
    public required Uri RequestUri { get; init; }
    public required HttpRequestHeaders Headers { get; init; }
    public string? Body { get; init; }
}

/// <summary>
/// Intercepts requests before anything touches a socket — the standard
/// .NET way to test HttpClient-based code without a real server. Records
/// every request it sees and replies according to <paramref name="handler"/>.
/// </summary>
public sealed class FakeHttpMessageHandler(Func<HttpRequestMessage, FakeResponse> handler) : HttpMessageHandler
{
    public List<RecordedRequest> Requests { get; } = [];

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var body = request.Content is null ? null : await request.Content.ReadAsStringAsync(cancellationToken);
        Requests.Add(new RecordedRequest
        {
            Method = request.Method,
            RequestUri = request.RequestUri!,
            Headers = request.Headers,
            Body = body,
        });

        var result = handler(request);
        var response = new HttpResponseMessage(result.Status)
        {
            Content = new StringContent(result.Body ?? string.Empty, Encoding.UTF8, result.ContentType),
        };
        return response;
    }
}

public sealed record FakeResponse(HttpStatusCode Status, string? Body = null, string ContentType = "application/json")
{
    public static FakeResponse Json(HttpStatusCode status, string body) => new(status, body);
    public static FakeResponse Text(HttpStatusCode status, string body) => new(status, body, "text/plain");
    public static FakeResponse NoContent() => new(HttpStatusCode.NoContent);
}

public static class RequestHeadersExtensions
{
    public static string? Get(this HttpRequestHeaders headers, string name) =>
        headers.TryGetValues(name, out var values) ? values.FirstOrDefault() : null;
}
