ALTER TABLE IF EXISTS categories
    ADD CONSTRAINT unique_categories_name UNIQUE (name);

ALTER TABLE IF EXISTS products
    ADD CONSTRAINT chk_products_price CHECK (price > 0);
    ADD CONSTRAINT chk_products_stock CHECK (stock >= 0);

ALTER TABLE IF EXISTS products
    ADD CONSTRAINT fk_products_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_products_name_category
    ON products(LOWER(name), category_id);