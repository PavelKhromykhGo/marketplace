package product

import (
	"context"
)

type Service interface {
	GetProduct(ctx context.Context, id int64) (*Product, error)
	ListProducts(ctx context.Context, offset, limit int, filter string) ([]*Product, error)
	CreateProduct(ctx context.Context, p *Product) (int64, error)
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id int64) error

	CreateCategory(ctx context.Context, c *Category) (int64, error)
	GetCategory(ctx context.Context, id int64) (*Category, error)
	ListCategories(ctx context.Context, offset, limit int, filter string) ([]*Category, error)
	UpdateCategory(ctx context.Context, c *Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type productService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &productService{repo: r}
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, offset, limit int, filter string) ([]*Product, error) {
	return s.repo.List(ctx, offset, limit, filter)
}

func (s *productService) CreateProduct(ctx context.Context, p *Product) (int64, error) {
	return s.repo.Create(ctx, p)
}

func (s *productService) UpdateProduct(ctx context.Context, p *Product) error {
	return s.repo.Update(ctx, p)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) CreateCategory(ctx context.Context, c *Category) (int64, error) {
	return s.repo.CreateCategory(ctx, c)
}

func (s *productService) GetCategory(ctx context.Context, id int64) (*Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *productService) ListCategories(ctx context.Context, offset, limit int, filter string) ([]*Category, error) {
	return s.repo.ListCategories(ctx, offset, limit, filter)
}

func (s *productService) UpdateCategory(ctx context.Context, c *Category) error {
	return s.repo.UpdateCategory(ctx, c)
}

func (s *productService) DeleteCategory(ctx context.Context, id int64) error {
	return s.repo.DeleteCategory(ctx, id)
}
