package service

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	"github.com/PIRSON21/mediasoft-go/internal/dto"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

type WarehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

func (s *WarehouseService) GetWarehouses(ctx context.Context) ([]*dto.WarehouseAtListResponse, error) {
	log := logger.GetLogger().With(zap.String("op", "service.WarehouseService.GetWarehouses"))

	warehouses, err := s.repo.GetWarehouses(ctx)
	if err != nil {
		log.Error("error while getting warehouses", zap.String("err", err.Error()))
		return []*dto.WarehouseAtListResponse{}, err
	}
	warehousesResp := make([]*dto.WarehouseAtListResponse, 0, len(warehouses))

	for _, v := range warehouses {
		warehousesResp = append(warehousesResp, &dto.WarehouseAtListResponse{
			ID:      v.ID,
			Address: v.Address,
		})
	}

	return warehousesResp, nil
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, request *dto.WarehouseRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.WarehouseService.CreateWarehouse"))

	warehouse := domain.Warehouse{
		Address: request.Address,
	}

	if err := s.repo.CreateWarehouse(ctx, &warehouse); err != nil {
		log.Error("error while creating warehouse", zap.String("err", err.Error()))
		return err
	}

	return nil
}
