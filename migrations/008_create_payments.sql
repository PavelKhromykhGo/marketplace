-- +goose UP
CREATE TABLE payment_intents (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'requires_confirmation', -- requires_confirmation | succeeded | cancelled | failed
    client_secret TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL Default now()
);

CREATE INDEX idx_payment_intents_order ON payment_intents(order_id);

-- +goose Down
DROP TABLE IF EXISTS payment_intents