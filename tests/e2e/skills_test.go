package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	agentclient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/agent"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
	"github.com/stretchr/testify/require"
)

const (
	testSkillNamePrefix            = "sdk-e2e-test-skill-"
	testSkillProvisioningUserAgent = "terraform-provider-groundcover/sdk-e2e"
)

func TestSkillsEndpoints(t *testing.T) {
	tc := NewTestClient(t)
	defer tc.Cleanup()
	ctx, apiClient := tc.BaseCtx, tc.Client

	createdName := fmt.Sprintf("%s%d", testSkillNamePrefix, time.Now().UnixNano())
	updatedName := createdName + "-updated"
	createdWhenToUse := "Use this skill during SDK E2E create checks"
	updatedWhenToUse := "Use this skill during SDK E2E update checks"
	description := "Created by SDK E2E test"
	updatedDescription := "Updated by SDK E2E test"
	isOrganizational := true

	var createdSkillID string

	t.Run("Create Skill", func(t *testing.T) {
		createParams := agentclient.NewAgentCreateSkillParams().
			WithContext(ctx).
			WithTimeout(defaultTimeout).
			WithBody(newSkillCreateRequest(createdName, createdWhenToUse, description, isOrganizational))

		createResp, err := apiClient.Agent.AgentCreateSkill(
			createParams,
			nil,
			transport.WithHeadersOverride(http.Header{"User-Agent": {testSkillProvisioningUserAgent}}),
		)
		if err == nil && createResp != nil && createResp.Payload != nil && createResp.Payload.Skill != nil &&
			createResp.Payload.Skill.ID != nil && *createResp.Payload.Skill.ID != "" {
			tc.TrackAgentSkill(*createResp.Payload.Skill.ID)
		}
		require.NoError(t, err, "Create skill failed")
		require.NotNil(t, createResp, "Create skill response should not be nil")
		require.NotNil(t, createResp.Payload, "Create skill response payload should not be nil")
		require.NotNil(t, createResp.Payload.Skill, "Create skill response skill should not be nil")
		require.NotNil(t, createResp.Payload.Skill.ID, "Created skill ID should not be nil")
		require.NotEmpty(t, *createResp.Payload.Skill.ID, "Created skill ID should not be empty")
		require.NotNil(t, createResp.Payload.Skill.IsOrganizational, "Created skill is_organizational should not be nil")
		require.True(t, *createResp.Payload.Skill.IsOrganizational, "Created skill should be organizational")
		requireSkillProvisioned(t, createResp.Payload.Skill.IsProvisioned, "created skill")

		createdSkillID = *createResp.Payload.Skill.ID
		t.Logf("Created agent skill with ID: %s", createdSkillID)
	})

	require.NotEmpty(t, createdSkillID, "Skill ID was not set after creation")

	t.Run("List Skills", func(t *testing.T) {
		limit := int64(250)
		query := createdName
		listParams := agentclient.NewAgentListSkillsParams().
			WithContext(ctx).
			WithTimeout(defaultTimeout).
			WithSearchQuery(&query).
			WithLimit(&limit)

		listResp, err := apiClient.Agent.AgentListSkills(listParams, nil)
		require.NoError(t, err, "List skills failed")
		require.NotNil(t, listResp, "List skills response should not be nil")
		require.NotNil(t, listResp.Payload, "List skills response payload should not be nil")

		var foundSkill *models.AgentSkillSummary
		for _, skill := range listResp.Payload.Skills {
			if skill != nil && skill.ID != nil && *skill.ID == createdSkillID {
				foundSkill = skill
				break
			}
		}
		require.NotNil(t, foundSkill, "Created skill %s not found in list response", createdSkillID)
		require.NotNil(t, foundSkill.Name, "Listed skill name should not be nil")
		require.Equal(t, createdName, *foundSkill.Name, "List skill name mismatch")
		requireSkillProvisioned(t, foundSkill.IsProvisioned, "listed skill")
	})

	t.Run("Get Skill", func(t *testing.T) {
		getParams := agentclient.NewAgentGetSkillParams().
			WithContext(ctx).
			WithTimeout(defaultTimeout).
			WithSkillID(createdSkillID)

		getResp, err := apiClient.Agent.AgentGetSkill(getParams, nil)
		require.NoError(t, err, "Get skill failed")
		require.NotNil(t, getResp, "Get skill response should not be nil")
		require.NotNil(t, getResp.Payload, "Get skill response payload should not be nil")
		require.NotNil(t, getResp.Payload.Skill, "Get skill response skill should not be nil")
		require.NotNil(t, getResp.Payload.Skill.Name, "Get skill name should not be nil")
		require.Equal(t, createdName, *getResp.Payload.Skill.Name, "Get skill name mismatch")
		requireSkillProvisioned(t, getResp.Payload.Skill.IsProvisioned, "fetched skill")
	})

	t.Run("Update Skill", func(t *testing.T) {
		updateParams := agentclient.NewAgentUpdateSkillParams().
			WithContext(ctx).
			WithTimeout(defaultTimeout).
			WithSkillID(createdSkillID).
			WithBody(newSkillUpdateRequest(updatedName, updatedWhenToUse, updatedDescription, isOrganizational))

		updateResp, err := apiClient.Agent.AgentUpdateSkill(updateParams, nil)
		require.NoError(t, err, "Update skill failed")
		require.NotNil(t, updateResp, "Update skill response should not be nil")
		require.NotNil(t, updateResp.Payload, "Update skill response payload should not be nil")
		require.NotNil(t, updateResp.Payload.Skill, "Update skill response skill should not be nil")
		require.NotNil(t, updateResp.Payload.Skill.Name, "Updated skill name should not be nil")
		require.Equal(t, updatedName, *updateResp.Payload.Skill.Name, "Updated skill name mismatch")
		requireSkillProvisioned(t, updateResp.Payload.Skill.IsProvisioned, "updated skill")
	})

	t.Run("Delete Skill", func(t *testing.T) {
		deleteParams := agentclient.NewAgentDeleteSkillParams().
			WithContext(ctx).
			WithTimeout(defaultTimeout).
			WithSkillID(createdSkillID)

		deleteResp, err := apiClient.Agent.AgentDeleteSkill(deleteParams, nil)
		require.NoError(t, err, "Delete skill failed")
		require.NotNil(t, deleteResp, "Delete skill response should not be nil")
		require.NotNil(t, deleteResp.Payload, "Delete skill response payload should not be nil")
		require.NotNil(t, deleteResp.Payload.SkillID, "Delete skill_id should not be nil")
		require.Equal(t, createdSkillID, *deleteResp.Payload.SkillID, "Delete skill ID mismatch")

		tc.UntrackAgentSkill(createdSkillID)
		t.Logf("Deleted agent skill %s", createdSkillID)
	})
}

func newSkillCreateRequest(name, whenToUse, description string, isOrganizational bool) *models.AgentSkillRequest {
	instructions := "Run the SDK E2E create skill instruction."

	return &models.AgentSkillRequest{
		Name:             &name,
		WhenToUse:        &whenToUse,
		Description:      description,
		Instructions:     &instructions,
		IsOrganizational: &isOrganizational,
	}
}

func newSkillUpdateRequest(name, whenToUse, description string, isOrganizational bool) *models.AgentSkillRequest {
	instructions := "Run the SDK E2E update skill instruction."

	return &models.AgentSkillRequest{
		Name:             &name,
		WhenToUse:        &whenToUse,
		Description:      description,
		Instructions:     &instructions,
		IsOrganizational: &isOrganizational,
	}
}

func requireSkillProvisioned(t *testing.T, isProvisioned *bool, context string) {
	t.Helper()
	require.NotNil(t, isProvisioned, "%s is_provisioned should not be nil", context)
	require.True(t, *isProvisioned, "%s should be provisioned", context)
}
