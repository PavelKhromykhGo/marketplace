package product

// CreateProductReq represents the request body for creating a new product.
// swagger:model CreateProductReq
type CreateProductReq struct {
	// Product name
	// required: true
	// min: 2
	// max: 200
	Name string `json:"name" binding:"required,min=2,max=200"`

	// Product description
	// max: 2000
	Description string `json:"description" binding:"max=2000"`

	// Price in copecks
	// required: true
	// min: 1
	Price int64 `json:"price" binding:"required,gt=0"`

	// Stock quantity
	// required: true
	// min: 0
	Stock int `json:"stock" binding:"required,gte=0"`

	// Category ID
	// required: true
	// min: 1
	CategoryID int64 `json:"category_id" binding:"required,gt=0"`
}

// UpdateProductReq represents the request body for updating an existing product.
// swagger:model UpdateProductReq
type UpdateProductReq struct {
	Name        string `json:"name" binding:"required, min=2, max=200"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int    `json:"stock" binding:"required,gte=0"`
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
}

// CreateCategoryReq represents the request body for creating a new category.
// swagger:model CreateCategoryReq
type CreateCategoryReq struct {
	// Category name
	// required: true
	// min: 2
	// max: 128
	Name string `json:"name" binding:"required,min=2,max=128"`
}

// UpdateCategoryReq represents the request body for updating an existing category.
// swagger:model UpdateCategoryReq
type UpdateCategoryReq struct {
	Name string `json:"name" binding:"required,min=2,max=128"`
}

// IDResponse represents a response containing an ID.
// swagger:model IDResponse
type IDResponse struct {
	ID int64 `json:"id"` // ID of the created or updated entity
}

// ErrorResponse represents a generic error response.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}
