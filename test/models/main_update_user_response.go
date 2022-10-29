// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MainUpdateUserResponse main update user response
//
// swagger:model main.UpdateUserResponse
type MainUpdateUserResponse struct {

	// email
	Email string `json:"email,omitempty"`

	// first name
	FirstName string `json:"firstName,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// last name
	LastName string `json:"lastName,omitempty"`

	// updated at
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// Validate validates this main update user response
func (m *MainUpdateUserResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this main update user response based on context it is used
func (m *MainUpdateUserResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MainUpdateUserResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MainUpdateUserResponse) UnmarshalBinary(b []byte) error {
	var res MainUpdateUserResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}