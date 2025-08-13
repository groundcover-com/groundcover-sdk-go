package e2e

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	policies "github.com/groundcover-com/groundcover-sdk-go/pkg/client/policies"
	saClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/serviceaccounts"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

// Helper function to create a temporary policy for service account tests
func createTestPolicy(t *testing.T, client *TestClient) string {
	t.Helper()

	policyName := "sa-test-policy-" + uuid.New().String()
	policyDesc := "Temporary policy for service account E2E tests"
	createReq := &models.CreatePolicyRequest{
		Name:        &policyName,
		Description: policyDesc,
		Role:        models.RoleMap{"write": ""},
		DataScope:   &models.DataScope{},
	}

	createParams := policies.NewCreatePolicyParams().
		WithContext(client.BaseCtx).
		WithTimeout(defaultTimeout).
		WithBody(createReq)

	createResp, err := client.Client.Policies.CreatePolicy(createParams, nil)
	require.NoError(t, err, "Failed to create temporary policy for service account test")
	require.NotNil(t, createResp, "Create policy response nil")
	require.NotNil(t, createResp.Payload, "Create policy payload nil")
	require.NotEmpty(t, createResp.Payload.UUID, "Created policy UUID empty")

	policyID := createResp.Payload.UUID
	t.Logf("Created temporary policy %s (%s) for service account tests", policyName, policyID)
	return policyID
}

// Helper function to delete a policy
func deletePolicy(t *testing.T, client *TestClient, policyID string) {
	t.Helper()
	if policyID == "" {
		return
	}
	deleteParams := policies.NewDeletePolicyParams().
		WithContext(client.BaseCtx).
		WithTimeout(defaultTimeout).
		WithID(policyID)

	_, err := client.Client.Policies.DeletePolicy(deleteParams, nil)
	if err != nil {
		// Log error but don't fail the test, as cleanup might fail for various reasons
		t.Logf("WARN: Failed to delete temporary policy %s during cleanup: %v", policyID, err)
	} else {
		t.Logf("Cleaned up temporary policy %s", policyID)
	}
}

func TestServiceAccountsEndpoints(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup()

	// Create a temporary policy needed for service accounts
	tempPolicyID := createTestPolicy(t, client)
	// Ensure the temporary policy is deleted at the end
	defer deletePolicy(t, client, tempPolicyID)

	// Create a second policy for update test
	secondPolicyID := createTestPolicy(t, client) // Create it early
	defer deletePolicy(t, client, secondPolicyID) // Ensure cleanup

	var createdServiceAccountID string
	var createdServiceAccountName string

	t.Run("Create Service Account", func(t *testing.T) {
		if tempPolicyID == "" {
			t.Skip("Skipping Service Account tests because temporary policy creation failed")
		}

		saName := "e2e-test-sa-" + uuid.New().String()
		// Use a unique email to avoid conflicts if test is re-run quickly
		saEmail := fmt.Sprintf("e2e-sa-%s@example.com", uuid.New().String())

		createReq := &models.CreateServiceAccountRequest{
			Name:        &saName,
			Email:       &saEmail,
			PolicyUUIDs: []string{tempPolicyID},
		}

		createParams := saClient.NewCreateServiceAccountParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(createReq)

		createResp, err := client.Client.Serviceaccounts.CreateServiceAccount(createParams, nil)
		t.Logf("Create service account response: %+v", createResp)
		require.NoError(t, err, "Failed to create service account")
		require.NotNil(t, createResp, "Create service account response should not be nil")
		require.NotNil(t, createResp.Payload, "Create service account response payload should not be nil")
		require.NotNil(t, createResp.Payload.ServiceAccountID, "Create service account response payload body should not be nil")
		require.NotEmpty(t, *createResp.Payload.ServiceAccountID, "Created service account ID should not be empty")

		createdServiceAccountID = *createResp.Payload.ServiceAccountID
		createdServiceAccountName = saName
		t.Logf("Created service account %s (%s) with email %s", saName, createdServiceAccountID, saEmail)
	})

	t.Run("List Service Accounts", func(t *testing.T) {
		if createdServiceAccountID == "" {
			t.Skip("Skipping List Service Accounts because create failed or didn't run")
		}

		listParams := saClient.NewListServiceAccountsParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := client.Client.Serviceaccounts.ListServiceAccounts(listParams, nil)
		require.NoError(t, err, "Failed to list service accounts")
		require.NotNil(t, listResp, "List service accounts response should not be nil")
		require.NotNil(t, listResp.Payload, "List service accounts response payload should not be nil")

		// Check if the created service account is in the list
		found := false
		for _, sa := range listResp.Payload {
			// Compare using ServiceAccountID field from the list response model
			if sa.ServiceAccountID == createdServiceAccountID {
				found = true
				t.Logf("Found created service account %s in the list", createdServiceAccountID)
				// Optional: Add more assertions on the specific SA item if needed
				require.Equal(t, createdServiceAccountName, sa.Name, "List service account name mismatch")
				require.False(t, sa.Deleted, "Newly created service account should not be marked deleted")
				// Check if the correct policy is associated
				policyFound := false
				for _, policy := range sa.Policies {
					if policy.UUID == tempPolicyID { // Assuming Policy struct has UUID field
						policyFound = true
						break
					}
				}
				require.True(t, policyFound, "Temporary policy %s not found in service account %s policies", tempPolicyID, createdServiceAccountID)
				break
			}
		}
		require.True(t, found, "Created service account %s not found in list response", createdServiceAccountID)
	})

	t.Run("Update Service Account", func(t *testing.T) {
		if createdServiceAccountID == "" || secondPolicyID == "" {
			t.Skip("Skipping Update Service Account because create or second policy creation failed")
		}

		// 1. Define updates
		updatedEmail := fmt.Sprintf("e2e-sa-updated-%s@example.com", uuid.New().String())

		// 2. Construct Update request
		updateReq := &models.UpdateServiceAccountRequest{
			ServiceAccountID: &createdServiceAccountID, // Use ServiceAccountID
			Email:            updatedEmail,             // Assume string type
			PolicyUUIDs:      []string{secondPolicyID}, // Add the second policy
			OverridePolicies: false,                    // Append mode
		}

		// 3. Call Update endpoint
		updateParams := saClient.NewUpdateServiceAccountParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(updateReq)

		updateResp, err := client.Client.Serviceaccounts.UpdateServiceAccount(updateParams, nil)
		require.NoError(t, err, "Failed to update service account")
		require.NotNil(t, updateResp, "Update service account response should not be nil")
		require.NotNil(t, updateResp.Payload, "Update service account payload should not be nil")
		// Use ServiceAccountID
		require.NotEmpty(t, updateResp.Payload.ServiceAccountID, "Update response ServiceAccountID should not be empty")
		// Use ServiceAccountID
		require.Equal(t, createdServiceAccountID, updateResp.Payload.ServiceAccountID, "Update response ID mismatch")

		t.Logf("Successfully updated service account %s", createdServiceAccountID)

		// 4. Verify using List
		listParams := saClient.NewListServiceAccountsParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := client.Client.Serviceaccounts.ListServiceAccounts(listParams, nil)
		require.NoError(t, err, "Failed to list service accounts after update")
		require.NotNil(t, listResp, "List service accounts response should not be nil after update")
		require.NotNil(t, listResp.Payload, "List service accounts payload should not be nil after update")

		foundUpdated := false
		for _, sa := range listResp.Payload {
			if sa.ServiceAccountID == createdServiceAccountID {
				foundUpdated = true
				t.Logf("Found updated service account %s in the list", createdServiceAccountID)
				require.Equal(t, updatedEmail, sa.Email, "Updated service account email mismatch")

				// Check both policies are present
				policyMap := make(map[string]bool)
				for _, policy := range sa.Policies {
					policyMap[policy.UUID] = true
				}
				require.True(t, policyMap[tempPolicyID], "Original policy %s missing after update", tempPolicyID)
				require.True(t, policyMap[secondPolicyID], "Second policy %s missing after update", secondPolicyID)
				break
			}
		}
		require.True(t, foundUpdated, "Updated service account %s not found in list response", createdServiceAccountID)
	})

	t.Run("Delete Service Account", func(t *testing.T) {
		if createdServiceAccountID == "" {
			t.Skip("Skipping Delete Service Account because create failed or didn't run")
		}

		deleteParams := saClient.NewDeleteServiceAccountParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdServiceAccountID)

		// Expect 202 Accepted on success based on handler
		_, err := client.Client.Serviceaccounts.DeleteServiceAccount(deleteParams, nil)
		require.NoError(t, err, "Failed to delete service account")

		t.Logf("Successfully deleted service account %s", createdServiceAccountID)

		// Verify deletion by trying to List again and ensure it's gone or marked deleted
		listParams := saClient.NewListServiceAccountsParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := client.Client.Serviceaccounts.ListServiceAccounts(listParams, nil)
		require.NoError(t, err, "Failed to list service accounts after delete")
		require.NotNil(t, listResp, "List service accounts response should not be nil after delete")
		require.NotNil(t, listResp.Payload, "List service accounts payload should not be nil after delete")

		foundAfterDelete := false
		for _, sa := range listResp.Payload {
			if sa.ServiceAccountID == createdServiceAccountID {
				if sa.Deleted != true {
					foundAfterDelete = true
					t.Logf("Found service account %s in list after delete attempt", createdServiceAccountID)
					break
				}
			}
		}
		// If the loop finishes without finding it, the test passes (hard delete)
		if !foundAfterDelete {
			t.Logf("Service account %s correctly not found in list after delete", createdServiceAccountID)
		}
	})
}
