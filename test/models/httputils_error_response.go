// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HttputilsErrorResponse httputils error response
//
// swagger:model httputils.ErrorResponse
type HttputilsErrorResponse struct {

	// error
	Error interface{} `json:"error,omitempty"`
}

// Validate validates this httputils error response
func (m *HttputilsErrorResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this httputils error response based on context it is used
func (m *HttputilsErrorResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HttputilsErrorResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HttputilsErrorResponse) UnmarshalBinary(b []byte) error {
	var res HttputilsErrorResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}