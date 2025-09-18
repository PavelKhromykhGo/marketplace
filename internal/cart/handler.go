package cart

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

	g := r.Group("/cart", auth.JWTAuth())
	{
		g.GET("", h.list)
		g.POST("/items", h.add)
		g.DELETE("/items/:product_id", h.remove)
		g.DELETE("", h.clear)
	}

}

type addReq struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

// @Summary Add item to cart
// @Description Add an item to the user's cart
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body addReq true "Item to add"
// @Success 201 {object} map[string]int64 "id"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/items [post]
func (h *Handler) add(c *gin.Context) {
	var req addReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := auth.GetUserID(c)
	id, err := h.svc.AddItem(c.Request.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary List cart items
// @Description Get all items in the user's cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {array} CartItem
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [get]
func (h *Handler) list(c *gin.Context) {
	userID := auth.GetUserID(c)
	items, err := h.svc.ListItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// @Summary Remove item from cart
// @Description Remove an item from the user's cart
// @Tags Cart
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/items/{product_id} [delete]
func (h *Handler) remove(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}
	userID := auth.GetUserID(c)
	if err = h.svc.RemoveItem(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Clear cart
// @Description Remove all items from the user's cart
// @Tags Cart
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [delete]
func (h *Handler) clear(c *gin.Context) {
	userID := auth.GetUserID(c)
	if err := h.svc.Clear(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
