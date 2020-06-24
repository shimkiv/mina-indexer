package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/coda-indexer/config"
)

// CORSMiddleware inject CORS headers into the response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

// RollbarMiddleware reports panics to rollback error tracker
func RollbarMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				config.TrackPanic(err)
				panic(err) // continue with default panic loger
			}
		}()
		c.Next()
	}
}
