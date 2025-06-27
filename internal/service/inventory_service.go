package service

import (
	"context"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	"github.com/PIRSON21/mediasoft-go/internal/dto"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/google/uuid"
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

	inventory, err := parseInventoryRequestToDomain(request)
	if err != nil {
		log.Error("error while parsing inventory request", zap.Error(err))
		return err
	}

	err = s.repo.CreateInventory(ctx, inventory)
	if err != nil {
		log.Error("error while creating inventory in repository", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func parseInventoryRequestToDomain(req *dto.InventoryCreateRequest) (*domain.Inventory, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, err
	}

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		return nil, err
	}

	return &domain.Inventory{
		ProductID:    productID,
		WarehouseID:  warehouseID,
		ProductCount: *req.Count,
		ProductPrice: *req.Price,
	}, nil
}

func (s *InventoryService) ChangeProductCount(ctx context.Context, request *dto.ChangeProductCountRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.InventoryService.ChangeProductCount"))

	inventory, err := parseChangeProductCountRequestToDomain(request)
	if err != nil {
		log.Error("error while parsing request", zap.Error(err))
		return err
	}

	err = s.repo.ChangeProductCount(ctx, inventory)
	if err != nil {
		log.Error("error while changing product count in repository", zap.Error(err))
		return err
	}

	return nil
}

func parseChangeProductCountRequestToDomain(req *dto.ChangeProductCountRequest) (*domain.Inventory, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, err
	}

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		return nil, err
	}

	return &domain.Inventory{
		ProductID:    productID,
		WarehouseID:  warehouseID,
		ProductCount: *req.Count,
	}, nil
}
