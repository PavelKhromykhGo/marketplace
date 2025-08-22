package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var pqe *pq.Error
		if errors.As(err, &pqe) {
			switch pqe.Code {
			case "23505": // unique_violation
				c.JSON(http.StatusConflict, gin.H{
					"error":  "duplicate key value violates unique constraint",
					"detail": pqe.Detail,
				})
				return
			case "23503": // foreign_key_violation
				c.JSON(http.StatusConflict, gin.H{
					"error":  "insert or update on table violates foreign key constraint",
					"detail": pqe.Detail,
				})
				return
			case "23514": // check_violation
				c.JSON(http.StatusBadRequest, gin.H{
					"error":  "check constraint violation",
					"detail": pqe.Detail,
				})
				return
			}
		}
		var verrs gin.H
		if errors.As(err, &verrs) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "validation failed",
				"detail": verrs,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
