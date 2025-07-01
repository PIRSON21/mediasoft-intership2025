package service

import (
	"context"
	"fmt"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	"github.com/PIRSON21/mediasoft-intership2025/internal/repository"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InventoryService struct {
	analytics *AnalyticsService
	repo      repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository, analytics *AnalyticsService) *InventoryService {
	return &InventoryService{
		analytics: analytics,
		repo:      repo,
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

func (s *InventoryService) CalculateCart(ctx context.Context, cartReq *dto.CartRequest) (*dto.CartResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.InventoryService.CalculateCart"),
	)
	_ = log

	cart, err := parseCartRequestToDomain(cartReq)
	if err != nil {
		log.Error("error while parsing cart request to domain", zap.Error(err))
		return nil, err
	}

	err = s.repo.GetPriceAndDiscount(ctx, cart)
	if err != nil {
		log.Error("error while getting price and discount from repository", zap.Error(err))
		return nil, err
	}

	resp := parseDomainToCartResponse(cart)

	return resp, nil
}

func parseCartRequestToDomain(req *dto.CartRequest) ([]*domain.Inventory, error) {
	var inv []*domain.Inventory

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		return nil, err
	}
	warehouse := &domain.Warehouse{
		ID: warehouseID,
	}

	for _, v := range req.Products {
		product, err := parseProductFromCartToDomain(v, warehouse)
		if err != nil {
			return nil, err
		}

		inv = append(inv, product)
	}

	return inv, nil
}

func parseProductFromCartToDomain(prod *dto.ProductInCartRequest, warehouse *domain.Warehouse) (*domain.Inventory, error) {
	productID, err := uuid.Parse(prod.ProductID)
	if err != nil {
		return nil, err
	}

	return &domain.Inventory{
		Warehouse: warehouse,
		Product: &domain.Product{
			ID: productID,
		},
		ProductCount: *prod.Count,
	}, nil
}

func parseDomainToCartResponse(invs []*domain.Inventory) *dto.CartResponse {
	var (
		resp               dto.CartResponse
		totalPrice         float64
		totalDiscountPrice float64
	)

	for _, inv := range invs {
		discountPrice := inv.ProductPrice
		if inv.ProductSale != 0 {
			discountPrice = inv.ProductPrice - (inv.ProductPrice * float64(inv.ProductSale) / 100)
		}

		fullPrice := inv.ProductPrice * float64(inv.ProductCount)
		discountFullPrice := discountPrice * float64(inv.ProductCount)

		prod := &dto.ProductInCartResponse{
			ProductID:         inv.Product.ID.String(),
			Count:             inv.ProductCount,
			FullPrice:         fullPrice,
			PriceWithDiscount: discountFullPrice,
		}
		resp.Products = append(resp.Products, prod)

		totalPrice += fullPrice
		totalDiscountPrice += discountFullPrice
	}

	resp.TotalProductPrice = totalPrice
	resp.TotalProductPriceWithDiscount = totalDiscountPrice

	return &resp
}

func (s *InventoryService) GetProductsAtWarehouse(ctx context.Context, params *dto.Pagination, warehouseID string) (*dto.ProductsResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.InventoryService.GetProducts"),
	)

	products, err := s.repo.GetProductsAtWarehouse(ctx, params, warehouseID)
	if err != nil {
		log.Error("error while getting products from repository", zap.Error(err))
		return nil, err
	}

	resp := parseProductsToResponse(products, params)

	return resp, nil
}

func parseProductsToResponse(prods []*domain.Inventory, params *dto.Pagination) *dto.ProductsResponse {
	resp := &dto.ProductsResponse{
		Page:     params.Page,
		Limit:    params.Limit,
		Products: make([]*dto.ProductAtList, 0),
	}

	for _, inv := range prods {
		discountPrice := inv.ProductPrice
		if inv.ProductSale != 0 {
			discountPrice = inv.ProductPrice - (inv.ProductPrice * float64(inv.ProductSale) / 100)
		}

		prod := dto.ProductAtList{
			ProductID:                inv.Product.ID.String(),
			ProductName:              inv.Product.Name,
			ProductPrice:             inv.ProductPrice,
			ProductPriceWithDiscount: discountPrice,
		}
		resp.Products = append(resp.Products, &prod)
	}

	return resp
}

func (s *InventoryService) BuyProducts(ctx context.Context, cart *dto.CartRequest) (*dto.CartResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "service.InventoryService.BuyProducts"),
	)

	invs, err := parseCartRequestToDomain(cart)
	if err != nil {
		log.Error("error while parsing cart to domain", zap.Error(err))
		return nil, err
	}

	err = s.repo.BuyProducts(ctx, invs)
	if err != nil {
		log.Error("error while changing products count in repository", zap.Error(err))
		return nil, err
	}

	go s.analytics.AddProductSell(invs)

	response := parseDomainToCartResponse(invs)

	return response, nil
}
