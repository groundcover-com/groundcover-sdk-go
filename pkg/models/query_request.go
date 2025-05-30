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

// QueryRequest query request
//
// swagger:model QueryRequest
type QueryRequest struct {

	// conditions
	Conditions []*Condition `json:"Conditions"`

	// end
	// Format: date-time
	End strfmt.DateTime `json:"End,omitempty"`

	// filters
	Filters string `json:"Filters,omitempty"`

	// pipeline
	Pipeline *PromqlPipeline `json:"Pipeline,omitempty"`

	// promql
	Promql string `json:"Promql,omitempty"`

	// query type
	QueryType string `json:"QueryType,omitempty"`

	// start
	// Format: date-time
	Start strfmt.DateTime `json:"Start,omitempty"`

	// step
	Step string `json:"Step,omitempty"`

	// sub pipelines
	SubPipelines KnownPipelines `json:"SubPipelines,omitempty"`
}

// Validate validates this query request
func (m *QueryRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConditions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnd(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePipeline(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubPipelines(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *QueryRequest) validateConditions(formats strfmt.Registry) error {
	if swag.IsZero(m.Conditions) { // not required
		return nil
	}

	for i := 0; i < len(m.Conditions); i++ {
		if swag.IsZero(m.Conditions[i]) { // not required
			continue
		}

		if m.Conditions[i] != nil {
			if err := m.Conditions[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Conditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Conditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *QueryRequest) validateEnd(formats strfmt.Registry) error {
	if swag.IsZero(m.End) { // not required
		return nil
	}

	if err := validate.FormatOf("End", "body", "date-time", m.End.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *QueryRequest) validatePipeline(formats strfmt.Registry) error {
	if swag.IsZero(m.Pipeline) { // not required
		return nil
	}

	if m.Pipeline != nil {
		if err := m.Pipeline.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Pipeline")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("Pipeline")
			}
			return err
		}
	}

	return nil
}

func (m *QueryRequest) validateStart(formats strfmt.Registry) error {
	if swag.IsZero(m.Start) { // not required
		return nil
	}

	if err := validate.FormatOf("Start", "body", "date-time", m.Start.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *QueryRequest) validateSubPipelines(formats strfmt.Registry) error {
	if swag.IsZero(m.SubPipelines) { // not required
		return nil
	}

	if m.SubPipelines != nil {
		if err := m.SubPipelines.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("SubPipelines")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("SubPipelines")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this query request based on the context it is used
func (m *QueryRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateConditions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePipeline(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSubPipelines(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *QueryRequest) contextValidateConditions(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Conditions); i++ {

		if m.Conditions[i] != nil {

			if swag.IsZero(m.Conditions[i]) { // not required
				return nil
			}

			if err := m.Conditions[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Conditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Conditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *QueryRequest) contextValidatePipeline(ctx context.Context, formats strfmt.Registry) error {

	if m.Pipeline != nil {

		if swag.IsZero(m.Pipeline) { // not required
			return nil
		}

		if err := m.Pipeline.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Pipeline")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("Pipeline")
			}
			return err
		}
	}

	return nil
}

func (m *QueryRequest) contextValidateSubPipelines(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.SubPipelines) { // not required
		return nil
	}

	if err := m.SubPipelines.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("SubPipelines")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("SubPipelines")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *QueryRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *QueryRequest) UnmarshalBinary(b []byte) error {
	var res QueryRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
