package e2e

import (
	"net/http"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/google/uuid"
	policies "github.com/groundcover-com/groundcover-sdk-go/pkg/client/policies"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestPoliciesEndpoints(t *testing.T) {
	client := NewTestClient(t)
	defer client.Cleanup()

	var createdPolicyID string
	var createdPolicyName string // Store name for later tests

	t.Run("Create Policy", func(t *testing.T) {
		policyName := "e2e-test-policy-" + uuid.New().String()
		policyDesc := "Policy created during E2E testing"
		createReq := &models.CreatePolicyRequest{
			Name:        &policyName,
			Description: policyDesc,
			Role: models.RoleMap{
				"write": "",
			},
			DataScope: &models.DataScope{},
		}

		createParams := policies.NewCreatePolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(createReq)

		createResp, err := client.Client.Policies.CreatePolicy(createParams, nil)
		require.NoError(t, err, "Failed to create policy")
		require.NotNil(t, createResp, "Create policy response should not be nil")
		require.NotNil(t, createResp.Payload, "Create policy response payload should not be nil")

		// Access the ID from the *models.Policy payload (field is UUID)
		require.NotEmpty(t, createResp.Payload.UUID, "Created policy UUID should not be empty")

		createdPolicyID = createResp.Payload.UUID // Use UUID
		createdPolicyName = policyName            // Save name for Get test
		t.Logf("Created policy with ID: %s", createdPolicyID)
	})

	// Add subsequent tests using createdPolicyID
	t.Run("Get Policy", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping Get Policy test because create failed or didn't run")
		}

		getParams := policies.NewGetPolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID)

		getResp, err := client.Client.Policies.GetPolicy(getParams, nil)
		require.NoError(t, err, "Failed to get policy")
		require.NotNil(t, getResp, "Get policy response should not be nil")
		require.NotNil(t, getResp.Payload, "Get policy response payload should not be nil")
		require.NotNil(t, getResp.Payload.UUID, "Get policy response payload PolicyID should not be nil")

		// Assert fields from the Get response (assume ID is PolicyID *string)
		require.Equal(t, createdPolicyID, getResp.Payload.UUID, "Get policy ID mismatch")
		require.Equal(t, createdPolicyName, *getResp.Payload.Name, "Get policy name mismatch")
		// Add more assertions as needed (e.g., description, role)
		t.Logf("Successfully retrieved policy with ID: %s", getResp.Payload.UUID)
	})

	t.Run("List Policies", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping List Policies test because create failed or didn't run")
		}

		listParams := policies.NewListPoliciesParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := client.Client.Policies.ListPolicies(listParams, nil)
		require.NoError(t, err, "Failed to list policies")
		require.NotNil(t, listResp, "List policies response should not be nil")
		require.NotNil(t, listResp.Payload, "List policies response payload should not be nil")

		// Check if the created policy is in the list
		found := false
		for _, policy := range listResp.Payload {
			if policy.UUID == createdPolicyID {
				found = true
				t.Logf("Found created policy %s in the list", createdPolicyID)
				// Optional: Add more assertions on the specific policy item if needed
				require.Equal(t, createdPolicyName, *policy.Name, "List policy name mismatch")
				break
			}
		}
		require.True(t, found, "Created policy %s not found in list response", createdPolicyID)
	})

	t.Run("Update Policy", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping Update Policy test because create failed or didn't run")
		}

		// Note: We assume the models.Policy returned by Get doesn't contain all fields
		// needed for update (like Revision). We'll construct the update request
		// based on the UpdatePolicyRequest definition.

		// 1. Define updates
		updatedDesc := "Policy updated during E2E testing"
		updatedRoleMap := models.RoleMap{"admin": ""} // Change role to admin
		// Hardcode initial revision - adjust if API requires different initial value
		// Alternatively, if createResp *did* return revision, use that.
		var assumedCurrentRevision int32 = 1 // Adjust this if needed

		// 2. Construct Update request using required fields
		updateReq := &models.UpdatePolicyRequest{
			Name:        &createdPolicyName,  // Keep original name (use stored name)
			Description: updatedDesc,         // Update description (assume *string)
			Role:        updatedRoleMap,      // Update role
			DataScope:   &models.DataScope{}, // Provide empty scope (like in create)
			// ClaimRole: nil, // Omit if optional
			CurrentRevision: assumedCurrentRevision, // Provide assumed current revision
		}

		// 3. Call Update endpoint
		updateParams := policies.NewUpdatePolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID).
			WithBody(updateReq)

		updateResp, err := client.Client.Policies.UpdatePolicy(updateParams, nil)
		require.NoError(t, err, "Failed to update policy")
		require.NotNil(t, updateResp, "Update policy response should not be nil")
		require.NotNil(t, updateResp.Payload, "Update policy response payload should not be nil")

		// 4. Verify response payload (assuming it *does* return updated fields)
		// If updateResp.Payload also lacks fields, these assertions will fail.
		t.Logf("Policy updated successfully. Update response received.")

		// 5. Get again to verify basic fields persisted
		getUpdatedParams := policies.NewGetPolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID)

		getUpdatedResp, err := client.Client.Policies.GetPolicy(getUpdatedParams, nil)
		require.NoError(t, err, "Failed to get policy after update")
		require.NotNil(t, getUpdatedResp, "Get updated policy response should not be nil")
		require.NotNil(t, getUpdatedResp.Payload, "Get updated policy response payload should not be nil")

		// Only verify fields known to be in models.Policy from Get response
		require.Equal(t, createdPolicyID, getUpdatedResp.Payload.UUID, "Get policy ID mismatch after update")
		require.Equal(t, createdPolicyName, *getUpdatedResp.Payload.Name, "Get policy name mismatch after update")
		t.Logf("Verified policy fields via Get after update.")

	})

	t.Run("Apply Policy", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping Apply Policy test because create failed or didn't run")
		}

		// Define a target email (replace with a valid test email if necessary)
		testEmail := "e2e-test@example.com"

		applyReq := &models.ApplyPolicyRequest{
			PolicyUUIDs: []string{createdPolicyID},
			Emails:      []string{testEmail},
			Override:    false, // Append mode
		}

		applyParams := policies.NewApplyPolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(applyReq)

		applyResp, err := client.Client.Policies.ApplyPolicy(applyParams, nil)

		// Expect 200 OK for success based on handler code c.Status(http.StatusOK)
		require.NoError(t, err, "Failed to apply policy")
		require.NotNil(t, applyResp, "Apply policy response should not be nil")

		// The response body is empty (NoContentResponse), so no payload check needed.
		t.Logf("Successfully called Apply Policy for policy %s to email %s", createdPolicyID, testEmail)

		// Optional TODO: Add verification step if there's an endpoint to check policies applied to a user
	})

	t.Run("Get Policy Audit Trail", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping Get Policy Audit Trail test because create failed or didn't run")
		}

		auditParams := policies.NewGetPolicyAuditTrailParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID)

		auditResp, err := client.Client.Policies.GetPolicyAuditTrail(auditParams, nil)
		require.NoError(t, err, "Failed to get policy audit trail")
		require.NotNil(t, auditResp, "Get policy audit trail response should not be nil")
		require.NotNil(t, auditResp.Payload, "Get policy audit trail response payload should not be nil")

		// Check that the payload is a slice and not empty
		require.IsType(t, []*models.Policy{}, auditResp.Payload, "Audit trail payload should be a slice of policies")
		require.NotEmpty(t, auditResp.Payload, "Audit trail should not be empty")

		t.Logf("Successfully retrieved policy audit trail for policy %s with %d entries", createdPolicyID, len(auditResp.Payload))
		// Optional: Further assertions on audit trail content if needed
		// e.g., check len(auditResp.Payload) >= 2 to account for create and update

	})

	t.Run("Delete Policy", func(t *testing.T) {
		if createdPolicyID == "" {
			t.Skip("Skipping Delete Policy test because create failed or didn't run")
		}

		deleteParams := policies.NewDeletePolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID)

		// Expect 200 No Content on success
		_, err := client.Client.Policies.DeletePolicy(deleteParams, nil)
		require.NoError(t, err, "Failed to delete policy")

		t.Logf("Successfully deleted policy %s", createdPolicyID)

		// Verify deletion by trying to Get the policy again
		getParams := policies.NewGetPolicyParams().
			WithContext(client.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdPolicyID)

		_, err = client.Client.Policies.GetPolicy(getParams, nil)
		require.Error(t, err, "Expected error when getting deleted policy, but got nil")

		// Check if it's the specific GetPolicyNotFound error
		_, ok := err.(*policies.GetPolicyNotFound)
		if ok {
			// It's the expected specific 404 error type
			t.Logf("Correctly received GetPolicyNotFound error for deleted policy %s", createdPolicyID)
		} else {
			// Fallback: Check if it's a generic APIError with 404 status
			apiError, ok := err.(*runtime.APIError)
			require.True(t, ok, "Expected GetPolicyNotFound or APIError, got %T", err)
			require.Equal(t, http.StatusNotFound, apiError.Code, "Expected status code 404 if APIError, got %d", apiError.Code)
			t.Logf("Received generic 404 APIError for deleted policy %s", createdPolicyID)
		}
	})
}
