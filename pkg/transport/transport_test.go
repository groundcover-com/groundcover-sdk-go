package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	metricsclient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
)

// WithHeadersOverride must be usable directly as a generated client option.
var _ metricsclient.ClientOption = WithHeadersOverride(nil)

// captureTransport is a base http.RoundTripper that records the request it
// receives. If next is set it forwards the request (exercising the real network
// path); otherwise it returns a canned successful response.
type captureTransport struct {
	req  *http.Request
	next http.RoundTripper
}

func (c *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.req = req.Clone(req.Context())
	if c.next != nil {
		return c.next.RoundTrip(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func newTestRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "https://example.com/api/test", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	return req
}

func TestRoundTripAppliesPerRequestHeaders(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req := newTestRequest(t)
	ctx := withRequestHeaders(req.Context(), http.Header{
		"X-Custom-Header": {"custom-value"},
		"X-Another":       {"another-value"},
	})
	req = req.WithContext(ctx)

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if got := cap.req.Header.Get("X-Custom-Header"); got != "custom-value" {
		t.Errorf("X-Custom-Header = %q, want %q", got, "custom-value")
	}
	if got := cap.req.Header.Get("X-Another"); got != "another-value" {
		t.Errorf("X-Another = %q, want %q", got, "another-value")
	}
}

func TestRoundTripMergesPerRequestHeaders(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req := newTestRequest(t)
	ctx := withRequestHeaders(req.Context(), http.Header{"X-First": {"1"}})
	ctx = withRequestHeaders(ctx, http.Header{"X-Second": {"2"}})
	req = req.WithContext(ctx)

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if got := cap.req.Header.Get("X-First"); got != "1" {
		t.Errorf("X-First = %q, want %q", got, "1")
	}
	if got := cap.req.Header.Get("X-Second"); got != "2" {
		t.Errorf("X-Second = %q, want %q", got, "2")
	}
}

func TestRoundTripPerRequestHeadersOverrideDefaults(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req := newTestRequest(t)
	ctx := withRequestHeaders(req.Context(), http.Header{headerUserAgent: {"caller-agent"}})
	req = req.WithContext(ctx)

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	// Override (not append) semantics: the default User-Agent must be replaced,
	// leaving exactly one value.
	if got := cap.req.Header.Values(headerUserAgent); len(got) != 1 || got[0] != "caller-agent" {
		t.Errorf("User-Agent = %v, want [caller-agent] (override, not append)", got)
	}
}

func TestRoundTripAppliesMultiValuePerRequestHeaders(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req := newTestRequest(t)
	ctx := withRequestHeaders(req.Context(), http.Header{"X-Multi": {"a", "b"}})
	req = req.WithContext(ctx)

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if got := cap.req.Header.Values("X-Multi"); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("X-Multi = %v, want [a b]", got)
	}
}

func TestRoundTripPerRequestHeadersOverrideWorkflowContentType(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req, err := http.NewRequest(http.MethodPost, "https://example.com/api/workflows/create", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	ctx := withRequestHeaders(req.Context(), http.Header{"Content-Type": {"application/json"}})
	req = req.WithContext(ctx)

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if got := cap.req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q (per-request header should override the workflow default)", got, "application/json")
	}
}

func TestRoundTripSetsAuthorizationWithAPIKey(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("secret-key", "backend", cap, 1, 0, 0, nil)

	if _, err := tr.RoundTrip(newTestRequest(t)); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if got := cap.req.Header.Get(headerAuthorization); got != "Bearer secret-key" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer secret-key")
	}
	if got := cap.req.Header.Get(headerBackendID); got != "backend" {
		t.Errorf("X-Backend-Id = %q, want %q", got, "backend")
	}
}

func TestRoundTripOmitsAuthorizationWithoutAPIKey(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("", "backend", cap, 1, 0, 0, nil)

	if _, err := tr.RoundTrip(newTestRequest(t)); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if _, ok := cap.req.Header[headerAuthorization]; ok {
		t.Errorf("Authorization header should be absent when no API key is set, got %q", cap.req.Header.Get(headerAuthorization))
	}
}

func TestRoundTripOmitsBackendIDWhenEmpty(t *testing.T) {
	cap := &captureTransport{}
	tr := NewTransport("key", "", cap, 1, 0, 0, nil)

	if _, err := tr.RoundTrip(newTestRequest(t)); err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}

	if _, ok := cap.req.Header[headerBackendID]; ok {
		t.Errorf("X-Backend-Id header should be absent when no backend ID is set, got %q", cap.req.Header.Get(headerBackendID))
	}
}

// NewSDKClient is the low-level primitive: it takes credentials as explicit
// arguments and never enforces them. The "required by default" policy lives at
// the user-facing NewClient layer.
func TestNewSDKClientWithoutCredentials(t *testing.T) {
	c, err := NewSDKClient("", "", "https://api.example.com")
	if err != nil {
		t.Fatalf("NewSDKClient without credentials returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewSDKClient returned nil client")
	}
}

func TestNewClientRequiresAPIKeyByDefault(t *testing.T) {
	t.Setenv("GC_API_KEY", "")
	t.Setenv("GC_BACKEND_ID", "")
	t.Setenv("GC_BASE_URL", "")

	if _, err := NewClient(option.WithBaseURL("https://api.example.com")); err == nil {
		t.Fatal("NewClient without API key should return an error by default")
	}
}

func TestNewClientAllowUnauthenticated(t *testing.T) {
	t.Setenv("GC_API_KEY", "")
	t.Setenv("GC_BACKEND_ID", "")
	t.Setenv("GC_BASE_URL", "")

	c, err := NewClient(
		option.WithBaseURL("https://api.example.com"),
		option.AllowUnauthenticated(),
	)
	if err != nil {
		t.Fatalf("NewClient with AllowUnauthenticated returned error: %v", err)
	}
	if c == nil {
		t.Fatal("NewClient returned nil client")
	}
}

func TestNewSDKClientRequiresBaseURL(t *testing.T) {
	if _, err := NewSDKClient("key", "backend", ""); err == nil {
		t.Fatal("NewSDKClient with empty base URL should return an error")
	}
}

// TestTransportEndToEnd verifies, against a real server, that a custom base
// transport is honored and that the default and custom per-request headers all
// reach the wire.
func TestTransportEndToEnd(t *testing.T) {
	var received http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cap := &captureTransport{next: http.DefaultTransport}
	tr := NewTransport("key", "backend", cap, 1, 0, 0, nil)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	ctx := withRequestHeaders(req.Context(), http.Header{"X-Intent": {"proxy-call"}})
	req = req.WithContext(ctx)

	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip returned error: %v", err)
	}
	resp.Body.Close()

	// The custom base transport must have observed the request...
	if cap.req == nil {
		t.Fatal("custom base transport was not used")
	}
	// ...and the request must have reached the server carrying all headers.
	if got := received.Get("X-Intent"); got != "proxy-call" {
		t.Errorf("custom header not seen by server: X-Intent = %q", got)
	}
	if got := received.Get(headerAuthorization); got != "Bearer key" {
		t.Errorf("Authorization not seen by server: %q", got)
	}
	if got := received.Get(headerBackendID); got != "backend" {
		t.Errorf("X-Backend-Id not seen by server: %q", got)
	}
}

func TestWithHeadersOverrideStoresHeadersOnContext(t *testing.T) {
	op := &runtime.ClientOperation{}
	WithHeadersOverride(http.Header{"X-Multi": {"a", "b"}})(op)

	stored, ok := op.Context.Value(requestHeadersKey).(http.Header)
	if !ok {
		t.Fatal("WithHeadersOverride did not store headers on the operation context")
	}
	if got := stored.Values("X-Multi"); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("X-Multi = %v, want [a b]", got)
	}
}

// TestWithHeadersOverrideThroughGeneratedClient drives a real generated client
// operation through NewSDKClient against a test server, verifying that the
// per-request option's headers and the default headers all reach the wire.
func TestWithHeadersOverrideThroughGeneratedClient(t *testing.T) {
	var received http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, err := NewSDKClient("key", "backend", server.URL)
	if err != nil {
		t.Fatalf("NewSDKClient returned error: %v", err)
	}

	params := metricsclient.NewMetricsQueryParams().
		WithContext(context.Background()).
		WithBody(&models.QueryRequest{
			Start:     strfmt.DateTime(time.Unix(0, 0)),
			End:       strfmt.DateTime(time.Unix(60, 0)),
			QueryType: "instant",
			Promql:    "up",
			Step:      "1m",
		})

	if _, err := c.Metrics.MetricsQuery(params, nil, WithHeadersOverride(http.Header{
		"X-Example": {"example-value"},
	})); err != nil {
		t.Fatalf("MetricsQuery returned error: %v", err)
	}

	if got := received.Get("X-Example"); got != "example-value" {
		t.Errorf("per-request header not seen by server: X-Example = %q", got)
	}
	if got := received.Get(headerAuthorization); got != "Bearer key" {
		t.Errorf("Authorization not seen by server: %q", got)
	}
}
