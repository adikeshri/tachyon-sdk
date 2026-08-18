package tachyon

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Collection is a handle scoped to one collection. Get one via
// Client.Collection(name); it does not verify the collection exists until
// you make a request.
type Collection struct {
	Documents *DocumentsService

	name      string
	transport *httpTransport
}

// Name returns the collection name this handle is scoped to.
func (c *Collection) Name() string { return c.name }

// Search sends GET /collections/{name}/search.
func (c *Collection) Search(ctx context.Context, params SearchParams) (SearchResponse, error) {
	query := map[string]string{}
	if params.Q != "" {
		query["q"] = params.Q
	}
	if len(params.QueryBy) > 0 {
		query["query_by"] = strings.Join(params.QueryBy, ",")
	}
	if params.Filter != "" {
		query["filter"] = params.Filter
	}
	if params.Sort != "" {
		query["sort"] = params.Sort
	}
	if len(params.Facet) > 0 {
		query["facet"] = strings.Join(params.Facet, ",")
	}
	if params.Limit != 0 {
		query["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Offset != 0 {
		query["offset"] = strconv.Itoa(params.Offset)
	}
	if params.Prefix != nil {
		query["prefix"] = strconv.FormatBool(*params.Prefix)
	}
	if params.TypoTolerance != nil {
		query["typo_tolerance"] = strconv.FormatBool(*params.TypoTolerance)
	}
	if params.MatchMode != "" {
		query["match_mode"] = string(params.MatchMode)
	}

	var out SearchResponse
	path := "/collections/" + url.PathEscape(c.name) + "/search"
	err := c.transport.do(ctx, http.MethodGet, path, requestOptions{query: query}, &out)
	return out, err
}

// Suggest sends GET /collections/{name}/suggest.
func (c *Collection) Suggest(ctx context.Context, params SuggestParams) (SuggestResponse, error) {
	query := map[string]string{}
	if params.Q != "" {
		query["q"] = params.Q
	}
	if len(params.QueryBy) > 0 {
		query["query_by"] = strings.Join(params.QueryBy, ",")
	}
	if params.Limit != 0 {
		query["limit"] = strconv.Itoa(params.Limit)
	}
	if params.TypoTolerance != nil {
		query["typo_tolerance"] = strconv.FormatBool(*params.TypoTolerance)
	}

	var out SuggestResponse
	path := "/collections/" + url.PathEscape(c.name) + "/suggest"
	err := c.transport.do(ctx, http.MethodGet, path, requestOptions{query: query}, &out)
	return out, err
}
