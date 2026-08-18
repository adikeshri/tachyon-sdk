package tachyon

// FieldType is one of the scalar field types Tachyon supports in a collection schema.
type FieldType string

const (
	FieldTypeText    FieldType = "text"
	FieldTypeKeyword FieldType = "keyword"
	FieldTypeInt     FieldType = "int"
	FieldTypeFloat   FieldType = "float"
	FieldTypeBool    FieldType = "bool"
	FieldTypeDate    FieldType = "date"
)

// FieldSchema describes one field in a collection schema. Pointer fields
// distinguish "not set, use the server's default" from an explicit false —
// use Ptr(false) or Ptr(true) to set them.
type FieldSchema struct {
	Name string    `json:"name"`
	Type FieldType `json:"type"`
	// Facet builds a facet column. Implies filterable. Default: false.
	Facet *bool `json:"facet,omitempty"`
	// Filter allows filter expressions on this field. Default: false.
	Filter *bool `json:"filter,omitempty"`
	// Sort allows sorting on this field. Default: false.
	Sort *bool `json:"sort,omitempty"`
	// Index includes `text` content in the inverted index. Default: true.
	Index *bool `json:"index,omitempty"`
	// Optional allows documents that omit the field. Default: true.
	Optional *bool `json:"optional,omitempty"`
	// Boost is a per-field relevance multiplier. Defaults to 10/6/2/1 by field name.
	Boost *float64 `json:"boost,omitempty"`
}

type TypoTolerance struct {
	Enabled       *bool `json:"enabled,omitempty"`
	OneTypoMinLen *int  `json:"one_typo_min_len,omitempty"`
	TwoTypoMinLen *int  `json:"two_typo_min_len,omitempty"`
	MaxTypos      *int  `json:"max_typos,omitempty"`
}

type CollectionSchema struct {
	Name                string         `json:"name"`
	Fields              []FieldSchema  `json:"fields"`
	TypoTolerance       *TypoTolerance `json:"typo_tolerance,omitempty"`
	DefaultSortingField string         `json:"default_sorting_field,omitempty"`
}

// CollectionInfo is a collection schema as returned by the server, with live counters.
type CollectionInfo struct {
	CollectionSchema
	NumDocuments int `json:"num_documents"`
	NumSegments  int `json:"num_segments"`
}

// Document is an arbitrary JSON object; "id" is always a string.
type Document = map[string]any

type DocumentIndexResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id,omitempty"`
	Code    string `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
}

type DocumentsIndexResponse struct {
	NumIndexed int                   `json:"num_indexed"`
	NumFailed  int                   `json:"num_failed"`
	Results    []DocumentIndexResult `json:"results"`
}

// MatchMode controls whether a search requires every token or just one.
type MatchMode string

const (
	MatchModeAll MatchMode = "all"
	MatchModeAny MatchMode = "any"
)

// SearchParams are the parameters for Collection.Search. All fields are
// optional; the zero value (empty string / 0 / nil slice) means "omit and
// let the server use its default" for every field except Prefix and
// TypoTolerance, whose zero value (false) would otherwise be indistinguishable
// from an explicit false — use a *bool there (see Ptr).
type SearchParams struct {
	// Q is the query text. Empty matches everything.
	Q string
	// QueryBy is the fields to search. Defaults to every `text` field.
	QueryBy []string
	// Filter is a filter expression, e.g. "brand:=Logitech && price:<5000".
	Filter string
	// Sort is a sort expression, e.g. "_text_match:desc,price:asc".
	Sort string
	// Facet is the fields to facet on.
	Facet []string
	// Limit is the page size. Default 10, max 250.
	Limit int
	// Offset must satisfy offset+limit <= 10,000.
	Offset int
	// Prefix prefix-matches the final token. Default true.
	Prefix *bool
	// TypoTolerance allows typo correction. Defaults to the collection's setting.
	TypoTolerance *bool
	// MatchMode: MatchModeAll requires every token, MatchModeAny requires one. Default MatchModeAll.
	MatchMode MatchMode
}

type SearchHit struct {
	Document  Document `json:"document"`
	TextMatch float64  `json:"text_match"`
}

type SearchResponse struct {
	Found int `json:"found"`
	// FoundIsExact is false once block-max WAND pruning has skipped part of
	// a term's postings for this query, at which point Found (and facet
	// counts) become a lower bound rather than an exact count.
	FoundIsExact bool                      `json:"found_is_exact"`
	SearchTimeMs int                       `json:"search_time_ms"`
	Hits         []SearchHit               `json:"hits"`
	Facets       map[string]map[string]int `json:"facets,omitempty"`
}

// SuggestParams are the parameters for Collection.Suggest.
type SuggestParams struct {
	// Q is the text being typed; only the final token is completed.
	Q string
	// QueryBy is the fields whose terms may be suggested. Defaults to every `text` field.
	QueryBy []string
	// Limit is the number of suggestions to return. Default 5, max 50.
	Limit int
	// TypoTolerance also suggests corrections. Defaults to the collection's setting.
	TypoTolerance *bool
}

type Suggestion struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
	Typos int    `json:"typos"`
}

type SuggestResponse struct {
	Suggestions  []Suggestion `json:"suggestions"`
	SearchTimeMs int          `json:"search_time_ms"`
}

// AnalyticsQueryParams are the parameters for AnalyticsService.Top and .ZeroResults.
type AnalyticsQueryParams struct {
	Collection string
	// Limit defaults to 20, max 500.
	Limit int
}

type AnalyticsQuery struct {
	Query           string  `json:"query"`
	Collection      string  `json:"collection"`
	Count           int     `json:"count"`
	ZeroResultCount int     `json:"zero_result_count"`
	LastResultCount int     `json:"last_result_count"`
	AvgLatencyMs    float64 `json:"avg_latency_ms"`
	LastSeen        int64   `json:"last_seen"`
}

type AnalyticsQueriesResponse struct {
	Queries        []AnalyticsQuery `json:"queries"`
	TrackedQueries int              `json:"tracked_queries"`
	DroppedQueries int              `json:"dropped_queries"`
}

type AnalyticsLatencyResponse struct {
	Count            int     `json:"count"`
	MeanMs           float64 `json:"mean_ms"`
	P50Ms            float64 `json:"p50_ms"`
	P95Ms            float64 `json:"p95_ms"`
	P99Ms            float64 `json:"p99_ms"`
	MaxMs            float64 `json:"max_ms"`
	TotalSearches    int     `json:"total_searches"`
	UptimeSeconds    int64   `json:"uptime_seconds"`
	QueriesPerSecond float64 `json:"queries_per_second"`
}

type HealthResponse struct {
	OK             bool   `json:"ok"`
	Version        string `json:"version"`
	UptimeSeconds  int64  `json:"uptime_seconds"`
	NumCollections int    `json:"num_collections"`
}

// Ptr returns a pointer to v — a convenience for setting optional struct
// fields like FieldSchema.Facet or SearchParams.Prefix, e.g. Ptr(true).
func Ptr[T any](v T) *T {
	return &v
}
