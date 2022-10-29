package main

import (
	"errors"
	mkerror "github.com/rengas/pdfgen/pkg/errors"
	"github.com/rengas/pdfgen/pkg/httputils"
	"github.com/rengas/pdfgen/pkg/logging"
	"github.com/rengas/pdfgen/pkg/user"
	"net/http"
	"time"
)

type UserAPI struct {
	userRepo UserRepository
}

func NewUserAPI(userRepo UserRepository) *UserAPI {
	return &UserAPI{
		userRepo: userRepo,
	}
}

// GetUser func for getting user account.
// @Description  Get User Profile
// @Summary      Get User Profile
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  GetUserResponse
// @Failure      400  {object}  httputils.ErrorResponse  "Bad Request"
// @Failure      401  {object}  httputils.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  httputils.ErrorResponse  "Not Found"
// @Failure      500  {object}  httputils.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /user [get]
func (u *UserAPI) GetUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		logging.WithContext(ctx).Info("cannot find user id")
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}
	usr, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		logging.Error("cannot find user id")
		if errors.Is(err, user.ErrUserNotFound) {
			httputils.NotFound(ctx, w, mkerror.ErrNotFound)
		}
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}
	httputils.OK(ctx, w, GetUserResponseFromUser(usr))

}

func (u UpdateUserRequest) GetUser() user.User {
	usr := user.User{}
	if u.Email != "" {
		usr.Email = u.Email
	}

	if u.FirstName == "" {
		usr.FirstName = nil
	} else {
		usr.FirstName = &u.FirstName
	}

	if u.LastName == "" {
		usr.LastName = nil
	} else {
		usr.LastName = &u.LastName
	}

	usr.UpdatedAt = time.Now().UTC()

	return usr
}

// UpdateUser func for updating user account.
// @Description  Update User Profile
// @Summary      Update User Profile
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        data  body      UpdateUserRequest         true  "Body"
// @Success      200   {object}  UpdateUserResponse   "Updated"
// @Failure      400   {object}  httputils.ErrorResponse   "Bad Request"
// @Failure      401   {object}  httputils.ErrorResponse   "Unauthorized"
// @Failure      403   {object}  httputils.ErrorResponse   "Forbidden"
// @Failure      422   {object}  httputils.ErrorResponse   "Validation errors"
// @Failure      500   {object}  httputils.ErrorResponse   "Internal Server Error"
// @Security     BearerAuth
// @Router       /user [patch]
func (u *UserAPI) UpdateUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		logging.Debug("unable to get userId from context")
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}

	usr, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			httputils.NotFound(ctx, w, mkerror.ErrNotFound)
		}
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}

	var ur UpdateUserRequest
	err = httputils.ReadJson(req, &ur)
	if err != nil {
		httputils.BadRequest(ctx, w, err)
		return
	}

	err = ur.Validate()
	if err != nil {
		httputils.UnProcessableEntity(ctx, w, err)
		return
	}

	mUser := ur.GetUser()
	usrEmail, err := u.userRepo.GetByEmail(ctx, ur.Email)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		httputils.UnAuthorized(ctx, w, err)
		return
	}

	if err == nil && usr.Id != usrEmail.Id {
		httputils.Conflict(ctx, w, ErrUserWithEmailExists)
		return
	}

	updateAt := time.Now().UTC()
	updateUser := user.User{
		Id:        userId,
		FirstName: mUser.FirstName,
		LastName:  mUser.LastName,
		Email:     mUser.Email,
		UpdatedAt: updateAt,
	}

	err = u.userRepo.Update(ctx, updateUser)
	if err != nil {
		httputils.InternalServerError(ctx, w, mkerror.ErrInternalError)
		return
	}
	user := &UpdateUserResponse{
		Id:        updateUser.Id,
		Email:     updateUser.Email,
		UpdatedAt: updateAt,
	}

	if mUser.FirstName != nil {
		user.FirstName = *mUser.FirstName
	}
	if mUser.LastName != nil {
		user.LastName = *mUser.LastName
	}

	httputils.OK(ctx, w, user)

}
