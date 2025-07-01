package service

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	"github.com/PIRSON21/mediasoft-intership2025/internal/repository"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

func (s *AnalyticsService) AddProductSell(invs []*domain.Inventory) {
	log := logger.GetLogger().With(
		zap.String("op", "service.AnalyticsService.AddProductSell"),
	)

	err := s.repo.AddProductSell(invs)
	if err != nil {
		log.Error("error while adding product sell to analytics", zap.Error(err))
		return
	}
}

func (s *AnalyticsService) GetWarehouseAnalytics(ctx context.Context, warehouseID string) (*dto.WarehouseAnalyticsResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.GetWarehouseAnalytics"),
	)

	analytics, err := s.repo.GetWarehouseAnalytics(ctx, warehouseID)
	if err != nil {
		log.Error("error while getting warehouse analytics from repository", zap.Error(err))
		return nil, err
	}

	response := parseWarehouseAnalyticsToResponse(warehouseID, analytics)

	return response, nil
}

func parseWarehouseAnalyticsToResponse(warehouseID string, analytics []*domain.Analytics) *dto.WarehouseAnalyticsResponse {
	analMap := make(map[uuid.UUID]*dto.ProductAnalytic)
	resp := dto.WarehouseAnalyticsResponse{
		WarehouseID: warehouseID,
	}

	for _, analytic := range analytics {
		anal, ok := analMap[analytic.Product.ID]
		if ok {
			anal.ProductCount += analytic.ProductCount
			anal.ProductPrice += analytic.ProductPrice
			resp.TotalSum += analytic.ProductPrice
			continue
		}
		anal = &dto.ProductAnalytic{
			ProductID:    analytic.Product.ID.String(),
			ProductName:  analytic.Product.Name,
			ProductCount: analytic.ProductCount,
			ProductPrice: analytic.ProductPrice,
		}
		analMap[analytic.Product.ID] = anal
		resp.TotalSum += anal.ProductPrice
		resp.Products = append(resp.Products, anal)
	}

	return &resp
}

func (s *AnalyticsService) GetTopWarehouses(ctx context.Context, limit int) ([]*dto.WarehouseAnalyticsAtListResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.AnalyticsService.GetTopWarehouse"),
	)

	response, err := s.repo.GetTopWarehouses(ctx, limit)
	if err != nil {
		log.Error("error while getting top warehouses from repository", zap.Error(err))
		return nil, err
	}

	return response, nil
}
