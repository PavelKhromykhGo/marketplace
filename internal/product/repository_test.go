package product

import (
	"context"
	"marketplace/internal/repository/postgres"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	cleanup := func() {
		if err = db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}

	xdb := sqlx.NewDb(db, "postgres")
	return xdb, mock, cleanup
}

func TestProductRepository_List(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at",
	}).AddRow(
		1, "iphone", "Description 1", int64(10000), 10, int64(2), time.Now(), time.Now(),
	)
	rawQuery := `
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
OFFSET $2 LIMIT $3`

	mock.ExpectQuery(regexp.QuoteMeta(rawQuery)).
		WithArgs("iphone", 0, 10).
		WillReturnRows(rows)

	got, err := repo.List(context.Background(), 0, 10, "iphone")
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "iphone", got[0].Name)
	assert.Equal(t, "Description 1", got[0].Description)
	assert.Equal(t, int64(10000), got[0].Price)
	assert.Equal(t, 10, got[0].Stock)
	assert.Equal(t, int64(2), got[0].CategoryID)
	assert.WithinDuration(t, time.Now(), got[0].CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), got[0].UpdatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Create(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	product := &Product{
		Name:        "New Product",
		Description: "New Description",
		Price:       15000,
		Stock:       3,
		CategoryID:  2,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
INSERT INTO products (name, description, price, stock, category_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`)).
		WithArgs(product.Name, product.Description, product.Price, product.Stock, product.CategoryID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.Create(context.Background(), product)
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetByID(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	expected := &Product{
		ID:          1,
		Name:        "TestProduct",
		Description: "TestDescription",
		Price:       1000,
		Stock:       5,
		CategoryID:  2,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "price", "stock", "category_id", "created_at", "updated_at"}).
		AddRow(expected.ID, expected.Name, expected.Description, expected.Price, expected.Stock,
			expected.CategoryID, expected.CreatedAt, expected.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
WHERE id = $1`)).
		WithArgs(expected.ID).
		WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), expected.ID)
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Description, got.Description)
	assert.Equal(t, expected.Price, got.Price)
	assert.Equal(t, expected.Stock, got.Stock)
	assert.Equal(t, expected.CategoryID, got.CategoryID)
	assert.WithinDuration(t, expected.CreatedAt, got.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, got.UpdatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Update(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	product := &Product{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       2000,
		Stock:       10,
		CategoryID:  3,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
UPDATE products
SET name = $1, description = $2, price = $3, stock = $4, category_id = $5
WHERE id = $6`)).
		WithArgs(product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Last insert ID is not used in UPDATE

	err := repo.Update(context.Background(), product)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Delete(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	productID := int64(1)

	mock.ExpectExec(regexp.QuoteMeta(`
DELETE FROM products
WHERE id = $1`)).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Last insert ID is not used in DELETE

	err := repo.Delete(context.Background(), productID)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetCategories(t *testing.T) {
	xdb, mock, cleanup := newMockDB(t)
	defer cleanup()
	repo := postgres.NewProductRepository(xdb)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Electronics").
		AddRow(2, "Books").
		AddRow(3, "Clothing")

	mock.ExpectQuery(regexp.QuoteMeta(`
SELECT id, name
FROM categories
ORDER BY name`)).
		WillReturnRows(rows)

	got, err := repo.GetCategories(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, int64(1), got[0].ID)
	assert.Equal(t, "Electronics", got[0].Name)
	assert.Equal(t, int64(2), got[1].ID)
	assert.Equal(t, "Books", got[1].Name)
	assert.Equal(t, int64(3), got[2].ID)
	assert.Equal(t, "Clothing", got[2].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}
