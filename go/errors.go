package tachyon

import (
	"errors"
	"fmt"
)

// ErrorKind categorizes an *Error by the HTTP status it came from.
type ErrorKind int

const (
	ErrKindUnknown ErrorKind = iota
	// ErrKindRequest is 400: invalid_schema, invalid_document, invalid_query, invalid_json.
	ErrKindRequest
	// ErrKindAuthentication is 401: missing or wrong API key.
	ErrKindAuthentication
	// ErrKindAuthorization is 403: a search key attempted a write.
	ErrKindAuthorization
	// ErrKindNotFound is 404: collection_not_found, document_not_found.
	ErrKindNotFound
	// ErrKindConflict is 409: collection_exists.
	ErrKindConflict
	// ErrKindServer is 5xx: corrupt_data, io_error, internal_error.
	ErrKindServer
)

// Error is returned for every non-2xx response from Tachyon's HTTP API.
// Code is the stable machine-readable string from the response body (see
// https://github.com/adikeshri/tachyon/blob/main/docs/api.md#errors);
// Status is the HTTP status code.
type Error struct {
	Message string
	Code    string
	Status  int
	Kind    ErrorKind
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (code=%s, status=%d)", e.Message, e.Code, e.Status)
}

func errorFromResponse(status int, message, code string) *Error {
	kind := ErrKindUnknown
	switch status {
	case 400:
		kind = ErrKindRequest
	case 401:
		kind = ErrKindAuthentication
	case 403:
		kind = ErrKindAuthorization
	case 404:
		kind = ErrKindNotFound
	case 409:
		kind = ErrKindConflict
	default:
		if status >= 500 {
			kind = ErrKindServer
		}
	}
	return &Error{Message: message, Code: code, Status: status, Kind: kind}
}

func hasKind(err error, kind ErrorKind) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == kind
	}
	return false
}

func IsRequestError(err error) bool        { return hasKind(err, ErrKindRequest) }
func IsAuthenticationError(err error) bool { return hasKind(err, ErrKindAuthentication) }
func IsAuthorizationError(err error) bool  { return hasKind(err, ErrKindAuthorization) }
func IsNotFoundError(err error) bool       { return hasKind(err, ErrKindNotFound) }
func IsConflictError(err error) bool       { return hasKind(err, ErrKindConflict) }
func IsServerError(err error) bool         { return hasKind(err, ErrKindServer) }

// ConnectionError means the request never reached the server, or the server
// never replied at all. Timeout is set when this was specifically caused by
// the request exceeding its configured timeout.
type ConnectionError struct {
	Message string
	Err     error
	Timeout bool
}

func (e *ConnectionError) Error() string { return e.Message }
func (e *ConnectionError) Unwrap() error { return e.Err }

func IsConnectionError(err error) bool {
	var e *ConnectionError
	return errors.As(err, &e)
}

func IsTimeoutError(err error) bool {
	var e *ConnectionError
	return errors.As(err, &e) && e.Timeout
}
