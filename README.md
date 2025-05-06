# groundcover Go SDK

This is the official Go SDK for interacting with the groundcover API. It provides convenient access to groundcover's services, including metrics queries and policy management.

## Prerequisites

*   Go 1.24 or higher.

## Installation

To use the SDK in your Go project, you can install it using `go get`:

```bash
go get github.com/groundcover-com/groundcover-sdk-go
```

## Configuration

### Environment Variables

The SDK requires the following environment variables to be set for authentication and endpoint configuration:

*   `GC_BASE_URL`: The base URL of the groundcover API (e.g., `https://api.groundcover.com`).
*   `GC_API_KEY`: Your groundcover API key.
*   `GC_BACKEND_ID`: Your groundcover Backend ID.

Optionally, you can set:

*   `GC_TRACEPARENT`: A default traceparent header value for distributed tracing.

### Client Initialization

The SDK client requires a configured transport stack that handles authentication, retries, and other custom behaviors. Here's an example of how to set up the client:

```go
package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	client "github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultRetryCount = 5
	minRetryWait      = 1 * time.Second
	maxRetryWait      = 5 * time.Second
)

func main() {
	baseURLStr := os.Getenv("GC_BASE_URL")
	if baseURLStr == "" {
		log.Fatal("GC_BASE_URL environment variable is required")
	}

	apiKey := os.Getenv("GC_API_KEY")
	if apiKey == "" {
		log.Fatal("GC_API_KEY environment variable is required")
	}

	backendID := os.Getenv("GC_BACKEND_ID")
	if backendID == "" {
		log.Fatal("GC_BACKEND_ID environment variable is required")
	}

	traceparent := os.Getenv("GC_TRACEPARENT") // Optional
	enableGzip := true                         // Example: Enable Gzip for requests by default

	// Parse baseURL for go-openapi transport config
	parsedURL, err := url.Parse(baseURLStr)
	if err != nil {
		log.Fatalf("Error parsing GC_BASE_URL: %v", err)
	}

	host := parsedURL.Host
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = client.DefaultBasePath // Use default if path is empty
	}
	if !strings.HasPrefix(basePath, "/") && basePath != "" {
		basePath = "/" + basePath
	}

	schemes := []string{parsedURL.Scheme}
	if len(schemes) == 0 || schemes[0] == "" {
		schemes = client.DefaultSchemes // Use default if scheme is missing
	}

	// --- Transport Stack Construction ---

	// 1. Base HTTP Transport (from go's net/http)
	baseHttpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// Configure other standard transport settings if needed (e.g., TLS, timeouts)
	}

	// 2. Custom Transport (handles auth, Gzip, retries)
	//    The retry logic is built into NewTransport.
	transportWrapper := transport.NewTransport(
		apiKey,
		backendID,
		traceparent,
		enableGzip,
		baseHttpTransport, // The underlying transport
		defaultRetryCount, // Max number of retries
		minRetryWait,      // Minimum wait time between retries
		maxRetryWait,      // Maximum wait time between retries
		// Specify statuses to retry on, or pass nil/empty for defaults
		[]int{http.StatusServiceUnavailable, http.StatusTooManyRequests, http.StatusGatewayTimeout, http.StatusBadGateway},
	)

	// 3. Final OpenAPI Runtime Transport
	finalRuntimeTransport := client.New(host, basePath, schemes)
	finalRuntimeTransport.Transport = transportWrapper // Set our custom transport

	// --- Client Initialization ---
	sdkClient := client.New(finalRuntimeTransport, strfmt.Default)

	// Now you can use sdkClient to make API calls
	// Example: sdkClient.Metrics.MetricsQuery(...)
}
```

## Usage

### Making an API Call

Here's an example of how to make a metrics query:

```go
	// (Continued from Client Initialization above)
	// --- API Call: Metrics Query ---
	// import models "github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	// import metrics "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	// import "github.com/sirupsen/logrus"
	// import "github.com/davecgh/go-spew/spew"

	baseCtx := context.Background()

	logrus.Info("--- Calling Metrics Query ---")

	// Prepare the request body for the metrics query
	startTime := strfmt.DateTime(time.Now().Add(-time.Hour))
	endTime := strfmt.DateTime(time.Now())
	step := "30s"
	queryType := "instant"
	promqlQuery := "avg(groundcover_container_cpu_limit_m_cpu)"

	queryRequestBody := &models.QueryRequest{
		Start:     startTime,
		End:       endTime,
		Step:      step,
		QueryType: queryType,
		Promql:    promqlQuery,
	}

	// Prepare the parameters for metrics query
	metricsParams := metrics.NewMetricsQueryParams().
		WithContext(baseCtx).
		WithTimeout(defaultTimeout). // Overall request timeout
		WithBody(queryRequestBody)

	// Execute the metrics query
	queryResponse, err := sdkClient.Metrics.MetricsQuery(metricsParams)
	if err != nil {
		// Handle errors (see Error Handling section)
		logrus.Errorf("Error executing metrics query: %v", err)
		return
	}

	// Handle the successful metrics response payload
	logrus.Info("Metrics Query Response:")
	spew.Dump(queryResponse.Payload) // queryResponse.Payload contains the data
```

### Context for Request Overrides

The `pkg/transport` module provides functions to override default client settings for individual requests using `context.Context`.

*   **Gzip Compression**: Enable or disable Gzip compression for a specific request.
    ```go
    // Disable Gzip for this specific request
    metricsCtx := transport.WithRequestGzip(baseCtx, false)
    // ... then use metricsCtx in NewMetricsQueryParams().WithContext(metricsCtx)
    ```

*   **Traceparent**: Set a specific `traceparent` header for a request.
    ```go
    // Set a specific traceparent for this request
    metricsCtx := transport.WithRequestTraceparent(baseCtx, "00-customtraceid-customspanid-01")
    // ... then use metricsCtx in NewMetricsQueryParams().WithContext(metricsCtx)
    ```

### Retry Mechanism

The SDK's custom transport has a built-in retry mechanism that automatically retries requests on transient server errors (e.g., `503 Service Unavailable`, `429 Too Many Requests`). This is configured during client initialization via `transport.NewTransport`.

### Error Handling

API calls can return errors. It's important to handle these appropriately. The SDK uses specific error types for different API responses, and also a generic `runtime.APIError`.

```go
	// import "github.com/go-openapi/runtime"
	// import metrics "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"

	// (inside an API call block like the metrics query example)
	// queryResponse, err := sdkClient.Metrics.MetricsQuery(metricsParams)
	if err != nil {
		switch e := err.(type) {
		case *metrics.MetricsQueryBadRequest: // Example specific error
			logrus.Errorf("Metrics API Error (Bad Request): %s, Payload: %v", e.Error(), e.Payload)
		case *metrics.MetricsQueryInternalServerError: // Example specific error
			logrus.Errorf("Metrics API Error (Internal Server Error): %s, Payload: %v", e.Error(), e.Payload)
		default:
			if apiErr, ok := err.(*runtime.APIError); ok {
				// This is a generic error from the go-openapi runtime
				// apiErr.Code gives the HTTP status code
				// apiErr.Response gives the raw response body (needs to be parsed or read)
				logrus.Errorf("Generic API Error: Code %d, Message: %s, Response: %v", apiErr.Code, apiErr.Error(), apiErr.Response)
			} else {
				// Other unexpected errors
				logrus.Errorf("Error executing API call: %v", err)
			}
		}
		return // Or handle as appropriate
	}
	// Process successful response: queryResponse.Payload
```

## Available Services

The SDK is organized by service, available under the `sdkClient` object. For example:

*   `sdkClient.Metrics`: For querying metrics.
*   `sdkClient.Policies`: For managing policies (example was commented out in `main.go` but shows the pattern).

Refer to the generated SDK code in the `pkg/client` directory for a full list of services and their operations.

## License

This SDK is distributed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for more information.

