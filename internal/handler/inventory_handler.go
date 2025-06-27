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
