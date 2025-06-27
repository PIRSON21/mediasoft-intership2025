package dto

type WarehouseRequest struct {
	Address string `json:"address"`
}

type WarehouseAtListResponse struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}
