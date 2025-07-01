CREATE TABLE IF NOT EXISTS inventory(
    inv_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES product(product_id),
    warehouse_id UUID REFERENCES warehouse(warehouse_id),
    product_count INT CONSTRAINT positive_count CHECK (product_count >= 0),
    product_price NUMERIC(10, 2) CONSTRAINT positive_price CHECK (product_price >= 0),
    product_sale INT CONSTRAINT positive_sale CHECK (product_sale >= 0)
);

CREATE UNIQUE INDEX idx_product_warehouse ON inventory(product_id, warehouse_id);
