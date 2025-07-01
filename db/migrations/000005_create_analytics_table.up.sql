CREATE TABLE IF NOT EXISTS analytics(
    analytic_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID REFERENCES warehouse(warehouse_id),
    product_id UUID REFERENCES product(product_id),
    product_count INT CONSTRAINT positive_count CHECK (product_count >= 0),
    product_price NUMERIC(10, 2) CONSTRAINT positive_price CHECK (product_price >= 0)
)
