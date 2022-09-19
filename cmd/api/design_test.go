package main

import (
	"github.com/google/uuid"
	"github.com/rengas/pdfgen/pkg/design"
	"reflect"
	"testing"
)

func TestGetDesignModel(t *testing.T) {

	tests := []struct {
		name       string
		input      CreateTemplateRequest
		want       design.Design
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "should return error when name is empty",
			input: CreateTemplateRequest{},

			wantErr:    true,
			wantErrMsg: "name is empty",
		},
		{
			name: "should return error when profileId is empty",
			input: CreateTemplateRequest{
				Name: "new template",
			},
			wantErr:    true,
			wantErrMsg: "profileId is empty",
		},
		{
			name: "should return error when fields are empty",
			input: CreateTemplateRequest{
				Name:      "new template",
				ProfileId: uuid.NewString(),
			},
			wantErr:    true,
			wantErrMsg: "fields are empty",
		},
		{
			name: "should return error when invalid field value types",
			input: CreateTemplateRequest{
				Name:      "new template",
				ProfileId: uuid.NewString(),
				Fields:    map[string]interface{}{"amount": struct{}{}},
			},
			wantErr:    true,
			wantErrMsg: "amount has unsupported type for value",
		},
		{
			name: "should return error when design is empty",
			input: CreateTemplateRequest{
				Name:      "new template",
				ProfileId: uuid.NewString(),
				Fields:    map[string]interface{}{"amount": 10.2},
			},
			wantErr:    true,
			wantErrMsg: "design is empty",
		},
		{
			name: "should return error when design is not a valid base64 string",
			input: CreateTemplateRequest{
				Name:      "new template",
				ProfileId: uuid.NewString(),
				Fields:    map[string]interface{}{"amount": 10.2},
				Design:    "XXXXXaGVsbG8=",
			},
			wantErr:    true,
			wantErrMsg: "invalid design",
		},
	}

	for _, tc := range tests {
		_, err := tc.input.GetDesignModel()
		if tc.wantErr == true {
			if !reflect.DeepEqual(tc.wantErrMsg, err.Error()) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErrMsg, err.Error())
			}
		}

	}
}
