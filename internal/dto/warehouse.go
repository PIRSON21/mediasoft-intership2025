package dto

// WarehouseRequest представляет запрос на создание или обновление склада.
type WarehouseRequest struct {
	Address string `json:"address"`
}

// WarehouseAtListResponse представляет склад в списке с его деталями.
type WarehouseAtListResponse struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}
