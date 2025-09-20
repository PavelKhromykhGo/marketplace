-- +goose UP
CREATE TABLE idempotency_keys (
    key TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    request_hash TEXT NOT NULL,
    responce_body JSONB,
    status_code INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose DOWN
DROP TABLE IF EXISTS idempotency_keys;