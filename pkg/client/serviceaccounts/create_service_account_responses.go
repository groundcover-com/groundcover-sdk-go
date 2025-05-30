// Code generated by go-swagger; DO NOT EDIT.

package serviceaccounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

// CreateServiceAccountReader is a Reader for the CreateServiceAccount structure.
type CreateServiceAccountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateServiceAccountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateServiceAccountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateServiceAccountBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateServiceAccountConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateServiceAccountInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /api/rbac/service-account/create] createServiceAccount", response, response.Code())
	}
}

// NewCreateServiceAccountOK creates a CreateServiceAccountOK with default headers values
func NewCreateServiceAccountOK() *CreateServiceAccountOK {
	return &CreateServiceAccountOK{}
}

/*
CreateServiceAccountOK describes a response with status code 200, with default header values.

CreateServiceAccountOK create service account o k
*/
type CreateServiceAccountOK struct {
	Payload *models.ServiceAccountCreatePayload
}

// IsSuccess returns true when this create service account o k response has a 2xx status code
func (o *CreateServiceAccountOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create service account o k response has a 3xx status code
func (o *CreateServiceAccountOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create service account o k response has a 4xx status code
func (o *CreateServiceAccountOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this create service account o k response has a 5xx status code
func (o *CreateServiceAccountOK) IsServerError() bool {
	return false
}

// IsCode returns true when this create service account o k response a status code equal to that given
func (o *CreateServiceAccountOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the create service account o k response
func (o *CreateServiceAccountOK) Code() int {
	return 200
}

func (o *CreateServiceAccountOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountOK %s", 200, payload)
}

func (o *CreateServiceAccountOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountOK %s", 200, payload)
}

func (o *CreateServiceAccountOK) GetPayload() *models.ServiceAccountCreatePayload {
	return o.Payload
}

func (o *CreateServiceAccountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServiceAccountCreatePayload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateServiceAccountBadRequest creates a CreateServiceAccountBadRequest with default headers values
func NewCreateServiceAccountBadRequest() *CreateServiceAccountBadRequest {
	return &CreateServiceAccountBadRequest{}
}

/*
CreateServiceAccountBadRequest describes a response with status code 400, with default header values.

ErrorResponse defines a common error response structure.
*/
type CreateServiceAccountBadRequest struct {
	Payload *CreateServiceAccountBadRequestBody
}

// IsSuccess returns true when this create service account bad request response has a 2xx status code
func (o *CreateServiceAccountBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create service account bad request response has a 3xx status code
func (o *CreateServiceAccountBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create service account bad request response has a 4xx status code
func (o *CreateServiceAccountBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this create service account bad request response has a 5xx status code
func (o *CreateServiceAccountBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this create service account bad request response a status code equal to that given
func (o *CreateServiceAccountBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the create service account bad request response
func (o *CreateServiceAccountBadRequest) Code() int {
	return 400
}

func (o *CreateServiceAccountBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountBadRequest %s", 400, payload)
}

func (o *CreateServiceAccountBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountBadRequest %s", 400, payload)
}

func (o *CreateServiceAccountBadRequest) GetPayload() *CreateServiceAccountBadRequestBody {
	return o.Payload
}

func (o *CreateServiceAccountBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateServiceAccountBadRequestBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateServiceAccountConflict creates a CreateServiceAccountConflict with default headers values
func NewCreateServiceAccountConflict() *CreateServiceAccountConflict {
	return &CreateServiceAccountConflict{}
}

/*
CreateServiceAccountConflict describes a response with status code 409, with default header values.

ErrorResponse defines a common error response structure.
*/
type CreateServiceAccountConflict struct {
	Payload *CreateServiceAccountConflictBody
}

// IsSuccess returns true when this create service account conflict response has a 2xx status code
func (o *CreateServiceAccountConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create service account conflict response has a 3xx status code
func (o *CreateServiceAccountConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create service account conflict response has a 4xx status code
func (o *CreateServiceAccountConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this create service account conflict response has a 5xx status code
func (o *CreateServiceAccountConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this create service account conflict response a status code equal to that given
func (o *CreateServiceAccountConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the create service account conflict response
func (o *CreateServiceAccountConflict) Code() int {
	return 409
}

func (o *CreateServiceAccountConflict) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountConflict %s", 409, payload)
}

func (o *CreateServiceAccountConflict) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountConflict %s", 409, payload)
}

func (o *CreateServiceAccountConflict) GetPayload() *CreateServiceAccountConflictBody {
	return o.Payload
}

func (o *CreateServiceAccountConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateServiceAccountConflictBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateServiceAccountInternalServerError creates a CreateServiceAccountInternalServerError with default headers values
func NewCreateServiceAccountInternalServerError() *CreateServiceAccountInternalServerError {
	return &CreateServiceAccountInternalServerError{}
}

/*
CreateServiceAccountInternalServerError describes a response with status code 500, with default header values.

ErrorResponse defines a common error response structure.
*/
type CreateServiceAccountInternalServerError struct {
	Payload *CreateServiceAccountInternalServerErrorBody
}

// IsSuccess returns true when this create service account internal server error response has a 2xx status code
func (o *CreateServiceAccountInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create service account internal server error response has a 3xx status code
func (o *CreateServiceAccountInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create service account internal server error response has a 4xx status code
func (o *CreateServiceAccountInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create service account internal server error response has a 5xx status code
func (o *CreateServiceAccountInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create service account internal server error response a status code equal to that given
func (o *CreateServiceAccountInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create service account internal server error response
func (o *CreateServiceAccountInternalServerError) Code() int {
	return 500
}

func (o *CreateServiceAccountInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountInternalServerError %s", 500, payload)
}

func (o *CreateServiceAccountInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/rbac/service-account/create][%d] createServiceAccountInternalServerError %s", 500, payload)
}

func (o *CreateServiceAccountInternalServerError) GetPayload() *CreateServiceAccountInternalServerErrorBody {
	return o.Payload
}

func (o *CreateServiceAccountInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateServiceAccountInternalServerErrorBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
CreateServiceAccountBadRequestBody create service account bad request body
swagger:model CreateServiceAccountBadRequestBody
*/
type CreateServiceAccountBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this create service account bad request body
func (o *CreateServiceAccountBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create service account bad request body based on context it is used
func (o *CreateServiceAccountBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateServiceAccountBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateServiceAccountBadRequestBody) UnmarshalBinary(b []byte) error {
	var res CreateServiceAccountBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
CreateServiceAccountConflictBody create service account conflict body
swagger:model CreateServiceAccountConflictBody
*/
type CreateServiceAccountConflictBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this create service account conflict body
func (o *CreateServiceAccountConflictBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create service account conflict body based on context it is used
func (o *CreateServiceAccountConflictBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateServiceAccountConflictBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateServiceAccountConflictBody) UnmarshalBinary(b []byte) error {
	var res CreateServiceAccountConflictBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
CreateServiceAccountInternalServerErrorBody create service account internal server error body
swagger:model CreateServiceAccountInternalServerErrorBody
*/
type CreateServiceAccountInternalServerErrorBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this create service account internal server error body
func (o *CreateServiceAccountInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create service account internal server error body based on context it is used
func (o *CreateServiceAccountInternalServerErrorBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CreateServiceAccountInternalServerErrorBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateServiceAccountInternalServerErrorBody) UnmarshalBinary(b []byte) error {
	var res CreateServiceAccountInternalServerErrorBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
