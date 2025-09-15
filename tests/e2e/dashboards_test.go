package e2e

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/dashboards"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestDashboardsEndpoints(t *testing.T) {
	testClient := NewTestClient(t)
	defer testClient.Cleanup()

	const defaultTimeout = 30 * time.Second

	var createdDashboardID string
	var createdDashboardName string
	var createdDashboard *models.View

	t.Run("Create Dashboard", func(t *testing.T) {
		dashboardName := "e2e-test-dashboard-" + uuid.New().String()
		description := "Dashboard created during E2E testing"

		// Dashboard preset for testing - needs to be a valid dashboard preset JSON
		preset := `{
  "duration": "Last 30 minutes",
  "layout": [
    {
      "id": "A",
      "x": 0,
      "y": 0,
      "w": 4,
      "h": 3,
      "minH": 2
    },
    {
      "id": "B",
      "x": 0,
      "y": 3,
      "w": 4,
      "h": 3,
      "minH": 1
    }
  ],
  "widgets": [
    {
      "id": "A",
      "type": "widget",
      "name": "avg(groundcover_node_rt_disk_space_used_percent{})",
      "queries": [
        {
          "id": "A",
          "expr": "avg(groundcover_node_rt_disk_space_used_percent{})",
          "dataType": "metrics",
          "step": null,
          "editorMode": "builder"
        }
      ],
      "visualizationConfig": {
        "type": "time-series"
      }
    },
    {
      "id": "B",
      "type": "text",
      "html": "<p>SDK Test Widget</p>"
    }
  ],
  "variables": {},
  "schemaVersion": 3
}`

		createReq := &models.CreateDashboardRequest{
			Name:          dashboardName,
			Description:   description,
			Preset:        preset,
			IsProvisioned: false,
		}

		// Create dashboard using SDK client
		createParams := dashboards.NewCreateDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithBody(createReq)

		createResp, err := testClient.Client.Dashboards.CreateDashboard(createParams, nil)
		require.NoError(t, err, "Failed to create dashboard")
		require.NotNil(t, createResp, "Create dashboard response should not be nil")

		t.Logf("✓ SDK Create Dashboard call succeeded")

		// Since CreateDashboard response doesn't contain dashboard data,
		// use GetDashboards to find the created dashboard
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards after create")
		require.NotNil(t, listResp.Payload, "List response payload should not be nil")

		// Find our created dashboard in the models.View payload
		var dashboardView *models.View
		for _, view := range listResp.Payload {
			if view.Name == dashboardName {
				dashboardView = view
				break
			}
		}

		require.NotNil(t, dashboardView, "Created dashboard not found in list")
		require.NotEmpty(t, dashboardView.UUID, "Created dashboard UUID should not be empty")
		require.Equal(t, dashboardName, dashboardView.Name, "Dashboard name mismatch")
		require.Equal(t, description, dashboardView.Description, "Dashboard description mismatch")
		require.Equal(t, "explore", dashboardView.Type, "Dashboard view type should be 'explore'")
		require.Equal(t, "active", dashboardView.Status, "Dashboard status should be 'active'")

		// Save for other tests - convert to our Dashboard struct for convenience
		createdDashboardID = dashboardView.UUID
		createdDashboardName = dashboardName
		createdDashboard = &models.View{
			UUID:           dashboardView.UUID,
			Name:           dashboardView.Name,
			Description:    dashboardView.Description,
			Owner:          dashboardView.Owner,
			Preset:         dashboardView.Preset,
			Type:           dashboardView.Type,
			RevisionNumber: dashboardView.RevisionNumber,
			Status:         dashboardView.Status,
			IsProvisioned:  dashboardView.IsProvisioned,
		}

		t.Logf("Created dashboard with ID: %s", createdDashboardID)
	})

	t.Run("Get Dashboard", func(t *testing.T) {
		if createdDashboardID == "" {
			t.Skip("Skipping Get Dashboard test because create failed or didn't run")
		}

		// Get dashboard using SDK client
		getParams := dashboards.NewGetDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdDashboardID)

		getResp, err := testClient.Client.Dashboards.GetDashboard(getParams, nil)
		require.NoError(t, err, "Failed to get dashboard")
		require.NotNil(t, getResp, "Get dashboard response should not be nil")

		t.Logf("✓ SDK Get Dashboard call succeeded")

		// Note: GetDashboard response is empty, but the call succeeded which means the dashboard exists
		// We can verify by listing dashboards and finding our dashboard
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards for verification")

		// Find our dashboard in the list to verify it exists
		found := false
		for _, view := range listResp.Payload {
			if view.UUID == createdDashboardID {
				require.Equal(t, createdDashboardName, view.Name, "Get dashboard name mismatch")
				require.Equal(t, "explore", view.Type, "Dashboard view type should be 'explore'")
				found = true
				break
			}
		}
		require.True(t, found, "Dashboard not found in list after Get call")

		t.Logf("Successfully retrieved dashboard with ID: %s", createdDashboardID)
	})

	t.Run("List Dashboards", func(t *testing.T) {
		if createdDashboardID == "" {
			t.Skip("Skipping List Dashboards test because create failed or didn't run")
		}

		// List dashboards using SDK client
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards")
		require.NotNil(t, listResp, "List dashboards response should not be nil")
		require.NotNil(t, listResp.Payload, "List dashboards response payload should not be nil")

		t.Logf("✓ SDK Get Dashboards call succeeded with %d dashboards", len(listResp.Payload))

		// Check if the created dashboard is in the models.View payload
		found := false
		for _, view := range listResp.Payload {
			if view.UUID == createdDashboardID {
				require.Equal(t, createdDashboardName, view.Name, "List dashboard name mismatch")
				found = true
				t.Logf("Found created dashboard %s in the list", createdDashboardID)
				break
			}
		}
		require.True(t, found, "Created dashboard %s not found in list response", createdDashboardID)
	})

	t.Run("Update Dashboard", func(t *testing.T) {
		if createdDashboardID == "" || createdDashboard == nil {
			t.Skip("Skipping Update Dashboard test because create failed or didn't run")
		}

		updatedName := createdDashboardName + "-updated"
		updatedDescription := "Updated dashboard description"

		// Updated dashboard preset for testing
		preset := `{
  "duration": "Last 1 hour",
  "layout": [
    {
      "id": "A",
      "x": 0,
      "y": 0,
      "w": 6,
      "h": 4,
      "minH": 2
    },
    {
      "id": "B",
      "x": 0,
      "y": 4,
      "w": 6,
      "h": 2,
      "minH": 1
    }
  ],
  "widgets": [
    {
      "id": "A",
      "type": "widget",
      "name": "Updated: avg(groundcover_node_rt_disk_space_used_percent{})",
      "queries": [
        {
          "id": "A",
          "expr": "avg(groundcover_node_rt_disk_space_used_percent{})",
          "dataType": "metrics",
          "step": null,
          "editorMode": "builder"
        }
      ],
      "visualizationConfig": {
        "type": "time-series"
      }
    },
    {
      "id": "B",
      "type": "text",
      "html": "<p>Updated SDK Test Widget</p>"
    }
  ],
  "variables": {},
  "schemaVersion": 3
}`

		updateReq := &models.UpdateDashboardRequest{
			Name:            updatedName,
			Description:     updatedDescription,
			Preset:          preset,
			CurrentRevision: createdDashboard.RevisionNumber,
			Override:        false,
			IsProvisioned:   false,
		}

		// Update dashboard using SDK client
		updateParams := dashboards.NewUpdateDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdDashboardID).
			WithBody(updateReq)

		updateResp, err := testClient.Client.Dashboards.UpdateDashboard(updateParams, nil)
		require.NoError(t, err, "Failed to update dashboard")
		require.NotNil(t, updateResp, "Update dashboard response should not be nil")

		t.Logf("✓ SDK Update Dashboard call succeeded")

		// Verify update by listing dashboards and finding our updated dashboard
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards after update")

		// Find our updated dashboard in the list
		var updatedDashboardView *models.View
		for _, view := range listResp.Payload {
			if view.UUID == createdDashboardID {
				updatedDashboardView = view
				break
			}
		}

		require.NotNil(t, updatedDashboardView, "Updated dashboard not found in list")
		require.Equal(t, createdDashboardID, updatedDashboardView.UUID, "Update dashboard ID mismatch")
		require.Equal(t, updatedName, updatedDashboardView.Name, "Update dashboard name mismatch")
		require.Equal(t, updatedDescription, updatedDashboardView.Description, "Update dashboard description mismatch")
		require.Greater(t, updatedDashboardView.RevisionNumber, createdDashboard.RevisionNumber, "Revision number should be incremented")

		// Update our reference for potential future tests
		createdDashboard = &models.View{
			UUID:           updatedDashboardView.UUID,
			Name:           updatedDashboardView.Name,
			Description:    updatedDashboardView.Description,
			Owner:          updatedDashboardView.Owner,
			Preset:         updatedDashboardView.Preset,
			Type:           updatedDashboardView.Type,
			RevisionNumber: updatedDashboardView.RevisionNumber,
			Status:         updatedDashboardView.Status,
			IsProvisioned:  updatedDashboardView.IsProvisioned,
		}
		createdDashboardName = updatedName

		t.Logf("Successfully updated dashboard with ID: %s", createdDashboardID)
	})

	t.Run("Archive Dashboard", func(t *testing.T) {
		if createdDashboardID == "" || createdDashboard == nil {
			t.Skip("Skipping Archive Dashboard test because create failed or didn't run")
		}

		// Archive dashboard using SDK client
		archiveParams := dashboards.NewArchiveDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdDashboardID).
			WithCurrentRevision(createdDashboard.RevisionNumber)

		archiveResp, err := testClient.Client.Dashboards.ArchiveDashboard(archiveParams, nil)
		require.NoError(t, err, "Failed to archive dashboard")
		require.NotNil(t, archiveResp, "Archive dashboard response should not be nil")

		t.Logf("✓ SDK Archive Dashboard call succeeded")

		// Verify archive by listing dashboards and checking status
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards after archive")

		// Find our archived dashboard in the list
		var archivedDashboardView *models.View
		for _, view := range listResp.Payload {
			if view.UUID == createdDashboardID {
				archivedDashboardView = view
				break
			}
		}

		require.NotNil(t, archivedDashboardView, "Archived dashboard not found in list")
		require.Equal(t, createdDashboardID, archivedDashboardView.UUID, "Archive dashboard ID mismatch")
		require.Equal(t, "archived", archivedDashboardView.Status, "Dashboard status should be 'archived'")

		// Update our reference for potential future tests
		createdDashboard = &models.View{
			UUID:           archivedDashboardView.UUID,
			Name:           archivedDashboardView.Name,
			Description:    archivedDashboardView.Description,
			Owner:          archivedDashboardView.Owner,
			Preset:         archivedDashboardView.Preset,
			Type:           archivedDashboardView.Type,
			RevisionNumber: archivedDashboardView.RevisionNumber,
			Status:         archivedDashboardView.Status,
			IsProvisioned:  archivedDashboardView.IsProvisioned,
		}

		t.Logf("Successfully archived dashboard with ID: %s", createdDashboardID)
	})

	t.Run("Restore Dashboard", func(t *testing.T) {
		if createdDashboardID == "" || createdDashboard == nil {
			t.Skip("Skipping Restore Dashboard test because archive failed or didn't run")
		}

		// Restore dashboard using SDK client
		restoreParams := dashboards.NewRestoreDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdDashboardID).
			WithCurrentRevision(createdDashboard.RevisionNumber)

		restoreResp, err := testClient.Client.Dashboards.RestoreDashboard(restoreParams, nil)
		require.NoError(t, err, "Failed to restore dashboard")
		require.NotNil(t, restoreResp, "Restore dashboard response should not be nil")

		t.Logf("✓ SDK Restore Dashboard call succeeded")

		// Verify restore by listing dashboards and checking status
		listParams := dashboards.NewGetDashboardsParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout)

		listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
		require.NoError(t, err, "Failed to list dashboards after restore")

		// Find our restored dashboard in the list
		var restoredDashboardView *models.View
		for _, view := range listResp.Payload {
			if view.UUID == createdDashboardID {
				restoredDashboardView = view
				break
			}
		}

		require.NotNil(t, restoredDashboardView, "Restored dashboard not found in list")
		require.Equal(t, createdDashboardID, restoredDashboardView.UUID, "Restore dashboard ID mismatch")
		require.Equal(t, "active", restoredDashboardView.Status, "Dashboard status should be 'active' after restore")

		// Update our reference for potential future tests
		createdDashboard = &models.View{
			UUID:           restoredDashboardView.UUID,
			Name:           restoredDashboardView.Name,
			Description:    restoredDashboardView.Description,
			Owner:          restoredDashboardView.Owner,
			Preset:         restoredDashboardView.Preset,
			Type:           restoredDashboardView.Type,
			RevisionNumber: restoredDashboardView.RevisionNumber,
			Status:         restoredDashboardView.Status,
			IsProvisioned:  restoredDashboardView.IsProvisioned,
		}

		t.Logf("Successfully restored dashboard with ID: %s", createdDashboardID)
	})

	t.Run("Delete Dashboard", func(t *testing.T) {
		if createdDashboardID == "" {
			t.Skip("Skipping Delete Dashboard test because create failed or didn't run")
		}

		// Delete dashboard using SDK client
		deleteParams := dashboards.NewDeleteDashboardParams().
			WithContext(testClient.BaseCtx).
			WithTimeout(defaultTimeout).
			WithID(createdDashboardID)

		deleteResp, err := testClient.Client.Dashboards.DeleteDashboard(deleteParams, nil)
		require.NoError(t, err, "Failed to delete dashboard")
		require.NotNil(t, deleteResp, "Delete dashboard response should not be nil")

		t.Logf("✓ SDK Delete Dashboard call succeeded")
		t.Logf("Successfully deleted dashboard with ID: %s", createdDashboardID)

		// Verify dashboard is deleted by trying to get it via SDK
		t.Run("Verify Dashboard Deleted", func(t *testing.T) {
			getParams := dashboards.NewGetDashboardParams().
				WithContext(testClient.BaseCtx).
				WithTimeout(defaultTimeout).
				WithID(createdDashboardID)

			_, err := testClient.Client.Dashboards.GetDashboard(getParams, nil)
			require.Error(t, err, "Expected error when getting deleted dashboard")

			// Also verify dashboard is not in list
			listParams := dashboards.NewGetDashboardsParams().
				WithContext(testClient.BaseCtx).
				WithTimeout(defaultTimeout)

			listResp, err := testClient.Client.Dashboards.GetDashboards(listParams, nil)
			require.NoError(t, err, "Failed to list dashboards after delete")

			// Ensure deleted dashboard is not in the list
			found := false
			for _, view := range listResp.Payload {
				if view.UUID == createdDashboardID {
					found = true
					break
				}
			}
			require.False(t, found, "Deleted dashboard should not be in list")

			t.Logf("Confirmed dashboard %s is deleted", createdDashboardID)
		})
	})
}
