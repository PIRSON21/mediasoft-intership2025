CREATE TABLE IF NOT EXISTS product(
  product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_name VARCHAR UNIQUE,
  product_description TEXT,
  product_weight FLOAT CONSTRAINT possitive_weight CHECK (product_weight >= 0),
  product_params JSONb,
  product_barcode VARCHAR
);

CREATE INDEX idx_product_name ON product(product_name);
