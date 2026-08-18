namespace Tachyon.Sdk;

/// <summary>One of the scalar field types Tachyon supports in a collection schema.</summary>
public enum FieldType
{
    Text,
    Keyword,
    Int,
    Float,
    Bool,
    Date,
}

/// <summary>Whether a search requires every token (default) or just one.</summary>
public enum MatchMode
{
    All,
    Any,
}

/// <summary>One field in a collection schema. See <see cref="CollectionSchema"/>.</summary>
public sealed record FieldSchema
{
    public required string Name { get; init; }
    public required FieldType Type { get; init; }

    /// <summary>Build a facet column. Implies filterable. Default: false.</summary>
    public bool? Facet { get; init; }

    /// <summary>Allow filter expressions on this field. Default: false.</summary>
    public bool? Filter { get; init; }

    /// <summary>Allow sorting on this field. Default: false.</summary>
    public bool? Sort { get; init; }

    /// <summary>Include `text` content in the inverted index. Default: true.</summary>
    public bool? Index { get; init; }

    /// <summary>Allow documents that omit the field. Default: true.</summary>
    public bool? Optional { get; init; }

    /// <summary>Per-field relevance multiplier. Defaults to 10/6/2/1 by field name.</summary>
    public double? Boost { get; init; }
}

public sealed record TypoToleranceConfig
{
    public bool? Enabled { get; init; }
    public int? OneTypoMinLen { get; init; }
    public int? TwoTypoMinLen { get; init; }
    public int? MaxTypos { get; init; }
}

public record CollectionSchema
{
    public required string Name { get; init; }
    public required IReadOnlyList<FieldSchema> Fields { get; init; }
    public TypoToleranceConfig? TypoTolerance { get; init; }
    public string? DefaultSortingField { get; init; }
}

/// <summary>A collection schema as returned by the server, with live counters.</summary>
public sealed record CollectionInfo : CollectionSchema
{
    public int NumDocuments { get; init; }
    public int NumSegments { get; init; }
}

public sealed record DocumentIndexResult
{
    public required bool Success { get; init; }
    public string? Id { get; init; }
    public string? Code { get; init; }
    public string? Error { get; init; }
}

public sealed record DocumentsIndexResponse
{
    public required int NumIndexed { get; init; }
    public required int NumFailed { get; init; }
    public required IReadOnlyList<DocumentIndexResult> Results { get; init; }
}

/// <summary>Parameters for <see cref="Collection.SearchAsync"/>. Every property is optional.</summary>
public sealed class SearchParams
{
    /// <summary>Query text. Null/empty matches everything.</summary>
    public string? Q { get; init; }

    /// <summary>Fields to search. Defaults to every `text` field.</summary>
    public IReadOnlyList<string>? QueryBy { get; init; }

    /// <summary>Filter expression, e.g. "brand:=Logitech &amp;&amp; price:&lt;5000".</summary>
    public string? Filter { get; init; }

    /// <summary>Sort expression, e.g. "_text_match:desc,price:asc".</summary>
    public string? Sort { get; init; }

    /// <summary>Fields to facet on.</summary>
    public IReadOnlyList<string>? Facet { get; init; }

    /// <summary>Page size. Default 10, max 250.</summary>
    public int? Limit { get; init; }

    /// <summary>Offset + limit must not exceed 10,000.</summary>
    public int? Offset { get; init; }

    /// <summary>Prefix-match the final token. Default true.</summary>
    public bool? Prefix { get; init; }

    /// <summary>Allow typo correction. Defaults to the collection's setting.</summary>
    public bool? TypoTolerance { get; init; }

    /// <summary>All requires every token, Any requires one. Default All.</summary>
    public MatchMode? MatchMode { get; init; }
}

public sealed record SearchHit
{
    public required Document Document { get; init; }
    public required double TextMatch { get; init; }
}

public sealed record SearchResponse
{
    public required int Found { get; init; }

    /// <summary>
    /// False once block-max WAND pruning has skipped part of a term's
    /// postings for this query, at which point Found (and facet counts)
    /// become a lower bound rather than an exact count.
    /// </summary>
    public required bool FoundIsExact { get; init; }

    public required int SearchTimeMs { get; init; }
    public required IReadOnlyList<SearchHit> Hits { get; init; }
    public IReadOnlyDictionary<string, IReadOnlyDictionary<string, int>>? Facets { get; init; }
}

/// <summary>Parameters for <see cref="Collection.SuggestAsync"/>.</summary>
public sealed class SuggestParams
{
    /// <summary>Text being typed; only the final token is completed.</summary>
    public required string Q { get; init; }

    /// <summary>Fields whose terms may be suggested. Defaults to every `text` field.</summary>
    public IReadOnlyList<string>? QueryBy { get; init; }

    /// <summary>Suggestions to return. Default 5, max 50.</summary>
    public int? Limit { get; init; }

    /// <summary>Also suggest corrections. Defaults to the collection's setting.</summary>
    public bool? TypoTolerance { get; init; }
}

public sealed record Suggestion
{
    public required string Text { get; init; }
    public required int Count { get; init; }
    public required int Typos { get; init; }
}

public sealed record SuggestResponse
{
    public required IReadOnlyList<Suggestion> Suggestions { get; init; }
    public required int SearchTimeMs { get; init; }
}

/// <summary>Parameters for <see cref="AnalyticsResource.TopAsync"/> and <see cref="AnalyticsResource.ZeroResultsAsync"/>.</summary>
public sealed class AnalyticsQueryParams
{
    public string? Collection { get; init; }

    /// <summary>Default 20, max 500.</summary>
    public int? Limit { get; init; }
}

public sealed record AnalyticsQuery
{
    public required string Query { get; init; }
    public required string Collection { get; init; }
    public required int Count { get; init; }
    public required int ZeroResultCount { get; init; }
    public required int LastResultCount { get; init; }
    public required double AvgLatencyMs { get; init; }
    public required long LastSeen { get; init; }
}

public sealed record AnalyticsQueriesResponse
{
    public required IReadOnlyList<AnalyticsQuery> Queries { get; init; }
    public required int TrackedQueries { get; init; }
    public required int DroppedQueries { get; init; }
}

public sealed record AnalyticsLatencyResponse
{
    public required int Count { get; init; }
    public required double MeanMs { get; init; }
    public required double P50Ms { get; init; }
    public required double P95Ms { get; init; }
    public required double P99Ms { get; init; }
    public required double MaxMs { get; init; }
    public required int TotalSearches { get; init; }
    public required long UptimeSeconds { get; init; }
    public required double QueriesPerSecond { get; init; }
}

public sealed record HealthResponse
{
    public required bool Ok { get; init; }
    public required string Version { get; init; }
    public required long UptimeSeconds { get; init; }
    public required int NumCollections { get; init; }
}
