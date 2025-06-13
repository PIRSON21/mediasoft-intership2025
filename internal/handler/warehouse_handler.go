package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PIRSON21/mediasoft-go/internal/service"
)

type WarehouseHandler struct {
	Service *service.WarehouseService
}

func (h *WarehouseHandler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	warehouses := h.Service.GetWarehouses()

	if err := json.NewEncoder(w).Encode(warehouses); err != nil {
		r.Response.StatusCode = http.StatusBadRequest
		json.NewEncoder(w).Encode(map[string]string{
			"error": "error while encoding warehouses",
		})
	}
}
