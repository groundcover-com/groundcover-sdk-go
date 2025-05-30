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
)

// NewListServiceAccountsParams creates a new ListServiceAccountsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListServiceAccountsParams() *ListServiceAccountsParams {
	return &ListServiceAccountsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListServiceAccountsParamsWithTimeout creates a new ListServiceAccountsParams object
// with the ability to set a timeout on a request.
func NewListServiceAccountsParamsWithTimeout(timeout time.Duration) *ListServiceAccountsParams {
	return &ListServiceAccountsParams{
		timeout: timeout,
	}
}

// NewListServiceAccountsParamsWithContext creates a new ListServiceAccountsParams object
// with the ability to set a context for a request.
func NewListServiceAccountsParamsWithContext(ctx context.Context) *ListServiceAccountsParams {
	return &ListServiceAccountsParams{
		Context: ctx,
	}
}

// NewListServiceAccountsParamsWithHTTPClient creates a new ListServiceAccountsParams object
// with the ability to set a custom HTTPClient for a request.
func NewListServiceAccountsParamsWithHTTPClient(client *http.Client) *ListServiceAccountsParams {
	return &ListServiceAccountsParams{
		HTTPClient: client,
	}
}

/*
ListServiceAccountsParams contains all the parameters to send to the API endpoint

	for the list service accounts operation.

	Typically these are written to a http.Request.
*/
type ListServiceAccountsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list service accounts params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListServiceAccountsParams) WithDefaults() *ListServiceAccountsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list service accounts params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListServiceAccountsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list service accounts params
func (o *ListServiceAccountsParams) WithTimeout(timeout time.Duration) *ListServiceAccountsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list service accounts params
func (o *ListServiceAccountsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list service accounts params
func (o *ListServiceAccountsParams) WithContext(ctx context.Context) *ListServiceAccountsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list service accounts params
func (o *ListServiceAccountsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list service accounts params
func (o *ListServiceAccountsParams) WithHTTPClient(client *http.Client) *ListServiceAccountsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list service accounts params
func (o *ListServiceAccountsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ListServiceAccountsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
