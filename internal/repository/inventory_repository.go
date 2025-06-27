package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
)

type InventoryRepository interface {
	CreateInventory(context.Context, *domain.Inventory) error
	ChangeProductCount(context.Context, *domain.Inventory) error
	AddDiscountToProducts(context.Context, []*domain.Inventory) error
}
