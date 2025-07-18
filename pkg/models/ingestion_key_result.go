// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IngestionKeyResult ingestion key result
//
// swagger:model IngestionKeyResult
type IngestionKeyResult struct {

	// created by
	CreatedBy string `json:"createdBy,omitempty"`

	// creation date
	// Format: date-time
	CreationDate strfmt.DateTime `json:"creationDate,omitempty"`

	// ID
	ID string `json:"id,omitempty"`

	// key
	Key string `json:"key,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// remote config
	RemoteConfig bool `json:"remoteConfig,omitempty"`

	// tags
	Tags []string `json:"tags"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this ingestion key result
func (m *IngestionKeyResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreationDate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IngestionKeyResult) validateCreationDate(formats strfmt.Registry) error {
	if swag.IsZero(m.CreationDate) { // not required
		return nil
	}

	if err := validate.FormatOf("creationDate", "body", "date-time", m.CreationDate.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this ingestion key result based on context it is used
func (m *IngestionKeyResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IngestionKeyResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IngestionKeyResult) UnmarshalBinary(b []byte) error {
	var res IngestionKeyResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
