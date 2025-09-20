package order

import (
	"context"
	"fmt"
)

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
	CreateFromCart(ctx context.Context, userID int64) (int64, error)
	ListOrders(ctx context.Context, userID int64, offset, limit int) ([]*Order, error)
	GetOrder(ctx context.Context, userID, orderID int64) (*Order, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateFromCart(ctx context.Context, userID int64) (int64, error) {
	// 1) забираем корзину
	cartItems, err := s.repo.GetCartItemsForUser(ctx, userID)
	if err != nil || len(cartItems) == 0 {
		return 0, fmt.Errorf("empty cart: %w", err)
	}
	// 2) получаем цены на товары
	productIDs := make([]int64, 0, len(cartItems))
	for _, item := range cartItems {
		productIDs = append(productIDs, item.ProductID)
	}
	prices, err := s.repo.GetProductsPrices(ctx, productIDs)
	if err != nil {
		return 0, fmt.Errorf("cannot get prices: %w", err)
	}
	// 3) считаем сумму
	var totalAmount int64
	orderItems := make([]OrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		price, ok := prices[item.ProductID]
		if !ok {
			return 0, fmt.Errorf("price not found for product %d", item.ProductID)
		}
		totalAmount += price * int64(item.Quantity)
		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}
	// 4) начинаем транзакцию
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot begin tx: %w", err)
	}
	defer tx.Rollback()
	// 5) резервируем товары
	for _, item := range orderItems {
		if err = s.repo.DecrementStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
			return 0, fmt.Errorf("stock not enough for product=%d: %w", item.ProductID, err)
		}
	}
	// 6) создаем заказ
	order := &Order{
		UserID:      userID,
		Status:      "new",
		TotalAmount: totalAmount,
	}
	orderID, err := s.repo.CreateOrder(ctx, tx, order)
	if err != nil {
		return 0, fmt.Errorf("cannot create order: %w", err)
	}
	// 7) создаем позиции заказа
	if err = s.repo.BulkInsertItems(ctx, tx, orderID, orderItems); err != nil {
		return 0, fmt.Errorf("cannot insert order items: %w", err)
	}
	// 8) очищаем корзину
	if err = s.repo.ClearCart(ctx, tx, userID); err != nil {
		return 0, fmt.Errorf("cannot clear cart: %w", err)
	}
	// 9) коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("cannot commit tx: %w", err)
	}
	return orderID, nil
}

func (s *service) ListOrders(ctx context.Context, userID int64, offset, limit int) ([]*Order, error) {
	return s.repo.GetUserOrders(ctx, userID, offset, limit)
}

func (s *service) GetOrder(ctx context.Context, userID, orderID int64) (*Order, error) {
	return s.repo.GetOrderWithItems(ctx, userID, orderID)
}
