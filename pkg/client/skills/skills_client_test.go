package skills

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (ClientService, func()) {
	t.Helper()

	server := httptest.NewServer(handler)
	parsedURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	transport := httptransport.New(parsedURL.Host, "", []string{parsedURL.Scheme})
	return New(transport, strfmt.Default), server.Close
}

func TestCreateSkillSendsOrgHeadersAndBody(t *testing.T) {
	client, cleanup := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/agent/skills", r.URL.Path)
		require.Equal(t, "tenant-1", r.Header.Get("X-Tenant-UUID"))
		require.Equal(t, "admin", r.Header.Get("X-Resolved-Role"))

		var body models.UserSkillCreateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "Incident Triage", body.Name)
		require.True(t, body.IsOrganizational)
		require.Equal(t, 1, body.Definition.SchemaVersion)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","skill":{"id":"skill-1","name":"Incident Triage","when_to_use":"during incidents","definition_schema_version":1,"revision":1,"is_organizational":true,"created_at":"2026-06-22T00:00:00Z","updated_at":"2026-06-22T00:00:00Z","definition":{"schemaVersion":1,"steps":[{"id":"step-1","type":"llm_markdown","markdown":"Do the thing"}]},"owner_user_id":"auth0|user"}}`))
	})
	defer cleanup()

	resp, err := client.CreateSkill(
		NewCreateSkillParams().
			WithTenantUUID("tenant-1").
			WithResolvedRole("admin").
			WithBody(&models.UserSkillCreateRequest{
				Name:             "Incident Triage",
				WhenToUse:        "during incidents",
				IsOrganizational: true,
				Definition: &models.UserSkillDefinitionV1{
					SchemaVersion: 1,
					Steps: []*models.UserSkillStep{
						{ID: "step-1", Type: "llm_markdown", Markdown: "Do the thing"},
					},
				},
			}),
		nil,
	)

	require.NoError(t, err)
	require.NotNil(t, resp.Payload.Skill)
	require.Equal(t, "skill-1", resp.Payload.Skill.ID)
	require.True(t, resp.Payload.Skill.IsOrganizational)
}

func TestListSkillsSendsTenantAndQuery(t *testing.T) {
	client, cleanup := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/agent/skills", r.URL.Path)
		require.Equal(t, "tenant-1", r.Header.Get("X-Tenant-UUID"))
		require.Equal(t, "incident", r.URL.Query().Get("q"))
		require.Equal(t, "25", r.URL.Query().Get("limit"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","skills":[{"id":"skill-1","name":"Incident Triage","when_to_use":"during incidents","definition_schema_version":1,"revision":1,"created_at":"2026-06-22T00:00:00Z","updated_at":"2026-06-22T00:00:00Z"}]}`))
	})
	defer cleanup()

	q := "incident"
	limit := int64(25)
	resp, err := client.ListSkills(
		NewListSkillsParams().
			WithTenantUUID("tenant-1").
			WithQ(&q).
			WithLimit(&limit),
		nil,
	)

	require.NoError(t, err)
	require.Len(t, resp.Payload.Skills, 1)
	require.Equal(t, "Incident Triage", resp.Payload.Skills[0].Name)
}
