package order

import (
	"context"
	"errors"
)

var ErrIdempotencyConflict = errors.New("idempotency conflict")

type IdempotencyRepository interface {
	TryStartIdempotent(ctx context.Context, userID int64, key string, reqHash string) (ok bool, savedStatus int, savedOrderID int64, err error)
	SaveIdempotentResult(ctx context.Context, key string, status int, orderID int64) error
}

type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)
	GetCartItemsForUser(ctx context.Context, userID int64) ([]CartItemLite, error)
	GetProductsPrices(ctx context.Context, productIDs []int64) (map[int64]int64, error)
	DecrementStock(ctx context.Context, tx Tx, productID int64, quantity int) error
	CreateOrder(ctx context.Context, tx Tx, order *Order) (int64, error)
	BulkInsertItems(ctx context.Context, tx Tx, orderID int64, items []OrderItem) error
	ClearCart(ctx context.Context, tx Tx, userID int64) error
	GetUserOrders(ctx context.Context, userID int64, offset, limit int) ([]*Order, error)
	GetOrderWithItems(ctx context.Context, userID, orderID int64) (*Order, error)
	GetOrderStatus(ctx context.Context, orderID int64) (string, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, from, to string) error
}

type Tx interface {
	Commit() error
	Rollback() error
}

type CartItemLite struct {
	ProductID int64
	Quantity  int
}

type Service interface {
	CreateFromCart(ctx context.Context, userID int64, idemKey string) (int64, error)
	ListOrders(ctx context.Context, userID int64, offset, limit int) ([]*Order, error)
	GetOrder(ctx context.Context, userID, orderID int64) (*Order, error)
}
