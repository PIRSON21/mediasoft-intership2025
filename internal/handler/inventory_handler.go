package handler

import (
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
	"github.com/PIRSON21/mediasoft-intership2025/internal/service"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.CreateInventory"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

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

func parseInventory(r io.Reader) (*dto.InventoryCreateRequest, error) {
	var invReq dto.InventoryCreateRequest

	if err := json.NewDecoder(r).Decode(&invReq); err != nil {
		return nil, err
	}

	return &invReq, nil
}

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

func (h *InventoryHandler) ChangeProductCount(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.ChangeProductCount"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

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

func parseChangeProductCountRequest(r io.Reader) (*dto.ChangeProductCountRequest, error) {
	var request dto.ChangeProductCountRequest

	if err := json.NewDecoder(r).Decode(&request); err != nil {
		return nil, err
	}

	return &request, nil
}

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

func (h *InventoryHandler) AddDiscountToProduct(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.AddDiscountToProduct"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

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

func parseDiscountRequest(r io.Reader) (*dto.DiscountToProductRequest, error) {
	var discounts dto.DiscountToProductRequest

	err := json.NewDecoder(r).Decode(&discounts)
	if err != nil {
		return nil, err
	}

	return &discounts, nil
}

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

func (h *InventoryHandler) GetProductFromWarehouse(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("product_id") != "" {
		h.GetOneProductFromWarehouse(w, r)
	} else {
		h.GetProductsAtWarehouse(w, r)
	}
}

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

func parseWarehouseIDFromURL(r *http.Request) (string, error) {
	splits := strings.Split(r.URL.Path, "/")

	warehouseID := splits[len(splits)-1]

	if err := uuid.Validate(warehouseID); err != nil {
		return "", fmt.Errorf("warehouse id is not valid")
	}

	return warehouseID, nil
}

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

func (h *InventoryHandler) CalculateCart(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.CalculateSum"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

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

func parseCartRequest(r io.Reader) (*dto.CartRequest, error) {
	var cart dto.CartRequest

	err := json.NewDecoder(r).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

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

func (h *InventoryHandler) BuyProducts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.InventoryHandler.BuyProducts"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

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

	// TODO: добавить аналитику
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
