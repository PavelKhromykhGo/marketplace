package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"marketplace/internal/product"

	"github.com/jmoiron/sqlx"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *product.Product) (int64, error) {
	query := `
INSERT INTO products (name, description, price, stock, category_id, created_at, updated_at)
VALUES (:name, :description, :price, :stock, :category_id, NOW(), NOW())
RETURNING id
`

	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()

	var id int64
	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
	}
	return id, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	query := `
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
WHERE id = :id
`

	rows, err := r.db.NamedQueryContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()

	var p product.Product
	if rows.Next() {
		if err = rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
		return &p, nil
	}
	return nil, sql.ErrNoRows
}

func (r *ProductRepo) List(ctx context.Context, offset, limit int, filter string) ([]*product.Product, error) {
	query := `
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
WHERE name ILIKE '%' || :filter || '%'
ORDER BY name
OFFSET :offset LIMIT :limit 
`

	rows, err := r.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()

	var products []*product.Product
	for rows.Next() {
		var p product.Product
		if err = rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *product.Product) error {
	query := `
UPDATE products
SET name = :name, description = :description, price = :price, stock = :stock, category_id = :category_id, updated_at = NOW()
WHERE id = :id
`

	_, err := r.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	query := `
DELETE FROM products
WHERE id = :id
`

	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return nil
}

func (r *ProductRepo) CreateCategory(ctx context.Context, c *product.Category) (int64, error) {
	query := `
INSERT INTO categories (name)
VALUES (:name)
RETURNING id
`
	rows, err := r.db.NamedQueryContext(ctx, query, c)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()
	var id int64
	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
	}
	return id, nil
}

func (r *ProductRepo) GetCategory(ctx context.Context, id int64) (*product.Category, error) {
	query := `
SELECT id, name
FROM categories
WHERE id = :id
`
	rows, err := r.db.NamedQueryContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()
	var c product.Category
	if rows.Next() {
		if err = rows.StructScan(&c); err != nil {
			return nil, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
		return &c, nil
	}
	return nil, fmt.Errorf("категория с id %d не найдена", id)
}

func (r *ProductRepo) ListCategories(ctx context.Context, offset, limit int, filter string) ([]*product.Category, error) {
	query := `
SELECT id, name
FROM categories
WHERE name ILIKE '%' || :filter || '%'
ORDER BY name
OFFSET :offset LIMIT :limit
`
	args := map[string]any{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("ошибка закрытия rows: %v", err)
		}
	}()
	var categories []*product.Category
	for rows.Next() {
		var c product.Category
		if err = rows.StructScan(&c); err != nil {
			return nil, fmt.Errorf("ошибка сканирования результата: %w", err)
		}
		categories = append(categories, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения результатов: %w", err)
	}
	return categories, nil
}

func (r *ProductRepo) UpdateCategory(ctx context.Context, c *product.Category) error {
	query := `
UPDATE categories
SET name = :name
WHERE id = :id
`
	_, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return nil
}

func (r *ProductRepo) DeleteCategory(ctx context.Context, id int64) error {
	query := `
DELETE FROM categories
WHERE id = :id
`
	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	return nil
}
