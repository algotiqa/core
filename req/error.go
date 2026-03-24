//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
