package order

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
	g := r.Group("/orders", auth.JWTAuth())
	{
		g.POST("", h.createFromCart)
		g.GET("", h.listOrders)
		g.GET("/:id", h.getOrder)
	}
}

func (h *Handler) createFromCart(c *gin.Context) {
	id, err := h.svc.CreateFromCart(c, auth.GetUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) listOrders(c *gin.Context) {
	orders, err := h.svc.ListOrders(c, auth.GetUserID(c), 0, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) getOrder(c *gin.Context) {
	oid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, err := h.svc.GetOrder(c, auth.GetUserID(c), oid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
