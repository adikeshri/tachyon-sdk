// Package tachyon is the official Go client for Tachyon, the typo-tolerant
// full-text search engine.
//
//	client := tachyon.NewClient("http://localhost:8108", tachyon.WithAPIKey("my-admin-key"))
//
//	client.Collections.Create(ctx, tachyon.CollectionSchema{
//		Name:   "products",
//		Fields: []tachyon.FieldSchema{{Name: "title", Type: tachyon.FieldTypeText}},
//	})
//	client.Collection("products").Documents.Index(ctx, tachyon.Document{"id": "1", "title": "Wireless Mouse"})
//	results, err := client.Collection("products").Search(ctx, tachyon.SearchParams{Q: "wireless mouse"})
package tachyon

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// ClientOption configures a Client. Pass zero or more to NewClient.
type ClientOption func(*clientConfig)

type clientConfig struct {
	apiKey     string
	timeout    time.Duration
	headers    map[string]string
	httpClient *http.Client
}

// WithAPIKey sets the key sent as X-TACHYON-API-KEY. Use an admin key for
// writes, a search key for read-only access.
func WithAPIKey(key string) ClientOption {
	return func(c *clientConfig) { c.apiKey = key }
}

// WithTimeout sets the per-request timeout. Default 15s. Ignored if
// WithHTTPClient supplies a client with its own Timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *clientConfig) { c.timeout = d }
}

// WithHeader adds an extra header sent on every request. May be called more than once.
func WithHeader(key, value string) ClientOption {
	return func(c *clientConfig) {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		c.headers[key] = value
	}
}

// WithHTTPClient overrides the *http.Client used for requests (mainly for
// testing, or to share transport/connection pooling with other clients).
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *clientConfig) { c.httpClient = hc }
}

// Client is a client for a single Tachyon server.
type Client struct {
	Collections *CollectionsService
	Analytics   *AnalyticsService

	transport *httpTransport
}

// NewClient builds a Client for the Tachyon server at baseURL, e.g.
// "http://localhost:8108".
func NewClient(baseURL string, opts ...ClientOption) *Client {
	cfg := clientConfig{timeout: 15 * time.Second}
	for _, opt := range opts {
		opt(&cfg)
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.timeout}
	}

	t := &httpTransport{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     cfg.apiKey,
		headers:    cfg.headers,
		httpClient: httpClient,
	}

	return &Client{
		transport:   t,
		Collections: &CollectionsService{transport: t},
		Analytics:   &AnalyticsService{transport: t},
	}
}

// Collection returns a handle scoped to one collection, for documents/search/suggest.
func (c *Client) Collection(name string) *Collection {
	return &Collection{
		name:      name,
		transport: c.transport,
		Documents: &DocumentsService{transport: c.transport, collectionName: name},
	}
}

// Health sends GET /health. Always reachable without an API key.
func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	var out HealthResponse
	err := c.transport.do(ctx, http.MethodGet, "/health", requestOptions{}, &out)
	return out, err
}

// Metrics sends GET /metrics. Prometheus exposition format, returned as plain text.
func (c *Client) Metrics(ctx context.Context) (string, error) {
	return c.transport.doText(ctx, http.MethodGet, "/metrics")
}
