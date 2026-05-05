package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metricsPipelineClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics_pipeline"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestRemoteConfigMetricsPipelineCrudE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	// 1. CREATE
	createBody := &models.CreateOrUpdateMetricsPipelineConfigRequest{
		Rules: &models.RelabelConfig{
			KeepRegex: []string{"http_requests_total", "process_cpu_seconds_total"},
			AddLabel:  map[string]string{"team": "platform"},
		},
	}

	createParams := metricsPipelineClient.NewCreateMetricsPipelineConfigParamsWithContext(ctx).WithBody(createBody)
	createResp, err := apiClient.MetricsPipeline.CreateMetricsPipelineConfig(createParams, nil)

	require.NoError(t, err, "Create config request failed")
	require.NotNil(t, createResp, "Create config response should not be nil")
	require.NotNil(t, createResp.Payload, "Create config response payload should not be nil")
	require.NotNil(t, createResp.Payload.Rules, "Created config should have rules")

	assert.Equal(t, []string{"http_requests_total", "process_cpu_seconds_total"}, createResp.Payload.Rules.KeepRegex)
	assert.Equal(t, map[string]string{"team": "platform"}, createResp.Payload.Rules.AddLabel)
	assert.NotEmpty(t, createResp.Payload.UUID, "Created config should have an ID")
	assert.NotEmpty(t, createResp.Payload.CreatedTimestamp, "Created config should have a creation timestamp")
	originalConfigCreationTimestamp := createResp.Payload.CreatedTimestamp

	// 2. READ
	getParams := metricsPipelineClient.NewGetMetricsPipelineConfigParamsWithContext(ctx)
	getRespOk, getRespNoContent, err := apiClient.MetricsPipeline.GetMetricsPipelineConfig(getParams, nil)

	require.NoError(t, err, "Get config request failed")
	require.Nil(t, getRespNoContent, "Shouldnt get a 204 No Content response")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload, "Get config response payload should not be nil")
	require.NotNil(t, getRespOk.Payload.Rules, "Retrieved config should have rules")

	assert.Equal(t, []string{"http_requests_total", "process_cpu_seconds_total"}, getRespOk.Payload.Rules.KeepRegex)
	assert.Equal(t, map[string]string{"team": "platform"}, getRespOk.Payload.Rules.AddLabel)

	// 3. UPDATE
	updateBody := &models.CreateOrUpdateMetricsPipelineConfigRequest{
		Rules: &models.RelabelConfig{
			KeepRegex: []string{"http_requests_total", "node_cpu_seconds_total"},
			DropRegex: []string{"go_.*"},
			AddLabel:  map[string]string{"team": "platform", "env": "staging"},
		},
	}

	updateParams := metricsPipelineClient.NewUpdateMetricsPipelineConfigParamsWithContext(ctx).WithBody(updateBody)
	updateResp, err := apiClient.MetricsPipeline.UpdateMetricsPipelineConfig(updateParams, nil)

	require.NoError(t, err, "Update config request failed")
	require.NotNil(t, updateResp, "Update config response should not be nil")
	require.NotNil(t, updateResp.Payload, "Update config response payload should not be nil")
	require.NotNil(t, updateResp.Payload.Rules, "Updated config should have rules")

	assert.Equal(t, []string{"http_requests_total", "node_cpu_seconds_total"}, updateResp.Payload.Rules.KeepRegex)
	assert.Equal(t, []string{"go_.*"}, updateResp.Payload.Rules.DropRegex)
	assert.Equal(t, map[string]string{"team": "platform", "env": "staging"}, updateResp.Payload.Rules.AddLabel)
	assert.NotEmpty(t, updateResp.Payload.CreatedTimestamp, "Updated config should have an update timestamp")

	getRespOk, getRespNoContent, err = apiClient.MetricsPipeline.GetMetricsPipelineConfig(getParams, nil)
	require.NoError(t, err, "Get updated config request failed")
	require.Nil(t, getRespNoContent, "Shouldnt get a 204 No Content response")
	require.NotNil(t, getRespOk, "Get config response should not be nil")
	require.NotNil(t, getRespOk.Payload.Rules, "Retrieved config should have rules")
	assert.Equal(t, []string{"go_.*"}, getRespOk.Payload.Rules.DropRegex)
	assert.Greater(t, updateResp.Payload.CreatedTimestamp, originalConfigCreationTimestamp, "Updated config should have a creation timestamp greater than the original one")

	// 4. DELETE
	deleteParams := metricsPipelineClient.NewDeleteMetricsPipelineConfigParamsWithContext(ctx)
	_, err = apiClient.MetricsPipeline.DeleteMetricsPipelineConfig(deleteParams, nil)

	require.NoError(t, err, "Delete config request failed")

	getRespOk, getRespNoContent, err = apiClient.MetricsPipeline.GetMetricsPipelineConfig(getParams, nil)
	require.NoError(t, err, "Get deleted config request failed")
	require.Nil(t, getRespNoContent, "Shouldnt get a 204 No Content response")
	require.Nil(t, getRespOk.Payload.Rules, "Deleted config should have nil rules")
}
