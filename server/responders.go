package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// jsonError renders an error response
func jsonError(c *gin.Context, status int, err interface{}) {
	var message interface{}

	switch v := err.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = v
	}

	c.AbortWithStatusJSON(status, gin.H{
		"status": status,
		"error":  message,
	})
}

// badRequest renders a HTTP 400 bad request response
func badRequest(c *gin.Context, err interface{}) {
	jsonError(c, http.StatusBadRequest, err)
}

// notFound renders a HTTP 404 not found response
func notFound(c *gin.Context, err interface{}) {
	jsonError(c, http.StatusNotFound, err)
}

// jsonOk renders a successful response
func jsonOk(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}
