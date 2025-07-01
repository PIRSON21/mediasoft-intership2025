package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Postgres) GetWarehouses(ctx context.Context) ([]*domain.Warehouse, error) {
	var warehouses []*domain.Warehouse
	log := logger.GetLogger().With(zap.String("op", "repository.postgres.GetWarehouses"))

	stmt := `SELECT warehouse_id, warehouse_address FROM warehouse`

	rows, err := db.pool.Query(ctx, stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return warehouses, nil
		}
		log.Error("error while getting warehouses", zap.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse domain.Warehouse
		err := rows.Scan(&warehouse.ID, &warehouse.Address)
		if err != nil {
			log.Error("error while parsing warehouse", zap.String("err", err.Error()))
			continue
		}
		warehouses = append(warehouses, &warehouse)
	}

	log.Debug("warehouses written successfully", zap.Any("warehouses", warehouses))

	return warehouses, nil
}

func (db *Postgres) CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) error {
	log := logger.GetLogger().With(zap.String("op", "repository.Postgres.CreateWarehouse"))

	stmt := fmt.Sprintf(
		`INSERT INTO warehouse(warehouse_address)
		VALUES ($1)
	`)

	_, err := db.pool.Exec(ctx, stmt, warehouse.Address)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return custErr.ErrWarehouseAlreadyExists
			}
		}
		log.Error("error while creating warehouse", zap.String("err", err.Error()))
		return err
	}

	return nil
}
