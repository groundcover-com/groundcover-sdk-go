package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	secretClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/secret"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const (
	testSecretNamePrefix = "sdk-e2e-test-secret-"
)

func TestSecretsE2E(t *testing.T) {
	ctx, apiClient := setupTestClient(t)

	var createdSecretID string
	secretName := fmt.Sprintf("%s%d", testSecretNamePrefix, time.Now().UnixNano())

	t.Run("Create Secret", func(t *testing.T) {
		secretType := "api_key"
		content := "test-secret-content"

		createReq := &models.CreateSecretRequest{
			Name:    &secretName,
			Type:    &secretType,
			Content: &content,
		}
		createParams := secretClient.NewCreateSecretParamsWithContext(ctx).WithBody(createReq)
		createResp, err := apiClient.Secret.CreateSecret(createParams, nil)
		require.NoError(t, err, "CreateSecret failed")
		require.NotNil(t, createResp.Payload)
		require.NotEmpty(t, createResp.Payload.ID)
		assert.Equal(t, secretName, createResp.Payload.Name)
		assert.Equal(t, secretType, createResp.Payload.Type)

		createdSecretID = createResp.Payload.ID
		t.Logf("Successfully created Secret with ID: %s", createdSecretID)
	})

	require.NotEmpty(t, createdSecretID, "Secret ID was not set after creation")

	t.Run("Update Secret", func(t *testing.T) {
		secretType := "api_key"
		updatedContent := "updated-secret-content"

		updateReq := &models.UpdateSecretRequest{
			Name:    &secretName,
			Type:    &secretType,
			Content: &updatedContent,
		}
		updateParams := secretClient.NewUpdateSecretParamsWithContext(ctx).
			WithID(createdSecretID).
			WithBody(updateReq)

		updateResp, err := apiClient.Secret.UpdateSecret(updateParams, nil)
		require.NoError(t, err, "UpdateSecret failed")
		require.NotNil(t, updateResp.Payload)
		assert.Equal(t, createdSecretID, updateResp.Payload.ID)
		assert.Equal(t, secretName, updateResp.Payload.Name)

		t.Logf("Successfully updated Secret ID: %s", createdSecretID)
	})

	t.Run("Delete Secret", func(t *testing.T) {
		deleteParams := secretClient.NewDeleteSecretParamsWithContext(ctx).WithID(createdSecretID)
		_, err := apiClient.Secret.DeleteSecret(deleteParams, nil)
		require.NoError(t, err, "DeleteSecret failed")
		t.Logf("Successfully deleted Secret ID: %s", createdSecretID)
	})

	t.Run("Create Secret with Different Types", func(t *testing.T) {
		secretTypes := []string{"api_key", "password", "basic_auth"}
		var createdIDs []string

		defer func() {
			for _, id := range createdIDs {
				deleteParams := secretClient.NewDeleteSecretParamsWithContext(ctx).WithID(id)
				_, err := apiClient.Secret.DeleteSecret(deleteParams, nil)
				if err != nil {
					t.Logf("Warning: Could not delete secret %s: %v", id, err)
				}
			}
		}()

		for _, secretType := range secretTypes {
			name := fmt.Sprintf("%s%d-%s", testSecretNamePrefix, time.Now().UnixNano(), secretType)
			content := "test-content"

			createReq := &models.CreateSecretRequest{
				Name:    &name,
				Type:    &secretType,
				Content: &content,
			}
			createParams := secretClient.NewCreateSecretParamsWithContext(ctx).WithBody(createReq)
			createResp, err := apiClient.Secret.CreateSecret(createParams, nil)
			require.NoError(t, err, "CreateSecret failed for type %s", secretType)
			require.NotNil(t, createResp.Payload)
			assert.Equal(t, secretType, createResp.Payload.Type)

			createdIDs = append(createdIDs, createResp.Payload.ID)
			t.Logf("Created Secret with type %s, ID: %s", secretType, createResp.Payload.ID)
		}
	})

	t.Run("Create Secret with Managed Provider", func(t *testing.T) {
		name := fmt.Sprintf("%s%d-terraform", testSecretNamePrefix, time.Now().UnixNano())
		secretType := "api_key"
		content := "terraform-managed-secret"
		managedBy := "terraform"

		createReq := &models.CreateSecretRequest{
			Name:              &name,
			Type:              &secretType,
			Content:           &content,
			ManagedByProvider: managedBy,
		}
		createParams := secretClient.NewCreateSecretParamsWithContext(ctx).WithBody(createReq)
		createResp, err := apiClient.Secret.CreateSecret(createParams, nil)
		require.NoError(t, err, "CreateSecret with managed provider failed")
		require.NotNil(t, createResp.Payload)
		require.NotEmpty(t, createResp.Payload.ID)

		createdID := createResp.Payload.ID
		t.Logf("Successfully created Secret with managed provider, ID: %s", createdID)

		// Cleanup
		defer func() {
			deleteParams := secretClient.NewDeleteSecretParamsWithContext(ctx).WithID(createdID)
			_, err := apiClient.Secret.DeleteSecret(deleteParams, nil)
			if err != nil {
				t.Logf("Warning: Could not delete secret %s: %v", createdID, err)
			}
		}()
	})

	t.Run("Update Non-Existent Secret", func(t *testing.T) {
		nonExistentID := "secretRef::store::00000000-0000-0000-0000-000000000000"
		name := "non-existent"
		secretType := "api_key"
		content := "test"

		updateReq := &models.UpdateSecretRequest{
			Name:    &name,
			Type:    &secretType,
			Content: &content,
		}
		updateParams := secretClient.NewUpdateSecretParamsWithContext(ctx).
			WithID(nonExistentID).
			WithBody(updateReq)

		_, err := apiClient.Secret.UpdateSecret(updateParams, nil)
		assert.Error(t, err, "Update non-existent secret should fail")
		t.Logf("Update non-existent secret correctly returned error")
	})

	t.Run("Delete Non-Existent Secret", func(t *testing.T) {
		nonExistentID := "secretRef::store::00000000-0000-0000-0000-000000000000"

		deleteParams := secretClient.NewDeleteSecretParamsWithContext(ctx).WithID(nonExistentID)
		_, err := apiClient.Secret.DeleteSecret(deleteParams, nil)
		assert.Error(t, err, "Delete non-existent secret should fail")
		t.Logf("Delete non-existent secret correctly returned error")
	})
}
