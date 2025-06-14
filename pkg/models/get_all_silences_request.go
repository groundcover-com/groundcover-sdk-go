// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetAllSilencesRequest get all silences request
//
// swagger:model GetAllSilencesRequest
type GetAllSilencesRequest struct {

	// Filter to show only active silences
	Active bool `json:"active,omitempty"`

	// backend ID
	BackendID string `json:"BackendID,omitempty"`

	// client ID
	ClientID string `json:"ClientID,omitempty"`

	// Maximum number of silences to return
	Limit int64 `json:"limit,omitempty"`

	// Number of silences to skip
	Skip int64 `json:"skip,omitempty"`
}

// Validate validates this get all silences request
func (m *GetAllSilencesRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get all silences request based on context it is used
func (m *GetAllSilencesRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetAllSilencesRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetAllSilencesRequest) UnmarshalBinary(b []byte) error {
	var res GetAllSilencesRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
