package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		reqID, _ := c.Get("request_id")

		if len(c.Errors) > 0 {
			zap.L().Error(
				"Request",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Int("size", size),
				zap.Any("request_id", reqID),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.ClientIP()),
				zap.String("error", c.Errors.String()),
			)
			return
		}
		zap.L().Info(
			"Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Int("size", size),
			zap.Any("request_id", reqID),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
