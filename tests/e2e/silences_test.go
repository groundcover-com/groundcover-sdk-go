package e2e

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	monitors "github.com/groundcover-com/groundcover-sdk-go/pkg/client/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSilencesEndpoints(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup()

	var createdSilenceID string
	var createdSilenceComment string

	t.Run("Create Silence", func(t *testing.T) {
		// Define test silence data
		silenceComment := "e2e-test-silence-" + uuid.New().String()
		startsAt := strfmt.DateTime(time.Now().Add(1 * time.Minute))           // Start in 1 minute
		endsAt := strfmt.DateTime(time.Now().Add(1*time.Hour + 1*time.Minute)) // End in 1 hour and 1 minute

		// Create matchers for the silence using SilenceMatcher format
		isEqualTrue := true
		isEqualFalse := false
		isRegexTrue := true
		isRegexFalse := false
		matchers := []*models.SilenceMatcher{
			{
				Name:    "service",
				Value:   "test-equal",
				IsEqual: &isEqualTrue,
				IsRegex: &isRegexFalse,
			},
			{
				Name:    "environment",
				Value:   "*test-not-equal-regex*",
				IsEqual: &isEqualFalse,
				IsRegex: &isRegexTrue,
			},
			{
				Name:  "workload",
				Value: "test-empty-equal",
				// IsEqual and IsRegex are nil (omitted from JSON, will default to true for IsEqual)
			},
		}

		createReq := &models.CreateSilenceRequest{
			StartsAt: &startsAt,
			EndsAt:   &endsAt,
			Comment:  silenceComment,
			Matchers: matchers,
		}

		createParams := monitors.NewCreateSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(createReq)

		createResp, err := client.Client.Monitors.CreateSilence(createParams, nil)
		if err != nil {
			t.Logf("Create silence error: %v", err)
			t.Logf("Error type: %T", err)
			require.NoError(t, err, "Failed to create silence")
		}
		require.NotNil(t, createResp, "Create silence response should not be nil")
		require.NotNil(t, createResp.Payload, "Create silence response payload should not be nil")

		// Access the UUID from the silence payload
		require.NotEmpty(t, createResp.Payload.UUID, "Created silence UUID should not be empty")

		createdSilenceID = string(createResp.Payload.UUID)
		createdSilenceComment = silenceComment
		t.Logf("Created silence with ID: %s", createdSilenceID)
	})

	// Only run subsequent tests if create succeeded
	if createdSilenceID == "" {
		t.Skip("Skipping remaining tests because create failed")
		return
	}

	t.Run("Get Silence", func(t *testing.T) {
		getParams := monitors.NewGetSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSilenceID)

		getResp, err := client.Client.Monitors.GetSilence(getParams, nil)
		require.NoError(t, err, "Failed to get silence")
		require.NotNil(t, getResp, "Get silence response should not be nil")
		require.NotNil(t, getResp.Payload, "Get silence response payload should not be nil")
		require.NotEmpty(t, getResp.Payload.UUID, "Get silence response payload UUID should not be empty")

		// Assert fields from the Get response
		require.Equal(t, createdSilenceID, string(getResp.Payload.UUID), "Get silence UUID mismatch")
		require.Equal(t, createdSilenceComment, getResp.Payload.Comment, "Get silence comment mismatch")
		require.NotEmpty(t, getResp.Payload.Matchers, "Get silence should have matchers")
		t.Logf("Successfully retrieved silence with ID: %s", string(getResp.Payload.UUID))
	})

	t.Run("Get All Silences", func(t *testing.T) {
		if createdSilenceID == "" {
			t.Skip("Skipping Get All Silences test because create failed or didn't run")
		}

		// Test without filters first
		getAllParams := monitors.NewGetAllSilencesParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout)

		getAllResp, err := client.Client.Monitors.GetAllSilences(getAllParams, nil)
		require.NoError(t, err, "Failed to get all silences")
		require.NotNil(t, getAllResp, "Get all silences response should not be nil")
		require.NotNil(t, getAllResp.Payload, "Get all silences response payload should not be nil")

		// Check if the created silence is in the list
		found := false
		for _, silence := range getAllResp.Payload {
			if string(silence.UUID) == createdSilenceID {
				found = true
				t.Logf("Found created silence %s in the list", createdSilenceID)
				require.Equal(t, createdSilenceComment, silence.Comment, "List silence comment mismatch")
				break
			}
		}
		require.True(t, found, "Created silence %s not found in list response", createdSilenceID)

		// Test with active filter
		activeFilter := true
		getAllActiveParams := monitors.NewGetAllSilencesParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithActive(&activeFilter)

		getAllActiveResp, err := client.Client.Monitors.GetAllSilences(getAllActiveParams, nil)
		require.NoError(t, err, "Failed to get active silences")
		require.NotNil(t, getAllActiveResp, "Get active silences response should not be nil")
		require.NotNil(t, getAllActiveResp.Payload, "Get active silences response payload should not be nil")

		t.Logf("Successfully retrieved %d total silences and %d active silences",
			len(getAllResp.Payload), len(getAllActiveResp.Payload))
	})

	t.Run("Update Silence", func(t *testing.T) {
		if createdSilenceID == "" {
			t.Skip("Skipping Update Silence test because create failed or didn't run test create silence")
		}

		// Define updates
		updatedComment := "Updated silence comment during E2E testing"
		newStartsAt := strfmt.DateTime(time.Now().Add(2 * time.Minute))           // Start in 2 minutes
		newEndsAt := strfmt.DateTime(time.Now().Add(2*time.Hour + 2*time.Minute)) // End in 2 hours and 2 minutes

		// Updated matchers
		isEqualTrue := true
		isRegexFalse := false
		updatedMatchers := []*models.SilenceMatcher{
			{
				Name:    "service",
				Value:   "updated-test-service",
				IsEqual: &isEqualTrue,
				IsRegex: &isRegexFalse,
			},
			{
				Name:    "environment",
				Value:   "production",
				IsEqual: &isEqualTrue,
				IsRegex: &isRegexFalse,
			},
		}

		updateReq := &models.UpdateSilenceRequest{
			StartsAt: newStartsAt,
			EndsAt:   newEndsAt,
			Comment:  updatedComment,
			Matchers: updatedMatchers,
		}

		updateParams := monitors.NewUpdateSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSilenceID).
			WithBody(updateReq)

		updateResp, err := client.Client.Monitors.UpdateSilence(updateParams, nil)
		require.NoError(t, err, "Failed to update silence")
		require.NotNil(t, updateResp, "Update silence response should not be nil")
		require.NotNil(t, updateResp.Payload, "Update silence response payload should not be nil")

		t.Logf("Silence updated successfully. Update response received.")

		// Get again to verify updates persisted
		getUpdatedParams := monitors.NewGetSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSilenceID)

		getUpdatedResp, err := client.Client.Monitors.GetSilence(getUpdatedParams, nil)
		require.NoError(t, err, "Failed to get silence after update")
		require.NotNil(t, getUpdatedResp, "Get updated silence response should not be nil")
		require.NotNil(t, getUpdatedResp.Payload, "Get updated silence response payload should not be nil")

		// Verify updated fields
		require.Equal(t, createdSilenceID, string(getUpdatedResp.Payload.UUID), "Get silence UUID mismatch after update")
		require.Equal(t, updatedComment, getUpdatedResp.Payload.Comment, "Get silence comment mismatch after update")
		require.NotEmpty(t, getUpdatedResp.Payload.Matchers, "Updated silence should have matchers")

		t.Logf("Verified silence fields via Get after update.")
	})

	t.Run("Delete Silence", func(t *testing.T) {
		deleteParams := monitors.NewDeleteSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSilenceID)

		// Expect 200 OK on success
		_, err := client.Client.Monitors.DeleteSilence(deleteParams, nil)
		require.NoError(t, err, "Failed to delete silence")

		t.Logf("Successfully deleted silence %s", createdSilenceID)

		// Verify deletion by trying to Get the silence again
		getParams := monitors.NewGetSilenceParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSilenceID)

		_, err = client.Client.Monitors.GetSilence(getParams, nil)
		require.Error(t, err, "Expected error when getting deleted silence, but got nil")

		// Check if it's a 404
		_, ok := err.(*monitors.GetSilenceNotFound)
		require.True(t, ok, "Expected 404, got %T", err)
		t.Logf("Received expected 404 for deleted silence %s", createdSilenceID)
	})
}
