//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package req

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

//=============================================================================

type AppError struct {
	Code    int
	Message string
}

//-----------------------------------------------------------------------------

func (e AppError) Error() string {
	return e.Message
}

//=============================================================================

func NewBadRequestError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusBadRequest,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewForbiddenError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusForbidden,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewNotFoundError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusNotFound,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewUnprocessableEntityError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusUnprocessableEntity,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewServerError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusInternalServerError,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewServiceUnavailableError(message string, params ...any) error {
	return AppError {
		Code:    http.StatusServiceUnavailable,
		Message: sprintf(message, params),
	}
}

//=============================================================================

func NewServerErrorByError(err error) error {
	if err == nil {
		return nil
	}

	return AppError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

//=============================================================================

func ReturnUnauthorizedError(c *gin.Context, message string) {
	writeError(c, http.StatusUnauthorized, message)
}

//=============================================================================

func ReturnForbiddenError(c *gin.Context, message string) {
	writeError(c, http.StatusForbidden, message)
}

//=============================================================================

func ReturnError(c *gin.Context, err error) {
	if err != nil {
		var ae AppError
		if errors.As(err, &ae) {
			writeError(c, ae.Code, ae.Message)
		} else {
			writeError(c, http.StatusInternalServerError, "Found non AppError object : "+ err.Error())
		}
	}
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

type errorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
}

//-----------------------------------------------------------------------------

func writeError(c *gin.Context, errorCode int, errorMessage string) {

	slog.Error(errorMessage,
		"client", c.ClientIP(),
		"code", errorCode)

	c.JSON(errorCode, &errorResponse{
		Code:    errorCode,
		Error:   errorMessage,
	})
}

//=============================================================================

func sprintf(message string, params ...any) string {
	if params != nil {
		message = fmt.Sprintf(message, params...)
	}

	return message
}

//=============================================================================
