package events

import (
	"context"
	"errors"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

const (
	ReasonOOMKilled    = "OOMKilled"
	TypeContainerCrash = "container_crash"
	ReasonFilter       = "reason"
	TypeFilter         = "type"
)

func GetOOMEvents(ctx context.Context, service *Service, request *EventsOverTimeRequest) (*EventsOverTimeResponse, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	if request.Conditions == nil {
		request.Conditions = []models.Condition{}
	}

	request.Conditions = append(request.Conditions, models.NewEqualStringCondition(ReasonFilter, ReasonOOMKilled), models.NewEqualStringCondition(TypeFilter, TypeContainerCrash))
	return service.EventsOverTime(ctx, request)
}
