package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8sClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/k8s"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func TestEventsSearchE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	t.Run("Execute Events Search", func(t *testing.T) {
		// Define search parameters
		startDateTime := strfmt.DateTime(time.Now().Add(-24 * time.Hour))
		endDateTime := strfmt.DateTime(time.Now())
		query := "* | stats count(*)"

		// Construct the request body
		body := &models.EventsSearchRequest{
			Start: &startDateTime,
			End:   &endDateTime,
			Query: query,
		}

		// Create parameters
		params := k8sClient.NewEventsSearchParamsWithContext(ctx).WithBody(body)

		// Execute the search
		resp, err := apiClient.K8s.EventsSearch(params, nil)

		// Assertions
		require.NoError(t, err, "Events search request failed")
		require.NotNil(t, resp, "Events search response should not be nil")
		require.NotNil(t, resp.Payload, "Events search response payload should not be nil")

		var result struct {
			Count int32 `json:"count()"`
		}

		payloadSlice, ok := resp.Payload.([]interface{})
		require.True(t, ok, "Payload should be a slice of interfaces")

		jsonBytes, err := json.Marshal(payloadSlice[0])
		require.NoError(t, err, "Failed to marshal events search response payload")

		err = json.Unmarshal(jsonBytes, &result)
		require.NoError(t, err, "Failed to unmarshal events search response payload")

		assert.Greater(t, result.Count, int32(0), "Events search should return a count greater than 0")
	})
}
