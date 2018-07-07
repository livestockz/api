package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is a positive http response structure
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type PageSuccessResponse struct {
	SuccessResponse
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
	Total int32 `json:"total"`
}

// FailureResponse is a negative http response structure
type FailureResponse struct {
	Error []error `json:"error"`
}

//Ok writes http response with status code 200 and json object with `data` property
func Ok(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{data})
}

//Ok writes http response with status code 200 and json object with `data` property
func Page(c *gin.Context, data interface{}, page, limit, total int32) {
	c.JSON(200, PageSuccessResponse{data, page, limit, total})
}

//Created writes http response with status code 201 and json object with `data` property
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, SuccessResponse{data})
}

//NoContent writes http response with status code 204 and empty content
func NoContent(c *gin.Context) {
	c.JSON(204, nil)
}

//BadRequest writes http response with status code 400 and json object with `error` property
func BadRequest(c *gin.Context, errors ...error) {
	c.JSON(400, FailureResponse{errors})
}

//Unauthorized writes http response with status code 401 and json object with `error` property
func Unauthorized(c *gin.Context, messages ...string) {
	errors := make([]error, len(messages))
	for i, msg := range messages {
		errors[i] = errors.New(msg)
	}
	c.JSON(401, FailureResponse{errors})
}

//Forbidden writes http response with status code 403 and json object with `error` property
func Forbidden(c *gin.Context, messages ...string) {
	errors := make([]error, len(messages))
	for i, msg := range messages {
		errors[i] = errors.New(msg)
	}
	c.JSON(403, FailureResponse{errors})
}

//NotFound writes http response with status code 404 and json object with `error` property
func NotFound(c *gin.Context, messages ...string) {
	errors := make([]error, len(messages))
	for i, msg := range messages {
		errors[i] = errors.New(msg)
	}
	c.JSON(403, FailureResponse{errors})
}

//Error writes http response with status code 500 and json object with `error` property
func Error(c *gin.Context, errors ...error) {
	c.JSON(500, FailureResponse{errors})
}
