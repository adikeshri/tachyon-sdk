package tachyon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// capturedRequest records what a test server actually received, so tests
// can assert on method/path/query/headers/body.
type capturedRequest struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   []byte
}

// newTestClient starts a local httptest server backed by handler, returning
// a Client pointed at it plus a slice that accumulates every request the
// server receives, in order.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *[]capturedRequest) {
	t.Helper()
	captured := []capturedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = append(captured, capturedRequest{
			Method: r.Method,
			// EscapedPath, not Path: Path is auto-decoded by net/http (so
			// "my%20products" would already read back as "my products"),
			// which would hide an encoding bug in the client. EscapedPath
			// reflects what was actually sent on the wire.
			Path:   r.URL.EscapedPath(),
			Query:  r.URL.Query(),
			Header: r.Header.Clone(),
			Body:   body,
		})
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	client := NewClient(server.URL, WithAPIKey("admin-key"))
	return client, &captured
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
