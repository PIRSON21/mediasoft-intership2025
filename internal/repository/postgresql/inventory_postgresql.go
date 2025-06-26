package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Postgres) CreateInventory(ctx context.Context, inventory *domain.Inventory) error {
	log := logger.GetLogger().With(zap.String("op", "repository.Postgres.CreateInventory"))

	stmt := `
	INSERT INTO inventory(product_id, warehouse_id, product_count, product_price)
	VALUES ($1, $2, $3, $4)
	`

	tag, err := db.pool.Exec(ctx, stmt, inventory.ProductID, inventory.WarehouseID, inventory.ProductCount, inventory.ProductPrice)
	if err != nil {
		pgError := new(pgconn.PgError)
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return custErr.ErrInventoryAlreadyExists
			case "23503":
				return custErr.ErrForeignKey
			}
		}
		log.Error("error while executing statement", zap.Error(err))
		return err
	}

	if tag.RowsAffected() < 1 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (db *Postgres) ChangeProductCount(ctx context.Context, inventory *domain.Inventory) error {
	log := logger.GetLogger().With(zap.String("op", "repository.Postgres.ChangeProductCount"))

	// используется пользовательская функция. код в миграции 000004
	stmt := `SELECT increase_product_count($1, $2, $3)`

	tag, err := db.pool.Exec(ctx, stmt, &inventory.ProductID, &inventory.WarehouseID, &inventory.ProductCount)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if errors.As(err, &pgErr) {
			if pgErr.Code == "P0002" {
				return custErr.ErrInventoryNotFound
			}
		}
		log.Error("error while executing statement", zap.String("stmt", stmt), zap.Error(err))
		return err
	}

	if tag.RowsAffected() < 1 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
