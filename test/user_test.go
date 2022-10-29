//go:build e2e

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/rengas/pdfgen/test/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name       string                         `json:"name"`
		WantErr    *models.HttputilsErrorResponse `json:"err"`
		StatusCode interface{}                    `json:"HTTPStatusCode"`
		Token      LoginFunc
	}{
		{
			Name:       "Get user with expired token",
			StatusCode: http.StatusOK,
			Token:      validLoginToken,
		},
	}

	for _, ts := range tests {

		t.Run(ts.Name, func(t *testing.T) {
			tkn := ts.Token()

			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/api/v1/user/me", httpPort), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tkn.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, ts.StatusCode, resp.StatusCode)

			if ts.WantErr != nil {
				var e *models.HttputilsErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}
			if ts.WantErr == nil {
				var ls *models.MainGetUserResponse
				err = json.NewDecoder(resp.Body).Decode(&ls)
				require.NoError(t, err)
				require.Equal(t, tkn.User.ID, ls.ID)
				require.Equal(t, tkn.User.Email, ls.Email)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	//TODO improve this test
	t.Parallel()
	email := StringPtr(gofakeit.Email())
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	tkn := validLoginToken()

	existingEmail := validLoginToken()

	tests := []struct {
		Name       string                         `json:"name"`
		Request    models.MainUpdateUserRequest   `json:"request"`
		Response   models.MainUpdateUserResponse  `json:"response"`
		WantErr    *models.HttputilsErrorResponse `json:"err"`
		StatusCode interface{}                    `json:"statusCode"`
	}{
		{
			Name: "Update user, without email",
			Request: models.MainUpdateUserRequest{
				Email: StringPtr(""),
			},
			WantErr:    &models.HttputilsErrorResponse{Error: "email is empty"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name: "Update user, update only first name",
			Request: models.MainUpdateUserRequest{
				Email:     StringPtr(tkn.User.Email),
				FirstName: firstName,
			},
			Response: models.MainUpdateUserResponse{
				Email:     tkn.User.Email,
				FirstName: firstName,
				LastName:  "",
			},
			StatusCode: http.StatusOK,
		},
		{
			Name: "Update user, update only last name",
			Request: models.MainUpdateUserRequest{
				Email:    StringPtr(tkn.User.Email),
				LastName: lastName,
			},
			Response: models.MainUpdateUserResponse{
				Email:     tkn.User.Email,
				FirstName: firstName,
				LastName:  lastName,
			},
			StatusCode: http.StatusOK,
		},
		{
			Name: "Update user, update both first name and last name",
			Request: models.MainUpdateUserRequest{
				Email:     StringPtr(tkn.User.Email),
				FirstName: firstName + "new",
				LastName:  lastName + "new",
			},
			Response: models.MainUpdateUserResponse{
				Email:     tkn.User.Email,
				FirstName: firstName + "new",
				LastName:  lastName + "new",
			},
			StatusCode: http.StatusOK,
		},
		{
			Name: "Update user, With email from a existing user",
			Request: models.MainUpdateUserRequest{
				Email:     StringPtr(existingEmail.User.Email),
				FirstName: firstName,
				LastName:  lastName,
			},
			WantErr:    &models.HttputilsErrorResponse{Error: "user with this email exists"},
			StatusCode: http.StatusConflict,
		},
		{
			Name: "Update user, With all details",
			Request: models.MainUpdateUserRequest{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			},
			Response: models.MainUpdateUserResponse{
				Email:     fromStringPtr(email),
				FirstName: firstName,
				LastName:  lastName,
			},
			StatusCode: http.StatusOK,
		},
	}

	for _, ts := range tests {
		t.Run(ts.Name, func(t *testing.T) {
			b, err := json.Marshal(ts.Request)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%s/api/v1/user/me", httpPort), bytes.NewReader(b))
			require.NoError(t, err)
			req.Header.Set("Authorization", tkn.AccessToken)
			client := &http.Client{
				Transport: &http.Transport{},
			}
			resp, err := client.Do(req)

			require.NoError(t, err)
			require.Equal(t, ts.StatusCode, resp.StatusCode)

			if ts.WantErr != nil {
				var e models.HttputilsErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}
			if ts.WantErr == nil {
				var ls models.MainUpdateUserResponse
				err = json.NewDecoder(resp.Body).Decode(&ls)
				require.NoError(t, err)
				require.Equal(t, ts.Request.FirstName, ls.FirstName)
				require.Equal(t, ts.Request.LastName, ls.LastName)
				require.Equal(t, tkn.User.ID, ls.ID)
				require.Equal(t, ts.Request.Email, ls.Email)
			}
		})
	}
}
