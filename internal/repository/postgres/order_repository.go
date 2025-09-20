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
