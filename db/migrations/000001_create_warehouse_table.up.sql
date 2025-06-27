CREATE TABLE IF NOT EXISTS warehouse (
    warehouse_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_address TEXT UNIQUE
);
