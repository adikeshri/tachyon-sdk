package tachyon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const apiKeyHeader = "X-TACHYON-API-KEY"

// httpTransport is the thin JSON-over-HTTP client shared by every service in the SDK.
type httpTransport struct {
	baseURL    string
	apiKey     string
	headers    map[string]string
	httpClient *http.Client
}

type requestOptions struct {
	query map[string]string
	body  any
}

type apiErrorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// do sends a request and decodes a JSON response into out (which may be nil
// for responses with no body, e.g. DELETE).
func (t *httpTransport) do(ctx context.Context, method, path string, opts requestOptions, out any) error {
	status, body, err := t.send(ctx, method, path, opts)
	if err != nil {
		return err
	}
	if status == 204 || len(body) == 0 {
		return nil
	}

	if status >= 400 {
		return errorFromBody(status, body)
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return &ConnectionError{
			Message: fmt.Sprintf("tachyon returned a response that could not be decoded: %v", err),
			Err:     err,
		}
	}
	return nil
}

// doText sends a request and returns the raw response body as text (used for /metrics).
func (t *httpTransport) doText(ctx context.Context, method, path string) (string, error) {
	status, body, err := t.send(ctx, method, path, requestOptions{})
	if err != nil {
		return "", err
	}
	if status >= 400 {
		return "", errorFromBody(status, body)
	}
	return string(body), nil
}

func errorFromBody(status int, body []byte) error {
	var payload apiErrorBody
	if err := json.Unmarshal(body, &payload); err != nil || payload.Error.Code == "" {
		return errorFromResponse(status, string(body), "internal_error")
	}
	return errorFromResponse(status, payload.Error.Message, payload.Error.Code)
}

func (t *httpTransport) send(ctx context.Context, method, path string, opts requestOptions) (int, []byte, error) {
	u, err := url.Parse(t.baseURL + path)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid URL %q: %w", t.baseURL+path, err)
	}
	if len(opts.query) > 0 {
		q := u.Query()
		for k, v := range opts.query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	if opts.body != nil {
		encoded, err := json.Marshal(opts.body)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		bodyReader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if opts.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if t.apiKey != "" {
		req.Header.Set(apiKeyHeader, t.apiKey)
	}
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, nil, &ConnectionError{
				Message: fmt.Sprintf("request to %s timed out: %v", u.String(), err),
				Err:     err,
				Timeout: true,
			}
		}
		return 0, nil, &ConnectionError{
			Message: fmt.Sprintf("failed to reach tachyon at %s: %v", u.String(), err),
			Err:     err,
		}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, &ConnectionError{
			Message: fmt.Sprintf("failed to read response from %s: %v", u.String(), err),
			Err:     err,
		}
	}

	return resp.StatusCode, data, nil
}
