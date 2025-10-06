package payment

import (
	"marketplace/internal/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(r *gin.Engine, svc *Service) {
	h := NewHandler(svc)
	g := r.Group("/payments", auth.JWTAuth())
	{
		g.POST("/intents/:id", h.createIntent)
		g.POST("/intents/:id/confirm", h.confirmIntent)
	}
}

// @Summary Create Payment Intent
// @Description Create a payment intent for the specified order
// @Tags payments
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Sucscess 201 {object} Intent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /orders/{id}/payments [post]
func (h *Handler) createIntent(c *gin.Context) {
	oid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pi, err := h.svc.CreateIntent(c, auth.GetUserID(c), oid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pi)
}

// @Summary Confirm Payment Intent
// @Description Confirm the payment intent for the specified order using the client secret
// @Tags payments
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param client_secret body map[string]string true "client Secret"
// @Sucscess 200 {object} Intent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /orders/{id}/payments/confirm [post]
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
