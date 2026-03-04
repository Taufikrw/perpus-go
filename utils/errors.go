package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Details any
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NewValidationError(message string, details any) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: message, Details: details}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func HandleError(c *gin.Context, err error) {
	var appError *AppError

	if errors.As(err, &appError) {
		SendErrorResponse(c, appError.Code, appError.Message, appError.Details)
		return
	}

	SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", nil)
}
