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
		INSERT INTO products (name, description, price, category_id, created_at, updated_at)
		VALUES (:name, :description, :price, :category_id, NOW(), NOW())
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
		SET name = :name,
		    description = :description,
		    price = :price,
		    stock = :stock,
		    category_id = :category_id,
		    updated_at = NOW()
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, p)
	return fmt.Errorf("ошибка выполнения запроса: %w", err)
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM products
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{"id": id})
	return fmt.Errorf("ошибка выполнения запроса: %w", err)
}

func (r *ProductRepo) GetCategories(ctx context.Context) ([]*product.Category, error) {
	query := `
		SELECT id, name
		FROM categories
		ORDER BY name
	`

	rows, err := r.db.QueryxContext(ctx, query)
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
	return categories, nil
}
