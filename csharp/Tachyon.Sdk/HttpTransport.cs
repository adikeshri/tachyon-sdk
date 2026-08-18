using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;

namespace Tachyon.Sdk;

/// <summary>Thin JSON-over-HTTP client shared by every resource in the SDK.</summary>
internal sealed class HttpTransport
{
    private readonly HttpClient _httpClient;
    private readonly string _baseUrl;
    private readonly string? _apiKey;
    private readonly IReadOnlyDictionary<string, string>? _headers;

    public HttpTransport(HttpClient httpClient, string baseUrl, string? apiKey, IReadOnlyDictionary<string, string>? headers)
    {
        _httpClient = httpClient;
        _baseUrl = baseUrl.TrimEnd('/');
        _apiKey = apiKey;
        _headers = headers;
    }

    public async Task<T> RequestAsync<T>(
        HttpMethod method,
        string path,
        IReadOnlyDictionary<string, string?>? query = null,
        object? body = null,
        CancellationToken cancellationToken = default)
    {
        var (status, responseBody) = await SendAsync(method, path, query, body, cancellationToken).ConfigureAwait(false);
        if (status >= 400)
        {
            throw ErrorFromBody(status, responseBody);
        }

        try
        {
            var result = JsonSerializer.Deserialize<T>(responseBody, JsonOptions.Default);
            if (result is null)
            {
                throw new TachyonConnectionException($"Tachyon returned an empty body for a {typeof(T).Name} response.");
            }
            return result;
        }
        catch (JsonException ex)
        {
            throw new TachyonConnectionException($"Tachyon returned a response that could not be decoded: {ex.Message}", ex);
        }
    }

    public async Task RequestAsync(
        HttpMethod method,
        string path,
        IReadOnlyDictionary<string, string?>? query = null,
        object? body = null,
        CancellationToken cancellationToken = default)
    {
        var (status, responseBody) = await SendAsync(method, path, query, body, cancellationToken).ConfigureAwait(false);
        if (status >= 400)
        {
            throw ErrorFromBody(status, responseBody);
        }
    }

    public async Task<string> RequestTextAsync(HttpMethod method, string path, CancellationToken cancellationToken = default)
    {
        var (status, responseBody) = await SendAsync(method, path, null, null, cancellationToken).ConfigureAwait(false);
        if (status >= 400)
        {
            throw ErrorFromBody(status, responseBody);
        }
        return responseBody;
    }

    private async Task<(int Status, string Body)> SendAsync(
        HttpMethod method,
        string path,
        IReadOnlyDictionary<string, string?>? query,
        object? body,
        CancellationToken cancellationToken)
    {
        var url = BuildUrl(path, query);
        using var request = new HttpRequestMessage(method, url);
        request.Headers.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
        if (!string.IsNullOrEmpty(_apiKey))
        {
            request.Headers.Add("X-TACHYON-API-KEY", _apiKey);
        }
        if (_headers is not null)
        {
            foreach (var (key, value) in _headers)
            {
                request.Headers.TryAddWithoutValidation(key, value);
            }
        }
        if (body is not null)
        {
            var json = JsonSerializer.Serialize(body, JsonOptions.Default);
            request.Content = new StringContent(json, Encoding.UTF8, "application/json");
        }

        HttpResponseMessage response;
        try
        {
            response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        }
        catch (OperationCanceledException ex) when (!cancellationToken.IsCancellationRequested)
        {
            // The caller's own token wasn't cancelled, so this was HttpClient.Timeout firing.
            throw new TachyonTimeoutException($"Request to {url} timed out: {ex.Message}", ex);
        }
        catch (HttpRequestException ex)
        {
            throw new TachyonConnectionException($"Failed to reach Tachyon at {url}: {ex.Message}", ex);
        }

        using (response)
        {
            var responseBody = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
            return ((int)response.StatusCode, responseBody);
        }
    }

    private string BuildUrl(string path, IReadOnlyDictionary<string, string?>? query)
    {
        var url = _baseUrl + path;
        if (query is null || query.Count == 0)
        {
            return url;
        }
        var parts = new List<string>();
        foreach (var (key, value) in query)
        {
            if (value is null)
            {
                continue;
            }
            parts.Add($"{Uri.EscapeDataString(key)}={Uri.EscapeDataString(value)}");
        }
        return parts.Count == 0 ? url : $"{url}?{string.Join('&', parts)}";
    }

    private static TachyonException ErrorFromBody(int status, string body)
    {
        try
        {
            var payload = JsonSerializer.Deserialize<ApiErrorEnvelope>(body, JsonOptions.Default);
            if (payload?.Error is { Code.Length: > 0 } error)
            {
                return TachyonExceptionFactory.FromResponse(status, error.Message ?? body, error.Code);
            }
        }
        catch (JsonException)
        {
            // Falls through to the generic mapping below.
        }
        return TachyonExceptionFactory.FromResponse(status, body, "internal_error");
    }

    private sealed record ApiErrorEnvelope
    {
        public ApiError? Error { get; init; }
    }

    private sealed record ApiError
    {
        public string? Code { get; init; }
        public string? Message { get; init; }
    }
}
