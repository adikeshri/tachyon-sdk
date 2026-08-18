package tachyon

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestErrorMappingByStatus(t *testing.T) {
	cases := []struct {
		status int
		code   string
		kind   ErrorKind
		check  func(error) bool
	}{
		{400, "invalid_query", ErrKindRequest, IsRequestError},
		{401, "unauthorized", ErrKindAuthentication, IsAuthenticationError},
		{403, "forbidden", ErrKindAuthorization, IsAuthorizationError},
		{404, "collection_not_found", ErrKindNotFound, IsNotFoundError},
		{409, "collection_exists", ErrKindConflict, IsConflictError},
		{500, "internal_error", ErrKindServer, IsServerError},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.code, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				writeJSON(w, tc.status, `{"error":{"code":"`+tc.code+`","message":"boom: `+tc.code+`"}}`)
			})

			_, err := client.Collections.Retrieve(context.Background(), "products")
			if err == nil {
				t.Fatal("expected an error")
			}
			var tachyonErr *Error
			if !isTachyonError(err, &tachyonErr) {
				t.Fatalf("expected a *tachyon.Error, got %T: %v", err, err)
			}
			if tachyonErr.Code != tc.code || tachyonErr.Status != tc.status || tachyonErr.Kind != tc.kind {
				t.Fatalf("unexpected error: %+v", tachyonErr)
			}
			if !tc.check(err) {
				t.Fatalf("predicate for %s returned false", tc.code)
			}
		})
	}
}

func TestConnectionErrorOnNetworkFailure(t *testing.T) {
	// Nothing listens on port 1 (a reserved, unused TCP port), so this
	// exercises a real network failure rather than a mocked one.
	client := NewClient("http://127.0.0.1:1", WithTimeout(2*time.Second))

	_, err := client.Health(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	if !IsConnectionError(err) {
		t.Fatalf("expected a connection error, got %T: %v", err, err)
	}
	if IsTimeoutError(err) {
		t.Fatalf("did not expect a timeout error: %v", err)
	}
}

func TestTimeoutErrorWhenServerIsSlow(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		writeJSON(w, 200, `{"ok":true,"version":"x","uptime_seconds":0,"num_collections":0}`)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Health(ctx)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !IsTimeoutError(err) {
		t.Fatalf("expected a timeout error, got %T: %v", err, err)
	}
	if !IsConnectionError(err) {
		t.Fatalf("expected IsConnectionError to also be true for a timeout: %v", err)
	}
}

func TestFallsBackToGenericErrorCodeWhenBodyIsNotTheDocumentedShape(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("upstream exploded"))
	})

	_, err := client.Health(context.Background())
	var tachyonErr *Error
	if !isTachyonError(err, &tachyonErr) {
		t.Fatalf("expected a *tachyon.Error, got %T: %v", err, err)
	}
	if tachyonErr.Code != "internal_error" || tachyonErr.Status != 500 {
		t.Fatalf("unexpected error: %+v", tachyonErr)
	}
}

func isTachyonError(err error, target **Error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	*target = e
	return true
}
