package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/logs"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

// ANSI color codes
const (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorYellow    = "\033[33m"
	ColorGreen     = "\033[32m"
	ColorBlue      = "\033[34m"
	ColorGray      = "\033[90m"
	ColorBrightRed = "\033[91m"
)

func getColorByLevel(level string) string {
	switch level {
	case "error":
		return ColorRed
	case "critical", "fatal":
		return ColorBrightRed
	case "warning", "warn":
		return ColorYellow
	case "info":
		return ColorGreen
	case "debug":
		return ColorGray
	case "trace":
		return ColorGray
	default:
		return ColorReset
	}
}

func main() {
	fmt.Println("groundcover Go SDK - Logs Example")
	fmt.Println("=================================")

	// Create client
	client, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	ctx := context.Background()

	// Example: Error log search
	fmt.Println("\nSearching for error logs...")
	if err := errorLogSearch(ctx, client); err != nil {
		fmt.Printf("Error searching error logs: %v\n", err)
	}
}

func errorLogSearch(ctx context.Context, gcClient *client.GroundcoverAPI) error {
	// Search for error-level logs using Query
	endTime := time.Now()
	startTime := endTime.Add(-2 * time.Hour)

	startDateTime := strfmt.DateTime(startTime)
	endDateTime := strfmt.DateTime(endTime)

	searchRequest := &models.LogsSearchRequest{
		Start: &startDateTime,
		End:   &endDateTime,
		// Search for error level logs, limit 10 results
		Query: "level:error | limit 10",
	}

	params := logs.NewSearchLogsParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(searchRequest)

	resp, err := gcClient.Logs.SearchLogs(params, nil)
	if err != nil {
		return fmt.Errorf("failed to search error logs: %w", err)
	}

	return printLogsResponse("Error Log Search", resp.Payload)
}

func printLogsResponse(searchName string, payload interface{}) error {
	fmt.Printf("  %s:\n", searchName)

	// Payload is directly an array of log objects
	logs, ok := payload.([]interface{})
	if !ok {
		fmt.Println("    Unexpected payload format")
		return nil
	}

	if len(logs) == 0 {
		fmt.Println("    No logs found")
		return nil
	}

	// Print each log in simple format
	for _, logItem := range logs {
		logEntry, ok := logItem.(map[string]interface{})
		if !ok {
			continue
		}

		timestamp := "unknown"
		if ts, hasTimestamp := logEntry["timestamp"]; hasTimestamp {
			timestamp = fmt.Sprintf("%v", ts)
		}

		content := "no content"
		if _, hasContent := logEntry["content"]; hasContent {
			content = fmt.Sprintf("%v", logEntry["content"])
		}

		level := "unknown"
		if l, hasLevel := logEntry["level"]; hasLevel {
			level = fmt.Sprintf("%v", l)
		}

		color := getColorByLevel(level)
		fmt.Printf("    %s[%s]%s %s - %s\n", color, level, ColorReset, timestamp, content)
	}

	return nil
}
