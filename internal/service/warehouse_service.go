package service

import "github.com/PIRSON21/mediasoft-go/internal/repository"

type WarehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

type WarehouseAtListResponse struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

func (s *WarehouseService) GetWarehouses() []*WarehouseAtListResponse {
	warehouses := s.repo.GetWarehouses()
	warehousesResp := make([]*WarehouseAtListResponse, 0, len(warehouses))

	for _, v := range warehouses {
		warehousesResp = append(warehousesResp, &WarehouseAtListResponse{
			ID:      v.ID,
			Address: v.Address,
		})
	}

	return warehousesResp
}
