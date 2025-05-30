// Code generated by go-swagger; DO NOT EDIT.

package policies

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

// GetPolicyAuditTrailReader is a Reader for the GetPolicyAuditTrail structure.
type GetPolicyAuditTrailReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPolicyAuditTrailReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPolicyAuditTrailOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetPolicyAuditTrailBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetPolicyAuditTrailNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetPolicyAuditTrailInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /api/rbac/policy/{id}/auditTrail] getPolicyAuditTrail", response, response.Code())
	}
}

// NewGetPolicyAuditTrailOK creates a GetPolicyAuditTrailOK with default headers values
func NewGetPolicyAuditTrailOK() *GetPolicyAuditTrailOK {
	return &GetPolicyAuditTrailOK{}
}

/*
GetPolicyAuditTrailOK describes a response with status code 200, with default header values.

PolicyAuditTrailResponse contains the audit trail for a policy.
*/
type GetPolicyAuditTrailOK struct {
	Payload []*models.Policy
}

// IsSuccess returns true when this get policy audit trail o k response has a 2xx status code
func (o *GetPolicyAuditTrailOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get policy audit trail o k response has a 3xx status code
func (o *GetPolicyAuditTrailOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get policy audit trail o k response has a 4xx status code
func (o *GetPolicyAuditTrailOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get policy audit trail o k response has a 5xx status code
func (o *GetPolicyAuditTrailOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get policy audit trail o k response a status code equal to that given
func (o *GetPolicyAuditTrailOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get policy audit trail o k response
func (o *GetPolicyAuditTrailOK) Code() int {
	return 200
}

func (o *GetPolicyAuditTrailOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailOK %s", 200, payload)
}

func (o *GetPolicyAuditTrailOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailOK %s", 200, payload)
}

func (o *GetPolicyAuditTrailOK) GetPayload() []*models.Policy {
	return o.Payload
}

func (o *GetPolicyAuditTrailOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPolicyAuditTrailBadRequest creates a GetPolicyAuditTrailBadRequest with default headers values
func NewGetPolicyAuditTrailBadRequest() *GetPolicyAuditTrailBadRequest {
	return &GetPolicyAuditTrailBadRequest{}
}

/*
GetPolicyAuditTrailBadRequest describes a response with status code 400, with default header values.

ErrorResponse defines a common error response structure.
*/
type GetPolicyAuditTrailBadRequest struct {
	Payload *GetPolicyAuditTrailBadRequestBody
}

// IsSuccess returns true when this get policy audit trail bad request response has a 2xx status code
func (o *GetPolicyAuditTrailBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get policy audit trail bad request response has a 3xx status code
func (o *GetPolicyAuditTrailBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get policy audit trail bad request response has a 4xx status code
func (o *GetPolicyAuditTrailBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get policy audit trail bad request response has a 5xx status code
func (o *GetPolicyAuditTrailBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get policy audit trail bad request response a status code equal to that given
func (o *GetPolicyAuditTrailBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the get policy audit trail bad request response
func (o *GetPolicyAuditTrailBadRequest) Code() int {
	return 400
}

func (o *GetPolicyAuditTrailBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailBadRequest %s", 400, payload)
}

func (o *GetPolicyAuditTrailBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailBadRequest %s", 400, payload)
}

func (o *GetPolicyAuditTrailBadRequest) GetPayload() *GetPolicyAuditTrailBadRequestBody {
	return o.Payload
}

func (o *GetPolicyAuditTrailBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetPolicyAuditTrailBadRequestBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPolicyAuditTrailNotFound creates a GetPolicyAuditTrailNotFound with default headers values
func NewGetPolicyAuditTrailNotFound() *GetPolicyAuditTrailNotFound {
	return &GetPolicyAuditTrailNotFound{}
}

/*
GetPolicyAuditTrailNotFound describes a response with status code 404, with default header values.

ErrorResponse defines a common error response structure.
*/
type GetPolicyAuditTrailNotFound struct {
	Payload *GetPolicyAuditTrailNotFoundBody
}

// IsSuccess returns true when this get policy audit trail not found response has a 2xx status code
func (o *GetPolicyAuditTrailNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get policy audit trail not found response has a 3xx status code
func (o *GetPolicyAuditTrailNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get policy audit trail not found response has a 4xx status code
func (o *GetPolicyAuditTrailNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get policy audit trail not found response has a 5xx status code
func (o *GetPolicyAuditTrailNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get policy audit trail not found response a status code equal to that given
func (o *GetPolicyAuditTrailNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get policy audit trail not found response
func (o *GetPolicyAuditTrailNotFound) Code() int {
	return 404
}

func (o *GetPolicyAuditTrailNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailNotFound %s", 404, payload)
}

func (o *GetPolicyAuditTrailNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailNotFound %s", 404, payload)
}

func (o *GetPolicyAuditTrailNotFound) GetPayload() *GetPolicyAuditTrailNotFoundBody {
	return o.Payload
}

func (o *GetPolicyAuditTrailNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetPolicyAuditTrailNotFoundBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPolicyAuditTrailInternalServerError creates a GetPolicyAuditTrailInternalServerError with default headers values
func NewGetPolicyAuditTrailInternalServerError() *GetPolicyAuditTrailInternalServerError {
	return &GetPolicyAuditTrailInternalServerError{}
}

/*
GetPolicyAuditTrailInternalServerError describes a response with status code 500, with default header values.

ErrorResponse defines a common error response structure.
*/
type GetPolicyAuditTrailInternalServerError struct {
	Payload *GetPolicyAuditTrailInternalServerErrorBody
}

// IsSuccess returns true when this get policy audit trail internal server error response has a 2xx status code
func (o *GetPolicyAuditTrailInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get policy audit trail internal server error response has a 3xx status code
func (o *GetPolicyAuditTrailInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get policy audit trail internal server error response has a 4xx status code
func (o *GetPolicyAuditTrailInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get policy audit trail internal server error response has a 5xx status code
func (o *GetPolicyAuditTrailInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get policy audit trail internal server error response a status code equal to that given
func (o *GetPolicyAuditTrailInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get policy audit trail internal server error response
func (o *GetPolicyAuditTrailInternalServerError) Code() int {
	return 500
}

func (o *GetPolicyAuditTrailInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailInternalServerError %s", 500, payload)
}

func (o *GetPolicyAuditTrailInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/rbac/policy/{id}/auditTrail][%d] getPolicyAuditTrailInternalServerError %s", 500, payload)
}

func (o *GetPolicyAuditTrailInternalServerError) GetPayload() *GetPolicyAuditTrailInternalServerErrorBody {
	return o.Payload
}

func (o *GetPolicyAuditTrailInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetPolicyAuditTrailInternalServerErrorBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetPolicyAuditTrailBadRequestBody get policy audit trail bad request body
swagger:model GetPolicyAuditTrailBadRequestBody
*/
type GetPolicyAuditTrailBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get policy audit trail bad request body
func (o *GetPolicyAuditTrailBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get policy audit trail bad request body based on context it is used
func (o *GetPolicyAuditTrailBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetPolicyAuditTrailBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPolicyAuditTrailBadRequestBody) UnmarshalBinary(b []byte) error {
	var res GetPolicyAuditTrailBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
GetPolicyAuditTrailInternalServerErrorBody get policy audit trail internal server error body
swagger:model GetPolicyAuditTrailInternalServerErrorBody
*/
type GetPolicyAuditTrailInternalServerErrorBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get policy audit trail internal server error body
func (o *GetPolicyAuditTrailInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get policy audit trail internal server error body based on context it is used
func (o *GetPolicyAuditTrailInternalServerErrorBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetPolicyAuditTrailInternalServerErrorBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPolicyAuditTrailInternalServerErrorBody) UnmarshalBinary(b []byte) error {
	var res GetPolicyAuditTrailInternalServerErrorBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
GetPolicyAuditTrailNotFoundBody get policy audit trail not found body
swagger:model GetPolicyAuditTrailNotFoundBody
*/
type GetPolicyAuditTrailNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get policy audit trail not found body
func (o *GetPolicyAuditTrailNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get policy audit trail not found body based on context it is used
func (o *GetPolicyAuditTrailNotFoundBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetPolicyAuditTrailNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPolicyAuditTrailNotFoundBody) UnmarshalBinary(b []byte) error {
	var res GetPolicyAuditTrailNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
