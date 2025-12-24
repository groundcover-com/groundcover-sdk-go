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
  "withInventoryDiscovery": true
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
  "withInventoryDiscovery": true
}`

func TestCloudwatch(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// 1. CREATE (without name/tags - backward compatibility)
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

func TestCloudwatchWithNameAndTags(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	testName := "e2e-test-cloudwatch"
	testTags := map[string]interface{}{
		"env":     "test",
		"purpose": "e2e-testing",
	}

	// 1. CREATE with name and tags
	createBody := &models.CreateDataIntegrationConfigRequest{
		Config: testCloudwatchConfig,
		Name:   testName,
		Tags:   testTags,
	}
	createParams := integrations.NewCreateDataIntegrationConfigParamsWithContext(ctx).WithType("cloudwatch").WithBody(createBody)
	createResp, err := apiClient.Integrations.CreateDataIntegrationConfig(createParams, nil)
	require.NoError(t, err, "Create config request failed")
	require.NotNil(t, createResp, "Create config response should not be nil")
	require.NotNil(t, createResp.Payload, "Create config response payload should not be nil")

	assert.Equal(t, testCloudwatchConfig, createResp.Payload.Config, "Created config value should match the request")
	// TODO: Uncomment after dev environment is updated to support name/tags (>= 1.9.743)
	// assert.Equal(t, testName, createResp.Payload.Name, "Created config name should match the request")
	// assert.Equal(t, testTags, createResp.Payload.Tags, "Created config tags should match the request")
	assert.NotEmpty(t, createResp.Payload.ID, "Created config should have an ID")
	assert.NotEmpty(t, createResp.Payload.UpdateTimestamp, "Created config should have an update timestamp")
	require.False(t, createResp.Payload.IsArchived, "Config should not be archived")
	originalConfigCreationTimestamp := createResp.Payload.UpdateTimestamp
	originalConfigID := createResp.Payload.ID

	// 2. READ and verify name/tags are returned
	getParams := integrations.NewGetDataIntegrationConfigParamsWithContext(ctx).WithID(originalConfigID).WithType("cloudwatch")
	getRespOk, err := apiClient.Integrations.GetDataIntegrationConfig(getParams, nil)

	require.NoError(t, err, "Get config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	require.False(t, getRespOk.Payload.IsArchived, "Config should not be archived")
	assert.Equal(t, testCloudwatchConfig, getRespOk.Payload.Config, "Retrieved config value should match the created one")
	// TODO: Uncomment after dev environment is updated to support name/tags (>= 1.9.743)
	// assert.Equal(t, testName, getRespOk.Payload.Name, "Retrieved config name should match the created one")
	// assert.Equal(t, testTags, getRespOk.Payload.Tags, "Retrieved config tags should match the created one")

	// 3. UPDATE with new name and tags
	updatedName := "e2e-test-cloudwatch-updated"
	updatedTags := map[string]interface{}{
		"env":     "test-updated",
		"purpose": "e2e-testing-updated",
	}
	updateBody := &models.CreateDataIntegrationConfigRequest{
		Config: testCloudwatchConfigUpdated,
		Name:   updatedName,
		Tags:   updatedTags,
	}
	updateParams := integrations.NewUpdateDataIntegrationConfigParamsWithContext(ctx).WithType("cloudwatch").WithID(originalConfigID).WithBody(updateBody)
	updateResp, err := apiClient.Integrations.UpdateDataIntegrationConfig(updateParams, nil)

	require.NoError(t, err, "Update config request failed")
	require.NotNil(t, updateResp, "Update config response should not be nil")
	require.NotNil(t, updateResp.Payload, "Update config response payload should not be nil")
	assert.Equal(t, testCloudwatchConfigUpdated, updateResp.Payload.Config, "Updated config value should match the request")
	// TODO: Uncomment after dev environment is updated to support name/tags (>= 1.9.743)
	// assert.Equal(t, updatedName, updateResp.Payload.Name, "Updated config name should match the request")
	// assert.Equal(t, updatedTags, updateResp.Payload.Tags, "Updated config tags should match the request")
	assert.NotEmpty(t, updateResp.Payload.UpdateTimestamp, "Updated config should have an update timestamp")

	getRespOk, err = apiClient.Integrations.GetDataIntegrationConfig(getParams, nil)
	require.NoError(t, err, "Get updated config request failed")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	require.False(t, getRespOk.Payload.IsArchived, "Config should not be archived")
	assert.Equal(t, testCloudwatchConfigUpdated, getRespOk.Payload.Config, "Retrieved config should have updated value")
	// TODO: Uncomment after dev environment is updated to support name/tags (>= 1.9.743)
	// assert.Equal(t, updatedName, getRespOk.Payload.Name, "Retrieved config name should be updated")
	// assert.Equal(t, updatedTags, getRespOk.Payload.Tags, "Retrieved config tags should be updated")
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
