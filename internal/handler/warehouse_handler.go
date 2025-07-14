package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/internal/middleware"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/render"
	"go.uber.org/zap"
)

// WarehouseService определяет методы для работы со складами.
//
//go:generate mockery init github.com/PIRSON21/mediasoft-intership2025/internal/handler
type WarehouseService interface {
	GetWarehouses(ctx context.Context) ([]*dto.WarehouseAtListResponse, error)
	CreateWarehouse(ctx context.Context, request *dto.WarehouseRequest) error
}

// WarehouseHandler обрабатывает запросы, связанные со складами.
type WarehouseHandler struct {
	Service WarehouseService
}

// NewWarehouseHandler создает новый экземпляр WarehouseHandler с заданным сервисом склада.
func NewWarehouseHandler(s WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		Service: s,
	}
}

// WarehousesHandler обрабатывает запросы на получение списка складов или создание нового склада.
func (h *WarehouseHandler) WarehousesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetWarehouses(w, r)
	case http.MethodPost:
		h.CreateWarehouse(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetWarehouses обрабатывает запросы на получение списка складов.
func (h *WarehouseHandler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handlers.WarehouseHandler.GetWarehouses"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	warehouses, err := h.Service.GetWarehouses(r.Context())
	if err != nil {
		log.Error("error while getting warehouses", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, fmt.Sprintf("error while getting warehouses: %q", err.Error()))
		return
	}

	if err := render.JSON(w, http.StatusOK, warehouses); err != nil {
		log.Error("error while encoding warehouses", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, fmt.Sprintf("error while encoding warehouses: %q", err.Error()))
		return
	}
}

// CreateWarehouse обрабатывает запросы на создание нового склада.
func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var request dto.WarehouseRequest
	log := logger.GetLogger().With(
		zap.String("op", "handlers.WarehouseHandler.CreateWarehouse"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Header.Get("Content-Type") != "application/json" {
		log.Error("bad format", zap.String("format", r.Header.Get("Content-Type")))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, "unsupported format")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("error while unmarshalling request", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, fmt.Sprintf("error while decoding warehouse: %q", err.Error()))
		return
	}

	if validErr := validateWarehouse(&request); validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err := h.Service.CreateWarehouse(context.TODO(), &request)
	if err != nil {
		if errors.Is(err, custErr.ErrWarehouseAlreadyExists) {
			custErr.UnnamedError(w, http.StatusConflict, "warehouse already exists")
			return
		}
		log.Error("error while creating warehouse", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, fmt.Sprintf("error while creating warehouse: %q", err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// validateWarehouse проверяет корректность данных склада.
func validateWarehouse(warehouse *dto.WarehouseRequest) map[string]string {
	validErr := make(map[string]string)
	if len(warehouse.Address) == 0 {
		validErr["address"] = "address cannot be empty"
	}

	if len(validErr) != 0 {
		return validErr
	}
	return nil
}
