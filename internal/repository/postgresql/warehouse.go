package postgresql

import "github.com/PIRSON21/mediasoft-go/internal/domain"

func (db *Postgres) GetWarehouses() []*domain.Warehouse {
	return []*domain.Warehouse{
		createWarehouse(1, "aboba"),
		createWarehouse(2, "boba"),
	}
}

func createWarehouse(id int, address string) *domain.Warehouse {
	return &domain.Warehouse{
		ID:      id,
		Address: address,
	}
}
