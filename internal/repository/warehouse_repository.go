package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
)

// WarehouseRepository - интерфейс для работы со складами.
type WarehouseRepository interface {
	GetWarehouses(context.Context) ([]*domain.Warehouse, error)
	CreateWarehouse(context.Context, *domain.Warehouse) error
}
