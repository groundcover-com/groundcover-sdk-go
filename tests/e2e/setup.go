// Package e2e contains end-to-end tests for the groundcover SDK
package e2e

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	client "github.com/groundcover-com/groundcover-sdk-go/pkg/client"
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

// FixContentTypeTransport wraps a RoundTripper and corrects the Content-Type for specific endpoints.
type FixContentTypeTransport struct {
	transport http.RoundTripper
	testing   *testing.T
}

var getMonitorPathRegex = regexp.MustCompile(`^/api/monitors/[^/]+/?$`) // Matches /api/monitors/{id} but not /api/monitors/silences

func (f *FixContentTypeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Fix request Content-Type for workflow create endpoint
	if req.Method == http.MethodPost && req.URL.Path == "/api/workflows/create" {
		originalContentType := req.Header.Get("Content-Type")
		req.Header.Set("Content-Type", "text/plain")
		if isDebugEnabled() {
			f.testing.Logf("Fixed request Content-Type for POST %s. Original: '%s', Set to: 'text/plain'",
				req.URL.Path, originalContentType)
		}
	}

	// Execute the request using the underlying transport
	resp, err := f.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Check if this is the response we need to fix (only for monitor GET endpoints, not silences)
	if req.Method == http.MethodGet && resp.StatusCode == http.StatusOK &&
		getMonitorPathRegex.MatchString(req.URL.Path) &&
		!strings.Contains(req.URL.Path, "silences") {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, "application/x-yaml") {
			if isDebugEnabled() {
				f.testing.Logf("Fixing Content-Type for GET %s. Original: '%s', Setting to '%s'", req.URL.Path, contentType, YamlContentType)
			}
			resp.Header.Set("Content-Type", YamlContentType)
		}
	}

	return resp, nil
}

// --- Custom YAML Consumer ---

// yamlByteConsumer consumes application/x-yaml as raw bytes
type yamlByteConsumer struct{}

// Consume reads the response body directly into data without YAML parsing.
// It expects 'data' to be a pointer to a []byte or similar slice type.
func (c *yamlByteConsumer) Consume(reader io.Reader, data interface{}) error {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Check if data is a pointer to []byte or *[]byte
	byteSlicePtr, ok := data.(*[]byte)
	if !ok {
		// Also handle *strfmt.Base64, which is essentially *[]byte
		base64Ptr, ok := data.(*strfmt.Base64)
		if !ok {
			// Fallback: Try assigning to []uint8 if that's the underlying type
			uint8SlicePtr, ok := data.(*[]uint8)
			if !ok {
				return fmt.Errorf("yamlByteConsumer requires data to be *[]byte, *[]uint8, or *strfmt.Base64, got %T for content type %s", data, YamlContentType)
			}
			*uint8SlicePtr = buf
			return nil
		}
		*base64Ptr = buf
		return nil
	}

	*byteSlicePtr = buf
	return nil
}

// --- End Custom YAML Consumer ---

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

	// Get environment variables
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

	// Parse baseURL for go-openapi transport config
	parsedURL, err := url.Parse(baseURLStr)
	if err != nil {
		t.Fatalf("Error parsing GC_BASE_URL: %v", err)
	}

	host := parsedURL.Host
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = client.DefaultBasePath
	}
	if !strings.HasPrefix(basePath, "/") && basePath != "" {
		basePath = "/" + basePath
	}

	schemes := []string{parsedURL.Scheme}
	if len(schemes) == 0 || schemes[0] == "" {
		schemes = client.DefaultSchemes
	}

	// Log detailed client configuration if debugging is enabled
	if debug {
		t.Logf("SDK Client Configuration:")
		t.Logf("- Host: %s", host)
		t.Logf("- Base Path: %s", basePath)
		t.Logf("- Schemes: %v", schemes)
		t.Logf("- Default BasePath: %s", client.DefaultBasePath)
		t.Logf("- Default Schemes: %v", client.DefaultSchemes)
	}

	// Transport setup
	baseHttpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	transportWrapper := transport.NewTransport(
		apiKey,
		opts.backendID,
		baseHttpTransport,
		defaultRetryCount,
		minRetryWait,
		maxRetryWait,
		[]int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway},
	)

	// Wrap with content type fixer
	contentTypeFixer := &FixContentTypeTransport{
		transport: transportWrapper,
		testing:   t,
	}

	// Wrap with debug transport (this should be the outermost wrapper if we want to see the final request/response)
	finalTransportLayer := http.RoundTripper(contentTypeFixer)
	if debug {
		finalTransportLayer = &DebugTransport{
			transport: contentTypeFixer,
			testing:   t,
		}
	}

	finalRuntimeTransport := httptransport.New(host, basePath, schemes)
	finalRuntimeTransport.Transport = finalTransportLayer

	// --- Register Custom Consumer ---
	// Add our custom consumer for YAML content type to prevent default parsing
	finalRuntimeTransport.Consumers[YamlContentType] = &yamlByteConsumer{}
	if debug {
		t.Logf("Registered custom YAML consumer for %s", YamlContentType)
	}
	// --- End Register Custom Consumer ---

	// Create client
	sdkClient := client.New(finalRuntimeTransport, strfmt.Default)

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

// this is a helper function to create the required env variables for NewTestClient() by code in your test (don't commit the apikey :-))
func createEnvVariablesForTest(apiUrl, apiKey, backendId, traceparent string) {
	os.Setenv("GC_BASE_URL", apiUrl)
	os.Setenv("GC_API_KEY", apiKey)
	os.Setenv("GC_BACKEND_ID", backendId)
	os.Setenv("GC_TRACEPARENT", traceparent)
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
