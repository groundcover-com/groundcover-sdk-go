// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ListAPIKeysResponseItem list Api keys response item
//
// swagger:model ListApiKeysResponseItem
type ListAPIKeysResponseItem struct {

	// Email of the user who created the key.
	CreatedBy string `json:"createdBy,omitempty"`

	// Timestamp when the key was created.
	// Format: date-time
	// Format: date-time
	CreationDate strfmt.DateTime `json:"creationDate,omitempty"`

	// Optional description for the API key.
	Description string `json:"description,omitempty"`

	// Timestamp when the key expires/expired (null if no expiration).
	// Format: date-time
	// Format: date-time
	ExpiredAt strfmt.DateTime `json:"expiredAt,omitempty"`

	// The UUID of the API key resource.
	ID string `json:"id,omitempty"`

	// Timestamp of the last activity detected for this key.
	// Format: date-time
	// Format: date-time
	LastActive strfmt.DateTime `json:"lastActive,omitempty"`

	// User-defined name for the API key.
	Name string `json:"name,omitempty"`

	// Policies associated with the service account owning this key.
	Policies []*PolicyRef `json:"policies"`

	// Timestamp when the key was revoked (null if active).
	// Format: date-time
	// Format: date-time
	RevokedAt strfmt.DateTime `json:"revokedAt,omitempty"`

	// The UUID of the service account this key belongs to.
	ServiceAccountID string `json:"serviceAccountId,omitempty"`

	// The name of the service account this key belongs to.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
}

// Validate validates this list Api keys response item
func (m *ListAPIKeysResponseItem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreationDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExpiredAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastActive(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePolicies(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRevokedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ListAPIKeysResponseItem) validateCreationDate(formats strfmt.Registry) error {
	if swag.IsZero(m.CreationDate) { // not required
		return nil
	}

	if err := validate.FormatOf("creationDate", "body", "date-time", m.CreationDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ListAPIKeysResponseItem) validateExpiredAt(formats strfmt.Registry) error {
	if swag.IsZero(m.ExpiredAt) { // not required
		return nil
	}

	if err := validate.FormatOf("expiredAt", "body", "date-time", m.ExpiredAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ListAPIKeysResponseItem) validateLastActive(formats strfmt.Registry) error {
	if swag.IsZero(m.LastActive) { // not required
		return nil
	}

	if err := validate.FormatOf("lastActive", "body", "date-time", m.LastActive.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ListAPIKeysResponseItem) validatePolicies(formats strfmt.Registry) error {
	if swag.IsZero(m.Policies) { // not required
		return nil
	}

	for i := 0; i < len(m.Policies); i++ {
		if swag.IsZero(m.Policies[i]) { // not required
			continue
		}

		if m.Policies[i] != nil {
			if err := m.Policies[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("policies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("policies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ListAPIKeysResponseItem) validateRevokedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.RevokedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("revokedAt", "body", "date-time", m.RevokedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this list Api keys response item based on the context it is used
func (m *ListAPIKeysResponseItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePolicies(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ListAPIKeysResponseItem) contextValidatePolicies(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Policies); i++ {

		if m.Policies[i] != nil {

			if swag.IsZero(m.Policies[i]) { // not required
				return nil
			}

			if err := m.Policies[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("policies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("policies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ListAPIKeysResponseItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ListAPIKeysResponseItem) UnmarshalBinary(b []byte) error {
	var res ListAPIKeysResponseItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
