package mocks

import "marketplace/internal/product"

type MockRepository struct {
	CreateFn        func(p *product.Product) (int64, error)
	GetByIDFn       func(id int64) (*product.Product, error)
	ListFn          func(offset, limit int, filter string) ([]*product.Product, error)
	UpdateFn        func(p *product.Product) error
	DeleteFn        func(id int64) error
	GetCategoriesFn func() ([]*product.Category, error)
}

func (m *MockRepository) Create(p *product.Product) (int64, error) {
	return m.CreateFn(p)
}

func (m *MockRepository) GetByID(id int64) (*product.Product, error) {
	return m.GetByIDFn(id)
}

func (m *MockRepository) List(offset, limit int, filter string) ([]*product.Product, error) {
	return m.ListFn(offset, limit, filter)
}

func (m *MockRepository) Update(p *product.Product) error {
	return m.UpdateFn(p)
}

func (m *MockRepository) Delete(id int64) error {
	return m.DeleteFn(id)
}

func (m *MockRepository) GetCategories() ([]*product.Category, error) {
	return m.GetCategoriesFn()
}
