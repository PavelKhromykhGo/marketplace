INSERT INTO categories (name) VALUES
    ('Books'), ('Electronics'), ('Clothes');

INSERT INTO products (name, description, price, stock, category_id)
VALUES
    ('Go Book', 'Learn Go programming', 2990, 10, 1),
    ('Headphones', 'Noise-cancelling', 10990, 5, 2),
    ('T-Shirt', '100% cotton', 1990, 20, 3);