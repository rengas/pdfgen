package errors

const (
	ErrUnableToGetUserIdFromContext ValidationError = "unable to get userId from context"
	ErrNotAdminUser                 ValidationError = "not an admin user"
	ErrAdminUserNotFound            ValidationError = "admin user not found"
)

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

const ErrInternalError InternalError = "Sorry! Something is broken"

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

const ErrUnAuthorizedError UnAuthorizedError = "not authorised"

type UnAuthorizedError string

func (e UnAuthorizedError) Error() string {
	return string(e)
}

const ErrForbidden ForbiddenErr = "invalid credentials"

type ForbiddenErr string

func (e ForbiddenErr) Error() string {
	return string(e)
}

const ErrNotFound NotFoundError = "not found"

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}
