package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"marketplace/internal/order"
	"marketplace/internal/payment"

	"github.com/jmoiron/sqlx"
)

type PaymentRepo struct {
	db *sqlx.DB
}

func NewPaymentRepo(db *sqlx.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func randSecret(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *PaymentRepo) CreateIntent(ctx context.Context, o *order.Order) (*payment.Intent, error) {
	secret, err := randSecret(16)
	if err != nil {
		return nil, err
	}
	var pi payment.Intent
	err = r.db.QueryRowxContext(ctx, `
		WITH upd AS (
			UPDATE orders SET status = 'awaiting_payment'
			WHERE id = $1 AND status IN ('new')
			RETURNING id, total_amount
		)
		INSERT INTO payment_intents (order_id, amount, status, client_secret)
			SELECT id, total_amount, 'requires_confirmation', $2 FROM upd
		RETURNING id, order_id, amount, status, client_secret, created_at, updated_at
	`, o.ID, secret).Scan(&pi.ID, &pi.OrderID, &pi.Amount, &pi.Status, &pi.ClientSecret, &pi.CreatedAt, &pi.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create intent failed: %w", err)
	}
	return &pi, nil
}

func (r *PaymentRepo) ConfirmIntent(ctx context.Context, orderID int64, clientSecret string) (*payment.Intent, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var pi payment.Intent
	err = tx.QueryRowxContext(ctx, `
		UPDATE payment_intents
		SET status = 'succeeded', updated_at = NOW()
		WHERE order_id = $1 AND client_secret = $2 AND status = 'requires_confirmation'
		RETURNING id, order_id, amount, status, client_secret, created_at, updated_at
	`, orderID, clientSecret).Scan(
		&pi.ID, &pi.OrderID, &pi.Amount, &pi.Status, &pi.ClientSecret, &pi.CreatedAt, &pi.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("confirm intent failed: %w", err)
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE orders SET status = 'paid', updated_at = NOW()
		WHERE id = $1 AND status = 'awaiting_payment'
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("update order status failed: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("order not found or not in awaiting_payment status")
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return &pi, nil
}
