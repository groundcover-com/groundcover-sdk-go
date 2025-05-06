package events

import (
	"encoding/json"
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type EventsOverTimeRequest struct {
	Start         time.Time          `form:"start" url:"start" binding:"required"`
	End           time.Time          `form:"end" url:"end" binding:"required,gtefield=Start"`
	Conditions    []models.Condition `form:"conditions" url:"conditions"`
	SortBy        string             `form:"sortBy" url:"sortBy" binding:"required,oneof=timestamp namespace instance object_kind firstSeen lastSeen type reason count workload cluster"`
	SortOrder     string             `form:"sortOrder" url:"sortOrder" binding:"required,oneof=asc desc"`
	Limit         uint32             `form:"limit" url:"limit"`
	Skip          uint32             `form:"skip" url:"skip"`
	Sources       []models.Condition `form:"sources" url:"sources"`
	WithRawEvents bool               `form:"withRawEvents" url:"withRawEvents"`
}

type EventsOverTimeResponse struct {
	Events           []EventOverTimeItem `json:"events"`
	IsLimitReached   bool                `json:"isLimitReached"`
	WarningIndicator bool                `json:"warningIndicator"`
}

type EventOverTimeItem struct {
	Uid              string          `json:"uid"`
	NormalizedUid    string          `json:"normalized_uid"`
	Timestamp        time.Time       `json:"timestamp"`
	ObjectUid        string          `json:"object_uid"`
	Namespace        string          `json:"namespace"`
	Instance         string          `json:"instance"`
	ObjectKind       string          `json:"object_kind"`
	Type             string          `json:"type"`
	Reason           string          `json:"reason"`
	Message          string          `json:"message"`
	ExitCode         string          `json:"exitCode"`
	CreatedTimestamp time.Time       `json:"createdTimestamp"`
	Workload         string          `json:"workload"`
	Env              string          `json:"env"`
	Cluster          string          `json:"cluster"`
	Raw              json.RawMessage `json:"raw,omitempty"`
}
