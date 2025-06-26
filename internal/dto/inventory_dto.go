package dto

type InventoryCreateRequest struct {
	WarehouseID *int     `json:"warehouse_id"`
	ProductID   *int     `json:"product_id"`
	Count       *int     `json:"product_count"`
	Price       *float64 `json:"product_price"`
}

type ChangeProductCountRequest struct {
	WarehouseID *int `json:"warehouse_id"`
	ProductID   *int `json:"product_id"`
	Count       *int `json:"product_count"`
}
