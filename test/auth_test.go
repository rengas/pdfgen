//go:build e2e

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rengas/pdfgen/test/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	email := Email()

	tests := []struct {
		Name       string                         `json:"name"`
		Request    models.MainRegisterRequest     `json:"request"`
		WantErr    *models.HttputilsErrorResponse `json:"err"`
		StatusCode interface{}                    `json:"statusCode"`
	}{
		{
			Name:       "Register new user, without email should return error",
			Request:    models.MainRegisterRequest{},
			WantErr:    &models.HttputilsErrorResponse{Error: "email is empty"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:       "Register new user, without password should return error",
			Request:    models.MainRegisterRequest{Email: StringPtr(Email())},
			WantErr:    &models.HttputilsErrorResponse{Error: "password is empty"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name: "Register new user, with invalid password password length should return error",
			Request: models.MainRegisterRequest{
				Email:    StringPtr(Email()),
				Password: StringPtr(Password(true, true, true, true, false, 2)),
			},
			WantErr:    &models.HttputilsErrorResponse{Error: "password is less than 8 characters"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name: "Register new user, with valid input should successfully register",
			Request: models.MainRegisterRequest{
				Email:    StringPtr(Email()),
				Password: StringPtr(Password(true, true, true, true, false, 8)),
			},
			WantErr:    nil,
			StatusCode: http.StatusOK,
		},
		{
			Name: "Register new user, with existing email  should thrown an error",
			Request: models.MainRegisterRequest{
				Email:    StringPtr(email),
				Password: StringPtr(Password(true, true, true, true, false, 8)),
			},
			WantErr:    &models.HttputilsErrorResponse{Error: "user with email exists"},
			StatusCode: http.StatusConflict,
		},
	}

	for _, ts := range tests {

		t.Run(ts.Name, func(t *testing.T) {

			b, err := json.Marshal(ts.Request)
			require.NoError(t, err)

			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/auth/register", httpPort), "application/json", bytes.NewReader(b))
			require.NoError(t, err)
			require.Equal(t, ts.StatusCode, resp.StatusCode)
			if ts.WantErr != nil {

				var e models.HttputilsErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}

			if ts.WantErr == nil {
				var rs models.MainRegisterResponse
				err = json.NewDecoder(resp.Body).Decode(&rs)
				require.NoError(t, err)
				require.NotEmpty(t, rs.ID)
			}
		})

	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
	email := StringPtr(Email())
	password := StringPtr(Password(true, true, true, true, false, 8))

	b, err := json.Marshal(models.MainRegisterRequest{Email: StringPtr(Email()), Password: password})
	require.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/auth/register", httpPort), "application/json", bytes.NewReader(b))
	require.NoError(t, err)

	var rs models.MainRegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&rs)

	tests := []struct {
		Name       string                         `json:"name"`
		Request    models.MainLoginRequest        `json:"request"`
		WantErr    *models.HttputilsErrorResponse `json:"err"`
		StatusCode interface{}                    `json:"statusCode"`
	}{
		{
			Name:       "Login, without email should return error",
			Request:    models.MainLoginRequest{},
			WantErr:    &models.HttputilsErrorResponse{Error: "email is empty"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:       "Login new user, without password should return error",
			Request:    models.MainLoginRequest{Email: email},
			WantErr:    &models.HttputilsErrorResponse{Error: "password is empty"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name: "Login  user, with valid input should successfully register",
			Request: models.MainLoginRequest{
				Email:    email,
				Password: password,
			},
			WantErr:    nil,
			StatusCode: http.StatusOK,
		},
	}

	for _, ts := range tests {

		t.Run(ts.Name, func(t *testing.T) {

			b, err := json.Marshal(ts.Request)
			require.NoError(t, err)

			resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/auth/login", httpPort), "application/json", bytes.NewReader(b))
			require.NoError(t, err)

			require.Equal(t, ts.StatusCode, resp.StatusCode)

			if ts.WantErr != nil {

				var e models.HttputilsErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}

			if ts.WantErr == nil {
				var ls models.MainLoginResponse
				err = json.NewDecoder(resp.Body).Decode(&ls)
				require.NoError(t, err)
				require.Equal(t, rs.ID, ls.User.ID)
				require.NotEmpty(t, ls.AccessToken)
				require.NotEmpty(t, ls.RefreshToken)
			}
		})
	}
}
