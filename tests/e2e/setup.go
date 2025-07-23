// Package e2e contains end-to-end tests for the groundcover SDK
package e2e

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultRetryCount = 5
	minRetryWait      = 1 * time.Second
	maxRetryWait      = 5 * time.Second
	YamlContentType   = "application/x-yaml" // Define YAML content type constant
)

// isDebugEnabled returns true if SDK_DEBUG environment variable is set to any value
func isDebugEnabled() bool {
	return os.Getenv("SDK_DEBUG") != ""
}

// DebugTransport wraps a RoundTripper and logs all requests and responses
type DebugTransport struct {
	transport http.RoundTripper
	testing   *testing.T
}

func (d *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	debug := isDebugEnabled()

	// Log the request if debug is enabled
	if debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			d.testing.Logf("Error dumping request: %v", err)
		} else {
			d.testing.Logf("REQUEST:\n%s", string(reqDump))
		}
	}

	// Execute the request
	resp, err := d.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the response if debug is enabled
	if debug {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			d.testing.Logf("Error dumping response: %v", err)
		} else {
			d.testing.Logf("RESPONSE:\n%s", string(respDump))
		}
	}

	// The response body is read and replaced here to ensure it remains readable
	// by subsequent handlers after the DebugTransport has processed it.
	// This is necessary because operations like httputil.DumpResponse (if debug is enabled)
	// or the act of reading the body itself for this buffering would consume the original stream.
	buf, readErr := io.ReadAll(resp.Body)
	resp.Body.Close() // Close the original body
	if readErr != nil {
		if debug {
			d.testing.Logf("Error reading response body after dump: %v", readErr)
		}
		// Return the response even if reading failed, might still be usable partially
		resp.Body = io.NopCloser(strings.NewReader("")) // Set an empty body
		return resp, err                                // Return the original transport error if any
	}

	// Set a new body with the same content
	resp.Body = io.NopCloser(strings.NewReader(string(buf)))

	return resp, err // Return the original transport error
}

// TestClient holds the client and context for testing
type TestClient struct {
	Client  *client.GroundcoverAPI
	BaseCtx context.Context
	Cleanup func()
	T       *testing.T
}

type testClientOptions struct {
	backendID string
}

type TestClientOption func(*testClientOptions)

func TestClientWithBackendID(backendID string) TestClientOption {
	return func(opts *testClientOptions) {
		opts.backendID = backendID
	}
}

// NewTestClient creates a new client for testing
func NewTestClient(t *testing.T, options ...TestClientOption) *TestClient {
	t.Helper()
	debug := isDebugEnabled()

	baseURLStr := os.Getenv("GC_BASE_URL")
	if baseURLStr == "" {
		t.Fatal("GC_BASE_URL environment variable is required")
	}

	apiKey := os.Getenv("GC_API_KEY")
	if apiKey == "" {
		t.Fatal("GC_API_KEY environment variable is required")
	}

	opts := &testClientOptions{
		backendID: os.Getenv("GC_BACKEND_ID"),
	}

	for _, option := range options {
		option(opts)
	}

	if opts.backendID == "" {
		t.Fatal("GC_BACKEND_ID environment variable is required")
	}

	traceparent := os.Getenv("GC_TRACEPARENT")
	if traceparent == "" {
		traceparent = generateTraceParent()
	}
	t.Logf("TraceID: %s", extractTraceID(traceparent))

	// Create debug transport wrapper if enabled
	var debugWrapper func(http.RoundTripper) http.RoundTripper
	if debug {
		debugWrapper = func(transport http.RoundTripper) http.RoundTripper {
			return &DebugTransport{
				transport: transport,
				testing:   t,
			}
		}
	}

	// Use our new simplified client creation with options
	var clientOptions []option.Option

	// Set the credentials explicitly (even though they're in env vars)
	clientOptions = append(clientOptions,
		option.WithAPIKey(apiKey),
		option.WithBackendID(opts.backendID),
		option.WithBaseURL(baseURLStr),
	)

	// Add custom HTTP transport
	baseHttpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	clientOptions = append(clientOptions, option.WithHTTPTransport(baseHttpTransport))

	// Add retry config
	clientOptions = append(clientOptions, option.WithRetryConfig(
		defaultRetryCount,
		minRetryWait,
		maxRetryWait,
		[]int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway},
	))

	// Add debug wrapper if enabled
	if debugWrapper != nil {
		clientOptions = append(clientOptions, option.WithTransportWrapper(debugWrapper))
	}

	// Create the client using our new simplified API
	sdkClient, err := groundcover.NewClient(clientOptions...)
	if err != nil {
		t.Fatalf("Failed to create SDK client: %v", err)
	}

	if debug {
		t.Logf("Created SDK client with automatic YAML consumer and content-type fixes")
	}

	// Create base context
	baseCtx := context.Background()

	// If a traceparent is provided via environment variable, add it to the base context for all test requests.
	if traceparent != "" {
		baseCtx = transport.WithRequestTraceparent(baseCtx, traceparent)
		if debug {
			t.Logf("- Applying default Traceparent to BaseCtx: %s", traceparent)
		}
	}

	// Create test client
	return &TestClient{
		Client:  sdkClient,
		BaseCtx: baseCtx,
		T:       t,
		Cleanup: func() {
			// Add cleanup logic here
		},
	}
}

// setupTestClient is a convenience wrapper around NewTestClient
// that returns the context and client directly for use in tests
func setupTestClient(t *testing.T, options ...TestClientOption) (context.Context, *client.GroundcoverAPI) {
	tc := NewTestClient(t, options...)
	return tc.BaseCtx, tc.Client
}

func generateTraceParent() string {
	// Generate 16 random bytes for the first hex section (32 hex chars)
	part1 := make([]byte, 16)
	rand.Read(part1)

	// Generate 8 random bytes for the second hex section (16 hex chars)
	part2 := make([]byte, 8)
	rand.Read(part2)

	// Format: 00-{32 hex chars}-{16 hex chars}-01
	return fmt.Sprintf("00-%x-%x-01", part1, part2)
}

func extractTraceID(traceParent string) string {
	// Split by hyphens and return the second part (index 1) - the 32-char trace ID
	parts := strings.Split(traceParent, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// this is a helper function to create the required env variables for NewTestClient() by code in your test (don't commit the apikey :-))
func createEnvVariablesForTest(apiUrl, apiKey, backendId, traceparent string) {
	os.Setenv("GC_BASE_URL", apiUrl)
	os.Setenv("GC_API_KEY", apiKey)
	os.Setenv("GC_BACKEND_ID", backendId)
	os.Setenv("GC_TRACEPARENT", traceparent)
}
