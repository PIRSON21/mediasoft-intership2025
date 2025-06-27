package dto

type InventoryCreateRequest struct {
	WarehouseID string   `json:"warehouse_id"`
	ProductID   string   `json:"product_id"`
	Count       *int     `json:"product_count"`
	Price       *float64 `json:"product_price"`
}

type ChangeProductCountRequest struct {
	WarehouseID string `json:"warehouse_id"`
	ProductID   string `json:"product_id"`
	Count       *int   `json:"product_count"`
}
