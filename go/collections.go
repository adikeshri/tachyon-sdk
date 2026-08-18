package tachyon

import (
	"context"
	"net/http"
	"net/url"
)

// CollectionsService is /collections — create, list, and remove collections.
type CollectionsService struct {
	transport *httpTransport
}

// Create sends POST /collections. Field types are immutable after creation.
func (s *CollectionsService) Create(ctx context.Context, schema CollectionSchema) (CollectionInfo, error) {
	var out CollectionInfo
	err := s.transport.do(ctx, http.MethodPost, "/collections", requestOptions{body: schema}, &out)
	return out, err
}

// List sends GET /collections.
func (s *CollectionsService) List(ctx context.Context) ([]CollectionInfo, error) {
	var out []CollectionInfo
	err := s.transport.do(ctx, http.MethodGet, "/collections", requestOptions{}, &out)
	return out, err
}

// Retrieve sends GET /collections/{name}.
func (s *CollectionsService) Retrieve(ctx context.Context, name string) (CollectionInfo, error) {
	var out CollectionInfo
	err := s.transport.do(ctx, http.MethodGet, "/collections/"+url.PathEscape(name), requestOptions{}, &out)
	return out, err
}

// Delete sends DELETE /collections/{name}. Removes the collection and all its data.
func (s *CollectionsService) Delete(ctx context.Context, name string) error {
	return s.transport.do(ctx, http.MethodDelete, "/collections/"+url.PathEscape(name), requestOptions{}, nil)
}
