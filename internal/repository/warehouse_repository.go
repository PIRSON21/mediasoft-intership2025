package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
)

type WarehouseRepository interface {
	GetWarehouses(context.Context) ([]*domain.Warehouse, error)
	CreateWarehouse(context.Context, *domain.Warehouse) error
}
