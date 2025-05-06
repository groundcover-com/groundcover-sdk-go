package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// Import the generated client and custom transport
	pkgClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	pkgMetrics "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"

	pkgModels "github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"

	"github.com/PuerkitoBio/rehttp"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
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

	traceparent := os.Getenv("GC_TRACEPARENT")
	// Example: Enable Gzip for requests (could be config-driven)
	enableGzip := true

	// Parse baseURL for go-openapi transport config
	parsedURL, err := url.Parse(baseURLStr)
	if err != nil {
		log.Fatalf("Error parsing GC_BASE_URL: %v", err)
	}

	host := parsedURL.Host
	basePath := parsedURL.Path
	if basePath == "" {
		basePath = pkgClient.DefaultBasePath // Use default if path is empty
	}
	// Ensure basePath starts with a slash if it's not just "/"
	if !strings.HasPrefix(basePath, "/") && basePath != "" {
		basePath = "/" + basePath
	}

	schemes := []string{parsedURL.Scheme}
	if len(schemes) == 0 || schemes[0] == "" {
		schemes = pkgClient.DefaultSchemes // Use default if scheme is missing
	}

	// --- Transport Stack Construction ---

	// 1. Base HTTP Transport (from go's net/http)
	//    Configure default timeout here.
	baseHttpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		// Add other standard transport settings if needed (e.g., TLS config)
	}

	// 2. Retry Transport (rehttp)
	//    Wraps the base transport to add retry logic.
	retryTransport := rehttp.NewTransport(
		baseHttpTransport, // Pass the base transport (RoundTripper)
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(defaultRetryCount),
			rehttp.RetryStatuses(http.StatusServiceUnavailable, http.StatusTooManyRequests),
		),
		rehttp.ExpJitterDelay(minRetryWait, maxRetryWait),
	)

	// 3. Custom Header Transport (our pkg/transport)
	//    Wraps the retry transport to add custom headers.
	customTransportWrapper := transport.NewCustomTransport(apiKey, backendID, traceparent, enableGzip, retryTransport)

	// 4. Final OpenAPI Runtime Transport
	//    Uses the fully wrapped transport stack.
	//    We need to create the runtime transport and then set its underlying RoundTripper.
	finalRuntimeTransport := httptransport.New(host, basePath, schemes)
	// Set the underlying RoundTripper (Transport field of http.Client)
	// to our custom wrapper which includes retries.
	// Note: We replace the default client/transport created by httptransport.New
	finalRuntimeTransport.Transport = customTransportWrapper

	// Ensure the underlying HTTP client used by the runtime transport has our settings
	// This step might be redundant if httptransport directly uses the provided RoundTripper, but it's safer.
	// If httptransport creates its own internal client, we might need to configure that client instead.
	// Testing would be required to confirm the exact behavior of httptransport.New.
	// For now, we assume setting finalRuntimeTransport.Transport is sufficient.

	// --- Client Initialization ---
	client := pkgClient.New(finalRuntimeTransport, strfmt.Default)

	// --- API Call 1: Metrics Query ---

	// Create a base context for requests
	baseCtx := context.Background()

	logrus.Info("--- Calling Metrics Query ---")
	// Example: Override settings using context for metrics query
	metricsCtx := transport.WithRequestGzip(baseCtx, false)                                      // Disable Gzip for this request
	metricsCtx = transport.WithRequestTraceparent(metricsCtx, "00-testtraceid-metricsspanid-01") // Specific traceparent

	// Prepare the request body for the metrics query
	startTime := strfmt.DateTime(time.Now().Add(-time.Hour))
	endTime := strfmt.DateTime(time.Now())
	step := "30s"
	queryType := "instant"
	promqlQuery := "avg(groundcover_container_cpu_limit_m_cpu)"

	queryRequestBody := &pkgModels.QueryRequest{
		Start:     startTime,
		End:       endTime,
		Step:      step,
		QueryType: queryType,
		Promql:    promqlQuery,
	}

	// Prepare the parameters for metrics query
	metricsParams := pkgMetrics.NewMetricsQueryParams().
		WithContext(metricsCtx). // Pass the specific context
		WithTimeout(defaultTimeout).
		WithBody(queryRequestBody)

	// Execute the metrics query
	queryResponse, err := client.Metrics.MetricsQuery(metricsParams)
	if err != nil {
		// Handle metrics query errors (using type switch as before)
		switch e := err.(type) {
		case *pkgMetrics.MetricsQueryBadRequest:
			logrus.Errorf("Metrics API Error (Bad Request): %s", e.Error())
		case *pkgMetrics.MetricsQueryInternalServerError:
			logrus.Errorf("Metrics API Error (Internal Server Error): %s", e.Error())
		default:
			if apiErr, ok := err.(*runtime.APIError); ok {
				logrus.Errorf("Metrics Generic API Error: Code %d, Response: %v", apiErr.Code, apiErr.Response)
			} else {
				logrus.Errorf("Error executing metrics query: %v", err)
			}
		}
		// Decide if you want to return or continue after a metrics error
		// return
	} else {
		// Handle the successful metrics response payload
		logrus.Info("Metrics Query Response:")
		spew.Dump(queryResponse) // queryResponse is the payload
	}

	// --- API Call 2: List Policies ---

	// logrus.Info("\n--- Calling List Policies ---")
	// // Use the base context (or create a new one with different overrides if needed)
	// policiesCtx := baseCtx
	// // Example: Use client defaults for Gzip/Traceparent, just set timeout.
	// policiesParams := pkgPolicies.NewListPoliciesParams().
	// 	WithContext(policiesCtx).
	// 	WithTimeout(defaultTimeout)
	// // Note: ListPolicies likely doesn't need a request body or many parameters,
	// // but check NewListPoliciesOperationParams definition if specific filters etc. are needed.

	// // Execute the list policies request
	// policiesResponse, err := client.Policies.ListPolicies(policiesParams)
	// if err != nil {
	// 	// Handle list policies errors
	// 	switch apiErr := err.(type) { // Assign asserted type to apiErr
	// 	case *pkgPolicies.ListPoliciesInternalServerError:
	// 		// apiErr is already *pkgPolicies.ListPoliciesOperationInternalServerError
	// 		logrus.Errorf("Policies API Error (Internal Server Error): %s", apiErr.Error())
	// 		// Potentially inspect apiErr payload if defined
	// 	default:
	// 		// Check for generic runtime.APIError within the default case
	// 		if genericAPIErr, ok := err.(*runtime.APIError); ok {
	// 			logrus.Errorf("Policies Generic API Error: Code %d, Response: %v", genericAPIErr.Code, genericAPIErr.Response)
	// 		} else {
	// 			// Log the original error if it's not a known API error type
	// 			logrus.Errorf("Error executing list policies: %v", err)
	// 		}
	// 	}
	// 	return // Exit after error
	// }

	// // Handle the successful list policies response payload
	// logrus.Info("List Policies Response:")
	// spew.Dump(policiesResponse) // policiesResponse IS the payload

}
