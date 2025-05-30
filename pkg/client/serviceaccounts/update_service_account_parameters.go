// Code generated by go-swagger; DO NOT EDIT.

package serviceaccounts

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

// NewUpdateServiceAccountParams creates a new UpdateServiceAccountParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateServiceAccountParams() *UpdateServiceAccountParams {
	return &UpdateServiceAccountParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateServiceAccountParamsWithTimeout creates a new UpdateServiceAccountParams object
// with the ability to set a timeout on a request.
func NewUpdateServiceAccountParamsWithTimeout(timeout time.Duration) *UpdateServiceAccountParams {
	return &UpdateServiceAccountParams{
		timeout: timeout,
	}
}

// NewUpdateServiceAccountParamsWithContext creates a new UpdateServiceAccountParams object
// with the ability to set a context for a request.
func NewUpdateServiceAccountParamsWithContext(ctx context.Context) *UpdateServiceAccountParams {
	return &UpdateServiceAccountParams{
		Context: ctx,
	}
}

// NewUpdateServiceAccountParamsWithHTTPClient creates a new UpdateServiceAccountParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateServiceAccountParamsWithHTTPClient(client *http.Client) *UpdateServiceAccountParams {
	return &UpdateServiceAccountParams{
		HTTPClient: client,
	}
}

/*
UpdateServiceAccountParams contains all the parameters to send to the API endpoint

	for the update service account operation.

	Typically these are written to a http.Request.
*/
type UpdateServiceAccountParams struct {

	/* Body.

	   Service account update information
	*/
	Body *models.UpdateServiceAccountRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update service account params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateServiceAccountParams) WithDefaults() *UpdateServiceAccountParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update service account params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateServiceAccountParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update service account params
func (o *UpdateServiceAccountParams) WithTimeout(timeout time.Duration) *UpdateServiceAccountParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update service account params
func (o *UpdateServiceAccountParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update service account params
func (o *UpdateServiceAccountParams) WithContext(ctx context.Context) *UpdateServiceAccountParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update service account params
func (o *UpdateServiceAccountParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update service account params
func (o *UpdateServiceAccountParams) WithHTTPClient(client *http.Client) *UpdateServiceAccountParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update service account params
func (o *UpdateServiceAccountParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update service account params
func (o *UpdateServiceAccountParams) WithBody(body *models.UpdateServiceAccountRequest) *UpdateServiceAccountParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update service account params
func (o *UpdateServiceAccountParams) SetBody(body *models.UpdateServiceAccountRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateServiceAccountParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
