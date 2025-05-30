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

// NewUpdateConfigParams creates a new UpdateConfigParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateConfigParams() *UpdateConfigParams {
	return &UpdateConfigParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateConfigParamsWithTimeout creates a new UpdateConfigParams object
// with the ability to set a timeout on a request.
func NewUpdateConfigParamsWithTimeout(timeout time.Duration) *UpdateConfigParams {
	return &UpdateConfigParams{
		timeout: timeout,
	}
}

// NewUpdateConfigParamsWithContext creates a new UpdateConfigParams object
// with the ability to set a context for a request.
func NewUpdateConfigParamsWithContext(ctx context.Context) *UpdateConfigParams {
	return &UpdateConfigParams{
		Context: ctx,
	}
}

// NewUpdateConfigParamsWithHTTPClient creates a new UpdateConfigParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateConfigParamsWithHTTPClient(client *http.Client) *UpdateConfigParams {
	return &UpdateConfigParams{
		HTTPClient: client,
	}
}

/*
UpdateConfigParams contains all the parameters to send to the API endpoint

	for the update config operation.

	Typically these are written to a http.Request.
*/
type UpdateConfigParams struct {

	/* Body.

	   The configuration to create or update
	*/
	Body *models.CreateOrUpdateLogsPipelineConfigRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateConfigParams) WithDefaults() *UpdateConfigParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update config params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateConfigParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update config params
func (o *UpdateConfigParams) WithTimeout(timeout time.Duration) *UpdateConfigParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update config params
func (o *UpdateConfigParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update config params
func (o *UpdateConfigParams) WithContext(ctx context.Context) *UpdateConfigParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update config params
func (o *UpdateConfigParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update config params
func (o *UpdateConfigParams) WithHTTPClient(client *http.Client) *UpdateConfigParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update config params
func (o *UpdateConfigParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update config params
func (o *UpdateConfigParams) WithBody(body *models.CreateOrUpdateLogsPipelineConfigRequest) *UpdateConfigParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update config params
func (o *UpdateConfigParams) SetBody(body *models.CreateOrUpdateLogsPipelineConfigRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateConfigParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
