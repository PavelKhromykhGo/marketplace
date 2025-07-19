package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AbortWithError(c *gin.Context, status int, err error) {
	if status >= http.StatusInternalServerError {
		_ = c.Error(err) // добавляет в gin.Context.Errors
	}

	c.AbortWithStatusJSON(status, gin.H{
		"error": err.Error(),
	})
}

func Abort400(c *gin.Context, err error) {
	AbortWithError(c, http.StatusBadRequest, err)
}

func Abort404(c *gin.Context, err error) {
	AbortWithError(c, http.StatusNotFound, err)
}

func Abort500(c *gin.Context, err error) {
	AbortWithError(c, http.StatusInternalServerError, err)
}
