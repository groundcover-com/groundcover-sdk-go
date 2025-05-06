package transport

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

// Define unexported types for context keys to prevent collisions
type contextKey int

const (
	gzipOverrideKey contextKey = iota
	traceparentOverrideKey
)

const (
	headerAuthorization   = "Authorization"
	headerBackendID       = "X-Backend-Id"
	headerUserAgent       = "User-Agent"
	headerTraceparent     = "traceparent"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	userAgent             = "groundcover-go-sdk" // Or a new user agent if preferred
	encodingGzip          = "gzip"
)

const (
	defaultRetryCount = 3
	minRetryWait      = 1 * time.Second
	maxRetryWait      = 30 * time.Second
)

// WithRequestGzip returns a new context with the Gzip override setting.
func WithRequestGzip(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, gzipOverrideKey, enabled)
}

// WithRequestTraceparent returns a new context with the Traceparent override.
func WithRequestTraceparent(ctx context.Context, traceparent string) context.Context {
	return context.WithValue(ctx, traceparentOverrideKey, traceparent)
}

// customTransport wraps an existing http.RoundTripper to add custom headers and optional request Gzip compression.
type customTransport struct {
	apiKey            string
	backendID         string
	clientTraceparent string // Renamed to avoid confusion with context override
	clientGzipEnabled bool   // Renamed to avoid confusion with context override
	retryTransport    http.RoundTripper
}

// NewCustomTransport creates a new customTransport.
// traceparent is optional and can be an empty string.
// retryCount, minWait, maxWait configure the retry mechanism.
func NewCustomTransport(
	apiKey, backendID, clientTraceparent string,
	clientGzipEnabled bool,
	baseHttpTransport http.RoundTripper, // This is the transport *before* retries
	retryCount int,
	minWait, maxWait time.Duration,
	retryStatuses []int,
) *customTransport {
	if baseHttpTransport == nil {
		baseHttpTransport = http.DefaultTransport
	}

	// Default retry statuses if not provided or empty
	if len(retryStatuses) == 0 {
		retryStatuses = []int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway}
	}

	// Configure retry transport
	rt := rehttp.NewTransport(
		baseHttpTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(retryCount),
			rehttp.RetryStatuses(retryStatuses...),
		),
		rehttp.ExpJitterDelay(minWait, maxWait),
	)

	return &customTransport{
		apiKey:            apiKey,
		backendID:         backendID,
		clientTraceparent: clientTraceparent,
		clientGzipEnabled: clientGzipEnabled,
		retryTransport:    rt,
	}
}

// RoundTrip executes a single HTTP transaction, checking context for overrides.
func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Determine effective settings based on context overrides or client defaults
	effectiveGzipEnabled := t.clientGzipEnabled
	if gzipVal, ok := ctx.Value(gzipOverrideKey).(bool); ok {
		effectiveGzipEnabled = gzipVal
	}

	effectiveTraceparent := t.clientTraceparent
	if traceVal, ok := ctx.Value(traceparentOverrideKey).(string); ok {
		effectiveTraceparent = traceVal
	}

	// Clone the request to avoid modifying the original passed to the base transport
	newReq := req.Clone(ctx)

	// --- Gzip Request Body (if effectiveGzipEnabled and body exists) ---
	if effectiveGzipEnabled && req.Body != nil && req.Body != http.NoBody {
		// Read the original body
		originalBodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading original request body for gzipping: %w", err)
		}
		req.Body.Close() // Close the original body reader

		// Compress the body
		var compressedBody bytes.Buffer
		gzw := gzip.NewWriter(&compressedBody)
		if _, err := gzw.Write(originalBodyBytes); err != nil {
			return nil, fmt.Errorf("error gzipping request body: %w", err)
		}
		if err := gzw.Close(); err != nil {
			return nil, fmt.Errorf("error closing gzip writer for request body: %w", err)
		}

		// Set the compressed body on the new request
		newReq.Body = io.NopCloser(&compressedBody)
		// Set Content-Encoding header
		newReq.Header.Set(headerContentEncoding, encodingGzip)
		// Remove Content-Length; http client will set it or use chunked encoding
		newReq.Header.Del(headerContentLength)
		newReq.ContentLength = int64(compressedBody.Len()) // Set ContentLength explicitly
	} else if req.Body != nil {
		// If not gzipping, rely on req.Clone providing a usable Body for newReq.
	}

	// --- Add Custom Headers ---
	newReq.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", t.apiKey))
	newReq.Header.Set(headerBackendID, t.backendID)
	newReq.Header.Set(headerUserAgent, userAgent)

	// Add traceparent header (using effective value from context or client default)
	if effectiveTraceparent != "" {
		newReq.Header.Set(headerTraceparent, effectiveTraceparent)
	}

	// Proceed with the request using the retry transport
	return t.retryTransport.RoundTrip(newReq)
}
