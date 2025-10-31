package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tracesClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/traces"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestTracesSearchE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	t.Run("Execute Traces Search", func(t *testing.T) {
		// Define search parameters
		startDateTime := strfmt.DateTime(time.Now().Add(-24 * time.Hour))
		endDateTime := strfmt.DateTime(time.Now())
		query := "* | stats count(*)"

		// Construct the request body
		body := &models.TracesSearchRequest{
			Start: &startDateTime,
			End:   &endDateTime,
			Query: query,
		}

		// Create parameters
		params := tracesClient.NewSearchTracesParamsWithContext(ctx).WithBody(body)

		// Execute the search
		resp, err := apiClient.Traces.SearchTraces(params, nil)

		// Assertions
		require.NoError(t, err, "Traces search request failed")
		require.NotNil(t, resp, "Traces search response should not be nil")
		require.NotNil(t, resp.Payload, "Traces search response payload should not be nil")

		var result struct {
			Count int64 `json:"count()"`
		}

		payloadSlice, ok := resp.Payload.([]interface{})
		require.True(t, ok, "Payload should be a slice of interfaces")

		jsonBytes, err := json.Marshal(payloadSlice[0])
		require.NoError(t, err, "Failed to marshal traces search response payload")

		err = json.Unmarshal(jsonBytes, &result)
		require.NoError(t, err, "Failed to unmarshal traces search response payload")

		assert.Greater(t, result.Count, int64(0), "Traces search should return a count greater than 0")
	})
}
