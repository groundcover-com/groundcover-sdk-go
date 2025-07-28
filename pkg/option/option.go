package option

import (
	"net/http"
	"time"
)

// Option represents a configuration option for the groundcover client.
type Option func(*Config)

// Config holds all configuration options for the client.
type Config struct {
	APIKey           string
	BackendID        string
	BaseURL          string
	HTTPTransport    http.RoundTripper
	RetryCount       int
	MinWait          time.Duration
	MaxWait          time.Duration
	RetryStatuses    []int
	TransportWrapper func(http.RoundTripper) http.RoundTripper
}

// WithAPIKey sets the API key for authentication.
// If not provided, defaults to the GC_API_KEY environment variable.
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithBackendID sets the backend ID for the client.
// If not provided, defaults to the GC_BACKEND_ID environment variable.
func WithBackendID(backendID string) Option {
	return func(c *Config) {
		c.BackendID = backendID
	}
}

// WithBaseURL sets the base URL for the groundcover API.
// If not provided, defaults to the GC_BASE_URL environment variable, or https://api.groundcover.com if not set.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithHTTPTransport sets a custom HTTP transport.
func WithHTTPTransport(transport http.RoundTripper) Option {
	return func(c *Config) {
		c.HTTPTransport = transport
	}
}

// WithRetryConfig sets custom retry configuration.
func WithRetryConfig(retryCount int, minWait, maxWait time.Duration, retryStatuses []int) Option {
	return func(c *Config) {
		c.RetryCount = retryCount
		c.MinWait = minWait
		c.MaxWait = maxWait
		c.RetryStatuses = retryStatuses
	}
}

// WithTransportWrapper allows wrapping the transport (e.g., for debugging).
func WithTransportWrapper(wrapper func(http.RoundTripper) http.RoundTripper) Option {
	return func(c *Config) {
		c.TransportWrapper = wrapper
	}
}
