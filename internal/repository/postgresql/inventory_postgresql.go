package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Postgres) CreateInventory(ctx context.Context, inventory *domain.Inventory) error {
	log := logger.GetLogger().With(zap.String("op", "repository.Postgres.CreateInventory"))

	stmt := `
	INSERT INTO inventory(product_id, warehouse_id, product_count, product_price)
	VALUES ($1, $2, $3, $4)
	`

	tag, err := db.pool.Exec(ctx, stmt, inventory.Product.ID, inventory.Warehouse.ID, inventory.ProductCount, inventory.ProductPrice)
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

	tag, err := db.pool.Exec(ctx, stmt, &inventory.Product.ID, &inventory.Warehouse.ID, &inventory.ProductCount)
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

func (db *Postgres) AddDiscountToProducts(ctx context.Context, inventory []*domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.AddDiscountToProduct"),
	)

	transaction, err := db.pool.Begin(ctx)
	if err != nil {
		log.Error("error while beginning transaction", zap.Error(err))
		return err
	}

	for _, discount := range inventory {
		err = addDiscount(ctx, transaction, discount)
		if err != nil {
			log.Error("error while adding discount", zap.Error(err))
			transaction.Rollback(ctx)
			return fmt.Errorf("error while adding discount: %w", err)
		}
	}

	return transaction.Commit(ctx)
}

func addDiscount(ctx context.Context, conn pgx.Tx, discount *domain.Inventory) error {
	stmt := `
		UPDATE inventory SET product_sale = $1 WHERE warehouse_id = $2 AND product_id = $3
	`

	tag, err := conn.Exec(ctx, stmt, discount.ProductSale, discount.Warehouse.ID, discount.Product.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return custErr.ErrInventoryNotFound
	}

	return nil
}

func (db *Postgres) GetProductFromWarehouse(ctx context.Context, inventory *domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetProductFromWarehouse"),
	)

	stmt := `
	SELECT p.product_name, p.product_description, p.product_weight, p.product_params, p.product_barcode, inv.product_count, inv.product_price, inv.product_sale
	FROM inventory inv
	JOIN product p USING (product_id)
	WHERE product_id = $1 AND warehouse_id = $2
	`

	err := db.pool.QueryRow(ctx, stmt, inventory.Product.ID, inventory.Warehouse.ID).Scan(
		&inventory.Product.Name,
		&inventory.Product.Description,
		&inventory.Product.Weight,
		&inventory.Product.Params,
		&inventory.Product.Barcode,
		&inventory.ProductCount,
		&inventory.ProductPrice,
		&inventory.ProductSale,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custErr.ErrProductNotFound
		}
		log.Error("error while getting rows", zap.Error(err))
		return err
	}

	return nil
}
