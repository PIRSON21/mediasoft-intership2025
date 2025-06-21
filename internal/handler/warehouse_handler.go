package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PIRSON21/mediasoft-go/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/internal/service"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

type WarehouseHandler struct {
	Service *service.WarehouseService
}

func (h *WarehouseHandler) WarehousesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetWarehouses(w, r)
	case http.MethodPost:
		h.CreateWarehouse(w, r)
	}
}

func (h *WarehouseHandler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log := logger.GetLogger().With(
		zap.String("op", "handlers.WarehouseHandler.GetWarehouses"),
		zap.String("remoteAddr", r.RemoteAddr),
	)

	warehouses, err := h.Service.GetWarehouses(r.Context())
	if err != nil {
		log.Error("error while getting warehouses", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "error while getting warehouses",
		})
		return
	}

	if err := json.NewEncoder(w).Encode(warehouses); err != nil {
		log.Error("error while encoding warehouses", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "error while encoding warehouses",
		})
		return
	}

	log.Info(fmt.Sprintf("GET 200 %s", time.Since(start)))
}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var request dto.WarehouseRequest
	log := logger.GetLogger().With(
		zap.String("op", "handlers.WarehouseHandler.CreateWarehouse"),
		zap.String("remoteAddr", r.RemoteAddr),
	)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error("error while unmarshalling request", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("error while decoding warehouse: %q", err.Error()),
		})
		return
	}

	if validErr := validateWarehouse(&request); validErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validErr)
		return
	}

	err := h.Service.CreateWarehouse(context.TODO(), &request)
	if err != nil {
		if errors.Is(err, custErr.ErrWarehouseAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "warehouse already exists",
			})
			return
		}
		log.Error("error while creating warehouse", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "error while creating warehouse",
		})
		return
	}

	log.Info(fmt.Sprintf("POST 201 %s", time.Since(start)))
	w.WriteHeader(http.StatusCreated)
}

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
