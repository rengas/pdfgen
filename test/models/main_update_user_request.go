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

// MainUpdateUserRequest main update user request
//
// swagger:model main.UpdateUserRequest
type MainUpdateUserRequest struct {

	// email
	// Example: John.doe@email.com
	// Required: true
	Email *string `json:"email"`

	// first name
	// Example: John
	FirstName string `json:"firstName,omitempty"`

	// last name
	// Example: Doe
	LastName string `json:"lastName,omitempty"`
}

// Validate validates this main update user request
func (m *MainUpdateUserRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MainUpdateUserRequest) validateEmail(formats strfmt.Registry) error {

	if err := validate.Required("email", "body", m.Email); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this main update user request based on context it is used
func (m *MainUpdateUserRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MainUpdateUserRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MainUpdateUserRequest) UnmarshalBinary(b []byte) error {
	var res MainUpdateUserRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
