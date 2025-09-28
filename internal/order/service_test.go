package order

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
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
