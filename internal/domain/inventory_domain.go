package domain

import "github.com/google/uuid"

type Inventory struct {
	ProductID    uuid.UUID
	WarehouseID  uuid.UUID
	ProductCount int
	ProductPrice float64
	ProductSale  int
}
