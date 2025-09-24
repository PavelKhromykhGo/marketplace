package order

import (
	"errors"
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

// @Summary Create Order from Cart
// @Description Create a new order based on the current user's cart
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Idempotency-Key header string false "Idempotency Key"
// @Success 201 {object} map[string]int64 "id"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string "idempotency conflict"
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *Handler) createFromCart(c *gin.Context) {
	idKey := c.GetHeader("Idempotency-Key")

	id, err := h.svc.CreateFromCart(c, auth.GetUserID(c), idKey)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrIdempotencyConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary List Orders
// @Description List orders for the current user
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} order.Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *Handler) listOrders(c *gin.Context) {
	orders, err := h.svc.ListOrders(c, auth.GetUserID(c), 0, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// @Summary Get Order
// @Description Get details of a specific order by ID
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} order.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (h *Handler) getOrder(c *gin.Context) {
	oid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, err := h.svc.GetOrder(c, auth.GetUserID(c), oid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
