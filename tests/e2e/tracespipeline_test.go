package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tracesPipelineClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/traces_pipeline"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const tracesTestConfigValue = `ottlRules:
- ruleName: example-rule
  conditions:
    - workload == "nginx"
  statements:
    - set(attributes["test.key"], "test-value")`

const tracesTestConfigValueUpdated = `ottlRules:
- ruleName: example-rule-updated
  conditions:
    - workload == "nginx"
  statements:
    - set(attributes["test.key"], "test-value-updated")`

func TestRemoteConfigTracesPipelineCrudE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// 1. CREATE - Create a new traces pipeline configuration
	createBody := &models.CreateOrUpdateTracesPipelineConfigRequest{
		Value: tracesTestConfigValue,
	}

	// Create parameters
	createParams := tracesPipelineClient.NewCreateTracesPipelineConfigParamsWithContext(ctx).WithBody(createBody)

	// Execute the create request
	createResp, err := apiClient.TracesPipeline.CreateTracesPipelineConfig(createParams, nil)

	// Assertions
	require.NoError(t, err, "Create config request failed")
	require.NotNil(t, createResp, "Create config response should not be nil")
	require.NotNil(t, createResp.Payload, "Create config response payload should not be nil")

	// Verify the created config
	assert.Equal(t, tracesTestConfigValue, createResp.Payload.Value, "Created config value should match the request")
	assert.NotEmpty(t, createResp.Payload.UUID, "Created config should have an ID")
	assert.NotEmpty(t, createResp.Payload.CreatedTimestamp, "Created config should have a creation timestamp")
	originalConfigCreationTimestamp := createResp.Payload.CreatedTimestamp

	// 2. READ - Get the traces pipeline configuration and verify it matches the created one
	getParams := tracesPipelineClient.NewGetTracesPipelineConfigParamsWithContext(ctx)

	// Execute the get request
	getRespOk, getRespNoContent, err := apiClient.TracesPipeline.GetTracesPipelineConfig(getParams, nil)

	// Assertions
	require.NoError(t, err, "Get config request failed")
	require.Nil(t, getRespNoContent, "Shouldnt get a 204 No Content response")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")

	// Verify the retrieved config
	assert.Equal(t, tracesTestConfigValue, getRespOk.Payload.Value, "Retrieved config value should match the created one")

	// 3. UPDATE - Update the traces pipeline configuration
	updateBody := &models.CreateOrUpdateTracesPipelineConfigRequest{
		Value: tracesTestConfigValueUpdated,
	}

	// Create parameters
	updateParams := tracesPipelineClient.NewUpdateTracesPipelineConfigParamsWithContext(ctx).WithBody(updateBody)

	// Execute the update request
	updateResp, err := apiClient.TracesPipeline.UpdateTracesPipelineConfig(updateParams, nil)

	// Assertions
	require.NoError(t, err, "Update config request failed")
	require.NotNil(t, updateResp, "Update config response should not be nil")
	require.NotNil(t, updateResp.Payload, "Update config response payload should not be nil")

	// Verify the updated config
	assert.Equal(t, tracesTestConfigValueUpdated, updateResp.Payload.Value, "Updated config value should match the request")
	assert.NotEmpty(t, updateResp.Payload.CreatedTimestamp, "Updated config should have an update timestamp")

	// Verify we can retrieve the updated version
	getRespOk, getRespNoContent, err = apiClient.TracesPipeline.GetTracesPipelineConfig(getParams, nil)
	require.NoError(t, err, "Get updated config request failed")
	require.Nil(t, getRespNoContent, "Shouldnt get a 204 No Content response")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	assert.Equal(t, tracesTestConfigValueUpdated, getRespOk.Payload.Value, "Retrieved config should have updated value")
	assert.Greater(t, updateResp.Payload.CreatedTimestamp, originalConfigCreationTimestamp, "Updated config should have a creation timestamp greater than the original one")

	// 4. DELETE - Delete the traces pipeline configuration
	deleteParams := tracesPipelineClient.NewDeleteTracesPipelineConfigParamsWithContext(ctx)
	_, err = apiClient.TracesPipeline.DeleteTracesPipelineConfig(deleteParams, nil)

	// Assertions
	require.NoError(t, err, "Delete config request failed")

	t.Log("Successfully deleted traces pipeline config")

	// Verify the config was deleted by trying to get it - should return 204 No Content
	getRespOk, getRespNoContent, err = apiClient.TracesPipeline.GetTracesPipelineConfig(getParams, nil)
	require.NoError(t, err, "Get deleted config request failed")
	require.Nil(t, getRespNoContent, "Should get a 204 No Content response")
	require.Equal(t, getRespOk.Payload.Value, "", "Should get the same config value")
}
