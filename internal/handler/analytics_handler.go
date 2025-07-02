package handler

import (
	"net/http"
	"strconv"

	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/internal/middleware"
	"github.com/PIRSON21/mediasoft-intership2025/internal/service"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/render"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) GetWarehouseAnalytics(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.AnalyticsHandler.GetWarehouseAnalytics"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	warehouseID, err := parseWarehouseIDFromURL(r)
	if err != nil {
		log.Error("error while parsing warehouseID", zap.Error(err))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response, err := h.service.GetWarehouseAnalytics(r.Context(), warehouseID)
	if err != nil {
		log.Error("error from service module", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while getting warehouse analytics")
		return
	}

	render.JSON(w, http.StatusOK, response)
}

func (h *AnalyticsHandler) GetTopWarehouses(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.AnalyticsHandler.GetTopWarehouses"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit, err := parseLimitParams(r)
	if err != nil {
		log.Error("error parsing limit", zap.Error(err))
		custErr.UnnamedError(w, http.StatusUnprocessableEntity, "wrong limit param")
		return
	}

	response, err := h.service.GetTopWarehouses(r.Context(), limit)
	if err != nil {
		log.Error("error from service module", zap.Error(err))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while getting top warehouses")
		return
	}

	render.JSON(w, http.StatusOK, response)
}

func parseLimitParams(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 10, nil
	}

	return strconv.Atoi(limitStr)
}
