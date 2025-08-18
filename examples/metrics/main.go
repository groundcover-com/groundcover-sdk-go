package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func main() {
	fmt.Println("groundcover Go SDK - Basic Usage Example")
	fmt.Println("=======================================")

	// Create client using environment variables
	// Required: GC_API_KEY, GC_BACKEND_ID, GC_BASE_URL
	gcClient, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	// Example 1: Simple instant query
	fmt.Println("\n1. Simple instant query...")
	err = simpleInstantQuery(gcClient)
	if err != nil {
		log.Printf("Error in instant query: %v", err)
	}

	// Example 2: Simple range query
	fmt.Println("\n2. Simple range query...")
	err = simpleRangeQuery(gcClient)
	if err != nil {
		log.Printf("Error in range query: %v", err)
	}

	fmt.Println("\n✓ Basic usage examples completed")
}

func simpleInstantQuery(client *client.GroundcoverAPI) error {
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

	resp, err := client.Metrics.MetricsQuery(params, nil)
	if err != nil {
		return fmt.Errorf("metrics query failed: %w", err)
	}

	// Marshal and print the response
	responseJSON, err := json.MarshalIndent(resp.Payload, "", "  ")
	if err != nil {
		fmt.Printf("✓ Instant query successful, but failed to marshal response: %v\n", err)
		fmt.Printf("Raw response type: %T\n", resp.Payload)
		return nil
	}

	fmt.Printf("✓ Instant query successful\n")
	fmt.Printf("Response:\n%s\n", string(responseJSON))
	return nil
}

func simpleRangeQuery(client *client.GroundcoverAPI) error {
	ctx := context.Background()

	queryRequest := &models.QueryRequest{
		Start:     strfmt.DateTime(time.Now().Add(-30 * time.Minute)),
		End:       strfmt.DateTime(time.Now()),
		QueryType: "range",
		Promql:    "avg(groundcover_container_cpu_limit_m_cpu)",
		Step:      "5m",
	}

	params := metrics.NewMetricsQueryParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(queryRequest)

	resp, err := client.Metrics.MetricsQuery(params, nil)
	if err != nil {
		return fmt.Errorf("metrics query failed: %w", err)
	}

	// Marshal and print the response
	responseJSON, err := json.MarshalIndent(resp.Payload, "", "  ")
	if err != nil {
		fmt.Printf("✓ Range query successful, but failed to marshal response: %v\n", err)
		fmt.Printf("Raw response type: %T\n", resp.Payload)
		return nil
	}

	fmt.Printf("✓ Range query successful\n")
	fmt.Printf("Response:\n%s\n", string(responseJSON))
	return nil
}
