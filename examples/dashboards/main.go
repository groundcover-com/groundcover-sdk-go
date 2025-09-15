package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/dashboards"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func main() {
	fmt.Println("groundcover Go SDK - Dashboards Example")
	fmt.Println("=======================================")

	// Create client using environment variables
	// Required: GC_API_KEY, GC_BACKEND_ID, GC_BASE_URL
	gcClient, err := groundcover.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("✓ Client created successfully")

	ctx := context.Background()

	// Run the complete dashboard lifecycle example
	err = runDashboardLifecycle(ctx, gcClient)
	if err != nil {
		log.Fatalf("Dashboard lifecycle example failed: %v", err)
	}

	fmt.Println("\n✓ Dashboards example completed successfully")
}

func runDashboardLifecycle(ctx context.Context, gcClient *client.GroundcoverAPI) error {
	dashboardName := "SDK Example Dashboard " + uuid.New().String()[:8]
	description := "Dashboard created by SDK example"

	// Step 1: Create Dashboard
	fmt.Printf("\n1. Creating dashboard: %s\n", dashboardName)
	dashboard, err := createDashboard(ctx, gcClient, dashboardName, description)
	if err != nil {
		return fmt.Errorf("failed to create dashboard: %w", err)
	}
	fmt.Printf("   ✓ Created dashboard with ID: %s\n", dashboard.UUID)

	// Step 2: Get Dashboard
	fmt.Printf("\n2. Getting dashboard by ID: %s\n", dashboard.UUID)
	retrievedDashboard, err := getDashboard(ctx, gcClient, dashboard.UUID)
	if err != nil {
		return fmt.Errorf("failed to get dashboard: %w", err)
	}
	fmt.Printf("   ✓ Retrieved dashboard: %s (Status: %s)\n", retrievedDashboard.Name, retrievedDashboard.Status)

	// Step 3: List Dashboards
	fmt.Println("\n3. Listing all dashboards")
	dashboards, err := listDashboards(ctx, gcClient)
	if err != nil {
		return fmt.Errorf("failed to list dashboards: %w", err)
	}
	fmt.Printf("   ✓ Found %d total dashboards\n", len(dashboards))

	// Find our dashboard in the list
	found := false
	for _, d := range dashboards {
		if d.UUID == dashboard.UUID {
			found = true
			fmt.Printf("   ✓ Confirmed our dashboard appears in the list\n")
			break
		}
	}
	if !found {
		return fmt.Errorf("created dashboard not found in list")
	}

	// Step 4: Update Dashboard
	updatedName := dashboardName + " (Updated)"
	updatedDescription := "Updated " + description
	fmt.Printf("\n4. Updating dashboard to: %s\n", updatedName)

	updatedDashboard, err := updateDashboard(ctx, gcClient, dashboard.UUID, updatedName, updatedDescription, dashboard.RevisionNumber)
	if err != nil {
		return fmt.Errorf("failed to update dashboard: %w", err)
	}
	fmt.Printf("   ✓ Updated dashboard: %s (Revision: %d -> %d)\n",
		updatedDashboard.Name, dashboard.RevisionNumber, updatedDashboard.RevisionNumber)

	// Step 5: Archive Dashboard
	fmt.Println("\n5. Archiving dashboard")
	archivedDashboard, err := archiveDashboard(ctx, gcClient, updatedDashboard.UUID, updatedDashboard.RevisionNumber)
	if err != nil {
		return fmt.Errorf("failed to archive dashboard: %w", err)
	}
	fmt.Printf("   ✓ Archived dashboard (Status: %s -> %s)\n", updatedDashboard.Status, archivedDashboard.Status)

	// Step 6: Restore Dashboard
	fmt.Println("\n6. Restoring dashboard")
	restoredDashboard, err := restoreDashboard(ctx, gcClient, archivedDashboard.UUID, archivedDashboard.RevisionNumber)
	if err != nil {
		return fmt.Errorf("failed to restore dashboard: %w", err)
	}
	fmt.Printf("   ✓ Restored dashboard (Status: %s -> %s)\n", archivedDashboard.Status, restoredDashboard.Status)

	// Step 7: Delete Dashboard
	fmt.Printf("\n7. Deleting dashboard: %s\n", restoredDashboard.UUID)
	err = deleteDashboard(ctx, gcClient, restoredDashboard.UUID)
	if err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}
	fmt.Printf("   ✓ Deleted dashboard successfully\n")

	// Step 8: Verify Deletion
	fmt.Println("\n8. Verifying dashboard deletion")
	_, err = getDashboard(ctx, gcClient, restoredDashboard.UUID)
	if err == nil {
		return fmt.Errorf("dashboard should have been deleted but was still found")
	}
	fmt.Printf("   ✓ Confirmed dashboard is deleted (SDK call returned error)\n")

	return nil
}

func createDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, name, description string) (*models.View, error) {
	preset := getExamplePreset()

	createReq := &models.CreateDashboardRequest{
		Name:          name,
		Description:   description,
		Preset:        preset,
		IsProvisioned: false,
	}

	// Create dashboard using SDK client
	createParams := dashboards.NewCreateDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithBody(createReq)

	_, err := gcClient.Dashboards.CreateDashboard(createParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK create dashboard call failed: %w", err)
	}

	// Since CreateDashboard response is empty, find the created dashboard by listing
	dashboardsList, err := listDashboards(ctx, gcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards after create: %w", err)
	}

	// Find our newly created dashboard
	for _, dashboard := range dashboardsList {
		if dashboard.Name == name {
			return dashboard, nil
		}
	}

	return nil, fmt.Errorf("created dashboard not found in list")
}

func getDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, id string) (*models.View, error) {
	// Get dashboard using SDK client
	getParams := dashboards.NewGetDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithID(id)

	_, err := gcClient.Dashboards.GetDashboard(getParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK get dashboard call failed: %w", err)
	}

	// Since GetDashboard response is empty, find the dashboard by listing
	dashboardsList, err := listDashboards(ctx, gcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards after get: %w", err)
	}

	// Find our dashboard in the list
	for _, dashboard := range dashboardsList {
		if dashboard.UUID == id {
			return dashboard, nil
		}
	}

	return nil, fmt.Errorf("dashboard with ID %s not found", id)
}

func listDashboards(ctx context.Context, gcClient *client.GroundcoverAPI) ([]*models.View, error) {
	// List dashboards using SDK client
	listParams := dashboards.NewGetDashboardsParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second)

	listResp, err := gcClient.Dashboards.GetDashboards(listParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK list dashboards call failed: %w", err)
	}

	if listResp.Payload == nil {
		return []*models.View{}, nil
	}

	return listResp.Payload, nil
}

func updateDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, id, name, description string, currentRevision int32) (*models.View, error) {
	preset := getUpdatedExamplePreset()

	updateReq := &models.UpdateDashboardRequest{
		Name:            name,
		Description:     description,
		Preset:          preset,
		CurrentRevision: currentRevision,
		Override:        false,
		IsProvisioned:   false,
	}

	// Update dashboard using SDK client
	updateParams := dashboards.NewUpdateDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithID(id).
		WithBody(updateReq)

	_, err := gcClient.Dashboards.UpdateDashboard(updateParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK update dashboard call failed: %w", err)
	}

	// Since UpdateDashboard response is empty, find the updated dashboard by listing
	dashboardsList, err := listDashboards(ctx, gcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards after update: %w", err)
	}

	// Find our updated dashboard in the list
	for _, dashboard := range dashboardsList {
		if dashboard.UUID == id {
			return dashboard, nil
		}
	}

	return nil, fmt.Errorf("updated dashboard with ID %s not found", id)
}

func archiveDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, id string, currentRevision int32) (*models.View, error) {
	// Archive dashboard using SDK client
	archiveParams := dashboards.NewArchiveDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithID(id).
		WithCurrentRevision(currentRevision)

	_, err := gcClient.Dashboards.ArchiveDashboard(archiveParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK archive dashboard call failed: %w", err)
	}

	// Since ArchiveDashboard response is empty, find the archived dashboard by listing
	dashboardsList, err := listDashboards(ctx, gcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards after archive: %w", err)
	}

	// Find our archived dashboard in the list
	for _, dashboard := range dashboardsList {
		if dashboard.UUID == id {
			return dashboard, nil
		}
	}

	return nil, fmt.Errorf("archived dashboard with ID %s not found", id)
}

func restoreDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, id string, currentRevision int32) (*models.View, error) {
	// Restore dashboard using SDK client
	restoreParams := dashboards.NewRestoreDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithID(id).
		WithCurrentRevision(currentRevision)

	_, err := gcClient.Dashboards.RestoreDashboard(restoreParams, nil)
	if err != nil {
		return nil, fmt.Errorf("SDK restore dashboard call failed: %w", err)
	}

	// Since RestoreDashboard response is empty, find the restored dashboard by listing
	dashboardsList, err := listDashboards(ctx, gcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards after restore: %w", err)
	}

	// Find our restored dashboard in the list
	for _, dashboard := range dashboardsList {
		if dashboard.UUID == id {
			return dashboard, nil
		}
	}

	return nil, fmt.Errorf("restored dashboard with ID %s not found", id)
}

func deleteDashboard(ctx context.Context, gcClient *client.GroundcoverAPI, id string) error {
	// Delete dashboard using SDK client
	deleteParams := dashboards.NewDeleteDashboardParams().
		WithContext(ctx).
		WithTimeout(30 * time.Second).
		WithID(id)

	_, err := gcClient.Dashboards.DeleteDashboard(deleteParams, nil)
	if err != nil {
		return fmt.Errorf("SDK delete dashboard call failed: %w", err)
	}

	return nil
}

func getExamplePreset() string {
	return `{
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
      "name": "Node CPU Usage",
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
      "html": "<p>SDK Example Dashboard Widget</p>"
    }
  ],
  "variables": {},
  "schemaVersion": 3
}`
}

func getUpdatedExamplePreset() string {
	return `{
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
      "name": "Updated Node CPU Usage",
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
      "html": "<p>Updated SDK Example Dashboard Widget</p>"
    }
  ],
  "variables": {},
  "schemaVersion": 3
}`
}
