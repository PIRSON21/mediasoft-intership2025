package repository

import "github.com/PIRSON21/mediasoft-go/internal/repository/postgresql"

type Repository interface {
	WarehouseRepository
}

func NewRepository() (Repository, error) {
	return postgresql.NewPostgres()
}
