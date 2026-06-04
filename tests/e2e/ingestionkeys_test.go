package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ingestionKeysClient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/ingestionkeys"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

const (
	testIngestionKeyNamePrefix   = "sdk-e2e-test-ingestion-key"
	testIngestionKeyCreatorEmail = "sdk-e2e-test@groundcover.com"
)

var (
	testIngestionKeyType = "sensor"
)

func TestIngestionKeysE2E(t *testing.T) {
	backendIDOverride := "groundcover-staging"
	tc := NewTestClient(t, TestClientWithBackendID(backendIDOverride))
	defer tc.Cleanup()
	ctx, apiClient := tc.BaseCtx, tc.Client

	ingestionKeyName := fmt.Sprintf("%s-%d", testIngestionKeyNamePrefix, time.Now().UnixNano())
	var createdKey string

	ingestionKeyReq := models.CreateIngestionKeyRequest{
		Name: &ingestionKeyName,
		Type: &testIngestionKeyType,
	}

	ingestionCreateParams := ingestionKeysClient.NewCreateIngestionKeyParamsWithContext(ctx).WithBody(&ingestionKeyReq)
	ingestionCreateResp, err := apiClient.Ingestionkeys.CreateIngestionKey(ingestionCreateParams, nil)
	// Register for cleanup as early as possible - before any aborting assertion
	// (ingestion keys are deleted by name)
	if err == nil {
		tc.TrackIngestionKey(ingestionKeyName)
	}
	require.NoError(t, err, "Failed to create Ingestion Key: %v", err)
	require.NotNil(t, ingestionCreateResp, "Ingestion Key creation response is nil")
	require.NotNil(t, ingestionCreateResp.Payload, "Ingestion Key response payload is nil")
	require.NotEmpty(t, ingestionCreateResp.Payload.Key, "Created Ingestion Key is empty")
	require.Equal(t, ingestionKeyName, ingestionCreateResp.Payload.Name, "Ingestion Key name does not match expected value")

	createdKey = ingestionCreateResp.Payload.Key
	t.Logf("Successfully created Ingestion Key: %s", createdKey)

	ingestionListParams := ingestionKeysClient.NewListIngestionKeysParamsWithContext(ctx).WithBody(&models.ListIngestionKeysRequest{
		Name: ingestionKeyName,
	})

	var ingestionListResp *ingestionKeysClient.ListIngestionKeysOK
	timeout := time.Now().Add(10 * time.Second)
	for {
		var err error
		ingestionListResp, err = apiClient.Ingestionkeys.ListIngestionKeys(ingestionListParams, nil)
		require.NoError(t, err, "Failed to list Ingestion Keys: %v", err)
		require.NotNil(t, ingestionListResp.Payload, "Ingestion Key list response payload is nil")

		if len(ingestionListResp.Payload) > 0 || time.Now().After(timeout) {
			break
		}

		t.Logf("Waiting for Ingestion Key %s to be listed, retrying...", ingestionKeyName)
		time.Sleep(1 * time.Second)
	}
	require.Len(t, ingestionListResp.Payload, 1, "Expected exactly one Ingestion Key in the list, found %d", len(ingestionListResp.Payload))
	actualKey := ingestionListResp.Payload[0]
	require.Equal(t, ingestionKeyName, actualKey.Name, "Ingestion Key name does not match expected value")
	require.Equal(t, testIngestionKeyType, actualKey.Type, "Ingestion Key type does not match expected value")
	require.Equal(t, createdKey, actualKey.Key, "Ingestion Key value is empty")

	deleteParams := ingestionKeysClient.NewDeleteIngestionKeyParamsWithContext(ctx).WithBody(&models.DeleteIngestionKeyRequest{
		Name: &ingestionKeyName,
	})
	_, err = apiClient.Ingestionkeys.DeleteIngestionKey(deleteParams, nil)
	require.NoError(t, err, "Failed to delete Ingestion Key: %v", err)

	// The key was deleted by the test itself - no need for Cleanup to delete it
	tc.UntrackIngestionKey(ingestionKeyName)

	t.Logf("Successfully deleted Ingestion Key: %s", ingestionKeyName)

}
