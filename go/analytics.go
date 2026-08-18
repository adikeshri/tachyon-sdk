package tachyon

import (
	"context"
	"net/http"
	"strconv"
)

// AnalyticsService is /analytics/* — recorded automatically from search
// traffic, in memory only.
type AnalyticsService struct {
	transport *httpTransport
}

func analyticsQuery(params AnalyticsQueryParams) map[string]string {
	query := map[string]string{}
	if params.Collection != "" {
		query["collection"] = params.Collection
	}
	if params.Limit != 0 {
		query["limit"] = strconv.Itoa(params.Limit)
	}
	return query
}

// Top sends GET /analytics/top.
func (s *AnalyticsService) Top(ctx context.Context, params AnalyticsQueryParams) (AnalyticsQueriesResponse, error) {
	var out AnalyticsQueriesResponse
	err := s.transport.do(ctx, http.MethodGet, "/analytics/top", requestOptions{query: analyticsQuery(params)}, &out)
	return out, err
}

// ZeroResults sends GET /analytics/zero-results. Ranks by how often a query
// came back empty.
func (s *AnalyticsService) ZeroResults(ctx context.Context, params AnalyticsQueryParams) (AnalyticsQueriesResponse, error) {
	var out AnalyticsQueriesResponse
	err := s.transport.do(ctx, http.MethodGet, "/analytics/zero-results", requestOptions{query: analyticsQuery(params)}, &out)
	return out, err
}

// Latency sends GET /analytics/latency.
func (s *AnalyticsService) Latency(ctx context.Context) (AnalyticsLatencyResponse, error) {
	var out AnalyticsLatencyResponse
	err := s.transport.do(ctx, http.MethodGet, "/analytics/latency", requestOptions{}, &out)
	return out, err
}
