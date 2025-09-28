package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTx struct {
	mock.Mock
}

func (m *mockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) BeginTx(ctx context.Context) (Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(Tx), args.Error(1)
}

func (m *mockRepo) GetCartItemsForUser(ctx context.Context, userID int64) ([]CartItemLite, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]CartItemLite), args.Error(1)
}

func (m *mockRepo) GetProductsPrices(ctx context.Context, productIDs []int64) (map[int64]int64, error) {
	args := m.Called(ctx, productIDs)
	return args.Get(0).(map[int64]int64), args.Error(1)
}

func (m *mockRepo) DecrementStock(ctx context.Context, tx Tx, productID int64, quantity int) error {
	args := m.Called(ctx, tx, productID, quantity)
	return args.Error(0)
}

func (m *mockRepo) CreateOrder(ctx context.Context, tx Tx, order *Order) (int64, error) {
	args := m.Called(ctx, tx, order)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) BulkInsertItems(ctx context.Context, tx Tx, orderID int64, items []OrderItem) error {
	args := m.Called(ctx, tx, orderID, items)
	return args.Error(0)
}

func (m *mockRepo) ClearCart(ctx context.Context, tx Tx, userID int64) error {
	args := m.Called(ctx, tx, userID)
	return args.Error(0)
}

func (m *mockRepo) GetUserOrders(ctx context.Context, userID int64, offset, limit int) ([]*Order, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]*Order), args.Error(1)
}
func (m *mockRepo) GetOrderWithItems(ctx context.Context, userID, orderID int64) (*Order, error) {
	args := m.Called(ctx, userID, orderID)
	return args.Get(0).(*Order), args.Error(1)
}

type mockIdemRepo struct {
	mock.Mock
}

func (m *mockIdemRepo) TryStartIdempotent(ctx context.Context, userID int64, key string, reqHash string) (ok bool, savedStatus int, savedOrderID int64, err error) {
	args := m.Called(ctx, userID, key, reqHash)
	return args.Get(0).(bool), args.Get(1).(int), args.Get(2).(int64), args.Error(3)
}
func (m *mockIdemRepo) SaveIdempotentResult(ctx context.Context, key string, status int, orderID int64) error {
	args := m.Called(ctx, key, status, orderID)
	return args.Error(0)
}

func anyTx() any {
	return mock.MatchedBy(func(tx Tx) bool {
		return true
	})
}

func TestCreateFromCart_Success_NoIdempotency(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)
	items := []CartItemLite{{ProductID: 10, Quantity: 2}, {ProductID: 20, Quantity: 1}}
	prices := map[int64]int64{10: 1000, 20: 2000}
	orderID := int64(777)

	tx := new(mockTx)

	repo.On("GetCartItemForUser", ctx, userID).Return(items, nil)
	repo.On("GetProductPrices", ctx, []int64{10, 20}).Return(prices, nil)
	repo.On("BeginTx", ctx).Return(tx, nil)
	repo.On("DecrementStock", ctx, tx, int64(10), 2).Return(nil)
	repo.On("DecrementStock", ctx, tx, int64(20), 1).Return(nil)
	repo.On("CreateOrder", ctx, tx, mock.MatchedBy(func(o *Order) bool {
		return o.UserID == userID && o.Status == "new" && o.TotalAmount == 3000
	})).Return(orderID, nil)
	repo.On("BulkInsertItems", ctx, tx, orderID, mock.Anything).Return(nil)
	repo.On("ClearCart", ctx, tx, userID).Return(nil)
	tx.On("Commit").Return(nil)

	gotID, err := svc.CreateFromCart(ctx, userID, "")
	assert.NoError(t, err)
	assert.Equal(t, orderID, gotID)

	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}
