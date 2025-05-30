// Code generated by go-swagger; DO NOT EDIT.

package apikeys

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
)

// NewCreateAPIKeyParams creates a new CreateAPIKeyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateAPIKeyParams() *CreateAPIKeyParams {
	return &CreateAPIKeyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateAPIKeyParamsWithTimeout creates a new CreateAPIKeyParams object
// with the ability to set a timeout on a request.
func NewCreateAPIKeyParamsWithTimeout(timeout time.Duration) *CreateAPIKeyParams {
	return &CreateAPIKeyParams{
		timeout: timeout,
	}
}

// NewCreateAPIKeyParamsWithContext creates a new CreateAPIKeyParams object
// with the ability to set a context for a request.
func NewCreateAPIKeyParamsWithContext(ctx context.Context) *CreateAPIKeyParams {
	return &CreateAPIKeyParams{
		Context: ctx,
	}
}

// NewCreateAPIKeyParamsWithHTTPClient creates a new CreateAPIKeyParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateAPIKeyParamsWithHTTPClient(client *http.Client) *CreateAPIKeyParams {
	return &CreateAPIKeyParams{
		HTTPClient: client,
	}
}

/*
CreateAPIKeyParams contains all the parameters to send to the API endpoint

	for the create Api key operation.

	Typically these are written to a http.Request.
*/
type CreateAPIKeyParams struct {

	/* Body.

	   API Key creation details
	*/
	Body *models.CreateAPIKeyRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create Api key params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateAPIKeyParams) WithDefaults() *CreateAPIKeyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create Api key params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateAPIKeyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create Api key params
func (o *CreateAPIKeyParams) WithTimeout(timeout time.Duration) *CreateAPIKeyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create Api key params
func (o *CreateAPIKeyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create Api key params
func (o *CreateAPIKeyParams) WithContext(ctx context.Context) *CreateAPIKeyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create Api key params
func (o *CreateAPIKeyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create Api key params
func (o *CreateAPIKeyParams) WithHTTPClient(client *http.Client) *CreateAPIKeyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create Api key params
func (o *CreateAPIKeyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create Api key params
func (o *CreateAPIKeyParams) WithBody(body *models.CreateAPIKeyRequest) *CreateAPIKeyParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create Api key params
func (o *CreateAPIKeyParams) SetBody(body *models.CreateAPIKeyRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateAPIKeyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
