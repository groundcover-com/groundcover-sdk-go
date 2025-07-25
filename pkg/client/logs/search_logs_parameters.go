// Code generated by go-swagger; DO NOT EDIT.

package logs

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

// NewSearchLogsParams creates a new SearchLogsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchLogsParams() *SearchLogsParams {
	return &SearchLogsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchLogsParamsWithTimeout creates a new SearchLogsParams object
// with the ability to set a timeout on a request.
func NewSearchLogsParamsWithTimeout(timeout time.Duration) *SearchLogsParams {
	return &SearchLogsParams{
		timeout: timeout,
	}
}

// NewSearchLogsParamsWithContext creates a new SearchLogsParams object
// with the ability to set a context for a request.
func NewSearchLogsParamsWithContext(ctx context.Context) *SearchLogsParams {
	return &SearchLogsParams{
		Context: ctx,
	}
}

// NewSearchLogsParamsWithHTTPClient creates a new SearchLogsParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchLogsParamsWithHTTPClient(client *http.Client) *SearchLogsParams {
	return &SearchLogsParams{
		HTTPClient: client,
	}
}

/*
SearchLogsParams contains all the parameters to send to the API endpoint

	for the search logs operation.

	Typically these are written to a http.Request.
*/
type SearchLogsParams struct {

	/* Body.

	   Logs search request
	*/
	Body *models.LogsSearchRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search logs params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchLogsParams) WithDefaults() *SearchLogsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search logs params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchLogsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the search logs params
func (o *SearchLogsParams) WithTimeout(timeout time.Duration) *SearchLogsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search logs params
func (o *SearchLogsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search logs params
func (o *SearchLogsParams) WithContext(ctx context.Context) *SearchLogsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search logs params
func (o *SearchLogsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search logs params
func (o *SearchLogsParams) WithHTTPClient(client *http.Client) *SearchLogsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search logs params
func (o *SearchLogsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the search logs params
func (o *SearchLogsParams) WithBody(body *models.LogsSearchRequest) *SearchLogsParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the search logs params
func (o *SearchLogsParams) SetBody(body *models.LogsSearchRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *SearchLogsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
