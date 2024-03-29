// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MainCreateDesignRequest main create design request
//
// swagger:model main.CreateDesignRequest
type MainCreateDesignRequest struct {

	// design
	Design string `json:"design,omitempty"`

	// fields
	Fields DesignAttrs `json:"fields,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// user Id
	UserID string `json:"userId,omitempty"`
}

// Validate validates this main create design request
func (m *MainCreateDesignRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this main create design request based on context it is used
func (m *MainCreateDesignRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MainCreateDesignRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MainCreateDesignRequest) UnmarshalBinary(b []byte) error {
	var res MainCreateDesignRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
