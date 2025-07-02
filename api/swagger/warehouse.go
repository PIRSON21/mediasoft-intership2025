package swagger

import "github.com/PIRSON21/mediasoft-intership2025/internal/dto"

// WarehousesResponse represents list of warehouses.
// swagger:response WarehousesResponse
type WarehousesResponse struct {
	// in: body
	Body []dto.WarehouseAtListResponse
}
