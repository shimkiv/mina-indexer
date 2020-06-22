package server

import (
	"github.com/gin-gonic/gin"

	"github.com/figment-networks/coda-indexer/config"
)

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
