package transport

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// HandleDBError возвращает true, если ошибка распознана и ответ уже отправлен.
func HandleDBError(c *gin.Context, err error) bool {
	var pqe *pq.Error
	if errors.As(err, &pqe) {
		switch pqe.Code {
		case "23505": // Unique violation
			c.JSON(http.StatusConflict, gin.H{
				"error":   "resource already exists",
				"details": pqe.Detail,
			})
		case "23503": // Foreign key violation
			c.JSON(http.StatusConflict, gin.H{
				"error":   "invalid foreign key reference",
				"details": pqe.Detail,
			})
			return true
		case "23502": // Not null violation
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "null constraint violated",
				"details": pqe.Column,
			})
			return true
		case "23514": // Check violation
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "check constraint violation",
				"details": pqe.Detail,
			})
			return true
		}
	}
	return false
}
