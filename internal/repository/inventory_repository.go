package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	"github.com/PIRSON21/mediasoft-go/internal/dto"
)

type InventoryRepository interface {
	CreateInventory(context.Context, *domain.Inventory) error
	ChangeProductCount(context.Context, *domain.Inventory) error
	AddDiscountToProducts(context.Context, []*domain.Inventory) error
	GetProductFromWarehouse(context.Context, *domain.Inventory) error
	GetPriceAndDiscount(context.Context, []*domain.Inventory) error
	GetProductsAtWarehouse(context.Context, *dto.Pagination, string) ([]*domain.Inventory, error)
}
