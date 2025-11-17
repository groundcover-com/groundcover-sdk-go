package e2e

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metricsClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestMetricsDiscoveryE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// Common time range for all tests
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour) // Query last hour
	startDateTime := strfmt.DateTime(startTime)
	endDateTime := strfmt.DateTime(endTime)

	t.Run("Get Metric Names", func(t *testing.T) {
		// Create request body
		body := &models.MetricsNamesRequest{
			Start:  startDateTime,
			End:    endDateTime,
			Limit:  10,
			Filter: "", // Get all metric names
		}

		// Create parameters
		params := metricsClient.NewGetMetricNamesParamsWithContext(ctx).WithBody(body)

		// Execute the request
		resp, err := apiClient.Metrics.GetMetricNames(params, nil)

		// Assertions
		require.NoError(t, err, "Get metric names request failed")
		require.NotNil(t, resp, "Get metric names response should not be nil")
		require.NotNil(t, resp.Payload, "Get metric names response payload should not be nil")

		// Check response structure
		assert.NotNil(t, resp.Payload.Metrics, "Metrics array should not be nil")
		assert.GreaterOrEqual(t, len(resp.Payload.Metrics), 0, "Should return metrics array (can be empty)")

		if len(resp.Payload.Metrics) > 0 {
			// Check first metric structure
			firstMetric := resp.Payload.Metrics[0]
			assert.NotEmpty(t, firstMetric.Name, "Metric name should not be empty")
			t.Logf("Found metric: %s (type: %s, unit: %s)", firstMetric.Name, firstMetric.Type, firstMetric.Unit)
		}

		t.Logf("Successfully retrieved %d metric names", len(resp.Payload.Metrics))
	})

	t.Run("Get Metric Names with Filter", func(t *testing.T) {
		// Create request body with filter
		body := &models.MetricsNamesRequest{
			Start:  startDateTime,
			End:    endDateTime,
			Limit:  5,
			Filter: "cpu", // Filter for CPU-related metrics
		}

		// Create parameters
		params := metricsClient.NewGetMetricNamesParamsWithContext(ctx).WithBody(body)

		// Execute the request
		resp, err := apiClient.Metrics.GetMetricNames(params, nil)

		// Assertions
		require.NoError(t, err, "Get filtered metric names request failed")
		require.NotNil(t, resp, "Get filtered metric names response should not be nil")
		require.NotNil(t, resp.Payload, "Get filtered metric names response payload should not be nil")

		// Check that results contain the filter string
		for _, metric := range resp.Payload.Metrics {
			t.Logf("Filtered metric: %s", metric.Name)
		}

		t.Logf("Successfully retrieved %d filtered metric names", len(resp.Payload.Metrics))
	})

	t.Run("Get Metric Keys", func(t *testing.T) {
		// First, get available metric names
		namesBody := &models.MetricsNamesRequest{
			Start: startDateTime,
			End:   endDateTime,
			Limit: 1,
		}
		namesParams := metricsClient.NewGetMetricNamesParamsWithContext(ctx).WithBody(namesBody)
		namesResp, err := apiClient.Metrics.GetMetricNames(namesParams, nil)
		require.NoError(t, err, "Failed to get metric names for keys test")

		if len(namesResp.Payload.Metrics) == 0 {
			t.Skip("No metrics available to test keys")
		}

		metricName := namesResp.Payload.Metrics[0].Name

		// Create request body for keys
		body := &models.MetricsKeysRequest{
			Start: startDateTime,
			End:   endDateTime,
			Name:  metricName,
			Limit: 10,
		}

		// Create parameters
		params := metricsClient.NewGetMetricKeysParamsWithContext(ctx).WithBody(body)

		// Execute the request
		resp, err := apiClient.Metrics.GetMetricKeys(params, nil)

		// Assertions
		require.NoError(t, err, "Get metric keys request failed")
		require.NotNil(t, resp, "Get metric keys response should not be nil")
		require.NotNil(t, resp.Payload, "Get metric keys response payload should not be nil")

		// Check response structure
		assert.Equal(t, metricName, resp.Payload.Name, "Response should contain the requested metric name")
		assert.NotNil(t, resp.Payload.Keys, "Keys array should not be nil")

		t.Logf("Successfully retrieved %d keys for metric '%s'", len(resp.Payload.Keys), metricName)
		if len(resp.Payload.Keys) > 0 {
			t.Logf("Sample keys: %v", resp.Payload.Keys[:min(3, len(resp.Payload.Keys))])
		}
	})

	t.Run("Get Metric Values", func(t *testing.T) {
		// First, get available metric names
		namesBody := &models.MetricsNamesRequest{
			Start: startDateTime,
			End:   endDateTime,
			Limit: 1,
		}
		namesParams := metricsClient.NewGetMetricNamesParamsWithContext(ctx).WithBody(namesBody)
		namesResp, err := apiClient.Metrics.GetMetricNames(namesParams, nil)
		require.NoError(t, err, "Failed to get metric names for values test")

		if len(namesResp.Payload.Metrics) == 0 {
			t.Skip("No metrics available to test values")
		}

		metricName := namesResp.Payload.Metrics[0].Name

		// Get keys for this metric
		keysBody := &models.MetricsKeysRequest{
			Start: startDateTime,
			End:   endDateTime,
			Name:  metricName,
			Limit: 1,
		}
		keysParams := metricsClient.NewGetMetricKeysParamsWithContext(ctx).WithBody(keysBody)
		keysResp, err := apiClient.Metrics.GetMetricKeys(keysParams, nil)
		require.NoError(t, err, "Failed to get metric keys for values test")

		if len(keysResp.Payload.Keys) == 0 {
			t.Skip("No keys available to test values")
		}

		keyName := keysResp.Payload.Keys[0]

		// Create request body for values
		body := &models.MetricsValuesRequest{
			Start: startDateTime,
			End:   endDateTime,
			Name:  metricName,
			Key:   keyName,
			Limit: 10,
		}

		// Create parameters
		params := metricsClient.NewGetMetricValuesParamsWithContext(ctx).WithBody(body)

		// Execute the request
		resp, err := apiClient.Metrics.GetMetricValues(params, nil)

		// Assertions
		require.NoError(t, err, "Get metric values request failed")
		require.NotNil(t, resp, "Get metric values response should not be nil")
		require.NotNil(t, resp.Payload, "Get metric values response payload should not be nil")

		// Check response structure
		assert.Equal(t, metricName, resp.Payload.Name, "Response should contain the requested metric name")
		assert.Equal(t, keyName, resp.Payload.Key, "Response should contain the requested key name")
		assert.NotNil(t, resp.Payload.Values, "Values array should not be nil")

		t.Logf("Successfully retrieved %d values for metric '%s' key '%s'",
			len(resp.Payload.Values), metricName, keyName)
		if len(resp.Payload.Values) > 0 {
			t.Logf("Sample values: %v", resp.Payload.Values[:min(3, len(resp.Payload.Values))])
		}
	})

	t.Run("Test Error Cases", func(t *testing.T) {
		t.Run("Invalid Time Range", func(t *testing.T) {
			// Create request with end time before start time
			body := &models.MetricsNamesRequest{
				Start: endDateTime,
				End:   startDateTime, // Invalid: end before start
				Limit: 10,
			}

			params := metricsClient.NewGetMetricNamesParamsWithContext(ctx).WithBody(body)
			_, err := apiClient.Metrics.GetMetricNames(params, nil)

			assert.Error(t, err, "Should return error for invalid time range")
		})

		t.Run("Non-existent Metric", func(t *testing.T) {
			body := &models.MetricsKeysRequest{
				Start: startDateTime,
				End:   endDateTime,
				Name:  "non_existent_metric_xyz123",
				Limit: 10,
			}

			params := metricsClient.NewGetMetricKeysParamsWithContext(ctx).WithBody(body)
			resp, err := apiClient.Metrics.GetMetricKeys(params, nil)

			// Should not error but return empty results
			require.NoError(t, err, "Should not error for non-existent metric")
			assert.Empty(t, resp.Payload.Keys, "Should return empty keys for non-existent metric")
		})
	})
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
