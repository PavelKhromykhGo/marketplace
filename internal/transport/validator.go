package transport

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			out := make([]map[string]string, 0, len(verrs))
			for _, fe := range verrs {
				out = append(out, map[string]string{
					"field": fe.Field(),
					"tag":   fe.Tag(),
				})
			}
			c.JSON(400, gin.H{
				"errors":  "validation failed",
				"details": out,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return false
	}
	return true
}
