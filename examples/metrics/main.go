package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

type MetricsResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string         `json:"resultType"`
		Result     []MetricSeries `json:"result"`
	} `json:"data"`
}

type MetricSeries struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

func main() {
	fmt.Println("groundcover Go SDK - Metrics Example")
	fmt.Println("===================================")

	// Create client
	client, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	// Example: Get node count
	fmt.Println("Getting node count...")
	promql := "count(groundcover_node_capacity_mem_bytes{})"
	fmt.Printf("PromQL: %s\n", promql)

	queryRequest := &models.QueryRequest{
		Start:     strfmt.DateTime(time.Now().Add(-5 * time.Minute)),
		End:       strfmt.DateTime(time.Now()),
		QueryType: "instant",
		Promql:    promql,
		Step:      "1m",
	}

	params := metrics.NewMetricsQueryParams().
		WithContext(context.Background()).
		WithTimeout(30 * time.Second).
		WithBody(queryRequest)

	resp, err := client.Metrics.MetricsQuery(params, nil)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}

	printNodeCount(resp.Payload)
}

func printNodeCount(payload interface{}) {
	// Convert payload to JSON bytes and unmarshal into struct
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal response")
		return
	}

	var metricsResp MetricsResponse
	if err := json.Unmarshal(jsonBytes, &metricsResp); err != nil {
		fmt.Println("Failed to unmarshal response")
		return
	}

	if len(metricsResp.Data.Result) == 0 {
		fmt.Println("No results found")
		return
	}

	// Get the count value
	if len(metricsResp.Data.Result[0].Value) >= 2 {
		fmt.Printf("Result: %v nodes\n", metricsResp.Data.Result[0].Value[1])
	} else {
		fmt.Println("Invalid result format")
	}
}
