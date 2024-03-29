// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MainValidateDesignResponse main validate design response
//
// swagger:model main.ValidateDesignResponse
type MainValidateDesignResponse struct {

	// id
	// Example: 99d15987-e06f-492c-a520-e54185e5b80b
	ID string `json:"id,omitempty"`
}

// Validate validates this main validate design response
func (m *MainValidateDesignResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this main validate design response based on context it is used
func (m *MainValidateDesignResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MainValidateDesignResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MainValidateDesignResponse) UnmarshalBinary(b []byte) error {
	var res MainValidateDesignResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
