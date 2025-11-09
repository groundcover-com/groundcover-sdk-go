package e2e

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/go-openapi/runtime"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

// MonitorYAML represents a YAML monitor definition (simplified for testing)
type MonitorYAML struct {
	Title string `yaml:"title"`
}

// MonitorYAMLWithSeverity represents a YAML monitor definition with severity
type MonitorYAMLWithSeverity struct {
	Title    string `yaml:"title"`
	Severity string `yaml:"severity"`
}

// MonitorYAMLWithEvaluation represents a YAML monitor definition with evaluation interval
type MonitorYAMLWithEvaluation struct {
	Title              string `yaml:"title"`
	Severity           string `yaml:"severity"`
	EvaluationInterval struct {
		Interval   string `yaml:"interval"`
		PendingFor string `yaml:"pendingFor"`
	} `yaml:"evaluationInterval"`
}

func TestMonitorsEndpoints(t *testing.T) {
	// Log the base URL to verify it's correct
	t.Logf("Using base URL: %s", os.Getenv("GC_BASE_URL"))

	client := NewTestClient(t)
	defer client.Cleanup()

	t.Run("Create, Get, Update, and Delete Monitor", func(t *testing.T) {
		// Generate a random monitor name using a local, thread-safe random source
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomName := fmt.Sprintf("E2E Test - K8s Pod Not Healthy Monitor - %d", rng.Intn(10_000_000))
		randomHeader := fmt.Sprintf("E2E Test - K8s Pod Not Healthy - %d", rng.Intn(10_000_000))

		// First, create a monitor
		yamlMonitor := fmt.Sprintf(`
title: %s
display:
  header: %s
  resourceHeaderLabels:
    - namespace
    - workload
  contextHeaderLabels:
    - cluster
  description: Pod has been in a non-running state for longer than 15 minutes
severity: critical
measurementType: state
model:
  queries:
    - dataType: metrics
      name: threshold_input_query
      pipeline:
        function:
          name: avg_over_time
          pipelines:
            - function:
                name: max_by
                pipelines:
                  - metric: groundcover_kube_pod_status_phase
                args:
                  - namespace
                  - workload
                  - cluster
          args:
            - "600"
  thresholds:
    - name: threshold_1
      inputName: threshold_input_query
      operator: gt
      values:
        - 0
labels:
  severity: critical
annotations:
  description: Pod {{ .Labels.namespace }}/{{ .Labels.pod }} has been in a non-running state for longer than 15 minutes
  summary: Kubernetes Pod not healthy
executionErrorState: OK
noDataState: OK
evaluationInterval:
  interval: 1m
  pendingFor: 0s`, randomName, randomHeader)

		// Unmarshal the YAML into a struct
		monitor := &models.CreateMonitorRequest{}
		err := yaml.Unmarshal([]byte(yamlMonitor), monitor)
		if err != nil {
			t.Fatalf("Error unmarshalling YAML: %v", err)
		}

		// Create the monitor
		createParams := monitors.NewCreateMonitorParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(monitor)

		createResp, err := client.Client.Monitors.CreateMonitor(createParams, nil, monitors.WithContentTypeApplicationxYaml, monitors.WithAcceptApplicationJSON)
		if err != nil {
			t.Fatalf("Failed to create monitor: %v", err)
		}

		monitorID := createResp.Payload.MonitorID
		t.Logf("Created monitor with ID: %s", monitorID)

		// Get the monitor - Use client defaults for Accept header
		getParams := monitors.NewGetMonitorParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(monitorID)

		getResp, err := client.Client.Monitors.GetMonitor(getParams, nil)
		require.NoError(t, err, "Failed to get monitor")
		require.NotNil(t, getResp, "Get response should not be nil")

		// Unmarshal the YAML response
		var receivedMonitor MonitorYAML
		// The Body field now holds the raw YAML bytes (likely as strfmt.Base64)
		err = yaml.Unmarshal(getResp.Payload, &receivedMonitor)
		require.NoError(t, err, "Failed to unmarshal get monitor response YAML")

		// Also unmarshal with evaluation interval details
		var receivedMonitorWithEval MonitorYAMLWithEvaluation
		err = yaml.Unmarshal(getResp.Payload, &receivedMonitorWithEval)
		require.NoError(t, err, "Failed to unmarshal get monitor response YAML with evaluation interval")

		// Assertions on the received monitor
		require.Equal(t, *monitor.Title, receivedMonitor.Title, "Monitor title mismatch after get")

		// Check if pendingFor is preserved as 0s
		if receivedMonitorWithEval.EvaluationInterval.PendingFor != "0s" {
			t.Errorf("Expected pendingFor to be '0s', but got: '%s'", receivedMonitorWithEval.EvaluationInterval.PendingFor)
		}

		// --- Update the monitor ---

		updatedSeverity := "warning"
		// Reuse the original monitor definition and modify a field
		updatedYamlMonitor := strings.Replace(yamlMonitor, "critical", updatedSeverity, 1)
		updateMonitor := &models.UpdateMonitorRequest{}
		err = yaml.Unmarshal([]byte(updatedYamlMonitor), updateMonitor)
		require.NoError(t, err, "Failed to unmarshal update monitor YAML")

		// Prepare update parameters with the full, modified monitor definition
		updateParams := monitors.NewUpdateMonitorParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(monitorID).
			WithBody(updateMonitor) // Pass the modified monitor

		// Execute the update request - Requires application/x-yaml content type
		_, err = client.Client.Monitors.UpdateMonitor(updateParams, nil, monitors.WithContentTypeApplicationxYaml, monitors.WithAcceptApplicationJSON)
		if err != nil {
			t.Fatalf("Failed to update monitor: %v", err)
		}

		// --- Verify the update ---
		getUpdatedResp, err := client.Client.Monitors.GetMonitor(getParams, nil)
		if err != nil {
			t.Fatalf("Failed to get updated monitor: %v", err)
		}

		// Parse updated YAML
		var updatedYaml MonitorYAMLWithSeverity // Use struct that includes severity
		err = yaml.Unmarshal(getUpdatedResp.Payload, &updatedYaml)
		if err != nil {
			t.Fatalf("Error unmarshalling updated YAML: %v", err)
		}

		// Verify the updated field (severity)
		require.Equal(t, updatedSeverity, updatedYaml.Severity, "Monitor severity mismatch after update")
		// Also check title hasn't changed unintentionally (if needed)
		require.Equal(t, *monitor.Title, updatedYaml.Title, "Monitor title mismatch after update")

		// Test creating a monitor with the same title (should fail with 409)
		t.Run("Create duplicate monitor", func(t *testing.T) {
			duplicateMonitor := &models.CreateMonitorRequest{}
			err := yaml.Unmarshal([]byte(yamlMonitor), duplicateMonitor)
			if err != nil {
				t.Fatalf("Error unmarshalling YAML: %v", err)
			}

			// Set the title to the updated title (which exists)
			duplicateMonitor.Title = &updatedYaml.Title

			duplicateParams := monitors.NewCreateMonitorParams().
				WithContext(client.BaseCtx).
				WithTimeout(defaultTimeout).
				WithBody(duplicateMonitor)

			_, err = client.Client.Monitors.CreateMonitor(duplicateParams, nil, monitors.WithContentTypeApplicationxYaml, monitors.WithAcceptApplicationJSON)

			// Check if it's a 409 Conflict error
			// Check if the error type is the specific one for 409 Conflict
			conflictError, ok := err.(*monitors.CreateMonitorConflict)
			if !ok {
				// Fallback to generic APIError check if specific type assertion fails
				apiError, ok := err.(*runtime.APIError)
				if !ok {
					t.Fatalf("Expected API error, got: %T - %v", err, err)
				}
				if apiError.Code != http.StatusConflict {
					t.Fatalf("Expected status code 409, got: %d", apiError.Code)
				}
				// If it was a generic 409, we can't easily verify the payload type here, log and proceed
				t.Logf("Received generic 409 error: %v", apiError.Response)
				// We might want to fail here if we strictly expect the specific error type
				// t.Fatalf("Received generic 409 error, expected *monitors.CreateMonitorConflict")
			} else {
				// Check the error message within the payload of the specific error type
				if conflictError.Payload == nil {
					t.Fatalf("Conflict error payload is nil")
				}
				expectedErrMsg := "monitor with the same title already exists"
				if !strings.Contains(conflictError.Payload.Message, expectedErrMsg) {
					t.Fatalf("Expected error message to contain '%s', got: '%s'", expectedErrMsg, conflictError.Payload.Message)
				}
				t.Logf("Correctly received 409 error for duplicate title: %s", conflictError.Payload.Message)
			}
		})

		// Delete the monitor
		deleteParams := monitors.NewDeleteMonitorParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(monitorID)

		_, err = client.Client.Monitors.DeleteMonitor(deleteParams, nil, monitors.WithAcceptApplicationJSON)
		// Original error check - client should now handle 200 OK as success after regeneration
		if err != nil {
			t.Fatalf("Failed to delete monitor: %v", err)
		}
	})
}
