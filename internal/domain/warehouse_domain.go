package domain

import "github.com/google/uuid"

type Warehouse struct {
	ID      uuid.UUID
	Address string
}
