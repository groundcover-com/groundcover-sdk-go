package e2e

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/synthetics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestSyntheticsEndpoints(t *testing.T) {
	testClient := NewTestClient(t)
	defer testClient.Cleanup()

	const defaultTimeout = 30 * time.Second

	var createdSyntheticID string
	syntheticName := "e2e-test-synthetic-" + uuid.New().String()

	t.Run("Create Synthetic Test", func(t *testing.T) {
		createReq := &models.SyntheticTestCreateRequest{
			Name:     syntheticName,
			Version:  1,
			Enabled:  true,
			Interval: "5m",
			CheckConfig: &models.WorkerRequest{
				Kind: "http",
				Metadata: &models.Metadata{
					SyntheticName: syntheticName,
				},
				Request: &models.Request{
					HTTP: &models.HTTPRequest{
						Kind:    "http",
						Method:  "GET",
						URL:     "https://httpbin.org/get",
						Timeout: "30s",
					},
				},
				ExecutionPolicy: &models.ExecutionPolicy{
					Assertions: []*models.Assertion{
						{
							Source:   "statusCode",
							Operator: "eq",
							Target:   "200",
						},
					},
				},
				Tracing: &models.Tracing{},
			},
		}

		createParams := synthetics.NewCreateSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(createReq)

		createResp, err := testClient.Client.Synthetics.CreateSyntheticTest(createParams, nil)
		require.NoError(t, err, "Failed to create synthetic test")
		require.NotNil(t, createResp, "Create synthetic test response should not be nil")
		require.NotNil(t, createResp.Payload, "Create synthetic test response payload should not be nil")
		require.NotEmpty(t, createResp.Payload.ID, "Created synthetic test ID should not be empty")

		createdSyntheticID = createResp.Payload.ID
		t.Logf("✓ Created synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("List Synthetic Tests", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping List Synthetic Tests because create failed or didn't run")
		}

		listParams := synthetics.NewListSyntheticTestsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout * 4)

		listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
		require.NoError(t, err, "Failed to list synthetic tests")
		require.NotNil(t, listResp, "List synthetic tests response should not be nil")
		require.NotNil(t, listResp.Payload, "List synthetic tests response payload should not be nil")

		t.Logf("✓ SDK List Synthetic Tests call succeeded with %d synthetics", len(listResp.Payload.Synthetics))

		found := false
		for _, item := range listResp.Payload.Synthetics {
			if item.ID == createdSyntheticID {
				found = true
				require.Equal(t, syntheticName, item.Name, "Synthetic test name mismatch")
				require.Equal(t, "http", item.SyntheticType, "Synthetic test type should be 'http'")
				require.NotEmpty(t, item.Status, "Synthetic test status should not be empty")
				t.Logf("Found created synthetic test %s in the list (status: %s)", createdSyntheticID, item.Status)
				break
			}
		}
		require.True(t, found, "Created synthetic test %s not found in list response", createdSyntheticID)
	})

	t.Run("Get Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping Get Synthetic Test because create failed or didn't run")
		}

		getParams := synthetics.NewGetSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		getResp, err := testClient.Client.Synthetics.GetSyntheticTest(getParams, nil)
		require.NoError(t, err, "Failed to get synthetic test")
		require.NotNil(t, getResp, "Get synthetic test response should not be nil")

		t.Logf("✓ Successfully retrieved synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Update Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping Update Synthetic Test because create failed or didn't run")
		}

		updatedName := syntheticName + "-updated"
		updateReq := &models.SyntheticTestCreateRequest{
			Name:     updatedName,
			Version:  1,
			Enabled:  true,
			Interval: "10m",
			CheckConfig: &models.WorkerRequest{
				Kind: "http",
				Metadata: &models.Metadata{
					SyntheticName: updatedName,
				},
				Request: &models.Request{
					HTTP: &models.HTTPRequest{
						Kind:    "http",
						Method:  "GET",
						URL:     "https://httpbin.org/get",
						Timeout: "30s",
					},
				},
				ExecutionPolicy: &models.ExecutionPolicy{
					Assertions: []*models.Assertion{
						{
							Source:   "statusCode",
							Operator: "eq",
							Target:   "200",
						},
					},
				},
				Tracing: &models.Tracing{},
			},
		}

		updateParams := synthetics.NewUpdateSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID).
			WithBody(updateReq)

		updateResp, err := testClient.Client.Synthetics.UpdateSyntheticTest(updateParams, nil)
		require.NoError(t, err, "Failed to update synthetic test")
		require.NotNil(t, updateResp, "Update synthetic test response should not be nil")

		t.Logf("✓ Successfully updated synthetic test with ID: %s", createdSyntheticID)

		// Verify update by listing and checking the name
		listParams := synthetics.NewListSyntheticTestsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
		require.NoError(t, err, "Failed to list synthetic tests after update")

		found := false
		for _, item := range listResp.Payload.Synthetics {
			if item.ID == createdSyntheticID {
				found = true
				require.Equal(t, updatedName, item.Name, "Updated synthetic test name mismatch")
				t.Logf("Verified synthetic test %s was updated with new name: %s", createdSyntheticID, item.Name)
				break
			}
		}
		require.True(t, found, "Updated synthetic test %s not found in list response", createdSyntheticID)

		syntheticName = updatedName
	})

	t.Run("Delete Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping Delete Synthetic Test because create failed or didn't run")
		}

		deleteParams := synthetics.NewDeleteSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		_, err := testClient.Client.Synthetics.DeleteSyntheticTest(deleteParams, nil)
		require.NoError(t, err, "Failed to delete synthetic test")

		t.Logf("✓ Successfully deleted synthetic test %s", createdSyntheticID)

		// Verify deletion by listing synthetic tests
		listParams := synthetics.NewListSyntheticTestsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
		require.NoError(t, err, "Failed to list synthetic tests after deletion")

		found := false
		for _, item := range listResp.Payload.Synthetics {
			if item.ID == createdSyntheticID {
				found = true
				break
			}
		}
		require.False(t, found, "Deleted synthetic test %s should not be found in list response", createdSyntheticID)
		t.Logf("Verified synthetic test %s is no longer in the list after deletion", createdSyntheticID)
	})
}
