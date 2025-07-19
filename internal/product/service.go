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
	GetCategories(ctx context.Context) ([]*Category, error)
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

func (s *productService) GetCategories(ctx context.Context) ([]*Category, error) {
	return s.repo.GetCategories(ctx)
}
