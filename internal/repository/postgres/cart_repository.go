package postgres

import (
	"context"
	"fmt"
	"marketplace/internal/cart"

	"github.com/jmoiron/sqlx"
)

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (c *CartRepo) AddItem(ctx context.Context, item *cart.CartItem) (int64, error) {
	query := `
INSERT INTO cart_items (user_id, product_id, quantity, created_at, updated_at)
VALUES (:user_id, :product_id, :quantity, NOW(), NOW())
RETURNING id
`

	rows, err := c.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки в бд: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ошибка получения id: %w", err)
		}
		return id, nil
	}
	return 0, fmt.Errorf("не удалось получить id вставленной записи")
}

func (c *CartRepo) ListItems(ctx context.Context, userID int64) ([]*cart.CartItem, error) {
	query := `
SELECT id, user_id, product_id, quantity, created_at, updated_at
FROM cart_items
WHERE user_id = $1
`
	var items []*cart.CartItem
	if err := c.db.SelectContext(ctx, &items, query, userID); err != nil {
		return nil, fmt.Errorf("ошибка получения из бд: %w", err)
	}
	return items, nil
}

func (c *CartRepo) RemoveItem(ctx context.Context, userID, productID int64) error {
	query := `
DELETE FROM cart_items
WHERE user_id = $1 AND product_id = $2
`
	_, err := c.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return fmt.Errorf("ошибка удаления из бд: %w", err)
	}
	return err
}

func (c *CartRepo) Clear(ctx context.Context, userID int64) error {
	query := `
DELETE FROM cart_items
WHERE user_id = $1
`
	_, err := c.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("ошибка очистки корзины в бд: %w", err)
	}
	return err
}
