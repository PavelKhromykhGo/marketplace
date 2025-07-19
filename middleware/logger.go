package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			if err.Err != nil {
				log.Printf("Error: %v, Status: %d, Method: %s, Path: %s",
					err.Err, c.Writer.Status(), c.Request.Method, c.Request.URL.Path)
			}
		}
	}
}
