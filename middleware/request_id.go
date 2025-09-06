package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Writer.Header().Set(requestIDHeader, reqID)
		c.Set("RequestID", reqID)
		c.Next()
	}
}
