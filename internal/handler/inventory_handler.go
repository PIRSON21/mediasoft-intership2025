package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/internal/middleware"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InventoryService определяет методы для работы с инвентаризацией товаров на складах.
//
//go:generate mockery init github.com/PIRSON21/mediasoft-intership2025/internal/handler
type InventoryService interface {
	CreateInventory(ctx context.Context, request *dto.InventoryCreateRequest) error
	ChangeProductCount(ctx context.Context, request *dto.ChangeProductCountRequest) error
	AddDiscountToProduct(ctx context.Context, request *dto.DiscountToProductRequest) error
	GetProductFromWarehouse(ctx context.Context, warehouseID, productID string) (*dto.ProductFromWarehouseResponse, error)
	GetProductsAtWarehouse(ctx context.Context, params *dto.Pagination, warehouseID string) (*dto.ProductsResponse, error)
	CalculateCart(ctx context.Context, request *dto.CartRequest) (*dto.CartResponse, error)
	BuyProducts(ctx context.Context, request *dto.CartRequest) (*dto.CartResponse, error)
}

// InventoryHandler обрабатывает запросы, связанные с инвентаризацией товаров на складах.
type InventoryHandler struct {
	service InventoryService
}

// NewInventoryHandler создает новый экземпляр InventoryHandler с заданным сервисом инвентаризации.
func NewInventoryHandler(service InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

// CreateInventory обрабатывает запросы на создание инвентаризации товара на складе.
func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.CreateInventory"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	invRequest, err := parseInventory(r.Body)
	if err != nil {
		custErr.UnnamedError(w, http.StatusBadRequest, "error while parsing body")
		return
	}

	validErr := validateInventoryCreateRequest(invRequest)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err = h.service.CreateInventory(r.Context(), invRequest)
	if err != nil {
		if errors.Is(err, custErr.ErrInventoryAlreadyExists) {
			custErr.UnnamedError(w, http.StatusConflict, "this inventory already exists")
			return
		}
		if errors.Is(err, custErr.ErrForeignKey) {
			custErr.UnnamedError(w, http.StatusBadRequest, "wrong product ID or warehouse ID")
			return
		}

		log.Error("error while creating inventory", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while creating inventory")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// parseInventory извлекает данные инвентаризации из запроса и возвращает их в виде dto.InventoryCreateRequest.
func parseInventory(r io.Reader) (*dto.InventoryCreateRequest, error) {
	var invReq dto.InventoryCreateRequest

	if err := json.NewDecoder(r).Decode(&invReq); err != nil {
		return nil, err
	}

	return &invReq, nil
}

// validateInventoryCreateRequest проверяет корректность данных инвентаризации.
func validateInventoryCreateRequest(req *dto.InventoryCreateRequest) map[string]string {
	validErr := make(map[string]string, 0)

	if req.ProductID == "" {
		validErr["product_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.ProductID); err != nil {
		validErr["product_id"] = "invalid product ID"
	}

	if req.WarehouseID == "" {
		validErr["warehouse_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.WarehouseID); err != nil {
		validErr["warehouse_id"] = "invalid warehouse ID"
	}

	if req.Count == nil {
		validErr["product_count"] = "this field cannot be empty"
	} else if *req.Count < 0 {
		validErr["product_count"] = "invalid product count"
	}

	if req.Price == nil {
		validErr["product_price"] = "this field cannot be empty"
	} else if *req.Price < 0 {
		validErr["product_price"] = "invalid product price"
	}

	if len(validErr) > 0 {
		return validErr
	}
	return nil
}

// ChangeProductCount обрабатывает запросы на изменение количества товара на складе.
func (h *InventoryHandler) ChangeProductCount(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.ChangeProductCount"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	prodReq, err := parseChangeProductCountRequest(r.Body)
	if err != nil {
		log.Error("error while parsing JSON", zap.Error(err))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, "cannot parse JSON")
		return
	}

	validErr := validateChangeProductCountRequest(prodReq)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err = h.service.ChangeProductCount(r.Context(), prodReq)
	if err != nil {
		if errors.Is(err, custErr.ErrInventoryNotFound) {
			custErr.UnnamedError(w, http.StatusNotFound, "there is no information about this product on warehouse")
			return
		}
		log.Error("error while change product count", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while changing product count")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseChangeProductCountRequest извлекает данные запроса на изменение количества товара и возвращает их в виде dto.ChangeProductCountRequest.
func parseChangeProductCountRequest(r io.Reader) (*dto.ChangeProductCountRequest, error) {
	var request dto.ChangeProductCountRequest

	if err := json.NewDecoder(r).Decode(&request); err != nil {
		return nil, err
	}

	return &request, nil
}

// validateChangeProductCountRequest проверяет корректность данных запроса на изменение количества товара.
func validateChangeProductCountRequest(req *dto.ChangeProductCountRequest) map[string]string {
	validErr := make(map[string]string, 0)

	if req.ProductID == "" {
		validErr["product_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.ProductID); err != nil {
		validErr["product_id"] = "invalid product ID"
	}

	if req.WarehouseID == "" {
		validErr["warehouse_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.WarehouseID); err != nil {
		validErr["warehouse_id"] = "invalid warehouse ID"
	}

	if req.Count == nil {
		validErr["product_count"] = "this field cannot be empty"
	} else if *req.Count < 0 {
		validErr["product_count"] = "invalid product count"
	}

	if len(validErr) != 0 {
		return validErr
	}

	return nil
}

// AddDiscountToProduct обрабатывает запросы на добавление скидок к товарам на складе.
func (h *InventoryHandler) AddDiscountToProduct(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.AddDiscountToProduct"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	discountReq, err := parseDiscountRequest(r.Body)
	if err != nil {
		log.Error("error while parsing discounts", zap.Error(err))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, "error while parsing request")
		return
	}

	validErr := validateDiscountRequest(discountReq)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err = h.service.AddDiscountToProduct(r.Context(), discountReq)
	if err != nil {
		if errors.Is(err, custErr.ErrInventoryNotFound) {
			custErr.UnnamedError(w, http.StatusBadRequest, "there is no some products in warehouse")
			return
		}
		log.Error("error while adding discounts", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error adding discounts")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// parseDiscountRequest извлекает данные запроса на скидки и возвращает их в виде dto.DiscountToProductRequest.
func parseDiscountRequest(r io.Reader) (*dto.DiscountToProductRequest, error) {
	var discounts dto.DiscountToProductRequest

	err := json.NewDecoder(r).Decode(&discounts)
	if err != nil {
		return nil, err
	}

	return &discounts, nil
}

// validateDiscountRequest проверяет корректность данных запроса на скидки.
func validateDiscountRequest(req *dto.DiscountToProductRequest) map[string]any {
	validErr := make(map[string]any)

	if req.WarehouseID == "" {
		validErr["warehouse_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.WarehouseID); err != nil {
		validErr["warehouse_id"] = "invalid warehouse ID"
	}

	if len(req.Discounts) == 0 {
		validErr["discounts"] = "there is no discounts"
	}

	discountsErr := validateDiscounts(req)

	if discountsErr != nil {
		validErr["discounts"] = discountsErr
	}

	if len(validErr) != 0 {
		return validErr
	}

	return nil
}

// validateDiscounts проверяет корректность каждого элемента в списке скидок.
func validateDiscounts(req *dto.DiscountToProductRequest) map[int]any {
	discountsErr := make(map[int]any)
	for idx, discount := range req.Discounts {
		discountErr := validateDiscount(discount)
		if discountErr != nil {
			discountsErr[idx] = discountErr
		}
	}

	if len(discountsErr) != 0 {
		return discountsErr
	}

	return nil
}

// validateDiscount проверяет корректность данных скидки.
func validateDiscount(discount *dto.Discount) map[string]string {
	discountErr := make(map[string]string)

	if discount.ProductID == "" {
		discountErr["product_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(discount.ProductID); err != nil {
		discountErr["product_id"] = "invalid product ID"
	}

	if discount.DiscountValue == nil {
		discountErr["discount"] = "this field cannot be empty"
	} else if *discount.DiscountValue < 0 || *discount.DiscountValue > 100 {
		discountErr["discount"] = "discount must be greater than 0 and less than 100"
	}

	if len(discountErr) != 0 {
		return discountErr
	}

	return nil
}

// GetProductFromWarehouse обрабатывает запросы на получение информации о товаре на складе.
func (h *InventoryHandler) GetProductFromWarehouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("product_id") != "" {
		h.GetOneProductFromWarehouse(w, r)
	} else {
		h.GetProductsAtWarehouse(w, r)
	}
}

// GetOneProductFromWarehouse обрабатывает запросы на получение информации о конкретном товаре на складе.
func (h *InventoryHandler) GetOneProductFromWarehouse(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.GetProductFromWarehouse"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	warehouseID, err := parseWarehouseIDFromURL(r)
	if err != nil {
		custErr.UnnamedError(w, http.StatusBadRequest, err.Error())
		return
	}

	productID, err := parseProductIDFromQuery(r)
	if err != nil {
		custErr.UnnamedError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.GetProductFromWarehouse(r.Context(), warehouseID, productID)
	if err != nil {
		if errors.Is(err, custErr.ErrProductNotFound) {
			custErr.UnnamedError(w, http.StatusNotFound, "there is no such product in the warehouse")
			return
		}
		log.Error("error while getting product", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while getting product")
		return
	}

	render.JSON(w, http.StatusOK, product)
}

// parseWarehouseIDFromURL извлекает идентификатор склада из URL запроса и возвращает его.
func parseWarehouseIDFromURL(r *http.Request) (string, error) {
	splits := strings.Split(r.URL.Path, "/")

	warehouseID := splits[len(splits)-1]

	if err := uuid.Validate(warehouseID); err != nil {
		return "", fmt.Errorf("warehouse id is not valid")
	}

	return warehouseID, nil
}

// parseProductIDFromQuery извлекает идентификатор продукта из параметров запроса и возвращает его.
func parseProductIDFromQuery(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	productID := r.FormValue("product_id")
	if productID == "" {
		return "", fmt.Errorf("empty product id")
	} else if err := uuid.Validate(productID); err != nil {
		return "", fmt.Errorf("product id is not valid")
	}

	return productID, nil
}

// GetProductsAtWarehouse обрабатывает запросы на получение списка продуктов на складе.
func (h *InventoryHandler) GetProductsAtWarehouse(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.GetProducts"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	warehouseID, err := parseWarehouseIDFromURL(r)
	if err != nil {
		log.Error("error while parsing warehouseID", zap.Error(err))
		custErr.UnnamedError(w, http.StatusBadRequest, "wrong warehouseID")
		return
	}

	params := parseParams(r)

	response, err := h.service.GetProductsAtWarehouse(r.Context(), params, warehouseID)
	if err != nil {
		log.Error("error while getting products", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while getting products")
		return
	}

	render.JSON(w, http.StatusOK, response)
}

// parseParams извлекает параметры пагинации из запроса и возвращает их в виде dto.Pagination.
func parseParams(r *http.Request) *dto.Pagination {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := limit * (page - 1)

	return &dto.Pagination{
		Page:   page,
		Offset: offset,
		Limit:  limit,
	}
}

// CalculateCart обрабатывает запросы на расчет стоимости товаров в корзине.
func (h *InventoryHandler) CalculateCart(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.CalculateSum"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cartReq, err := parseCartRequest(r.Body)
	if err != nil {
		log.Error("error while parsing cart", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while parsing cart")
		return
	}

	validErr := validateCartRequest(cartReq)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	resp, err := h.service.CalculateCart(r.Context(), cartReq)
	if err != nil {
		log.Error("error while calculating cart", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while calculating cart")
		return
	}

	render.JSON(w, http.StatusOK, resp)
}

// parseCartRequest извлекает данные корзины из запроса и возвращает их в виде dto.CartRequest.
func parseCartRequest(r io.Reader) (*dto.CartRequest, error) {
	var cart dto.CartRequest

	err := json.NewDecoder(r).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// validateCartRequest проверяет корректность данных корзины.
func validateCartRequest(req *dto.CartRequest) map[string]any {
	validErr := make(map[string]any)
	var productsID []string

	if req.WarehouseID == "" {
		validErr["warehouse_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(req.WarehouseID); err != nil {
		validErr["warehouse_id"] = "invalid warehouse ID"
	}

	if len(req.Products) == 0 {
		validErr["products"] = "there is no products in cart"
	} else {
		productsErr := make(map[int]any)
		for idx, product := range req.Products {
			if slices.Contains(productsID, product.ProductID) {
				productsErr[idx] = map[string]string{"product_id": "product ID must be unique"}
				continue
			}
			productsID = append(productsID, product.ProductID)
			productErr := validateProductInCart(product)
			if productErr != nil {
				productsErr[idx] = productErr
			}
		}

		if len(productsErr) != 0 {
			validErr["products"] = productsErr
		}

	}

	if len(validErr) != 0 {
		return validErr
	}

	return nil
}

// validateProductInCart проверяет корректность данных продукта в корзине.
func validateProductInCart(product *dto.ProductInCartRequest) map[string]string {
	productErr := make(map[string]string)
	if product.ProductID == "" {
		productErr["product_id"] = "this field cannot be empty"
	} else if err := uuid.Validate(product.ProductID); err != nil {
		productErr["product_id"] = "invalid product ID"
	}

	if product.Count == nil {
		productErr["product_count"] = "this field cannot be empty"
	} else if *product.Count <= 0 {
		productErr["product_count"] = "product count must be greater than 0"
	}

	if len(productErr) != 0 {
		return productErr
	}

	return nil
}

// BuyProducts обрабатывает запросы на покупку товаров из корзины.
func (h *InventoryHandler) BuyProducts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.BuyProducts"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cart, err := parseCartRequest(r.Body)
	if err != nil {
		log.Error("error while parsing cart", zap.Error(err))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, "wrong request body")
		return
	}

	validErr := validateCartRequest(cart)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	response, err := h.service.BuyProducts(r.Context(), cart)
	if err != nil {
		if custErr.Any(err, custErr.ErrNotEnoughProductCount, custErr.ErrNotFoundProductAtWarehouse) {
			custErr.UnnamedError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("error in service module", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while buying products")
		return
	}

	render.JSON(w, http.StatusOK, response)
}
