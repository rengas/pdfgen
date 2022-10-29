package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	mkerror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"github.com/rengas/pdfgen/pkg/logging"
	"github.com/rengas/pdfgen/pkg/user"
	"net/http"
	"time"
)

type AuthAPI struct {
	userRepo UserRepository
	bcrypt   Bcrypt
	jwt      JWTToken
}

func NewAuthAPI(userRepo UserRepository,
	bcrypt Bcrypt,
	jwt JWTToken) *AuthAPI {
	return &AuthAPI{
		userRepo: userRepo,
		bcrypt:   bcrypt,
		jwt:      jwt,
	}
}

// Register func for register.
// @Description  Register a new user.
// @Summary      Register
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        register_details  body  RegisterRequest  true  "register details"
// @Success      200           {object}  RegisterResponse
// @Failure      400           {object}  httputils.ErrorResponse "Bad Request"
// @Failure      422           {object}  httputils.ErrorResponse "Validation errors"
// @Failure      500           {object}  httputils.ErrorResponse  "Internal Server Error"
// @Router       /register [post]
func (a *AuthAPI) Register(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var rr RegisterRequest
	err := httputils.ReadJson(req, &rr)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(ctx, w, err)
		return
	}

	err = rr.Validate()
	if err != nil {
		logging.WithField(logging.Field{Label: "validation", Value: err.Error()}).Debug("validation failed")
		httputils.UnProcessableEntity(ctx, w, err)
		return
	}

	u, err := a.userRepo.GetByEmail(ctx, rr.Email)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}
	if u.Id != "" {
		httputils.Conflict(ctx, w, ErrAuthEmailExists)
		return
	}

	b, err := a.bcrypt.GetHashedPassword(rr.Password)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}

	id := uuid.NewString()
	us := user.User{
		Id:           id,
		Email:        rr.Email,
		PasswordHash: string(b),
		Role:         user.RoleNormal,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err = a.userRepo.SaveNewUser(ctx, us)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}

	httputils.OK(ctx, w, &RegisterResponse{Id: id})
}

// Login func for login.
// @Description  Login as a user.
// @Summary      Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data  body      LoginRequest               true  "Credentials"
// @Success      200   {object}  LoginResponse              "User account and token pair"
// @Failure      400   {object}  httputils.ErrorResponse  "		 "Bad Request"
// @Failure      422   {object}  httputils.ErrorResponse         "Validation errors"
// @Failure      403   {object}  httputils.ErrorResponse         "Forbidden"
// @Router       /login [post]
func (a *AuthAPI) Login(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var lr LoginRequest
	err := httputils.ReadJson(req, &lr)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.BadRequest(ctx, w, err)
		return
	}
	err = lr.Validate()
	if err != nil {
		logging.WithField(logging.Field{Label: "validation", Value: err.Error()}).Debug("validation failed")
		httputils.UnProcessableEntity(ctx, w, err)
		return
	}

	u, err := a.userRepo.GetByEmail(ctx, lr.Email)
	if err != nil && errors.Is(err, user.ErrUserNotFound) {
		logging.WithContext(ctx).WithError(err)
		httputils.UnAuthorized(ctx, w, err)
		return
	}

	err = a.bcrypt.CompareHashedPassword(u.PasswordHash, lr.Password)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.Forbidden(ctx, w, mkerror.ErrForbidden)
		return
	}

	claims := make(map[string]interface{}, 0)
	claims["userId"] = u.Id
	token, err := a.jwt.TokePair(claims)
	if err != nil {
		logging.WithContext(ctx).WithError(err)
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}

	rs := LoginResponse{
		User: User{
			Id:    u.Id,
			Email: u.Email,
		},
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	httputils.OK(ctx, w, rs)
}

func (p *AuthAPI) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I'm ok")
}
