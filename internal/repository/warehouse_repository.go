package repository

import "github.com/PIRSON21/mediasoft-go/internal/domain"

type WarehouseRepository interface {
	GetWarehouses() []*domain.Warehouse
}
