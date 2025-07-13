package repository

import (
	"context"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
)

// ProductRepository - интерфейс для работы с продуктами.
type ProductRepository interface {
	GetProducts(context.Context) ([]*domain.Product, error)
	AddProduct(context.Context, *domain.Product) error
	UpdateProduct(context.Context, *domain.Product) error
}
