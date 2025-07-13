package domain

import "github.com/google/uuid"

// Warehouse представляет склад с его деталями.
type Warehouse struct {
	ID      uuid.UUID
	Address string
}
