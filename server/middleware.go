package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/config"
)

// corsMiddleware inject CORS headers into the response
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

// rollbarMiddleware reports panics to rollback error tracker
func rollbarMiddleware() gin.HandlerFunc {
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

func timeBucketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeBucket, err := getTimeBucket(c)
		if err != nil {
			badRequest(c, err)
			return
		}
		if err := timeBucket.validate(); err != nil {
			badRequest(c, err)
			return
		}
		c.Set("timebucket", timeBucket)
	}
}

func requestLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		msg := ""

		field := logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"client":   c.ClientIP(),
			"status":   status,
			"duration": duration.Milliseconds(),
			"path":     c.Request.URL.Path,
			"params":   c.Request.URL.Query(),
		})

		if err := c.Errors.Last(); err != nil {
			field = field.WithError(err)
		}

		switch {
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			field.Warn(msg)
		case status >= http.StatusInternalServerError:
			field.Error(msg)
		default:
			field.Info(msg)
		}
	}
}
