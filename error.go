package err

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func Error(errM ...*ErrorMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastError := c.Errors.Last()
		if lastError == nil {
			return
		}

		for _, err := range errM {
			if err.matchError(lastError.Err) {
				err.response(c)
			}
		}
	}
}

type ErrorMap struct {
	errors []error

	statusCode int
	response   func(c *gin.Context)
}

func (e *ErrorMap) StatusCode(statusCode int) *ErrorMap {
	e.statusCode = statusCode
	e.response = func(c *gin.Context) {
		c.Status(statusCode)
	}
	return e
}

func (e *ErrorMap) Response(response func(c *gin.Context)) *ErrorMap {
	e.response = response
	return e
}

func (e *ErrorMap) matchError(actual error) bool {
	for _, expected := range e.errors {
		if errors.Is(actual, expected) {
			return true
		}
	}
	return false
}

func NewErrMap(err ...error) *ErrorMap {
	return &ErrorMap{errors: err}
}
