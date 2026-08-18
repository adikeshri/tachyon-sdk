"""
Typed shapes for Tachyon's JSON API. These are `TypedDict`s, not classes —
values coming back from `HttpClient` are plain `dict`s parsed straight from
the response body; the types here exist for editor/mypy support, not runtime
validation.
"""

from typing import Any, Dict, List, Literal, TypedDict

FieldType = Literal["text", "keyword", "int", "float", "bool", "date"]


class _FieldSchemaRequired(TypedDict):
    name: str
    type: FieldType


class FieldSchema(_FieldSchemaRequired, total=False):
    """One field in a collection schema. See `CollectionSchema`."""

    facet: bool
    """Build a facet column. Implies filterable. Default: False."""
    filter: bool
    """Allow filter expressions on this field. Default: False."""
    sort: bool
    """Allow sorting on this field. Default: False."""
    index: bool
    """Include `text` content in the inverted index. Default: True."""
    optional: bool
    """Allow documents that omit the field. Default: True."""
    boost: float
    """Per-field relevance multiplier. Defaults to 10/6/2/1 by field name."""


class TypoToleranceConfig(TypedDict, total=False):
    enabled: bool
    one_typo_min_len: int
    two_typo_min_len: int
    max_typos: int


class _CollectionSchemaRequired(TypedDict):
    name: str
    fields: List[FieldSchema]


class CollectionSchema(_CollectionSchemaRequired, total=False):
    typo_tolerance: TypoToleranceConfig
    default_sorting_field: str


class _CollectionInfoRequired(_CollectionSchemaRequired):
    num_documents: int
    num_segments: int


class CollectionInfo(_CollectionInfoRequired, total=False):
    """A collection schema as returned by the server, with live counters."""

    typo_tolerance: TypoToleranceConfig
    default_sorting_field: str


TachyonDocument = Dict[str, Any]
"""A document is an arbitrary JSON object; `id` is always a string."""


class _DocumentIndexResultRequired(TypedDict):
    success: bool


class DocumentIndexResult(_DocumentIndexResultRequired, total=False):
    id: str
    code: str
    error: str


class DocumentsIndexResponse(TypedDict):
    num_indexed: int
    num_failed: int
    results: List[DocumentIndexResult]


MatchMode = Literal["all", "any"]


class SearchHit(TypedDict):
    document: TachyonDocument
    text_match: float


class _SearchResponseRequired(TypedDict):
    found: int
    found_is_exact: bool
    """
    False once block-max WAND pruning has skipped part of a term's postings
    for this query, at which point `found` (and facet counts) become a
    lower bound rather than an exact count.
    """
    search_time_ms: int
    hits: List[SearchHit]


class SearchResponse(_SearchResponseRequired, total=False):
    facets: Dict[str, Dict[str, int]]


class Suggestion(TypedDict):
    text: str
    count: int
    typos: int


class SuggestResponse(TypedDict):
    suggestions: List[Suggestion]
    search_time_ms: int


class AnalyticsQuery(TypedDict):
    query: str
    collection: str
    count: int
    zero_result_count: int
    last_result_count: int
    avg_latency_ms: float
    last_seen: int


class AnalyticsQueriesResponse(TypedDict):
    queries: List[AnalyticsQuery]
    tracked_queries: int
    dropped_queries: int


class AnalyticsLatencyResponse(TypedDict):
    count: int
    mean_ms: float
    p50_ms: float
    p95_ms: float
    p99_ms: float
    max_ms: float
    total_searches: int
    uptime_seconds: int
    queries_per_second: float


class HealthResponse(TypedDict):
    ok: bool
    version: str
    uptime_seconds: int
    num_collections: int
