package product

import (
	"errors"
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

	p := r.Group("/products")
	{
		p.GET("", h.listProducts)
		p.GET("/:id", h.getProduct)
		p.POST("", h.createProduct)
		p.PUT("/:id", h.updateProduct)
		p.DELETE("/:id", h.deleteProduct)
	}
	cg := r.Group("/categories")
	{
		cg.GET("", h.listCategories)
		cg.GET("/:id", h.getCategory)
		cg.POST("", h.createCategory)
		cg.PUT("/:id", h.updateCategory)
		cg.DELETE("/:id", h.deleteCategory)
	}
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return 0, false
	}
	return id, true
}

func parsePaging(c *gin.Context) (offset, limit int, filter string) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return 0, 0, ""
	}
	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return 0, 0, ""
	}
	return offset, limit, ""
}

// listProducts godoc
// @Summary List products
// @Description Get a list of products with pagination
// @Tags products
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Param filter query string false "Name filter"
// @Success 200 {array} Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]

func (h *Handler) listProducts(c *gin.Context) {
	offset, limit, filter := parsePaging(c)

	products, err := h.service.ListProducts(c.Request.Context(), offset, limit, filter)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, products)
}

// getProduct godoc
// @Summary Get a product by ID
// @Description Get a single product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]

func (h *Handler) getProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}
	if product == nil {
		_ = c.Error(errors.New("product not found")).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, product)
}

// createProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductReq true "Product payload"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *Handler) createProduct(c *gin.Context) {
	var req CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	id, err := h.service.CreateProduct(c.Request.Context(), &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateProduct godoc
// @Summary Update an existing product
// @Description Update the details of an existing product by its ID
// @Tags products
// @Accept json
// @Param id path int true "Product ID"
// @Param product body UpdateProductReq true "Product payload"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *Handler) updateProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	if err := h.service.UpdateProduct(c.Request.Context(), &Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a single product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.Status(http.StatusNoContent)
}

// listCategories godoc
// @Summary List categories
// @Description Get a list of categories with pagination
// @Tags categories
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Param filter query string false "Name filter"
// @Success 200 {array} Category
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (h *Handler) listCategories(c *gin.Context) {
	offset, limit, filter := parsePaging(c)

	categories, err := h.service.ListCategories(c.Request.Context(), offset, limit, filter)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// getCategory godoc
// @Summary Get a category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *Handler) getCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	category, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}
	if category == nil {
		_ = c.Error(errors.New("category not found")).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, category)
}

// createCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryReq true "Category payload"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [post]
func (h *Handler) createCategory(c *gin.Context) {
	var req CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	id, err := h.service.CreateCategory(c.Request.Context(), &Category{
		Name: req.Name,
	})
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateCategory godoc
// @Summary Update an existing category
// @Description Update the details of an existing category by its ID
// @Tags categories
// @Accept json
// @Param id path int true "Category ID"
// @Param category body UpdateCategoryReq true "Category payload"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *Handler) updateCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	var req UpdateCategoryReq

	if err := h.service.UpdateCategory(c.Request.Context(), &Category{
		ID:   id,
		Name: req.Name,
	}); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) deleteCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return // err уже в c.Errors
	}

	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.Status(http.StatusNoContent)
}
