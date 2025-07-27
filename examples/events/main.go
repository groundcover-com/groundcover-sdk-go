package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/events"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/utils"
)

func main() {
	fmt.Println("groundcover Go SDK - Events Example")
	fmt.Println("===================================")

	// Create client
	client, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	ctx := context.Background()

	// Search for OOM events in the last 24 hours (example of event searching)
	fmt.Println("\nSearching for events...")
	if err := searchOOMEvents(ctx, client); err != nil {
		fmt.Printf("Error searching events: %v\n", err)
	}
}

func searchOOMEvents(ctx context.Context, gcClient *client.GroundcoverAPI) error {
	// Search for OOM events in the last 24 hours (example event type)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	startDateTime := strfmt.DateTime(startTime)
	endDateTime := strfmt.DateTime(endTime)

	// Use the built-in OOM event conditions
	cs := utils.NewConditionSet()
	cs.AddOOMEventConditions()

	searchRequest := &models.EventsSearchRequest{
		Start: &startDateTime,
		End:   &endDateTime,
		// Search for OOM events, limit 10 results
		Query: "type:container_crash reason:OOMKilled | limit 10",
	}

	params := events.NewSearchEventsParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(searchRequest)

	resp, err := gcClient.Events.SearchEvents(params, nil)
	if err != nil {
		return fmt.Errorf("failed to search events: %w", err)
	}

	return printOOMEvents(resp.Payload)
}

func printOOMEvents(payload interface{}) error {
	events, ok := payload.([]interface{})
	if !ok {
		fmt.Printf("  Unexpected payload type: %T\n", payload)
		return nil
	}

	if len(events) == 0 {
		fmt.Println("  ✓ No OOM events found (that's good!)")
		return nil
	}

	fmt.Printf("  ⚠ Found %d event(s):\n", len(events))
	fmt.Println("  " + strings.Repeat("─", 60))

	for i, eventItem := range events {
		event := eventItem.(map[string]interface{})

		fmt.Printf("  #%d\n", i+1)

		// Main event fields
		fmt.Printf("    Time:      %s\n", getString(event, "timestamp"))
		fmt.Printf("    Namespace: %s\n", getString(event, "entity_namespace"))
		fmt.Printf("    Workload:  %s\n", getString(event, "entity_workload"))
		fmt.Printf("    Container: %s\n", getString(event, "entity_name"))
		fmt.Printf("    Reason:    %s\n", getString(event, "reason"))

		// Extract details from string_attributes
		if attrs, ok := event["string_attributes"].(map[string]interface{}); ok {
			fmt.Printf("    Pod:       %s\n", getString(attrs, "podName"))
			fmt.Printf("    Exit Code: %s\n", getString(attrs, "exitCode"))
			fmt.Printf("    Image:     %s\n", getString(attrs, "imageName"))
		}

		// Extract memory limit from float_attributes and convert to readable format
		if floatAttrs, ok := event["float_attributes"].(map[string]interface{}); ok {
			if memLimit, ok := floatAttrs["memoryLimit"].(float64); ok {
				fmt.Printf("    Memory Limit: %.0f MB\n", memLimit/1024/1024)
			}
		}

		if i < len(events)-1 {
			fmt.Println()
		}
	}

	return nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return "unknown"
}
