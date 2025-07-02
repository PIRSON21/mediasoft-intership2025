package swagger

import "github.com/PIRSON21/mediasoft-intership2025/internal/dto"

// WarehouseAnalyticsResponse swagger response
// swagger:response WarehouseAnalyticsResponse
type WarehouseAnalyticsResponseWrapper struct {
	// in: body
	Body dto.WarehouseAnalyticsResponse
}

// WarehouseAnalyticsAtListResponse swagger response
// swagger:response WarehouseAnalyticsAtListResponse
type WarehouseAnalyticsAtListResponseWrapper struct {
	// in: body
	Body []dto.WarehouseAnalyticsAtListResponse
}
