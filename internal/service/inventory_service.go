package service

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	"github.com/PIRSON21/mediasoft-go/internal/dto"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

type InventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) CreateInventory(ctx context.Context, request *dto.InventoryCreateRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.InventoryService.CreateInventory"))

	inventory := parseInventoryRequestToDomain(request)

	err := s.repo.CreateInventory(ctx, inventory)
	if err != nil {
		log.Error("error while creating inventory in repository", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func parseInventoryRequestToDomain(req *dto.InventoryCreateRequest) *domain.Inventory {
	return &domain.Inventory{
		ProductID:    *req.ProductID,
		WarehouseID:  *req.WarehouseID,
		ProductCount: *req.Count,
		ProductPrice: *req.Price,
	}
}

func (s *InventoryService) ChangeProductCount(ctx context.Context, request *dto.ChangeProductCountRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.InventoryService.ChangeProductCount"))

	inventory := parseChangeProductCountRequestToDomain(request)

	err := s.repo.ChangeProductCount(ctx, inventory)
	if err != nil {
		log.Error("error while changing product count in repository", zap.Error(err))
		return err
	}

	return nil
}

func parseChangeProductCountRequestToDomain(req *dto.ChangeProductCountRequest) *domain.Inventory {
	return &domain.Inventory{
		ProductID:    *req.ProductID,
		WarehouseID:  *req.WarehouseID,
		ProductCount: *req.Count,
	}
}
