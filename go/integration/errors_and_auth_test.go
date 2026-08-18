//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	tachyon "github.com/adikeshri/tachyon-sdk/go"
)

func TestNoAPIKeyIsRejectedWhenAuthIsEnabled(t *testing.T) {
	_, err := anonymousClient().Collections.List(context.Background())
	if !tachyon.IsAuthenticationError(err) {
		t.Fatalf("expected an authentication error, got %v", err)
	}
	var tErr *tachyon.Error
	if e, ok := err.(*tachyon.Error); ok {
		tErr = e
	}
	if tErr == nil || tErr.Code != "unauthorized" || tErr.Status != 401 {
		t.Fatalf("unexpected error: %+v", tErr)
	}
}

func TestSearchOnlyKeyCanRead(t *testing.T) {
	_, err := searchOnlyClient().Collections.List(context.Background())
	if err != nil {
		t.Fatalf("expected read to succeed, got %v", err)
	}
}

func TestSearchOnlyKeyCannotWrite(t *testing.T) {
	_, err := searchOnlyClient().Collections.Create(context.Background(), tachyon.CollectionSchema{
		Name: uniqueName("coll"), Fields: []tachyon.FieldSchema{},
	})
	if !tachyon.IsAuthorizationError(err) {
		t.Fatalf("expected an authorization error, got %v", err)
	}
	var tErr *tachyon.Error
	if e, ok := err.(*tachyon.Error); ok {
		tErr = e
	}
	if tErr == nil || tErr.Code != "forbidden" || tErr.Status != 403 {
		t.Fatalf("unexpected error: %+v", tErr)
	}
}

func TestRealNetworkFailureRaisesConnectionError(t *testing.T) {
	// Nothing listens on port 1 (a reserved, unused TCP port), so this
	// exercises a real network failure rather than a mocked one.
	unreachable := tachyon.NewClient("http://127.0.0.1:1", tachyon.WithTimeout(2*time.Second))
	_, err := unreachable.Health(context.Background())
	if !tachyon.IsConnectionError(err) {
		t.Fatalf("expected a connection error, got %v", err)
	}
}

func TestAdminKeyWorksEndToEnd(t *testing.T) {
	_, err := adminClient().Collections.List(context.Background())
	if err != nil {
		t.Fatalf("expected admin key to work, got %v", err)
	}
}
