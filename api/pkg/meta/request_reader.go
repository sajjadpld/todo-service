package meta

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"microservice/pkg/validator"
)

type IRequestConvertible[D any] interface {
	ToDomain() *D
}

// ReqBodyToDomain binds the request body and evaluates that by `GO Validator`
func ReqBodyToDomain[R IRequestConvertible[D], D any](c *gin.Context) (entity *D, err error) {
	var body R
	if err = c.ShouldBind(&body); err != nil {
		return
	}

	if err = validator.ValidateRequestDto(c.Request.Context(), body); err != nil {
		return
	}

	entity = body.ToDomain()
	return
}

// ReqRouteParamsToDomain binds the request body and evaluates that by `GO Validator`
func ReqRouteParamsToDomain[R IRequestConvertible[D], D any](c *gin.Context) (entity *D, err error) {
	var params R

	// Create a new instance of the params struct
	params = reflect.New(reflect.TypeOf(params).Elem()).Interface().(R)
	bindRouteParams(c, params)

	// Validate the parameters
	if err = validator.ValidateRequestDto(c.Request.Context(), params); err != nil {
		return
	}

	// Convert to domain entity
	entity = params.ToDomain()
	return
}

func ReqQryParamToDomain[R IRequestConvertible[D], D any](c *gin.Context) (entity *D, err error) {
	var qry R
	qry = reflect.New(reflect.TypeOf(qry).Elem()).Interface().(R)

	if err = c.ShouldBindQuery(qry); err != nil {
		return
	}

	if err = validator.ValidateRequestDto(c.Request.Context(), qry); err != nil {
		return
	}

	entity = qry.ToDomain()
	return
}

func ReqHeaderToDomain[T, K any](c *gin.Context, fn func(*T) *K) (entity *K, err error) {
	var headers T

	err = mapHeadersToStruct(c, &headers)
	if err != nil {
		return
	}

	if err = validator.ValidateRequestDto(c.Request.Context(), headers); err != nil {
		return
	}

	entity = fn(&headers)
	return
}

// HELPERS

// bindRouteParams to bind route parameters to a struct
func bindRouteParams(c *gin.Context, target interface{}) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	// Create case-insensitive map of route parameters
	paramMap := make(map[string]string)
	for _, p := range c.Params {
		paramMap[strings.ToLower(p.Key)] = p.Value
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Get param tag or use lowercase field name
		paramName := field.Tag.Get("param")
		if paramName == "" {
			paramName = strings.ToLower(field.Name)
		}

		// Case-insensitive match
		if value, exists := paramMap[strings.ToLower(paramName)]; exists {
			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(value)
			}
		}
	}
}

// mapHeadersToStruct: match and map the headers' values to the referenced struct JSON tags
func mapHeadersToStruct(c *gin.Context, dest interface{}) error {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		headerValue := c.GetHeader(jsonTag)

		// Handle Bearer token
		if strings.HasPrefix(headerValue, "Bearer ") {
			headerValue = strings.TrimPrefix(headerValue, "Bearer ")
		} else if strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
			headerValue = strings.TrimPrefix(strings.ToLower(headerValue), "bearer ")
		}

		if headerValue != "" {
			v.Field(i).SetString(headerValue)
		}
	}
	return nil
}
