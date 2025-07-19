package product

import (
	"marketplace/internal/transport"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/products")
	{
		group.GET("", h.listProducts)
		group.GET("/:id", h.getProduct)
		group.POST("", h.createProduct)
		group.PUT("/:id", h.updateProduct)
		group.DELETE("/:id", h.deleteProduct)
	}

}

func (h *Handler) listProducts(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		transport.Abort400(c, err)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		transport.Abort400(c, err)
		return
	}
	filter := c.Query("filter")

	products, err := h.service.ListProducts(c.Request.Context(), offset, limit, filter)
	if err != nil {
		transport.Abort500(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) getProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		transport.Abort400(c, err)
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		transport.Abort500(c, err)
		return
	}
	if product == nil {
		transport.Abort404(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) createProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		transport.Abort400(c, err)
		return
	}

	id, err := h.service.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		transport.Abort500(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
func (h *Handler) updateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		transport.Abort400(c, err)
		return
	}

	var product Product
	if err = c.ShouldBindJSON(&product); err != nil {
		transport.Abort400(c, err)
		return
	}
	product.ID = id

	if err = h.service.UpdateProduct(c.Request.Context(), &product); err != nil {
		transport.Abort500(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *Handler) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		transport.Abort400(c, err)
		return
	}

	if err = h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		transport.Abort500(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
