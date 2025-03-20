package helpers

import (
	"fmt"
	"net/http"
	"strings"
)

var ErrAppBadID = NewAppError(http.StatusBadRequest, "Invalid ID.", nil)

func ErrAppGeneric(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "Something went wrong.", err)
}

func NewAppError(status int, message string, e error) *AppError {
	return &AppError{
		status:  status,
		message: message,
		err:     e,
	}
}

type AppError struct {
	status  int
	message string
	err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", strings.ToLower(e.message), e.err)
}
func (e *AppError) HTTPStatus() int {
	return e.status
}
func (e *AppError) Message() string {
	return e.message
}
func (e *AppError) Inner() error {
	return e.err
}
