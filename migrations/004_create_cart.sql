-- +Goose Up
CREATE TABLE cart_items (
	id SERIAL PRIMARY KEY,
	user_id SERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +Goose Down
DROP TABLE IF EXISTS cart_items;

