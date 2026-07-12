package transport

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	client "github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
)

type contextKey int

const (
	traceparentOverrideKey contextKey = iota
	requestHeadersKey
)

const (
	headerAuthorization = "Authorization"
	headerBackendID     = "X-Backend-Id"
	headerUserAgent     = "User-Agent"
	headerTraceparent   = "traceparent"
	userAgent           = "groundcover-go-sdk"
	yamlContentType     = "application/x-yaml"
)

const (
	defaultRetryCount = 3
	minRetryWait      = 1 * time.Second
	maxRetryWait      = 30 * time.Second
)

// idempotentRetryMethods are the HTTP methods that are safe to replay after the
// server has already returned a response.
var idempotentRetryMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodTrace,
}

var getMonitorPathRegex = regexp.MustCompile(`^/api/monitors/[^/]+/?$`) // Matches /api/monitors/{id} but not /api/monitors/silences

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
				return fmt.Errorf("yamlByteConsumer requires data to be *[]byte, *[]uint8, or *strfmt.Base64, got %T for content type %s", data, yamlContentType)
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

// NewYamlByteConsumer returns a new instance of the YAML byte consumer
// for use with go-openapi runtime transport consumers
func NewYamlByteConsumer() *yamlByteConsumer {
	return &yamlByteConsumer{}
}

// ConfigureRuntimeTransport configures the provided runtime transport with
// the SDK's standard consumers and settings
func ConfigureRuntimeTransport(rt *httptransport.Runtime) {
	// Register the YAML byte consumer for application/x-yaml content type
	rt.Consumers[yamlContentType] = NewYamlByteConsumer()
}

// NewConfiguredRuntimeTransport creates a new runtime transport with
// the SDK's standard configuration applied
func NewConfiguredRuntimeTransport(host, basePath string, schemes []string) *httptransport.Runtime {
	rt := httptransport.New(host, basePath, schemes)
	ConfigureRuntimeTransport(rt)
	return rt
}

// ClientOption allows customization of the SDK client
type ClientOption func(*clientConfig)

type clientConfig struct {
	httpTransport    http.RoundTripper
	retryCount       int
	minWait          time.Duration
	maxWait          time.Duration
	retryStatuses    []int
	transportWrapper func(http.RoundTripper) http.RoundTripper
}

// WithHTTPTransport sets a custom HTTP transport
func WithHTTPTransport(transport http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		c.httpTransport = transport
	}
}

// WithRetryConfig sets custom retry configuration
func WithRetryConfig(retryCount int, minWait, maxWait time.Duration, retryStatuses []int) ClientOption {
	return func(c *clientConfig) {
		c.retryCount = retryCount
		c.minWait = minWait
		c.maxWait = maxWait
		c.retryStatuses = retryStatuses
	}
}

// WithTransportWrapper allows wrapping the transport (e.g., for debugging)
func WithTransportWrapper(wrapper func(http.RoundTripper) http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		c.transportWrapper = wrapper
	}
}

// NewSDKClient creates a fully configured groundcover SDK client with all
// standard configurations applied automatically. Use options to customize behavior.
func NewSDKClient(apiKey, backendID, baseURL string, options ...ClientOption) (*client.GroundcoverAPI, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}

	// Create default config
	config := &clientConfig{
		retryCount:    defaultRetryCount,
		minWait:       minRetryWait,
		maxWait:       maxRetryWait,
		retryStatuses: []int{503, 429}, // Service Unavailable, Too Many Requests
	}

	// Apply options
	for _, option := range options {
		option(config)
	}

	// Normalize and parse baseURL
	normalizedURL := normalizeBaseURL(baseURL)
	parsedURL, err := url.Parse(normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
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

	// Create transport with SDK functionality
	sdkTransport := NewTransport(
		apiKey,
		backendID,
		config.httpTransport,
		config.retryCount,
		config.minWait,
		config.maxWait,
		config.retryStatuses,
	)

	// Apply custom transport wrapper if provided
	finalTransport := http.RoundTripper(sdkTransport)
	if config.transportWrapper != nil {
		finalTransport = config.transportWrapper(sdkTransport)
	}

	// Create runtime transport with SDK configurations
	runtimeTransport := NewConfiguredRuntimeTransport(host, basePath, schemes)
	runtimeTransport.Transport = finalTransport

	// Create and return client
	return client.New(runtimeTransport, strfmt.Default), nil
}

// WithRequestTraceparent returns a new context with the Traceparent override.
func WithRequestTraceparent(ctx context.Context, traceparent string) context.Context {
	return context.WithValue(ctx, traceparentOverrideKey, traceparent)
}

// withRequestHeaders returns a new context carrying custom headers to apply to
// the outgoing request. It is the internal mechanism behind WithHeadersOverride;
// the transport reads these headers in RoundTrip and applies them after the
// default headers using override (not append) semantics.
func withRequestHeaders(ctx context.Context, headers http.Header) context.Context {
	merged := http.Header{}
	if existing, ok := ctx.Value(requestHeadersKey).(http.Header); ok {
		for k, v := range existing {
			merged[k] = append([]string(nil), v...)
		}
	}
	for k, v := range headers {
		merged[http.CanonicalHeaderKey(k)] = append([]string(nil), v...)
	}
	return context.WithValue(ctx, requestHeadersKey, merged)
}

// WithHeadersOverride returns a per-request option that attaches custom headers
// to a single operation call. It is accepted by the generated client methods
// alongside their other options, for example:
//
//	resp, err := client.Metrics.MetricsQuery(params, nil, transport.WithHeadersOverride(headers))
//
// Headers use override (not append) semantics: each header provided here replaces
// any existing values for that name, multi-valued headers are supported, and the
// per-request headers take precedence over the SDK's default headers. The return
// type matches the generated clients' option type (func(*runtime.ClientOperation)).
func WithHeadersOverride(headers http.Header) func(*runtime.ClientOperation) {
	return func(op *runtime.ClientOperation) {
		ctx := op.Context
		if ctx == nil {
			ctx = context.Background()
		}
		op.Context = withRequestHeaders(ctx, headers)
	}
}

// transport wraps an existing http.RoundTripper to add custom headers.
type transport struct {
	apiKey         string
	backendID      string
	retryTransport http.RoundTripper
}

// NewTransport creates a new transport.
// traceparent is optional and can be an empty string.
// retryCount, minWait, maxWait configure the retry mechanism.
// If retryCount, minWait, or maxWait are provided as 0, package defaults will be used.
func NewTransport(
	apiKey, backendID string,
	baseHttpTransport http.RoundTripper, // This is the transport *before* retries
	retryCount int,
	minWait, maxWait time.Duration,
	retryStatuses []int,
) *transport {
	if baseHttpTransport == nil {
		baseHttpTransport = http.DefaultTransport
	}

	// Default retry statuses if not provided or empty
	if len(retryStatuses) == 0 {
		retryStatuses = []int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway}
	}

	// Apply package defaults if parameters are zero
	if retryCount <= 0 {
		retryCount = defaultRetryCount
	}

	if minWait <= 0 {
		minWait = minRetryWait
	}

	if maxWait <= 0 {
		maxWait = maxRetryWait
	}

	// Configure retry transport. Status-based retries are restricted to
	// idempotent methods (see idempotentRetryMethods).
	rt := rehttp.NewTransport(
		baseHttpTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(retryCount),
			rehttp.RetryHTTPMethods(idempotentRetryMethods...),
			rehttp.RetryAny(
				rehttp.RetryStatuses(retryStatuses...),
				rehttp.RetryIsErr(func(err error) bool { return err != nil }),
			),
		),
		rehttp.ExpJitterDelay(minWait, maxWait),
	)

	return &transport{
		apiKey:         apiKey,
		backendID:      backendID,
		retryTransport: rt,
	}
}

// RoundTrip executes a single HTTP transaction, checking context for overrides.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	var effectiveTraceparent string
	if traceVal, ok := ctx.Value(traceparentOverrideKey).(string); ok {
		effectiveTraceparent = traceVal
	}

	// Clone the request to avoid modifying the original passed to the base transport
	newReq := req.Clone(ctx)

	// --- Add Default Headers ---
	// Authorization and backend ID are only set when configured, so the SDK can
	// be used without an API key (for example, when routing through a transport
	// that supplies its own credentials).
	if t.apiKey != "" {
		newReq.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", t.apiKey))
	}
	if t.backendID != "" {
		newReq.Header.Set(headerBackendID, t.backendID)
	}
	newReq.Header.Set(headerUserAgent, userAgent)

	if effectiveTraceparent != "" {
		newReq.Header.Set(headerTraceparent, effectiveTraceparent)
	}

	// --- Fix Content-Type for specific endpoints ---
	// Fix request Content-Type for workflow create endpoint
	if newReq.Method == http.MethodPost && newReq.URL.Path == "/api/workflows/create" {
		newReq.Header.Set("Content-Type", "text/plain")
	}

	// --- Apply Custom Per-Request Headers ---
	// Applied last so callers can override any default or endpoint-specific
	// header set above. Override (not append) semantics: each provided header
	// replaces any existing values for that name. Keys are already canonical
	// (see withRequestHeaders), matching the default headers set above.
	if headers, ok := ctx.Value(requestHeadersKey).(http.Header); ok {
		for k, vals := range headers {
			newReq.Header[k] = append([]string(nil), vals...)
		}
	}

	// Execute the request
	resp, err := t.retryTransport.RoundTrip(newReq)
	if err != nil {
		return nil, err
	}

	// Fix response Content-Type for monitor GET endpoints
	if newReq.Method == http.MethodGet && resp.StatusCode == http.StatusOK &&
		getMonitorPathRegex.MatchString(newReq.URL.Path) &&
		!strings.Contains(newReq.URL.Path, "silences") {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, yamlContentType) {
			resp.Header.Set("Content-Type", yamlContentType)
		}
	}

	return resp, nil
}

func normalizeBaseURL(baseURL string) string {
	if baseURL == "" {
		return ""
	}

	// Handle URLs starting with // (protocol-relative URLs)
	if strings.HasPrefix(baseURL, "//") {
		return "https:" + baseURL
	}

	// Handle URLs starting with / (relative paths) - assume https and remove leading slash
	if strings.HasPrefix(baseURL, "/") && !strings.HasPrefix(baseURL, "//") {
		return "https://" + strings.TrimPrefix(baseURL, "/")
	}

	// Handle URLs without scheme
	if !strings.Contains(baseURL, "://") {
		return "https://" + baseURL
	}

	return baseURL
}

// NewClient creates a new groundcover SDK client with a simplified API.
// It automatically reads configuration from environment variables unless overridden by options.
//
// Environment variables:
//   - GC_API_KEY: Your groundcover API key (required, unless option.AllowUnauthenticated is used)
//   - GC_BACKEND_ID: Your groundcover Backend ID (required, unless option.AllowUnauthenticated is used)
//   - GC_BASE_URL: The base URL of the groundcover API (optional, defaults to https://api.groundcover.com)
//
// The API key and backend ID are required by default. Pass option.AllowUnauthenticated
// to create a client without them; their headers are then not set and the client
// can be used with a custom transport that supplies its own credentials.
//
// Example usage:
//
//	client := transport.NewClient()
//	client := transport.NewClient(option.WithAPIKey("custom-key"))
func NewClient(options ...option.Option) (*client.GroundcoverAPI, error) {
	// Set up default configuration from environment variables
	config := &option.Config{
		APIKey:    os.Getenv("GC_API_KEY"),
		BackendID: os.Getenv("GC_BACKEND_ID"),
		BaseURL:   os.Getenv("GC_BASE_URL"),
	}

	// Apply provided options
	for _, opt := range options {
		opt(config)
	}

	// Set default base URL if not provided
	if config.BaseURL == "" {
		config.BaseURL = "https://api.groundcover.com"
	}

	// Validate required fields unless the caller opts out. When opted out, the
	// API key and backend ID headers are simply not set, allowing the SDK to be
	// used with a custom transport that supplies its own credentials; the server
	// will reject the request if it requires authentication that was not provided.
	if !config.AllowUnauthenticated {
		if config.APIKey == "" {
			return nil, fmt.Errorf("API key is required: set GC_API_KEY environment variable or use option.WithAPIKey()")
		}
		if config.BackendID == "" {
			return nil, fmt.Errorf("backend ID is required: set GC_BACKEND_ID environment variable or use option.WithBackendID()")
		}
	}

	// Convert to legacy ClientOption format
	var clientOptions []ClientOption

	if config.HTTPTransport != nil {
		clientOptions = append(clientOptions, WithHTTPTransport(config.HTTPTransport))
	}

	if config.RetryCount > 0 || config.MinWait > 0 || config.MaxWait > 0 || len(config.RetryStatuses) > 0 {
		clientOptions = append(clientOptions, WithRetryConfig(
			config.RetryCount,
			config.MinWait,
			config.MaxWait,
			config.RetryStatuses,
		))
	}

	if config.TransportWrapper != nil {
		clientOptions = append(clientOptions, WithTransportWrapper(config.TransportWrapper))
	}

	// Use the existing NewSDKClient function
	return NewSDKClient(config.APIKey, config.BackendID, config.BaseURL, clientOptions...)
}
