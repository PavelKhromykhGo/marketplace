package product

type CreateProductReq struct {
	Name        string `json:"name" binding:"required,min=2,max=200"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int    `json:"stock" binding:"required,gte=0"`
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
}

type UpdateProductReq struct {
	Name        string `json:"name" binding:"required, min=2, max=200"`
	Description string `json:"description" binding:"max=2000"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int    `json:"stock" binding:"required,gte=0"`
	CategoryID  int64  `json:"category_id" binding:"required,gt=0"`
}

type CreateCategoryReq struct {
	Name string `json:"name" binding:"required,min=2,max=128"`
}

type UpdateCategoryReq struct {
	Name string `json:"name" binding:"required,min=2,max=128"`
}
