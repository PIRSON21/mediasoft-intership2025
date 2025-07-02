package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
)

type AnalyticsRepository interface {
	AddProductSell([]*domain.Inventory) error
	GetWarehouseAnalytics(context.Context, string) ([]*domain.Analytics, error)
	GetTopWarehouses(context.Context, int) ([]*dto.WarehouseAnalyticsAtListResponse, error)
}
