namespace Tachyon.Sdk;

/// <summary>
/// A handle scoped to one collection. Get one via <see cref="TachyonClient.Collection"/>;
/// it does not verify the collection exists until you make a request.
/// </summary>
public sealed class Collection
{
    private readonly HttpTransport _transport;

    public string Name { get; }
    public DocumentsResource Documents { get; }

    internal Collection(HttpTransport transport, string name)
    {
        _transport = transport;
        Name = name;
        Documents = new DocumentsResource(transport, name);
    }

    /// <summary><c>GET /collections/{name}/search</c>.</summary>
    public Task<SearchResponse> SearchAsync(SearchParams? searchParams = null, CancellationToken cancellationToken = default)
    {
        searchParams ??= new SearchParams();
        var query = new Dictionary<string, string?>
        {
            ["q"] = searchParams.Q,
            ["query_by"] = Join(searchParams.QueryBy),
            ["filter"] = searchParams.Filter,
            ["sort"] = searchParams.Sort,
            ["facet"] = Join(searchParams.Facet),
            ["limit"] = searchParams.Limit?.ToString(),
            ["offset"] = searchParams.Offset?.ToString(),
            ["prefix"] = FormatBool(searchParams.Prefix),
            ["typo_tolerance"] = FormatBool(searchParams.TypoTolerance),
            ["match_mode"] = searchParams.MatchMode switch
            {
                MatchMode.All => "all",
                MatchMode.Any => "any",
                null => null,
                _ => throw new ArgumentOutOfRangeException(nameof(searchParams)),
            },
        };
        return _transport.RequestAsync<SearchResponse>(
            HttpMethod.Get,
            $"/collections/{Uri.EscapeDataString(Name)}/search",
            query: query,
            cancellationToken: cancellationToken);
    }

    /// <summary><c>GET /collections/{name}/suggest</c>.</summary>
    public Task<SuggestResponse> SuggestAsync(SuggestParams suggestParams, CancellationToken cancellationToken = default)
    {
        var query = new Dictionary<string, string?>
        {
            ["q"] = suggestParams.Q,
            ["query_by"] = Join(suggestParams.QueryBy),
            ["limit"] = suggestParams.Limit?.ToString(),
            ["typo_tolerance"] = FormatBool(suggestParams.TypoTolerance),
        };
        return _transport.RequestAsync<SuggestResponse>(
            HttpMethod.Get,
            $"/collections/{Uri.EscapeDataString(Name)}/suggest",
            query: query,
            cancellationToken: cancellationToken);
    }

    private static string? Join(IReadOnlyList<string>? values) =>
        values is null || values.Count == 0 ? null : string.Join(',', values);

    private static string? FormatBool(bool? value) =>
        value is null ? null : value.Value ? "true" : "false";
}
