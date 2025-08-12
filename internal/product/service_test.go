package product

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(ctx context.Context, p *Product) (int64, error) {
	args := m.Called(ctx, p)
	if id, ok := args.Get(0).(int64); ok {
		return id, args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*Product, error) {
	args := m.Called(ctx, id)
	if product, ok := args.Get(0).(*Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) List(ctx context.Context, offset, limit int, filter string) ([]*Product, error) {
	args := m.Called(ctx, offset, limit, filter)
	if products, ok := args.Get(0).([]*Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, p *Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockRepo) GetCategories(ctx context.Context) ([]*Category, error) {
	args := m.Called(ctx)
	if categories, ok := args.Get(0).([]*Category); ok {
		return categories, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestListProducts(t *testing.T) {
	ctx := context.Background()
	fakeRepo := new(mockRepo)
	svc := NewService(fakeRepo)

	expected := []*Product{
		{ID: 1, Name: "TestProduct_1", Price: 1000},
		{ID: 2, Name: "TestProduct_2", Price: 2000},
	}

	fakeRepo.On("List", ctx, 0, 10, "").Return(expected, nil)

	// Call the service method
	products, err := svc.ListProducts(ctx, 0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, products)

	fakeRepo.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	ctx := context.Background()
	fakeRepo := new(mockRepo)
	svc := NewService(fakeRepo)

	expected := &Product{ID: 1, Name: "TestProduct", Price: 1000}

	t.Run("успешно", func(t *testing.T) {
		fakeRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		product, err := svc.GetProduct(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, product)

		fakeRepo.AssertExpectations(t)
	})
	t.Run("не найден", func(t *testing.T) {
		fakeRepo.On("GetByID", ctx, int64(99999)).Return(nil, errors.New("не найдено"))

		product, err := svc.GetProduct(ctx, 2)
		assert.Error(t, err)
		assert.Nil(t, product)

		fakeRepo.AssertExpectations(t)
	})
}
