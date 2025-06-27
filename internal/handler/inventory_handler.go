package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/PIRSON21/mediasoft-go/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/internal/middleware"
	"github.com/PIRSON21/mediasoft-go/internal/service"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/PIRSON21/mediasoft-go/pkg/render"
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
