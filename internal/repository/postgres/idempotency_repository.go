package postgres

import (
	"context"
	"database/sql"
	"marketplace/internal/order"

	"github.com/jmoiron/sqlx"
)

type IdemRepo struct {
	db *sqlx.DB
}

func NewIdempotencyRepository(db *sqlx.DB) *IdemRepo {
	return &IdemRepo{db: db}
}

func (r *IdemRepo) TryStartIdempotent(ctx context.Context, userID int64, key string, reqHash string) (ok bool, savedStatus int, savedOrderID int64, err error) {
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO idempotency_keys (user_id, key, request_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, userID, key, reqHash)
	if err != nil {
		return false, 0, 0, err
	}

	var status sql.NullInt32
	var orderID sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT status_code, (response_body->>'order_id')::bigint
		FROM idempotency_keys
		WHERE key = $1
	`, key).Scan(&status, &orderID)
	if err != nil {
		return false, 0, 0, err
	}

	if orderID.Valid {
		return false, int(status.Int32), orderID.Int64, nil
	}
	return false, int(status.Int32), 0, order.ErrIdempotencyConflict
}

func (r *IdemRepo) SaveIdempotentResult(ctx context.Context, key string, status int, orderID int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE idempotency_keys
		SET status_code = $1, response_body = jsonb_build_object('order_id', $2)
		WHERE key = $3
	`, status, orderID, key)
	return err
}
