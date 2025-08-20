package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/logs"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const (
	QueryBatchSize = 1000
)

type ParsedLogLine struct {
	Timestamp        time.Time
	Cluster          string
	Env              string
	Workload         string
	Namespace        string
	Container        string
	Level            string
	Message          string
	Tags             []string
	StringAttributes map[string]interface{}
	FloatAttributes  map[string]interface{}
}

func main() {
	fmt.Println("groundcover Go SDK - Logs Example")
	fmt.Println("=================================")

	// Define command line flags
	var format string
	flag.StringVar(&format, "format", "csv", "Output format: csv or json")
	flag.Parse()

	// Check command line arguments
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: logs-fetcher [-format=<csv|json>] <query> <output_file>")
		fmt.Println("  -format: Output format (default: csv)")
		fmt.Println("Example: logs-fetcher 'level:error' results.csv")
		fmt.Println("Example: logs-fetcher -format=json 'level:error' results.json")
		os.Exit(1)
	}

	// Validate format
	if format != "json" && format != "csv" {
		fmt.Printf("Invalid format: %s. Must be 'json' or 'csv'\n", format)
		os.Exit(1)
	}

	query := flag.Arg(0)
	outputFile := flag.Arg(1)

	// Create client
	client, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	ctx := context.Background()

	// Example: Fetch logs in chunks
	fmt.Printf("\nFetching logs with query: %s\n", query)
	fmt.Printf("Output format: %s\n", format)
	fmt.Printf("Results will be saved to: %s\n", outputFile)

	logsWritten, err := fetchLogs(ctx, client, query, outputFile, format)
	if err != nil {
		fmt.Printf("Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if logsWritten == 0 {
		os.Exit(0)
	}

	fmt.Printf("✓ %d logs successfully saved to %s in %s format\n", logsWritten, outputFile, format)
}

func fetchLogs(ctx context.Context, gcClient *client.GroundcoverAPI, query, outputFile, format string) (int, error) {
	// Search for logs using provided query
	endTime := time.Now()
	startTime := endTime.Add(-20 * time.Minute)

	count, err := getLogsCount(ctx, gcClient, startTime, endTime, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get logs count: %w", err)
	}

	if count == 0 {
		fmt.Println("No logs found")
		return 0, nil
	}

	fmt.Printf("Fetching %d logs\n", count)

	logs, err := getLogs(ctx, gcClient, startTime, endTime, query, int(count), outputFile, format)
	if err != nil {
		return 0, fmt.Errorf("failed to get logs: %w", err)
	}

	fmt.Printf("Successfully fetched %d logs\n", logs)

	return logs, nil
}

func getLogsCount(ctx context.Context, gcClient *client.GroundcoverAPI, startTime time.Time, endTime time.Time, query string) (int64, error) {
	startDateTime := strfmt.DateTime(startTime)
	endDateTime := strfmt.DateTime(endTime)

	countQuery := fmt.Sprintf("%s | stats count() as count", query)

	searchRequest := models.LogsSearchRequest{
		Start: &startDateTime,
		End:   &endDateTime,
		Query: countQuery,
	}

	params := logs.NewSearchLogsParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(&searchRequest)

	resp, err := gcClient.Logs.SearchLogs(params, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to search logs: %w", err)
	}

	response, ok := resp.Payload.([]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected payload type: %T", resp.Payload)
	}

	if len(response) == 0 {
		// No logs in the window; treat as zero, not an error.
		return 0, nil
	}
	countElement := response[0]

	countElementMap, ok := countElement.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get count element")
	}

	countValue, ok := countElementMap["count"]
	if !ok {
		return 0, fmt.Errorf("failed to get count value")
	}

	countValueJsonNumber, ok := countValue.(json.Number)
	if !ok {
		return 0, fmt.Errorf("failed to get count value")
	}

	return countValueJsonNumber.Int64()
}

func getLogs(ctx context.Context, gcClient *client.GroundcoverAPI, startTime time.Time, endTime time.Time, query string, logCount int, outputFile, format string) (int, error) {
	startDateTime := strfmt.DateTime(startTime)
	endDateTime := strfmt.DateTime(endTime)

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return 0, fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	var csvWriter *csv.Writer
	if format == "csv" {
		csvWriter = csv.NewWriter(file)
		defer csvWriter.Flush()

		// Write CSV header
		header := []string{"Timestamp", "Cluster", "Env", "Workload", "Namespace", "Container", "Level", "Message"}
		if err := csvWriter.Write(header); err != nil {
			return 0, fmt.Errorf("failed to write CSV header: %w", err)
		}
	}

	offset := 0
	totalLogs := 0

	for offset < logCount {
		currentCount := QueryBatchSize
		if offset+currentCount > logCount {
			currentCount = logCount - totalLogs
		}

		log.Printf("Getting logs from %d to %d", offset, offset+currentCount)
		queryWithPagination := fmt.Sprintf("%s | offset %d | limit %d", query, offset, currentCount)

		searchRequest := models.LogsSearchRequest{
			Start: &startDateTime,
			End:   &endDateTime,
			Query: queryWithPagination,
		}

		params := logs.NewSearchLogsParams().
			WithContext(ctx).
			WithTimeout(30 * time.Second).
			WithBody(&searchRequest)

		resp, err := gcClient.Logs.SearchLogs(params, nil)

		if err != nil {
			return 0, fmt.Errorf("failed to search logs: %w", err)
		}

		// Process response and write to file
		response := resp.Payload.([]interface{})
		parsedLogs := make([]ParsedLogLine, 0)

		for _, logEntry := range response {
			logEntryMap, ok := logEntry.(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("failed to convert log entry to map")
			}

			timestamp, ok := logEntryMap["timestamp"].(string)

			if !ok {
				return 0, fmt.Errorf("failed to get timestamp")
			}

			timestampTime, err := time.Parse(time.RFC3339, timestamp)

			if err != nil {
				return 0, fmt.Errorf("failed to parse timestamp: %w", err)
			}

			level, ok := logEntryMap["level"].(string)

			if !ok {
				return 0, fmt.Errorf("failed to get level")
			}

			message, ok := logEntryMap["content"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get message")
			}

			stringAttributes, ok := logEntryMap["string_attributes"].(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("failed to get string attributes")
			}

			floatAttributes, ok := logEntryMap["float_attributes"].(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("failed to get float attributes")
			}

			cluster, ok := logEntryMap["cluster"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get cluster")
			}

			env, ok := logEntryMap["env"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get env")
			}

			workload, ok := logEntryMap["workload"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get workload")
			}

			namespace, ok := logEntryMap["namespace"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get namespace")
			}

			container, ok := logEntryMap["container_name"].(string)
			if !ok {
				return 0, fmt.Errorf("failed to get container")
			}

			parsedLog := ParsedLogLine{
				Timestamp:        timestampTime,
				Level:            level,
				Message:          message,
				StringAttributes: stringAttributes,
				FloatAttributes:  floatAttributes,
				Cluster:          cluster,
				Env:              env,
				Workload:         workload,
				Namespace:        namespace,
				Container:        container,
			}

			// Write to file based on format
			if format == "csv" {
				// Write as CSV row
				row := []string{
					parsedLog.Timestamp.Format(time.RFC3339),
					parsedLog.Cluster,
					parsedLog.Env,
					parsedLog.Workload,
					parsedLog.Namespace,
					parsedLog.Container,
					parsedLog.Level,
					parsedLog.Message,
				}
				if err := csvWriter.Write(row); err != nil {
					return 0, fmt.Errorf("failed to write CSV row: %w", err)
				}
			} else {
				// Write as JSON
				jsonLog, err := json.Marshal(parsedLog)
				if err != nil {
					return 0, fmt.Errorf("failed to marshal log: %w", err)
				}
				file.Write(jsonLog)
				file.WriteString("\n")
			}

			parsedLogs = append(parsedLogs, parsedLog)

			totalLogs++
		}

		offset += len(parsedLogs)
	}

	return totalLogs, nil
}
