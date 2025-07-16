package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/internal/middleware"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/render"
	"go.uber.org/zap"
)

// AnalyticsService определяет методы для работы с аналитикой складов.
//
//go:generate mockery init github.com/PIRSON21/mediasoft-intership2025/internal/handler
type AnalyticsService interface {
	AddProductSell(invs []*domain.Inventory)
	GetWarehouseAnalytics(ctx context.Context, warehouseID string) (*dto.WarehouseAnalyticsResponse, error)
	GetTopWarehouses(ctx context.Context, limit int) ([]*dto.WarehouseAnalyticsAtListResponse, error)
}

// AnalyticsHandler обрабатывает запросы, связанные с аналитикой складов.
type AnalyticsHandler struct {
	service AnalyticsService
}

// NewAnalyticsHandler создает новый экземпляр AnalyticsHandler с заданным сервисом аналитики.
func NewAnalyticsHandler(service AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// GetWarehouseAnalytics обрабатывает запросы на получение аналитики по складу.
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

// GetTopWarehouses обрабатывает запросы на получение списка топ-складов.
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

// parseLimitParams извлекает параметр limit из запроса и возвращает его значение.
func parseLimitParams(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 10, nil
	}

	return strconv.Atoi(limitStr)
}
