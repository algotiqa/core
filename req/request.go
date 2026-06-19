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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//=============================================================================

const MaxQueryLimit = 5000

//=============================================================================
//===
//=== Parameter retrieval
//===
//=============================================================================

func GetPagingParams(c *gin.Context) (offset int, limit int, errV error) {

	offset, err1 := GetParamAsInt(c, "offset", 0)

	if err1 != nil || offset < 0 {
		return 0, 0, NewBadRequestError("Invalid 'offset' param: %v", offset)
	}

	//--- Extract limit

	limit, err2 := GetParamAsInt(c, "limit", MaxQueryLimit)

	if err2 != nil || limit < 1 || limit > MaxQueryLimit {
			return 0, 0, NewBadRequestError("Invalid 'limit' param: %v", limit)
		}

	return offset, limit, nil
}

//=============================================================================

func GetParamAsBool(c *gin.Context, name string, defValue bool) (bool, error) {
	params := c.Request.URL.Query()

	if ! params.Has(name) {
		return defValue, nil
	}

	value := params.Get(name)

	res, err := strconv.ParseBool(value)

	if err == nil {
		return res, nil
	}

	return false, NewBadRequestError("Parameter '%v' has not a boolean value: %v", name, value)
}

//=============================================================================

func GetParamAsInt(c *gin.Context, name string, defValue int) (int, error) {
	params := c.Request.URL.Query()

	if ! params.Has(name) {
		return defValue, nil
	}

	value := params.Get(name)

	res, err := strconv.ParseInt(value, 10, 32)

	if err == nil {
		return int(res), nil
	}

	return 0, NewBadRequestError("Parameter '%v' has not an integer value: %v", name, value)
}

//=============================================================================

func GetParamAsInts(c *gin.Context, name string) ([]int64, error) {
	params := c.Request.URL.Query()
	values, ok := params[name]
	if !ok {
		return nil, nil
	}

	var res []int64

	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, NewBadRequestError("Parameter '%v' has not an integer value: %v", name, value)
		}

		res = append(res, v)
	}

	return res, nil
}

//=============================================================================

func GetParamAsString(c *gin.Context, name string, defValue string) string {
	params := c.Request.URL.Query()

	if ! params.Has(name) {
		return defValue
	}

	value := params.Get(name)

	if value == "" {
		return defValue
	}

	return value
}

//=============================================================================

func GetParamAsStrings(c *gin.Context, name string) []string  {
	params := c.Request.URL.Query()
	values, ok := params[name]
	if !ok {
		return nil
	}

	var res []string

	for _, value := range values {
		res = append(res, value)
	}

	return res
}

//=============================================================================

func BindParamsFromQuery(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		message := parseError(err)
		return NewBadRequestError(message)
	}

	return nil
}

//=============================================================================

func BindParamsFromBody(c *gin.Context, obj any) error {
	if err := c.ShouldBind(obj); err != nil {
		message := parseError(err)
		return NewBadRequestError(message)
	}

	return nil
}

//=============================================================================

func GetIdFromUrl(c *gin.Context) (uint, error) {
	sId := c.Param("id")
	iId, err := strconv.ParseInt(sId, 10, 64)

	if err != nil || iId<0 {
		return 0, NewBadRequestError("Invalid ID in url: %v", sId)
	}

	return uint(iId), nil
}

//=============================================================================

func GetIdsFromUrl(c *gin.Context) ([]uint, error) {
	params := c.Request.URL.Query()
	values, ok := params["id"]
	if !ok {
		return nil, nil
	}

	var res []uint

	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, NewBadRequestError("Parameter 'id' has not an integer value: %v", value)
		}

		res = append(res, uint(v))
	}

	return res, nil
}

//=============================================================================

func GetId2FromUrl(c *gin.Context) (uint, error) {
	sId := c.Param("id2")
	iId, err := strconv.ParseInt(sId, 10, 64)

	if err != nil || iId<0 {
		return 0, NewBadRequestError("Invalid ID in url: %v", sId)
	}

	return uint(iId), nil
}

//=============================================================================

type listResponse struct {
	Offset   int  `json:"offset"`
	Limit    int  `json:"limit"`
	Overflow bool `json:"overflow"`
	Result   any  `json:"result"`
}

//-----------------------------------------------------------------------------

func ReturnList(c *gin.Context, result any, offset int, limit int, size int) error {
	c.JSON(http.StatusOK, &listResponse{
		Offset:   offset,
		Limit:    limit,
		Overflow: size == MaxQueryLimit,
		Result:   result,
	})

	return nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func parseError(err error) string {
	switch typedError := any(err).(type) {
	case validator.ValidationErrors:
		for _, e := range typedError {
			return parseFieldError(e)
		}

	case *json.UnmarshalTypeError:
		return parseMarshallingError(*typedError)

	case *strconv.NumError:
		return parseConvertError(*typedError)
	}

	return err.Error()
}

//=============================================================================

func parseFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	fieldPrefix := fmt.Sprintf("The field %s", field)
	tag := strings.Split(e.Tag(), "|")[0]

	switch tag {
	case "required":
		return fmt.Sprintf("Missing the '%s' parameter", field)

	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())

	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)

	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)

	default:
		return fmt.Errorf("%v", e).Error()
	}
}

//=============================================================================

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("Invalid type: '%s' must be a %s", strings.ToLower(e.Field), e.Type.String())
}

//=============================================================================

func parseConvertError(e strconv.NumError) string {
	return fmt.Sprintf("Parameter must be an integer: %s", e.Num)
}

//=============================================================================
