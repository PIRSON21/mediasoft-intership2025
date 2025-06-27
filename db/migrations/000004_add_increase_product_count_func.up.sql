CREATE OR REPLACE FUNCTION increase_product_count(
    in_product_id UUID,
    in_warehouse_id UUID,
    in_delta INT
) RETURNS VOID AS $$
DECLARE
    updated_rows INT;
BEGIN
    UPDATE inventory
    SET product_count = product_count + in_delta
    WHERE product_id = in_product_id AND warehouse_id = in_warehouse_id;

    GET DIAGNOSTICS updated_rows = ROW_COUNT;

    IF updated_rows = 0 THEN
        RAISE EXCEPTION 'Inventory row not found for product_id=%, warehouse_id=%', in_product_id, in_warehouse_id
            USING ERRCODE = 'P0002';
    END IF;
END;
$$ LANGUAGE plpgsql;
