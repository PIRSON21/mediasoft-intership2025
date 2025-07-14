package service

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	"github.com/PIRSON21/mediasoft-intership2025/internal/handler"
	"github.com/PIRSON21/mediasoft-intership2025/internal/repository"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

// warehouseService предоставляет методы для работы с складами.
type warehouseService struct {
	repo repository.WarehouseRepository
}

// NewWarehouseService создает новый экземпляр warehouseService.
func NewWarehouseService(repo repository.WarehouseRepository) handler.WarehouseService {
	return &warehouseService{
		repo: repo,
	}
}

// GetWarehouses возвращает список складов.
func (s *warehouseService) GetWarehouses(ctx context.Context) ([]*dto.WarehouseAtListResponse, error) {
	log := logger.GetLogger().With(zap.String("op", "service.warehouseService.GetWarehouses"))

	warehouses, err := s.repo.GetWarehouses(ctx)
	if err != nil {
		log.Error("error while getting warehouses", zap.String("err", err.Error()))
		return []*dto.WarehouseAtListResponse{}, err
	}
	warehousesResp := createWarehouseListResponse(warehouses)

	return warehousesResp, nil
}

// createWarehouseListResponse преобразует список складов в ответ с параметрами.
func createWarehouseListResponse(warehouses []*domain.Warehouse) []*dto.WarehouseAtListResponse {
	warehousesResp := make([]*dto.WarehouseAtListResponse, 0, len(warehouses))

	for _, v := range warehouses {
		warehousesResp = append(warehousesResp, &dto.WarehouseAtListResponse{
			ID:      v.ID.String(),
			Address: v.Address,
		})
	}
	return warehousesResp
}

// CreateWarehouse создает новый склад в репозитории.
func (s *warehouseService) CreateWarehouse(ctx context.Context, request *dto.WarehouseRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.warehouseService.CreateWarehouse"))

	warehouse := domain.Warehouse{
		Address: request.Address,
	}

	if err := s.repo.CreateWarehouse(ctx, &warehouse); err != nil {
		log.Error("error while creating warehouse", zap.String("err", err.Error()))
		return err
	}

	return nil
}
