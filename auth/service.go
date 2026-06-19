//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package auth

import (
	"log/slog"
	"net/http"

	"github.com/algotiqa/core/req"
	"github.com/gin-gonic/gin"
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

func (c *Context) GetParamAsInts(name string) ([]int64, error) {
	return req.GetParamAsInts(c.Gin, name)
}

//=============================================================================

func (c *Context) GetParamAsString(name string, defValue string) string {
	return req.GetParamAsString(c.Gin, name, defValue)
}

//=============================================================================

func (c *Context) GetParamAsStrings(name string) []string {
	return req.GetParamAsStrings(c.Gin, name)
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

func (c *Context) GetIdsFromUrl() ([]uint, error) {
	return req.GetIdsFromUrl(c.Gin)
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
