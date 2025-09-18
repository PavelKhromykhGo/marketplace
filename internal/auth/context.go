package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) int64 {
	if c == nil {
		return 0
	}
	v, ok := c.Get("userID")
	if !ok || v == nil {
		return 0
	}
	switch v := v.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			return id
		}
	case fmt.Stringer:
		if id, err := strconv.ParseInt(v.String(), 10, 64); err == nil {
			return id
		}
	}
	return 0
}
