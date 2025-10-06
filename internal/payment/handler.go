package payment

import (
	"marketplace/internal/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, svc Service) {
	h := NewHandler(svc)
	g := r.Group("/payments", auth.JWTAuth())
	{
		g.POST("/intents/:id", h.createIntent)
		g.POST("/intents/:id/confirm", h.confirmIntent)
	}
}

func (h *Handler) createIntent(c *gin.Context) {
	oid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pi, err := h.svc.CreateIntent(c, auth.GetUserID(c), oid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pi)
}

func (h *Handler) confirmIntent(c *gin.Context) {
	oid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var body struct {
		ClientSecret string `json:"client_secret" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pi, err := h.svc.Confirm(c, auth.GetUserID(c), oid, body.ClientSecret)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pi)
}
