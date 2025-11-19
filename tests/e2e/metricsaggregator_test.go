package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aggregationsMetricsClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/aggregations_metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const testMetricsAggregatorConfigValue = `content: |
  - ignore_old_samples: true
    match: '{__name__=~"test_metric_counter"}'
    without: [instance]
    interval: 30s
    outputs: [total_prometheus]
  - match: '{__name__=~"test_metric_latency"}'
    without: [instance]
    interval: 30s
    outputs: [avg]`

const testMetricsAggregatorConfigValueUpdated = `content: |
  - ignore_old_samples: true
    match: '{__name__=~"test_metric_counter_updated"}'
    without: [instance]
    interval: 60s
    outputs: [total_prometheus]
  - match: '{__name__=~"test_metric_latency_updated"}'
    without: [instance]
    interval: 60s
    outputs: [avg]`

func TestRemoteConfigMetricsAggregatorCrudE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// 1. CREATE - Create a new metrics aggregator configuration
	createBody := &models.CreateOrUpdateMetricsAggregatorConfigRequest{
		Value: testMetricsAggregatorConfigValue,
	}

	// Create parameters
	createParams := aggregationsMetricsClient.NewCreateMetricsAggregatorConfigParamsWithContext(ctx).WithBody(createBody)

	// Execute the create request
	createResp, err := apiClient.AggregationsMetrics.CreateMetricsAggregatorConfig(createParams, nil)

	// Assertions
	require.NoError(t, err, "Create config request failed")
	require.NotNil(t, createResp, "Create config response should not be nil")
	require.NotNil(t, createResp.Payload, "Create config response payload should not be nil")

	// Verify the created config
	assert.Equal(t, testMetricsAggregatorConfigValue, createResp.Payload.Value, "Created config value should match the request")
	assert.NotEmpty(t, createResp.Payload.UUID, "Created config should have an ID")
	assert.NotEmpty(t, createResp.Payload.CreatedTimestamp, "Created config should have a creation timestamp")
	originalConfigCreationTimestamp := createResp.Payload.CreatedTimestamp

	// 2. READ - Get the metrics aggregator configuration and verify it matches the created one
	getParams := aggregationsMetricsClient.NewGetMetricsAggregatorConfigParamsWithContext(ctx)

	// Execute the get request
	getRespOk, err := apiClient.AggregationsMetrics.GetMetricsAggregatorConfig(getParams, nil)

	// Assertions
	require.NoError(t, err, "Get config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")

	// Verify the retrieved config
	assert.Equal(t, testMetricsAggregatorConfigValue, getRespOk.Payload.Value, "Retrieved config value should match the created one")

	// 3. UPDATE - Update the metrics aggregator configuration
	updateBody := &models.CreateOrUpdateMetricsAggregatorConfigRequest{
		Value: testMetricsAggregatorConfigValueUpdated,
	}

	// Create parameters
	updateParams := aggregationsMetricsClient.NewUpdateMetricsAggregatorConfigParamsWithContext(ctx).WithBody(updateBody)

	// Execute the update request
	updateResp, err := apiClient.AggregationsMetrics.UpdateMetricsAggregatorConfig(updateParams, nil)

	// Assertions
	require.NoError(t, err, "Update config request failed")
	require.NotNil(t, updateResp, "Update config response should not be nil")
	require.NotNil(t, updateResp.Payload, "Update config response payload should not be nil")

	// Verify the updated config
	assert.Equal(t, testMetricsAggregatorConfigValueUpdated, updateResp.Payload.Value, "Updated config value should match the request")
	assert.NotEmpty(t, updateResp.Payload.CreatedTimestamp, "Updated config should have an update timestamp")

	// Verify we can retrieve the updated version
	getRespOk, err = apiClient.AggregationsMetrics.GetMetricsAggregatorConfig(getParams, nil)
	require.NoError(t, err, "Get updated config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	assert.Equal(t, testMetricsAggregatorConfigValueUpdated, getRespOk.Payload.Value, "Retrieved config should have updated value")
	assert.Greater(t, updateResp.Payload.CreatedTimestamp, originalConfigCreationTimestamp, "Updated config should have a creation timestamp greater than the original one")

	// 4. DELETE - Delete the metrics aggregator configuration
	deleteParams := aggregationsMetricsClient.NewDeleteMetricsAggregatorConfigParamsWithContext(ctx)
	_, err = apiClient.AggregationsMetrics.DeleteMetricsAggregatorConfig(deleteParams, nil)

	// Assertions
	require.NoError(t, err, "Delete config request failed")

	t.Log("Successfully deleted metrics aggregator config")

	// Verify the config was deleted by trying to get it - should return 204 No Content
	getRespOk, err = apiClient.AggregationsMetrics.GetMetricsAggregatorConfig(getParams, nil)
	require.NoError(t, err, "Get deleted config request failed")
	require.Equal(t, getRespOk.Payload.Value, "", "Should get an empty config value")
}
