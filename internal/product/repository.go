package product

import "context"

type Repository interface {
	Create(ctx context.Context, p *Product) (int64, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, offset, limit int, filter string) ([]*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error

	CreateCategory(ctx context.Context, c *Category) (int64, error)
	GetCategory(ctx context.Context, id int64) (*Category, error)
	ListCategories(ctx context.Context, offset, limit int, filter string) ([]*Category, error)
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id int64) error
}
