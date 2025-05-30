// Code generated by go-swagger; DO NOT EDIT.

package logs_pipeline

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

// NewCreateConfigParams creates a new CreateConfigParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateConfigParams() *CreateConfigParams {
	return &CreateConfigParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateConfigParamsWithTimeout creates a new CreateConfigParams object
// with the ability to set a timeout on a request.
func NewCreateConfigParamsWithTimeout(timeout time.Duration) *CreateConfigParams {
	return &CreateConfigParams{
		timeout: timeout,
	}
}

// NewCreateConfigParamsWithContext creates a new CreateConfigParams object
// with the ability to set a context for a request.
func NewCreateConfigParamsWithContext(ctx context.Context) *CreateConfigParams {
	return &CreateConfigParams{
		Context: ctx,
	}
}

// NewCreateConfigParamsWithHTTPClient creates a new CreateConfigParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateConfigParamsWithHTTPClient(client *http.Client) *CreateConfigParams {
	return &CreateConfigParams{
		HTTPClient: client,
	}
}

/*
CreateConfigParams contains all the parameters to send to the API endpoint

	for the create config operation.

	Typically these are written to a http.Request.
*/
type CreateConfigParams struct {

	/* Body.

	   The configuration to create or update
	*/
	Body *models.CreateOrUpdateLogsPipelineConfigRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateConfigParams) WithDefaults() *CreateConfigParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateConfigParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create config params
func (o *CreateConfigParams) WithTimeout(timeout time.Duration) *CreateConfigParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create config params
func (o *CreateConfigParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create config params
func (o *CreateConfigParams) WithContext(ctx context.Context) *CreateConfigParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create config params
func (o *CreateConfigParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create config params
func (o *CreateConfigParams) WithHTTPClient(client *http.Client) *CreateConfigParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create config params
func (o *CreateConfigParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create config params
func (o *CreateConfigParams) WithBody(body *models.CreateOrUpdateLogsPipelineConfigRequest) *CreateConfigParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create config params
func (o *CreateConfigParams) SetBody(body *models.CreateOrUpdateLogsPipelineConfigRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateConfigParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
