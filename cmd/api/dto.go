package main

import (
	"errors"
	"fmt"
	"github.com/rengas/pdfgen/pkg/design"
	pgerrror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/pagination"
	"github.com/rengas/pdfgen/pkg/user"
	"reflect"
	"time"
)

const (
	ErrAuthPasswordIsEmpty               pgerrror.ValidationError = "password is empty"
	ErrAuthPasswordInvalidLength         pgerrror.ValidationError = "password is less than 8 characters"
	ErrAuthEmailIsEmpty                  pgerrror.ValidationError = "email is empty"
	ErrAuthEmailExists                   pgerrror.ValidationError = "user with email exists"
	ErrUserEmailIsEmpty                  pgerrror.ValidationError = "email is empty"
	ErrUserWithEmailExists               pgerrror.ValidationError = "user with this email exists"
	ErrDesignNameIsEmpty                 pgerrror.ValidationError = "name is empty"
	ErrDesignDesignIsEmpty               pgerrror.ValidationError = "design is empty"
	ErrDesignUserIdIsEmpty               pgerrror.ValidationError = "userId is empty"
	ErrDesignUnsupportedFieldType        pgerrror.ValidationError = "unsupported field type for field"
	ErrDesignDesignIdIsEmpty             pgerrror.ValidationError = "designId is empty"
	ErrDesignCountIsEmpty                pgerrror.ValidationError = "count is empty"
	ErrDesignCountInvalid                pgerrror.ValidationError = "count is invalid"
	ErrDesignPageIsEmpty                 pgerrror.ValidationError = "page is empty"
	ErrDesignPageInvalid                 pgerrror.ValidationError = "page is invalid"
	ErrDesignUnableToGetDesigns          pgerrror.ValidationError = "unable to get designs"
	ErrDesignMustBeBase64Encoded         pgerrror.ValidationError = "design must be base64 encoded"
	ErrDesignInvalidHTML                 pgerrror.ValidationError = "invalid html design"
	ErrDesignUnableToMinify              pgerrror.ValidationError = "unable to minify design"
	ErrDesignUnableToUpdate              pgerrror.ValidationError = "unable to update design"
	ErrDesignUnableToReadRequest         pgerrror.ValidationError = "unable to read request body"
	ErrDesignUnableToGetDesign           pgerrror.ValidationError = "unable to get design"
	ErrDesignUnableToParseDesign         pgerrror.ValidationError = "unable to parse design"
	ErrDesignUnableToMatchFieldsToDesign pgerrror.ValidationError = "unable to match fields to design"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required" example:"John@email.com" `
	Password string `json:"password" example:"your password" validate:"required"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" {
		return ErrAuthEmailIsEmpty
	}
	if r.Password == "" {
		return ErrAuthPasswordIsEmpty
	}
	return nil
}

type User struct {
	Id    string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
	Email string `json:"email" example:"John@email.com"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken" example:"JWT token format"`
	RefreshToken string `json:"refreshToken"  example:"JWT token format"`
}

type RegisterRequest struct {
	Email    string `json:"email"  validate:"required" example:"John@email.com"`
	Password string `json:"password" minLength:"8"  validate:"required" example:"random_string"`
}

func (r RegisterRequest) Validate() error {
	if r.Email == "" {
		return ErrAuthEmailIsEmpty
	}
	if r.Password == "" {
		return ErrAuthPasswordIsEmpty
	}
	if len(r.Password) < 8 {
		return ErrAuthPasswordInvalidLength
	}
	return nil
}

type RegisterResponse struct {
	Id string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
}

type ValidateDesignRequest struct {
	Name   string       `json:"name"`
	Design string       `json:"design"`
	Fields design.Attrs `json:"fields"`
}

func (c ValidateDesignRequest) Validate() error {

	if c.Name == "" {
		return ErrDesignNameIsEmpty
	}

	if c.Design == "" {
		return ErrDesignUserIdIsEmpty
	}

	if c.Fields != nil {
		for _, v := range c.Fields {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Slice,
				reflect.Array, reflect.Map:
				continue
			default:
				return ErrDesignUnsupportedFieldType
			}
		}
	}

	return nil
}

type ValidateDesignResponse struct {
	Message string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
}

type CreateDesignRequest struct {
	Name   string       `json:"name"`
	UserId string       `json:"userId"`
	Design string       `json:"design"`
	Fields design.Attrs `json:"fields"`
}

func (c CreateDesignRequest) Validate() error {

	if c.Name == "" {
		return ErrDesignNameIsEmpty
	}

	if c.UserId == "" {
		return ErrDesignUserIdIsEmpty
	}

	if c.Design == "" {
		return ErrDesignDesignIsEmpty
	}

	if c.Fields != nil {
		for _, v := range c.Fields {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Slice,
				reflect.Array, reflect.Map:
				continue
			default:
				return ErrDesignUnsupportedFieldType
			}
		}
	}

	return nil
}

type CreateDesignResponse struct {
	Id string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
}

type UpdateDesignRequest struct {
	Name   string       `json:"name"`
	Design string       `json:"design"`
	Fields design.Attrs `json:"fields"`
}

func (c UpdateDesignRequest) Validate() error {

	if c.Name == "" {
		return errors.New("name is empty")
	}

	if c.Design == "" {
		return errors.New("design is empty")
	}

	if c.Fields != nil {
		for k, v := range c.Fields {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Slice,
				reflect.Array, reflect.Map:
				continue
			default:
				return errors.New(fmt.Sprintf("%s has unsupported type for value", k))
			}
		}
	}

	return nil
}

type UpdateDesignResponse struct {
	Id string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
}

type GetDesignResponse design.Design

type GetUserResponse struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetUserResponseFromUser(u user.User) GetUserResponse {
	usr := GetUserResponse{
		Id:        u.Id,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.FirstName != nil {
		usr.FirstName = *u.FirstName
	}

	if u.LastName != nil {
		usr.LastName = *u.LastName
	}

	return usr
}

type ListDesignResponse struct {
	Designs    []design.Design       `json:"designs,omitempty"`
	Pagination pagination.Pagination `json:"pagination"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"required" example:"John.doe@email.com"`
	FirstName string `json:"firstName"  example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}

func (u UpdateUserRequest) Validate() error {
	if u.Email == "" {
		return ErrUserEmailIsEmpty
	}
	return nil
}

type UpdateUserResponse struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteDesignResponse struct {
	Id string `json:"id" example:"99d15987-e06f-492c-a520-e54185e5b80b"`
}

type GeneratePDFRequest struct {
	DesignId string       `json:"DesignId"`
	Fields   design.Attrs `json:"fields"`
}

func (g GeneratePDFRequest) Validate() error {
	if g.DesignId == "" {
		return ErrDesignDesignIdIsEmpty
	}

	if g.Fields != nil {
		for _, v := range g.Fields {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Slice,
				reflect.Array, reflect.Map:
				continue
			default:
				return ErrDesignUnsupportedFieldType
			}
		}
	}

	return nil
}
