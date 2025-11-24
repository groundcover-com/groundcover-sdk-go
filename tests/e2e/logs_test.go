package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logsClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/logs"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestLogsSearchE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	t.Run("Execute Logs Search", func(t *testing.T) {
		// Define search parameters
		startDateTime := strfmt.DateTime(time.Now().Add(-24 * time.Hour))
		endDateTime := strfmt.DateTime(time.Now())
		query := "* | stats count(*)"

		// Construct the request body
		body := &models.LogsSearchRequest{
			Start: &startDateTime,
			End:   &endDateTime,
			Query: query,
		}

		// Create parameters
		params := logsClient.NewSearchLogsParamsWithContext(ctx).WithBody(body)

		// Execute the search
		resp, err := apiClient.Logs.SearchLogs(params, nil)

		// Assertions
		require.NoError(t, err, "Logs search request failed")
		require.NotNil(t, resp, "Logs search response should not be nil")
		require.NotNil(t, resp.Payload, "Logs search response payload should not be nil")

		var result struct {
			Count int64 `json:"count"`
		}

		payloadSlice, ok := resp.Payload.([]interface{})
		require.True(t, ok, "Payload should be a slice of interfaces")

		jsonBytes, err := json.Marshal(payloadSlice[0])
		require.NoError(t, err, "Failed to marshal logs search response payload")

		err = json.Unmarshal(jsonBytes, &result)
		require.NoError(t, err, "Failed to unmarshal logs search response payload")

		assert.Greater(t, result.Count, int64(0), "Logs search should return a count greater than 0")
	})
}
