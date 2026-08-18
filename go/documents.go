package tachyon

import (
	"context"
	"net/http"
	"net/url"
)

// DocumentsService is /collections/{name}/documents — index, fetch, and delete documents.
type DocumentsService struct {
	transport      *httpTransport
	collectionName string
}

// Index sends POST /collections/{name}/documents, upserting one or more
// documents by id. Individual documents can fail without failing their
// neighbours — check NumFailed and Results on the response.
func (s *DocumentsService) Index(ctx context.Context, docs ...Document) (DocumentsIndexResponse, error) {
	var out DocumentsIndexResponse
	path := "/collections/" + url.PathEscape(s.collectionName) + "/documents"
	err := s.transport.do(ctx, http.MethodPost, path, requestOptions{body: docs}, &out)
	return out, err
}

// Retrieve sends GET /collections/{name}/documents/{id}.
func (s *DocumentsService) Retrieve(ctx context.Context, id string) (Document, error) {
	var out Document
	path := "/collections/" + url.PathEscape(s.collectionName) + "/documents/" + url.PathEscape(id)
	err := s.transport.do(ctx, http.MethodGet, path, requestOptions{}, &out)
	return out, err
}

// Delete sends DELETE /collections/{name}/documents/{id}.
func (s *DocumentsService) Delete(ctx context.Context, id string) error {
	path := "/collections/" + url.PathEscape(s.collectionName) + "/documents/" + url.PathEscape(id)
	return s.transport.do(ctx, http.MethodDelete, path, requestOptions{}, nil)
}
