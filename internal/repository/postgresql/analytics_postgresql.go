package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

// AddProductSell добавляет информацию о продаже продуктов в аналитику.
func (db *Postgres) AddProductSell(invs []*domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.AddProductSell"),
	)

	stmt, values := getAddProductSellStatement(invs)

	tag, err := db.pool.Exec(context.Background(), stmt, values...)
	if err != nil {
		log.Error("executing statement", zap.Error(err), zap.String("stmt", stmt))
		return err
	}
	if int(tag.RowsAffected()) != len(invs) {
		log.Error("not all product was added to analytics", zap.Int("want", len(invs)), zap.Int64("actual", tag.RowsAffected()))
	}

	return nil
}

// getAddProductSellStatement формирует SQL-запрос для добавления информации
// о продаже продуктов в аналитику.
func getAddProductSellStatement(invs []*domain.Inventory) (string, []any) {
	var (
		cursor = 1
		rows   []string
		values []any
	)

	query := `INSERT INTO analytics(warehouse_id, product_id, product_count, product_price) VALUES `

	for _, inv := range invs {
		price := inv.ProductPrice
		if inv.ProductSale != 0 {
			price = price - (price * float64(inv.ProductSale) / 100)
		}
		row := fmt.Sprintf("($%d, $%d, $%d, $%d)", cursor, cursor+1, cursor+2, cursor+3)
		rows = append(rows, row)
		values = append(values, inv.Warehouse.ID.String(), inv.Product.ID.String(), inv.ProductCount, price)

		cursor += 4
	}

	stmt := query + strings.Join(rows, ", ")

	return stmt, values
}

func (db *Postgres) GetWarehouseAnalytics(ctx context.Context, warehouseID string) ([]*domain.Analytics, error) {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetWarehouseAnalytics"),
	)
	var res []*domain.Analytics

	stmt := `
	SELECT inv.warehouse_id, p.product_id, p.product_name, a.product_count, a.product_price
	FROM inventory inv
	JOIN product p USING (product_id)
	JOIN analytics a USING (warehouse_id, product_id)
	WHERE warehouse_id = $1
	`

	rows, err := db.pool.Query(ctx, stmt, warehouseID)
	if err != nil {
		log.Error("error while getting analytics rows", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		anal := &domain.Analytics{
			Product:   &domain.Product{},
			Warehouse: &domain.Warehouse{},
		}

		err = rows.Scan(&anal.Warehouse.ID, &anal.Product.ID, &anal.Product.Name, &anal.ProductCount, &anal.ProductPrice)
		if err != nil {
			log.Error("error while scanning row", zap.Error(err))
			continue
		}

		res = append(res, anal)

	}

	if rows.Err() != nil {
		log.Error("error after scanning rows", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return res, nil
}

// GetTopWarehouses возвращает топ limit складов по сумме продаж продуктов.
func (db *Postgres) GetTopWarehouses(ctx context.Context, limit int) ([]*dto.WarehouseAnalyticsAtListResponse, error) {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetTopWarehouses"),
	)

	stmt := `
	SELECT
	w.warehouse_id,
	w.warehouse_address,
	COALESCE(SUM(a.product_price), 0) AS warehouse_total_sum
	FROM warehouse w
	LEFT JOIN analytics a USING (warehouse_id)
	GROUP BY warehouse_id
	ORDER BY warehouse_total_sum DESC
	LIMIT $1
	`

	rows, err := db.pool.Query(ctx, stmt, limit)
	if err != nil {
		log.Error("error while getting top warehouses", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var res []*dto.WarehouseAnalyticsAtListResponse
	for rows.Next() {
		warehouse := &dto.WarehouseAnalyticsAtListResponse{
			WarehouseID:       "",
			WarehouseAddress:  "",
			WarehouseTotalSum: 0,
		}

		err = rows.Scan(&warehouse.WarehouseID, &warehouse.WarehouseAddress, &warehouse.WarehouseTotalSum)
		if err != nil {
			log.Error("error while scanning row", zap.Error(err))
			continue
		}

		res = append(res, warehouse)
	}
	if rows.Err() != nil {
		log.Error("error after scanning rows", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	if len(res) == 0 {
		log.Warn("no warehouses found")
		return nil, nil
	}
	return res, nil
}
