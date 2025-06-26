CREATE TABLE IF NOT EXISTS inventory(
    inv_id SERIAL PRIMARY KEY,
    product_id INT REFERENCES product(product_id),
    warehouse_id INT REFERENCES warehouse(warehouse_id),
    product_count INT CONSTRAINT positive_count CHECK (product_count >= 0),
    product_price NUMERIC(10, 2) CONSTRAINT positive_price CHECK (product_price >= 0),
    product_sale INT CONSTRAINT positive_sale CHECK (product_sale >= 0)
);

CREATE UNIQUE INDEX idx_product_warehouse ON inventory(product_id, warehouse_id);
