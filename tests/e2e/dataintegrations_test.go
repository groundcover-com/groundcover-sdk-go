package e2e

import (
	"testing"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/integrations"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCloudwatchConfig = `{
  "version": 1,
  "name": "test-cloudwatch",
  "exporters": ["prometheus"],
  "scrapeInterval": "5m",
  "stsRegion": "us-east-1",
  "regions": ["us-east-1"],
  "roleArn": "arn:aws:iam::123456789012:role/test-role",
  "awsMetrics": [
    {
      "namespace": "AWS/EC2",
      "metrics": [
        {
          "name": "CPUUtilization",
          "statistics": ["Average"],
          "period": 300,
          "length": 300,
          "nullAsZero": false
        }
      ]
    }
  ],
  "apiConcurrencyLimits": {
    "listMetrics": 3,
    "getMetricData": 5,
    "getMetricStatistics": 5,
    "listInventory": 10
  },
  "withContextTagsOnInfoMetrics": false,
  "withInventoryDiscovery": false
}`

const testCloudwatchConfigUpdated = `{
  "version": 1,
  "name": "test-cloudwatch",
  "exporters": ["prometheus"],
  "scrapeInterval": "5m",
  "stsRegion": "us-east-2",
  "regions": ["us-east-2"],
  "roleArn": "arn:aws:iam::123456789012:role/test-role",
  "awsMetrics": [
    {
      "namespace": "AWS/EC2",
      "metrics": [
        {
          "name": "CPUUtilization",
          "statistics": ["Average"],
          "period": 300,
          "length": 300,
          "nullAsZero": false
        }
      ]
    }
  ],
  "apiConcurrencyLimits": {
    "listMetrics": 1,
    "getMetricData": 5,
    "getMetricStatistics": 5,
    "listInventory": 10
  },
  "withContextTagsOnInfoMetrics": false,
  "withInventoryDiscovery": false
}`

func TestCloudwatch(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// 1. CREATE
	createBody := &models.CreateDataIntegrationConfigRequest{
		Config: testCloudwatchConfig,
	}
	createParams := integrations.NewCreateDataIntegrationConfigParamsWithContext(ctx).WithType("cloudwatch").WithBody(createBody)
	createResp, err := apiClient.Integrations.CreateDataIntegrationConfig(createParams, nil)
	require.NoError(t, err, "Create config request failed")
	require.NotNil(t, createResp, "Create config response should not be nil")
	require.NotNil(t, createResp.Payload, "Create config response payload should not be nil")

	assert.Equal(t, testCloudwatchConfig, createResp.Payload.Config, "Created config value should match the request")
	assert.NotEmpty(t, createResp.Payload.ID, "Created config should have an ID")
	assert.NotEmpty(t, createResp.Payload.UpdateTimestamp, "Created config should have an update timestamp")
	require.False(t, createResp.Payload.IsArchived, "Config should not be archived")
	originalConfigCreationTimestamp := createResp.Payload.UpdateTimestamp
	originalConfigID := createResp.Payload.ID

	// 2. READ
	getParams := integrations.NewGetDataIntegrationConfigParamsWithContext(ctx).WithID(originalConfigID).WithType("cloudwatch")
	getRespOk, err := apiClient.Integrations.GetDataIntegrationConfig(getParams, nil)

	require.NoError(t, err, "Get config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	require.False(t, getRespOk.Payload.IsArchived, "Config should not be archived")
	assert.Equal(t, testCloudwatchConfig, getRespOk.Payload.Config, "Retrieved config value should match the created one")

	// 3. UPDATE
	updateBody := &models.CreateDataIntegrationConfigRequest{
		Config: testCloudwatchConfigUpdated,
	}
	updateParams := integrations.NewUpdateDataIntegrationConfigParamsWithContext(ctx).WithType("cloudwatch").WithID(originalConfigID).WithBody(updateBody)
	updateResp, err := apiClient.Integrations.UpdateDataIntegrationConfig(updateParams, nil)

	require.NoError(t, err, "Update config request failed")
	require.NotNil(t, updateResp, "Update config response should not be nil")
	require.NotNil(t, updateResp.Payload, "Update config response payload should not be nil")
	assert.Equal(t, testCloudwatchConfigUpdated, updateResp.Payload.Config, "Updated config value should match the request")
	assert.NotEmpty(t, updateResp.Payload.UpdateTimestamp, "Updated config should have an update timestamp")

	getRespOk, err = apiClient.Integrations.GetDataIntegrationConfig(getParams, nil)
	require.NoError(t, err, "Get updated config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	require.False(t, getRespOk.Payload.IsArchived, "Config should not be archived")
	assert.Equal(t, testCloudwatchConfigUpdated, getRespOk.Payload.Config, "Retrieved config should have updated value")
	assert.Greater(t, updateResp.Payload.UpdateTimestamp, originalConfigCreationTimestamp, "Updated config should have an update timestamp greater than the original one")

	// 4. DELETE
	deleteParams := integrations.NewDeleteDataIntegrationConfigParamsWithContext(ctx).WithID(originalConfigID).WithType("cloudwatch")
	_, err = apiClient.Integrations.DeleteDataIntegrationConfig(deleteParams, nil)
	require.NoError(t, err, "Delete config request failed")
	t.Log("Successfully deleted data integration config")

	// Verify the config was deleted by trying to get it without includeArchived - should return 404
	getRespOk, err = apiClient.Integrations.GetDataIntegrationConfig(getParams, nil)
	require.Error(t, err, "Get deleted config request should fail")
	require.Nil(t, getRespOk, "Should not get a config response")
}
