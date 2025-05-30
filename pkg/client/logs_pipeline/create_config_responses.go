// Code generated by go-swagger; DO NOT EDIT.

package logs_pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

// CreateConfigReader is a Reader for the CreateConfig structure.
type CreateConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateConfigCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateConfigBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateConfigInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 503:
		result := NewCreateConfigServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /api/pipelines/logs/config] createConfig", response, response.Code())
	}
}

// NewCreateConfigCreated creates a CreateConfigCreated with default headers values
func NewCreateConfigCreated() *CreateConfigCreated {
	return &CreateConfigCreated{}
}

/*
CreateConfigCreated describes a response with status code 201, with default header values.

logsPipelineConfigResponse contains a logs pipeline configuration entry
*/
type CreateConfigCreated struct {
	Payload *models.LogsPipelineConfig
}

// IsSuccess returns true when this create config created response has a 2xx status code
func (o *CreateConfigCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create config created response has a 3xx status code
func (o *CreateConfigCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create config created response has a 4xx status code
func (o *CreateConfigCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this create config created response has a 5xx status code
func (o *CreateConfigCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this create config created response a status code equal to that given
func (o *CreateConfigCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the create config created response
func (o *CreateConfigCreated) Code() int {
	return 201
}

func (o *CreateConfigCreated) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigCreated %s", 201, payload)
}

func (o *CreateConfigCreated) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigCreated %s", 201, payload)
}

func (o *CreateConfigCreated) GetPayload() *models.LogsPipelineConfig {
	return o.Payload
}

func (o *CreateConfigCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.LogsPipelineConfig)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateConfigBadRequest creates a CreateConfigBadRequest with default headers values
func NewCreateConfigBadRequest() *CreateConfigBadRequest {
	return &CreateConfigBadRequest{}
}

/*
CreateConfigBadRequest describes a response with status code 400, with default header values.

emptyLogsPipelineConfigResponse is used for empty responses
*/
type CreateConfigBadRequest struct {
	Payload interface{}
}

// IsSuccess returns true when this create config bad request response has a 2xx status code
func (o *CreateConfigBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create config bad request response has a 3xx status code
func (o *CreateConfigBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create config bad request response has a 4xx status code
func (o *CreateConfigBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this create config bad request response has a 5xx status code
func (o *CreateConfigBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this create config bad request response a status code equal to that given
func (o *CreateConfigBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the create config bad request response
func (o *CreateConfigBadRequest) Code() int {
	return 400
}

func (o *CreateConfigBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigBadRequest %s", 400, payload)
}

func (o *CreateConfigBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigBadRequest %s", 400, payload)
}

func (o *CreateConfigBadRequest) GetPayload() interface{} {
	return o.Payload
}

func (o *CreateConfigBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateConfigInternalServerError creates a CreateConfigInternalServerError with default headers values
func NewCreateConfigInternalServerError() *CreateConfigInternalServerError {
	return &CreateConfigInternalServerError{}
}

/*
CreateConfigInternalServerError describes a response with status code 500, with default header values.

emptyLogsPipelineConfigResponse is used for empty responses
*/
type CreateConfigInternalServerError struct {
	Payload interface{}
}

// IsSuccess returns true when this create config internal server error response has a 2xx status code
func (o *CreateConfigInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create config internal server error response has a 3xx status code
func (o *CreateConfigInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create config internal server error response has a 4xx status code
func (o *CreateConfigInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create config internal server error response has a 5xx status code
func (o *CreateConfigInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create config internal server error response a status code equal to that given
func (o *CreateConfigInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create config internal server error response
func (o *CreateConfigInternalServerError) Code() int {
	return 500
}

func (o *CreateConfigInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigInternalServerError %s", 500, payload)
}

func (o *CreateConfigInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigInternalServerError %s", 500, payload)
}

func (o *CreateConfigInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *CreateConfigInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateConfigServiceUnavailable creates a CreateConfigServiceUnavailable with default headers values
func NewCreateConfigServiceUnavailable() *CreateConfigServiceUnavailable {
	return &CreateConfigServiceUnavailable{}
}

/*
CreateConfigServiceUnavailable describes a response with status code 503, with default header values.

emptyLogsPipelineConfigResponse is used for empty responses
*/
type CreateConfigServiceUnavailable struct {
	Payload interface{}
}

// IsSuccess returns true when this create config service unavailable response has a 2xx status code
func (o *CreateConfigServiceUnavailable) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create config service unavailable response has a 3xx status code
func (o *CreateConfigServiceUnavailable) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create config service unavailable response has a 4xx status code
func (o *CreateConfigServiceUnavailable) IsClientError() bool {
	return false
}

// IsServerError returns true when this create config service unavailable response has a 5xx status code
func (o *CreateConfigServiceUnavailable) IsServerError() bool {
	return true
}

// IsCode returns true when this create config service unavailable response a status code equal to that given
func (o *CreateConfigServiceUnavailable) IsCode(code int) bool {
	return code == 503
}

// Code gets the status code for the create config service unavailable response
func (o *CreateConfigServiceUnavailable) Code() int {
	return 503
}

func (o *CreateConfigServiceUnavailable) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigServiceUnavailable %s", 503, payload)
}

func (o *CreateConfigServiceUnavailable) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/pipelines/logs/config][%d] createConfigServiceUnavailable %s", 503, payload)
}

func (o *CreateConfigServiceUnavailable) GetPayload() interface{} {
	return o.Payload
}

func (o *CreateConfigServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
