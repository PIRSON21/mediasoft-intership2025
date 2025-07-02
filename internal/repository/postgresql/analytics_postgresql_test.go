package postgresql

import (
	"testing"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAddProductSellStatement(t *testing.T) {
	invs := []*domain.Inventory{
		{
			Warehouse:    &domain.Warehouse{ID: uuid.New()},
			Product:      &domain.Product{ID: uuid.New()},
			ProductCount: 10,
			ProductPrice: 100.0,
			ProductSale:  15,
		},
		{
			Warehouse:    &domain.Warehouse{ID: uuid.New()},
			Product:      &domain.Product{ID: uuid.New()},
			ProductCount: 5,
			ProductPrice: 200.0,
			ProductSale:  0,
		},
	}

	stmt, values := getAddProductSellStatement(invs)
	expectedStmt := `INSERT INTO analytics(warehouse_id, product_id, product_count, product_price) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`
	expectedValues := []any{
		invs[0].Warehouse.ID.String(),
		invs[0].Product.ID.String(),
		invs[0].ProductCount,
		invs[0].ProductPrice - (invs[0].ProductPrice * float64(invs[0].ProductSale) / 100),
		invs[1].Warehouse.ID.String(),
		invs[1].Product.ID.String(),
		invs[1].ProductCount,
		invs[1].ProductPrice,
	}

	require.Equal(t, expectedStmt, stmt)
	require.Equal(t, expectedValues, values)
}
