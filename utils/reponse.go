package utils

import (
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
	Error []string `json:"error"`
}

//Ok writes http response with status code 200 and json object with `data` property
func Ok(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{data})
}

//Ok writes http response with status code 200 and json object with `data` property
func Page(c *gin.Context, data interface{}, page, limit, total int32) {
	c.JSON(200, &PageSuccessResponse{SuccessResponse{data}, page, limit, total})
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
	msg := make([]string, len(errors))
	for i, err := range errors {
		msg[i] = err.Error()
	}
	c.JSON(400, FailureResponse{msg})
}

//Unauthorized writes http response with status code 401 and json object with `error` property
func Unauthorized(c *gin.Context, messages ...string) {
	c.JSON(401, FailureResponse{messages})
}

//Forbidden writes http response with status code 403 and json object with `error` property
func Forbidden(c *gin.Context, messages ...string) {
	c.JSON(403, FailureResponse{messages})
}

//NotFound writes http response with status code 404 and json object with `error` property
func NotFound(c *gin.Context, messages ...string) {
	c.JSON(403, FailureResponse{messages})
}

//Error writes http response with status code 500 and json object with `error` property
func Error(c *gin.Context, errors ...error) {
	msg := make([]string, len(errors))
	for i, err := range errors {
		msg[i] = err.Error()
	}
	c.JSON(500, FailureResponse{msg})
}
