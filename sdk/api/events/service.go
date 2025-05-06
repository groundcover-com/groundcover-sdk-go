package events

import (
	"context"
	"net/http"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
)

type Service struct {
	client httpclient.Doer
}

func New(client httpclient.Doer) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) EventsOverTime(ctx context.Context, request *EventsOverTimeRequest, opts ...httpclient.RequestOption) (*EventsOverTimeResponse, error) {
	urlPath := "/api/k8s/v2/events-over-time"
	response := &EventsOverTimeResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodPost,
		Path:     urlPath,
		Body:     request,
		Response: response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}
