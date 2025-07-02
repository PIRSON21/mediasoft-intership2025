package dto

type WarehouseAnalyticsResponse struct {
	WarehouseID string             `json:"warehouse_id"`
	Products    []*ProductAnalytic `json:"products"`
	TotalSum    float64            `json:"total_sum"`
}

type ProductAnalytic struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductCount int     `json:"total_product_count"`
	ProductPrice float64 `json:"total_product_price"`
}

type WarehouseAnalyticsAtListResponse struct {
	WarehouseID       string  `json:"warehouse_id"`
	WarehouseAddress  string  `json:"warehouse_address"`
	WarehouseTotalSum float64 `json:"warehouse_total_sum"`
}
