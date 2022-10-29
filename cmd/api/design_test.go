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
		input      CreateDesignRequest
		want       design.Design
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "should return error when name is empty",
			input: CreateDesignRequest{},

			wantErr:    true,
			wantErrMsg: "name is empty",
		},
		{
			name: "should return error when profileId is empty",
			input: CreateDesignRequest{
				Name: "new Design",
			},
			wantErr:    true,
			wantErrMsg: "profileId is empty",
		},
		{
			name: "should return error when design is empty",
			input: CreateDesignRequest{
				Name:   "new Design",
				UserId: uuid.NewString(),
				Fields: map[string]interface{}{"amount": 10.2},
			},
			wantErr:    true,
			wantErrMsg: "design is empty",
		},
		{
			name: "should return error when invalid field value types",
			input: CreateDesignRequest{
				Name:   "new Design",
				UserId: uuid.NewString(),
				Fields: map[string]interface{}{"amount": struct{}{}},
				Design: "PCFET0NUWVBFIGh0bWw+CjxodG1sPgogICA8aGVhZD4KICAgICAgPHRpdGxlPnt7LmFtb3VudH19PC90aXRsZT4KICAgPC9oZWFkPgogICA8Ym9keT4KICAgICAgPGgxPnt7LmFtb3VudH19IDwvaDE+CiAgICAgIDxoMT57ey5uYW1lfX0gPC9oMT4KICAgICAgPGgxPnt7LmFkZHJlc3N9fSA8L2gxPgogICAgICA8dWwgPgogICAgICAgICB7e3JhbmdlICRpLCAkYSA6PSAuaXRlbXN9fQogICAgICAgICA8bGk+e3skYX19PC9saT4KICAgICAgICAge3tlbmR9fQogICAgICA8L3VsPgogICAgICA8dWwgPgogICAgICAgICB7e3JhbmdlICRpLCAkYSA6PSAuaXRlbU1hcH19CiAgICAgICAgIDxsaT57eyRhfX08L2xpPgogICAgICAgICB7e2VuZH19CiAgICAgIDwvdWw+CiAgIDwvYm9keT4KPC9odG1sPg=="},
			wantErr:    true,
			wantErrMsg: "amount has unsupported type for value",
		},
	}

	for _, tc := range tests {
		err := tc.input.Validate()
		if tc.wantErr == true {
			if !reflect.DeepEqual(tc.wantErrMsg, err.Error()) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErrMsg, err.Error())
			}
		}

	}
}
