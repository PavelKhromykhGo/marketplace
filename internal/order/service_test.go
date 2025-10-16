package order

import (
	"context"
	"errors"
	"net/http"
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

func (m *mockRepo) GetOrderStatus(ctx context.Context, orderID int64) (string, error) {
	args := m.Called(ctx, orderID)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) UpdateOrderStatus(ctx context.Context, orderID int64, from, to string) error {
	args := m.Called(ctx, orderID, from, to)
	return args.Error(0)
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

	repo.On("GetCartItemsForUser", ctx, userID).Return(items, nil)
	repo.On("GetProductsPrices", ctx, []int64{10, 20}).Return(prices, nil)
	repo.On("BeginTx", ctx).Return(tx, nil)
	repo.On("DecrementStock", ctx, tx, int64(10), 2).Return(nil)
	repo.On("DecrementStock", ctx, tx, int64(20), 1).Return(nil)
	repo.On("CreateOrder", ctx, tx, mock.MatchedBy(func(o *Order) bool {
		return o.UserID == userID && o.Status == "new" && o.TotalAmount == 4000
	})).Return(orderID, nil)
	repo.On("BulkInsertItems", ctx, tx, orderID, mock.Anything).Return(nil)
	repo.On("ClearCart", ctx, tx, userID).Return(nil)
	tx.On("Rollback").Return(nil)
	tx.On("Commit").Return(nil)

	gotID, err := svc.CreateFromCart(ctx, userID, "")
	assert.NoError(t, err)
	assert.Equal(t, orderID, gotID)

	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestCreateFromCart_Success_WithIdempotency_FirstCall(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)
	idemKey := "abc-123"

	items := []CartItemLite{{ProductID: 10, Quantity: 2}, {ProductID: 20, Quantity: 1}}
	prices := map[int64]int64{10: 1000, 20: 2000}
	orderID := int64(555)

	tx := new(mockTx)

	idem.On("TryStartIdempotent", ctx, userID, idemKey, mock.AnythingOfType("string")).Return(true, 0, int64(0), nil)

	repo.On("GetCartItemsForUser", ctx, userID).Return(items, nil)
	repo.On("GetProductsPrices", ctx, []int64{10, 20}).Return(prices, nil)
	repo.On("BeginTx", ctx).Return(tx, nil)
	repo.On("DecrementStock", ctx, tx, int64(10), 2).Return(nil)
	repo.On("DecrementStock", ctx, tx, int64(20), 1).Return(nil)
	repo.On("CreateOrder", ctx, tx, mock.MatchedBy(func(o *Order) bool {
		return o.UserID == userID && o.TotalAmount == 4000
	})).Return(orderID, nil)
	repo.On("BulkInsertItems", ctx, tx, orderID, mock.Anything).Return(nil)
	repo.On("ClearCart", ctx, tx, userID).Return(nil)
	tx.On("Rollback").Return(nil)
	tx.On("Commit").Return(nil)

	idem.On("SaveIdempotentResult", ctx, idemKey, http.StatusCreated, orderID).Return(nil)

	gotID, err := svc.CreateFromCart(ctx, userID, idemKey)
	assert.NoError(t, err)
	assert.Equal(t, orderID, gotID)

	repo.AssertExpectations(t)
	idem.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestCreateFromCart_SecondCall_ReturnsSaved(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)
	idemKey := "abc-123"
	orderID := int64(555)

	idem.On("TryStartIdempotent", ctx, userID, idemKey, mock.AnythingOfType("string")).Return(false, http.StatusCreated, orderID, nil)

	gotID, err := svc.CreateFromCart(ctx, userID, idemKey)
	assert.NoError(t, err)
	assert.Equal(t, orderID, gotID)

	idem.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateFromCart_Conflict_WhenKeyBusyWithoutResult(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)
	idemKey := "abc-123"

	idem.On("TryStartIdempotent", ctx, userID, idemKey, mock.AnythingOfType("string")).Return(false, 0, int64(0), nil)

	gotID, err := svc.CreateFromCart(ctx, userID, idemKey)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrIdempotencyConflict))
	assert.Equal(t, int64(0), gotID)

	idem.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateFromCart_EmptyCart_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)

	repo.On("GetCartItemsForUser", ctx, userID).Return([]CartItemLite{}, nil)

	gotID, err := svc.CreateFromCart(ctx, userID, "")
	assert.Error(t, err)
	assert.Equal(t, int64(0), gotID)

	repo.AssertExpectations(t)
}

func TestCreateFromCart_StockNotEnough_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	idem := new(mockIdemRepo)
	svc := NewService(repo, idem)

	userID := int64(1)
	items := []CartItemLite{{ProductID: 10, Quantity: 2}, {ProductID: 20, Quantity: 1}}
	prices := map[int64]int64{10: 1000, 20: 2000}
	tx := new(mockTx)

	repo.On("GetCartItemsForUser", ctx, userID).Return(items, nil)
	repo.On("GetProductsPrices", ctx, []int64{10, 20}).Return(prices, nil)
	repo.On("BeginTx", ctx).Return(tx, nil)

	repo.On("DecrementStock", ctx, tx, int64(10), 2).Return(errors.New("not enough stock"))
	tx.On("Rollback").Return(nil)

	gotID, err := svc.CreateFromCart(ctx, userID, "")
	assert.Error(t, err)
	assert.Equal(t, int64(0), gotID)

	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}
