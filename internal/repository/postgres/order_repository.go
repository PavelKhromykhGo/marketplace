package postgres

import (
	"context"
	"database/sql"
	"marketplace/internal/order"

	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

type txWrap struct {
	*sqlx.Tx
}

func (t *txWrap) Commit() error {
	return t.Tx.Commit()
}

func (t *txWrap) Rollback() error {
	return t.Tx.Rollback()
}

func (r *OrderRepo) BeginTx(ctx context.Context) (order.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	return &txWrap{tx}, nil
}

func (r *OrderRepo) GetCartItemsForUser(ctx context.Context, userID int64) ([]order.CartItemLite, error) {
	var items []order.CartItemLite
	err := r.db.SelectContext(ctx, &items, `
		SELECT product_id, quantity
		FROM cart_items
		WHERE user_id=$1
	`, userID)
	return items, err
}

func (r *OrderRepo) GetProductsPrices(ctx context.Context, productIDs []int64) (map[int64]int64, error) {
	query, args, err := sqlx.In(`
		SELECT id, price
		FROM products
		WHERE id IN (?)
	`, productIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	type row struct {
		ID    int64 `db:"id"`
		Price int64 `db:"price"`
	}
	var rows []row
	if err = r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	m := make(map[int64]int64, len(rows))
	for _, v := range rows {
		m[v.ID] = v.Price
	}
	return m, nil
}

func (r *OrderRepo) DecrementStock(ctx context.Context, tx order.Tx, productID int64, quantity int) error {
	xtx := tx.(*txWrap)
	result, err := xtx.ExecContext(ctx, `
		UPDATE products
		SET stock = stock - $1
		WHERE id=$2 AND stock >= $1
	`, quantity, productID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *OrderRepo) CreateOrder(ctx context.Context, tx order.Tx, o *order.Order) (int64, error) {
	xtx := tx.(*txWrap)
	var id int64
	err := xtx.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, status, total_amount, created_at, updated_at)
		VALUES (:user_id, :status, :total_amount, now(), now())
		RETURNING id
	`, o).Scan(&id)

	return id, err
}

func (r *OrderRepo) BulkInsertItems(ctx context.Context, tx order.Tx, orderID int64, items []order.OrderItem) error {
	xtx := tx.(*txWrap)
	q := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
	for _, item := range items {
		_, err := xtx.ExecContext(ctx, q, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepo) ClearCart(ctx context.Context, tx order.Tx, userID int64) error {
	xtx := tx.(*txWrap)
	_, err := xtx.ExecContext(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1
	`, userID)
	return err
}

func (r *OrderRepo) GetUserOrders(ctx context.Context, userID int64, offset, limit int) ([]*order.Order, error) {
	var orders []*order.Order
	err := r.db.SelectContext(ctx, &orders, `
		SELECT id, user_id, status, total_amount, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`, userID, offset, limit)
	return orders, err
}

func (r *OrderRepo) GetOrderWithItems(ctx context.Context, userID, orderID int64) (*order.Order, error) {
	var o order.Order
	err := r.db.GetContext(ctx, &o, `
		SELECT id, user_id, status, total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`, orderID, userID)
	if err != nil {
		return nil, err
	}
	var items []order.OrderItem
	err = r.db.SelectContext(ctx, &items, `
		SELECT product_id, quantity, price
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (r *OrderRepo) GetOrderStatus(ctx context.Context, orderID int64) (string, error) {
	var status string
	if err := r.db.GetContext(ctx, &status, `
		SELECT status
		FROM orders
		WHERE id = $1
	`, orderID); err != nil {
		return "", err
	}
	return status, nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderID int64, from, to string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND status = $3
	`, to, orderID, from)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
