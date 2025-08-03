package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

func main() {
	fmt.Println("groundcover Go SDK - Client Options Example")
	fmt.Println("==========================================")

	// Example 1: Simple client creation (using environment variables)
	fmt.Println("\n1. Simple client creation (environment variables)...")
	simpleClient()

	// Example 2: Client with custom options
	fmt.Println("\n2. Client with custom options...")
	clientWithOptions()

	// Example 3: Legacy transport-based client
	fmt.Println("\n3. Legacy transport-based client...")
	legacyClient()

	fmt.Println("\n✓ All client configuration examples completed")
}

func simpleClient() {
	// Create client with environment variables (GC_API_KEY, GC_BACKEND_ID, GC_BASE_URL)
	client, err := groundcover.NewClient()
	if err != nil {
		log.Printf("Failed to create simple client: %v", err)
		return
	}

	fmt.Println("✓ Simple client created successfully")

	// Test with a simple query
	err = testClient(client)
	if err != nil {
		log.Printf("Error testing simple client: %v", err)
	}
}

func clientWithOptions() {
	// Create client with custom options
	client, err := groundcover.NewClient(
		option.WithAPIKey("your-api-key"),
		option.WithBackendID("your-backend-id"),
		option.WithBaseURL("https://api.groundcover.com"),
		option.WithRetryConfig(
			3,              // retry count
			1*time.Second,  // min wait
			10*time.Second, // max wait
			[]int{http.StatusServiceUnavailable, http.StatusTooManyRequests},
		),
		// Note: No timeout option available in current SDK version
	)
	if err != nil {
		log.Printf("Failed to create client with options: %v", err)
		return
	}

	fmt.Println("✓ Client with custom options created successfully")

	// Test with a simple query
	err = testClient(client)
	if err != nil {
		log.Printf("Error testing client with options: %v", err)
	}
}

func legacyClient() {
	// Legacy method using transport directly
	// Note: You would typically get these from environment variables
	baseURL := "https://api.groundcover.com"
	apiKey := "your-api-key"
	backendID := "your-backend-id"

	// Create a fully configured client - handles auth, retries, content-type fixes, etc.
	sdkClient, err := transport.NewSDKClient(apiKey, backendID, baseURL)
	if err != nil {
		log.Printf("Failed to create legacy SDK client: %v", err)
		return
	}

	fmt.Println("✓ Legacy transport-based client created successfully")

	// Test with a simple query
	err = testClient(sdkClient)
	if err != nil {
		log.Printf("Error testing legacy client: %v", err)
	}
}

func testClient(client *client.GroundcoverAPI) error {
	ctx := context.Background()

	queryRequest := &models.QueryRequest{
		Start:     strfmt.DateTime(time.Now().Add(-5 * time.Minute)),
		End:       strfmt.DateTime(time.Now()),
		QueryType: "instant",
		Promql:    "avg(groundcover_container_cpu_limit_m_cpu)",
		Step:      "1m",
	}

	params := metrics.NewMetricsQueryParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(queryRequest)

	_, err := client.Metrics.MetricsQuery(params, nil)
	if err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	fmt.Println("✓ Test query successful")
	return nil
}
