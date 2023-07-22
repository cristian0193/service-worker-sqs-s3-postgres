package exceptions

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ErrNotFound error custom.
var (
	ErrNotFound           = errors.New("not found")
	ErrInternalError      = errors.New("internal error")
	ErrEntityAlreadyExist = errors.New("entity already exists")
	ErrInvalidEntity      = errors.New("invalid Entity")
)

// NewError returns a new controlled error.
func NewError(code int, err error) *echo.HTTPError {
	return echo.NewHTTPError(code, err.Error())
}

// HandleServiceError returns an error according to its type.
func HandleServiceError(err error) *echo.HTTPError {
	switch err.Error() {
	case ErrNotFound.Error():
		return NewError(http.StatusNotFound, err)
	case ErrInternalError.Error():
		return NewError(http.StatusInternalServerError, err)
	case ErrEntityAlreadyExist.Error():
		return NewError(http.StatusConflict, err)
	case ErrInvalidEntity.Error():
		return NewError(http.StatusUnprocessableEntity, err)
	default:
		return NewError(http.StatusInternalServerError, err)
	}
}
