package service

import (
	"context"
	"fmt"

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
		Product: &domain.Product{
			ID: productID,
		},
		Warehouse: &domain.Warehouse{
			ID: warehouseID,
		},
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
		Product: &domain.Product{
			ID: productID,
		},
		Warehouse: &domain.Warehouse{
			ID: warehouseID,
		},
		ProductCount: *req.Count,
	}, nil
}

func (s *InventoryService) AddDiscountToProduct(ctx context.Context, request *dto.DiscountToProductRequest) error {
	log := logger.GetLogger().With(
		zap.String("op", "service.InventoryService.AddDiscountToProduct"),
	)

	inventory, err := parseDiscountToInventory(request)
	if err != nil {
		log.Error("error while parsing discounts to inventory", zap.Error(err))
		return err
	}

	err = s.repo.AddDiscountToProducts(ctx, inventory)
	if err != nil {
		log.Error("error while adding discounts to repository", zap.Error(err))
		return err
	}

	return nil
}

func parseDiscountToInventory(req *dto.DiscountToProductRequest) ([]*domain.Inventory, error) {
	var inventory []*domain.Inventory

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		return nil, err
	}

	for _, discount := range req.Discounts {
		productID, err := uuid.Parse(discount.ProductID)
		if err != nil {
			continue
		}

		inv := &domain.Inventory{
			Product: &domain.Product{
				ID: productID,
			},
			Warehouse: &domain.Warehouse{
				ID: warehouseID,
			},
			ProductSale: *discount.DiscountValue,
		}

		inventory = append(inventory, inv)
	}

	if len(inventory) == 0 {
		return nil, fmt.Errorf("no one discount were parsed")
	}

	return inventory, nil
}

func (s *InventoryService) GetProductFromWarehouse(ctx context.Context, warehouseID, productID string) (*dto.ProductFromWarehouseResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.InventoryService.GetProductFromWarehouse"),
	)

	inventory, err := parseProductRequestToInventory(warehouseID, productID)
	if err != nil {
		log.Error("error while parsing productID and warehouseID", zap.Error(err))
		return nil, err
	}

	err = s.repo.GetProductFromWarehouse(ctx, inventory)
	if err != nil {
		log.Error("error while getting product from warehouse", zap.Error(err))
		return nil, err
	}

	response := parseProductFromWarehouseToResponse(inventory)

	return response, nil
}

func parseProductRequestToInventory(warehouseIDStr, productIDStr string) (*domain.Inventory, error) {
	warehouseID, err := uuid.Parse(warehouseIDStr)
	if err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, err
	}

	return &domain.Inventory{
		Product: &domain.Product{
			ID: productID,
		},
		Warehouse: &domain.Warehouse{
			ID: warehouseID,
		},
	}, nil
}

func parseProductFromWarehouseToResponse(inv *domain.Inventory) *dto.ProductFromWarehouseResponse {
	response := &dto.ProductFromWarehouseResponse{
		ProductID:          inv.Product.ID.String(),
		ProductName:        inv.Product.Name,
		ProductDescription: inv.Product.Description,
		ProductWeight:      inv.Product.Weight,
		ProductBarcode:     inv.Product.Barcode,
		ProductCount:       inv.ProductCount,
		ProductPrice:       inv.ProductPrice,
	}

	response.ProductParams = copyMap(inv.Product.Params)

	if inv.ProductSale == 0 {
		response.ProductPriceWithSale = inv.ProductPrice
	} else {
		response.ProductPriceWithSale = inv.ProductPrice * (1 - float64(inv.ProductSale)/100)
	}

	return response
}
