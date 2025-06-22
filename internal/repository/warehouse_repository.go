package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
)

type WarehouseRepository interface {
	GetWarehouses(context.Context) ([]*domain.Warehouse, error)
	CreateWarehouse(context.Context, *domain.Warehouse) error
}
