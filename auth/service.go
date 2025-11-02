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

package auth

import (
	"github.com/tradalia/core/req"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

//=============================================================================

type RestService func(c *Context)

//=============================================================================

type Context struct {
	Gin     *gin.Context
	Session *UserSession
	Log     *slog.Logger
	Config  any
	Token   string
}

//=============================================================================

func (c *Context) ReturnError(err error) {
	req.ReturnError(c.Gin, err)
}

//=============================================================================

func (c *Context) ReturnList(result any, offset int, limit int, size int) error {
	return req.ReturnList(c.Gin, result, offset, limit, size)
}

//=============================================================================

func (c *Context) ReturnObject(data any) error {
	c.Gin.JSON(http.StatusOK, data)
	return nil
}

//=============================================================================

func (c *Context) ReturnData(contentType string, data []byte) error {
	c.Gin.Data(http.StatusOK, contentType, data)
	return nil
}

//=============================================================================

func (c *Context) GetPagingParams() (offset int, limit int, errV error) {
	return req.GetPagingParams(c.Gin)
}

//=============================================================================

func (c *Context) GetParamAsBool(name string, defValue bool) (bool, error) {
	return req.GetParamAsBool(c.Gin, name, defValue)
}

//=============================================================================

func (c *Context) GetParamAsInt(name string, defValue int) (int, error) {
	return req.GetParamAsInt(c.Gin, name, defValue)
}

//=============================================================================

func (c *Context) GetParamAsString(name string, defValue string) string {
	return req.GetParamAsString(c.Gin, name, defValue)
}

//=============================================================================

func (c *Context) BindParamsFromQuery(obj any) (err error) {
	return req.BindParamsFromQuery(c.Gin, obj)
}

//=============================================================================

func (c *Context) BindParamsFromBody(obj any) (err error) {
	return req.BindParamsFromBody(c.Gin, obj)
}

//=============================================================================

func (c *Context) GetIdFromUrl() (uint, error) {
	return req.GetIdFromUrl(c.Gin)
}

//=============================================================================

func (c *Context) GetId2FromUrl() (uint, error) {
	return req.GetId2FromUrl(c.Gin)
}

//=============================================================================

func (c *Context) GetCodeFromUrl() string {
	return c.Gin.Param("code")
}

//=============================================================================
