package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/synthetics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPSyntheticsEndpoints(t *testing.T) {
	t.Parallel()
	testClient := NewTestClient(t)
	defer testClient.Cleanup()

	const defaultTimeout = 30 * time.Second

	var createdSyntheticID string
	syntheticName := "e2e-test-tcp-synthetic-" + uuid.New().String()

	t.Run("Create TCP Synthetic Test", func(t *testing.T) {
		createReq := &models.SyntheticTestCreateRequest{
			Name:     syntheticName,
			Version:  1,
			Enabled:  true,
			Interval: "5m",
			CheckConfig: &models.WorkerRequest{
				Kind: "tcp",
				Metadata: &models.Metadata{
					SyntheticName: syntheticName,
				},
				Request: &models.Request{
					TCP: &models.TCPRequest{
						Kind: "tcp",
						Host: "google.com",
						Port: 80,
					},
				},
				ExecutionPolicy: &models.ExecutionPolicy{
					Assertions: []*models.Assertion{
						{
							Source:   "tcp",
							Operator: "exists",
							Target:   "true",
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
		// Register for cleanup as early as possible - before any aborting assertion
		if err == nil && createResp != nil && createResp.Payload != nil && createResp.Payload.ID != "" {
			testClient.TrackSyntheticTest(createResp.Payload.ID)
		}
		require.NoError(t, err, "Failed to create TCP synthetic test")
		require.NotNil(t, createResp, "Create TCP synthetic test response should not be nil")
		require.NotNil(t, createResp.Payload, "Create TCP synthetic test response payload should not be nil")
		require.NotEmpty(t, createResp.Payload.ID, "Created TCP synthetic test ID should not be empty")

		createdSyntheticID = createResp.Payload.ID
		t.Logf("Created TCP synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Get TCP Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		getParams := synthetics.NewGetSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		getResp, err := testClient.Client.Synthetics.GetSyntheticTest(getParams, nil)
		require.NoError(t, err, "Failed to get TCP synthetic test")
		require.NotNil(t, getResp, "Get TCP synthetic test response should not be nil")
		require.NotNil(t, getResp.Payload, "Get TCP synthetic test response payload should not be nil")

		require.Equal(t, syntheticName, getResp.Payload.Name, "TCP synthetic test name mismatch")
		require.NotNil(t, getResp.Payload.CheckConfig, "TCP synthetic test CheckConfig should not be nil")
		require.Equal(t, models.WorkerRequestKind("tcp"), getResp.Payload.CheckConfig.Kind, "TCP synthetic test kind should be 'tcp'")
		require.NotNil(t, getResp.Payload.CheckConfig.Request, "TCP synthetic test Request should not be nil")
		require.NotNil(t, getResp.Payload.CheckConfig.Request.TCP, "TCP synthetic test TCP request should not be nil")
		require.Equal(t, "google.com", getResp.Payload.CheckConfig.Request.TCP.Host, "TCP synthetic test host mismatch")
		require.Equal(t, int64(80), getResp.Payload.CheckConfig.Request.TCP.Port, "TCP synthetic test port mismatch")

		t.Logf("Successfully retrieved TCP synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Delete TCP Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		deleteParams := synthetics.NewDeleteSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		_, err := testClient.Client.Synthetics.DeleteSyntheticTest(deleteParams, nil)
		require.NoError(t, err, "Failed to delete TCP synthetic test")

		// The synthetic test was deleted by the test itself - no need for Cleanup to delete it
		testClient.UntrackSyntheticTest(createdSyntheticID)

		t.Logf("Successfully deleted TCP synthetic test %s", createdSyntheticID)
	})
}

func TestSSLSyntheticsEndpoints(t *testing.T) {
	t.Parallel()
	testClient := NewTestClient(t)
	defer testClient.Cleanup()

	const defaultTimeout = 30 * time.Second

	var createdSyntheticID string
	syntheticName := "e2e-test-ssl-synthetic-" + uuid.New().String()

	t.Run("Create SSL Synthetic Test", func(t *testing.T) {
		createReq := &models.SyntheticTestCreateRequest{
			Name:     syntheticName,
			Version:  1,
			Enabled:  true,
			Interval: "5m",
			CheckConfig: &models.WorkerRequest{
				Kind: "ssl",
				Metadata: &models.Metadata{
					SyntheticName: syntheticName,
				},
				Request: &models.Request{
					Ssl: &models.SslRequest{
						Kind: "ssl",
						Host: "google.com",
						Port: 443,
					},
				},
				ExecutionPolicy: &models.ExecutionPolicy{
					Assertions: []*models.Assertion{
						{
							Source:   "ssl",
							Operator: "eq",
							Target:   "true",
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
		// Register for cleanup as early as possible - before any aborting assertion
		if err == nil && createResp != nil && createResp.Payload != nil && createResp.Payload.ID != "" {
			testClient.TrackSyntheticTest(createResp.Payload.ID)
		}
		require.NoError(t, err, "Failed to create SSL synthetic test")
		require.NotNil(t, createResp, "Create SSL synthetic test response should not be nil")
		require.NotNil(t, createResp.Payload, "Create SSL synthetic test response payload should not be nil")
		require.NotEmpty(t, createResp.Payload.ID, "Created SSL synthetic test ID should not be empty")

		createdSyntheticID = createResp.Payload.ID
		t.Logf("Created SSL synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Get SSL Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		getParams := synthetics.NewGetSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		getResp, err := testClient.Client.Synthetics.GetSyntheticTest(getParams, nil)
		require.NoError(t, err, "Failed to get SSL synthetic test")
		require.NotNil(t, getResp, "Get SSL synthetic test response should not be nil")
		require.NotNil(t, getResp.Payload, "Get SSL synthetic test response payload should not be nil")

		require.Equal(t, syntheticName, getResp.Payload.Name, "SSL synthetic test name mismatch")
		require.NotNil(t, getResp.Payload.CheckConfig, "SSL synthetic test CheckConfig should not be nil")
		require.Equal(t, models.WorkerRequestKind("ssl"), getResp.Payload.CheckConfig.Kind, "SSL synthetic test kind should be 'ssl'")
		require.NotNil(t, getResp.Payload.CheckConfig.Request, "SSL synthetic test Request should not be nil")
		require.NotNil(t, getResp.Payload.CheckConfig.Request.Ssl, "SSL synthetic test Ssl request should not be nil")
		require.Equal(t, "google.com", getResp.Payload.CheckConfig.Request.Ssl.Host, "SSL synthetic test host mismatch")
		require.Equal(t, int64(443), getResp.Payload.CheckConfig.Request.Ssl.Port, "SSL synthetic test port mismatch")

		t.Logf("Successfully retrieved SSL synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Delete SSL Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		deleteParams := synthetics.NewDeleteSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		_, err := testClient.Client.Synthetics.DeleteSyntheticTest(deleteParams, nil)
		require.NoError(t, err, "Failed to delete SSL synthetic test")

		// The synthetic test was deleted by the test itself - no need for Cleanup to delete it
		testClient.UntrackSyntheticTest(createdSyntheticID)

		t.Logf("Successfully deleted SSL synthetic test %s", createdSyntheticID)
	})
}

func TestSyntheticsEndpoints(t *testing.T) {
	t.Parallel()
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
		// Register for cleanup as early as possible - before any aborting assertion
		if err == nil && createResp != nil && createResp.Payload != nil && createResp.Payload.ID != "" {
			testClient.TrackSyntheticTest(createResp.Payload.ID)
		}
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

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			listParams := synthetics.NewListSyntheticTestsParams().
				WithContext(testClient.BaseCtx).
				WithTimeout(defaultTimeout)

			listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
			if !assert.NoError(collect, err, "ListSyntheticTests failed") {
				return
			}
			if !assert.NotNil(collect, listResp, "ListSyntheticTests response is nil") || !assert.NotNil(collect, listResp.Payload, "ListSyntheticTests payload is nil") {
				return
			}
			var found bool
			for _, item := range listResp.Payload.Synthetics {
				if item.ID == createdSyntheticID {
					found = true
					assert.Equal(collect, syntheticName, item.Name, "Synthetic test name mismatch")
					assert.Equal(collect, "http", item.SyntheticType, "Synthetic test type should be 'http'")
					assert.NotEmpty(collect, item.Status, "Synthetic test status should not be empty")
					break
				}
			}
			assert.True(collect, found, "Created synthetic test %s not found in list response", createdSyntheticID)
		}, 2*time.Minute, 2*time.Second, "Created synthetic test %s not found in list response with expected name and type", createdSyntheticID)
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

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			listParams := synthetics.NewListSyntheticTestsParams().
				WithContext(testClient.BaseCtx).
				WithTimeout(defaultTimeout)

			listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
			if !assert.NoError(collect, err, "ListSyntheticTests failed") {
				return
			}
			if !assert.NotNil(collect, listResp, "ListSyntheticTests response is nil") || !assert.NotNil(collect, listResp.Payload, "ListSyntheticTests payload is nil") {
				return
			}
			var found bool
			for _, item := range listResp.Payload.Synthetics {
				if item.ID == createdSyntheticID && item.Name == updatedName {
					found = true
					break
				}
			}
			assert.True(collect, found, "Updated synthetic test %s not found in list with expected name %q", createdSyntheticID, updatedName)
		}, 2*time.Minute, 2*time.Second, "Updated synthetic test %s not found in list with expected name", createdSyntheticID)

		getParams := synthetics.NewGetSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		getResp, err := testClient.Client.Synthetics.GetSyntheticTest(getParams, nil)
		require.NoError(t, err, "Failed to get updated synthetic test")
		require.NotNil(t, getResp.Payload, "Get updated synthetic test payload should not be nil")

		require.Equal(t, updatedName, getResp.Payload.Name, "Updated synthetic test name mismatch")
		require.Equal(t, "10m", getResp.Payload.Interval, "Updated synthetic test interval mismatch")
		require.Equal(t, models.WorkerRequestKind("http"), getResp.Payload.CheckConfig.Kind, "Updated synthetic test kind mismatch")
		require.NotNil(t, getResp.Payload.CheckConfig.Request.HTTP, "Updated synthetic test HTTP request should not be nil")
		require.Equal(t, "GET", getResp.Payload.CheckConfig.Request.HTTP.Method, "Updated synthetic test HTTP method mismatch")
		require.Equal(t, "https://httpbin.org/get", getResp.Payload.CheckConfig.Request.HTTP.URL, "Updated synthetic test HTTP URL mismatch")
		require.Equal(t, "30s", getResp.Payload.CheckConfig.Request.HTTP.Timeout, "Updated synthetic test HTTP timeout mismatch")
		require.Equal(t, models.HTTPRequestKind("http"), getResp.Payload.CheckConfig.Request.HTTP.Kind, "Updated synthetic test HTTP kind mismatch")

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

		// The synthetic test was deleted by the test itself - no need for Cleanup to delete it
		testClient.UntrackSyntheticTest(createdSyntheticID)

		t.Logf("✓ Successfully deleted synthetic test %s", createdSyntheticID)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			listParams := synthetics.NewListSyntheticTestsParams().
				WithContext(testClient.BaseCtx).
				WithTimeout(defaultTimeout)

			listResp, err := testClient.Client.Synthetics.ListSyntheticTests(listParams, nil)
			if !assert.NoError(collect, err, "ListSyntheticTests failed") {
				return
			}
			if !assert.NotNil(collect, listResp, "ListSyntheticTests response is nil") || !assert.NotNil(collect, listResp.Payload, "ListSyntheticTests payload is nil") {
				return
			}
			for _, item := range listResp.Payload.Synthetics {
				if item.ID == createdSyntheticID {
					assert.Fail(collect, "Deleted synthetic test %s still found in list response", createdSyntheticID)
					return
				}
			}
		}, 2*time.Minute, 2*time.Second, "Deleted synthetic test %s should not be found in list response", createdSyntheticID)
		t.Logf("Verified synthetic test %s is no longer in the list after deletion", createdSyntheticID)
	})
}

func TestDNSSyntheticsEndpoints(t *testing.T) {
	t.Parallel()
	testClient := NewTestClient(t)
	defer testClient.Cleanup()

	const defaultTimeout = 30 * time.Second

	var createdSyntheticID string
	syntheticName := "e2e-test-dns-synthetic-" + uuid.New().String()

	t.Run("Create DNS Synthetic Test", func(t *testing.T) {
		createReq := &models.SyntheticTestCreateRequest{
			Name:     syntheticName,
			Version:  1,
			Enabled:  true,
			Interval: "5m",
			CheckConfig: &models.WorkerRequest{
				Kind: "dns",
				Metadata: &models.Metadata{
					SyntheticName: syntheticName,
				},
				Request: &models.Request{
					DNS: &models.DNSRequest{
						Kind:       "dns",
						Domain:     "google.com",
						Resolver:   "8.8.8.8",
						Port:       53,
						RecordType: "A",
						Timeout:    "30s",
					},
				},
				ExecutionPolicy: &models.ExecutionPolicy{
					Assertions: []*models.Assertion{
						{
							Source:   "dnsAnswer",
							Operator: "exists",
							Target:   "true",
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
		// Register for cleanup as early as possible - before any aborting assertion
		if err == nil && createResp != nil && createResp.Payload != nil && createResp.Payload.ID != "" {
			testClient.TrackSyntheticTest(createResp.Payload.ID)
		}
		require.NoError(t, err, "Failed to create DNS synthetic test")
		require.NotNil(t, createResp, "Create DNS synthetic test response should not be nil")
		require.NotNil(t, createResp.Payload, "Create DNS synthetic test response payload should not be nil")
		require.NotEmpty(t, createResp.Payload.ID, "Created DNS synthetic test ID should not be empty")

		createdSyntheticID = createResp.Payload.ID
		t.Logf("Created DNS synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Get DNS Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		getParams := synthetics.NewGetSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		getResp, err := testClient.Client.Synthetics.GetSyntheticTest(getParams, nil)
		require.NoError(t, err, "Failed to get DNS synthetic test")
		require.NotNil(t, getResp, "Get DNS synthetic test response should not be nil")
		require.NotNil(t, getResp.Payload, "Get DNS synthetic test response payload should not be nil")

		require.Equal(t, syntheticName, getResp.Payload.Name, "DNS synthetic test name mismatch")
		require.NotNil(t, getResp.Payload.CheckConfig, "DNS synthetic test CheckConfig should not be nil")
		require.Equal(t, models.WorkerRequestKind("dns"), getResp.Payload.CheckConfig.Kind, "DNS synthetic test kind should be 'dns'")
		require.NotNil(t, getResp.Payload.CheckConfig.Request, "DNS synthetic test Request should not be nil")
		require.NotNil(t, getResp.Payload.CheckConfig.Request.DNS, "DNS synthetic test DNS request should not be nil")
		require.Equal(t, "google.com", getResp.Payload.CheckConfig.Request.DNS.Domain, "DNS synthetic test domain mismatch")
		require.Equal(t, "8.8.8.8", getResp.Payload.CheckConfig.Request.DNS.Resolver, "DNS synthetic test resolver mismatch")
		require.Equal(t, int64(53), getResp.Payload.CheckConfig.Request.DNS.Port, "DNS synthetic test port mismatch")
		require.Equal(t, models.DNSRequestRecordType("A"), getResp.Payload.CheckConfig.Request.DNS.RecordType, "DNS synthetic test record type mismatch")

		t.Logf("Successfully retrieved DNS synthetic test with ID: %s", createdSyntheticID)
	})

	t.Run("Delete DNS Synthetic Test", func(t *testing.T) {
		if createdSyntheticID == "" {
			t.Skip("Skipping because create failed or didn't run")
		}

		deleteParams := synthetics.NewDeleteSyntheticTestParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdSyntheticID)

		_, err := testClient.Client.Synthetics.DeleteSyntheticTest(deleteParams, nil)
		require.NoError(t, err, "Failed to delete DNS synthetic test")

		// The synthetic test was deleted by the test itself - no need for Cleanup to delete it
		testClient.UntrackSyntheticTest(createdSyntheticID)

		t.Logf("Successfully deleted DNS synthetic test %s", createdSyntheticID)
	})
}

// TestHTTPSyntheticFollowRedirectsFalse verifies that setting followRedirects=false
// via the SDK is correctly serialized and persisted. This catches the go-swagger bug
// where bool fields with omitempty silently drop false values
// (github.com/go-swagger/go-swagger/issues/1601).
func TestHTTPSyntheticFollowRedirectsFalse(t *testing.T) {
	httpReq := &models.HTTPRequest{
		Kind:            "http",
		URL:             "https://httpbin.org/redirect/1",
		Method:          "GET",
		Timeout:         "30s",
		FollowRedirects: swag.Bool(false),
		AllowInsecure:   swag.Bool(false),
	}

	data, err := json.Marshal(httpReq)
	require.NoError(t, err)
	require.Contains(t, string(data), `"followRedirects":false`,
		"SDK must preserve followRedirects=false in JSON serialization, got: %s", string(data))
	require.Contains(t, string(data), `"allowInsecure":false`,
		"SDK must preserve allowInsecure=false in JSON serialization, got: %s", string(data))
}
