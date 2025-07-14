package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PIRSON21/mediasoft-intership2025/internal/domain"
	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// CreateInventory создает новую запись в таблице inventory.
//
// Если запись с таким product_id и warehouse_id уже существует, то возвращает ошибку ErrInventoryAlreadyExists.
//
// Если warehouse_id или product_id не существует, то возвращает ошибку ErrForeignKey.
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

// ChangeProductCount изменяет количество продукта на складе.
//
// Если количество меньше нуля, то возвращает ошибку ErrNotEnoughProductCount.
//
// Если запись не найдена, то возвращает ErrInventoryNotFound.
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

// AddDiscountToProducts добавляет скидку на продукты в инвентаре.
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

// addDiscount добавляет скидку на продукт в инвентаре.
//
// Если запись не найдена, то возвращает ErrInventoryNotFound.
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

// GetProductFromWarehouse получает информацию о продукте на складе.
//
// Если продукт не найден, то возвращает ErrProductNotFound.
func (db *Postgres) GetProductFromWarehouse(ctx context.Context, inventory *domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetProductFromWarehouse"),
	)

	inv := struct {
		ProductName        string
		ProductDescription string
		ProductWeight      sql.NullFloat64
		ProductParams      map[string]any
		ProductBarcode     string
		ProductCount       sql.NullInt64
		ProductPrice       sql.NullFloat64
		ProductSale        sql.NullInt64
	}{}

	stmt := `
	SELECT p.product_name, p.product_description, p.product_weight, p.product_params, p.product_barcode, inv.product_count, inv.product_price, inv.product_sale
	FROM inventory inv
	JOIN product p USING (product_id)
	WHERE product_id = $1 AND warehouse_id = $2
	`

	err := db.pool.QueryRow(ctx, stmt, inventory.Product.ID, inventory.Warehouse.ID).Scan(
		&inv.ProductName,
		&inv.ProductDescription,
		&inv.ProductWeight,
		&inv.ProductParams,
		&inv.ProductBarcode,
		&inv.ProductCount,
		&inv.ProductPrice,
		&inv.ProductSale,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custErr.ErrProductNotFound
		}
		log.Error("error while getting rows", zap.Error(err))
		return err
	}

	inventory.Product.Name = inv.ProductName
	inventory.Product.Description = inv.ProductDescription
	if inv.ProductWeight.Valid {
		inventory.Product.Weight = inv.ProductWeight.Float64
	}
	if inv.ProductParams != nil {
		inventory.Product.Params = inv.ProductParams
	}
	inventory.Product.Barcode = inv.ProductBarcode
	if inv.ProductCount.Valid {
		inventory.ProductCount = int(inv.ProductCount.Int64)
	} else {
		inventory.ProductCount = 0
	}
	if inv.ProductPrice.Valid {
		inventory.ProductPrice = inv.ProductPrice.Float64
	} else {
		inventory.ProductPrice = 0
	}
	if inv.ProductSale.Valid {
		inventory.ProductSale = int(inv.ProductSale.Int64)
	} else {
		inventory.ProductSale = 0
	}

	return nil
}

// GetPriceAndDiscount получает цену и скидку для продуктов в инвентаре.
//
// Если запись не найдена, то возвращает ErrInventoryNotFound.
func (db *Postgres) GetPriceAndDiscount(ctx context.Context, invs []*domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetPriceAndDiscount"),
	)

	if len(invs) == 0 {
		return nil
	}

	invMap := make(map[string]*domain.Inventory)
	var productsID []string
	warehouseID := invs[0].Warehouse.ID

	for _, inv := range invs {
		productID := inv.Product.ID.String()
		invMap[productID] = inv
		productsID = append(productsID, productID)
	}

	stmt := `
		SELECT product_id, product_price, product_sale, product_count FROM inventory
		WHERE warehouse_id = $1 AND product_id = ANY($2)
	`

	rows, err := db.pool.Query(ctx, stmt, warehouseID, productsID)
	if err != nil {
		log.Error("error while getting rows from DB", zap.Error(err))
		return err
	}
	defer rows.Close()

	err = scanRows(rows, invMap)
	if err != nil {
		log.Error("error while scanning rows", zap.Error(err))
		return err
	}

	if rows.Err() != nil {
		log.Error("error from rows", zap.Error(rows.Err()))
		return rows.Err()
	}

	return nil
}

// scanRows сканирует строки из результата запроса и заполняет информацию о цене и скидке.
func scanRows(rows pgx.Rows, invMap map[string]*domain.Inventory) error {
	for rows.Next() {
		var (
			productID string
			price     sql.NullFloat64
			discount  sql.NullInt64
			count     sql.NullInt64
		)

		err := rows.Scan(&productID, &price, &discount, &count)
		if err != nil {
			return err
		}

		if !price.Valid {
			price.Float64 = 0
		}
		if !discount.Valid {
			discount.Int64 = 0
		}
		if !count.Valid {
			count.Int64 = 0
		}

		if inv, ok := invMap[productID]; ok {
			if inv.ProductCount > int(count.Int64) {
				inv.ProductCount = int(count.Int64)
			}
			inv.ProductPrice = price.Float64
			inv.ProductSale = int(discount.Int64)
		}
	}
	return nil
}

// GetProductsAtWarehouse получает продукты на складе с пагинацией.
func (db *Postgres) GetProductsAtWarehouse(ctx context.Context, params *dto.Pagination, warehouseID string) ([]*domain.Inventory, error) {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.GetProducts"),
	)

	var products []*domain.Inventory

	stmt := `
	SELECT p.product_id, p.product_name, inv.product_price, inv.product_sale
	FROM inventory inv
	JOIN product p USING (product_id)
	WHERE inv.warehouse_id = $1
	OFFSET $2
	LIMIT $3
	`

	rows, err := db.pool.Query(ctx, stmt, warehouseID, params.Offset, params.Limit)
	if err != nil {
		log.Error("error while executing statement", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var (
			id    string
			name  string
			price sql.NullFloat64
			sale  sql.NullInt64
		)

		err = rows.Scan(&id, &name, &price, &sale)
		if err != nil {
			continue
		}

		productID, _ := uuid.Parse(id)

		prod := &domain.Inventory{
			Product: &domain.Product{
				ID:   productID,
				Name: name,
			},
		}

		if price.Valid {
			prod.ProductPrice = float64(price.Float64)
		}
		if sale.Valid {
			prod.ProductSale = int(sale.Int64)
		}

		products = append(products, prod)
	}

	if rows.Err() != nil {
		log.Error("error after scanning rows", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return products, nil
}

// BuyProducts вычитает количество продуктов из инвентаря.
//
// Если продуктов нет на складе, то возвращает ErrNotEnoughProductCount.
func (db *Postgres) BuyProducts(ctx context.Context, inventories []*domain.Inventory) error {
	log := logger.GetLogger().With(
		zap.String("op", "repository.Postgres.BuyProducts"),
	)

	if len(inventories) == 0 {
		return nil
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		log.Error("error while beginning transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	err = validateProductCount(ctx, tx, inventories)
	if err != nil {
		log.Error("error while validating product count", zap.Error(err))
		return err
	}

	err = updateProductCount(ctx, tx, inventories)
	if err != nil {
		log.Error("error while updating product count", zap.Error(err))
		return err
	}

	return tx.Commit(ctx)
}

// validateProductCount проверяет, что количество продуктов на складе достаточно для покупки.
//
// Если количество продуктов меньше, чем нужно, то возвращает ErrNotEnoughProductCount.
func validateProductCount(ctx context.Context, tx pgx.Tx, invs []*domain.Inventory) error {
	warehouseID := invs[0].Warehouse.ID.String()
	products := make([]string, 0, len(invs))
	countMap := make(map[string]int, len(invs))
	invMap := make(map[string]*domain.Inventory, len(invs))

	for _, inv := range invs {
		productID := inv.Product.ID.String()
		products = append(products, productID)
		countMap[productID] = inv.ProductCount
		invMap[productID] = inv
	}

	stmt := `
	SELECT product_id, product_count, product_price, product_sale
	FROM inventory
	WHERE warehouse_id = $1 AND product_id = ANY($2)
	FOR UPDATE
	`

	rows, err := tx.Query(ctx, stmt, warehouseID, products)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = processRows(rows, invMap)
	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	rows.Close()
	tag := rows.CommandTag()
	if int(tag.RowsAffected()) != len(invs) {
		return custErr.ErrNotFoundProductAtWarehouse
	}

	return nil
}

// processRows обрабатывает строки из результата запроса и проверяет количество продуктов.
func processRows(rows pgx.Rows, invMap map[string]*domain.Inventory) error {
	for rows.Next() {
		var (
			dbProductID string
			dbCount     sql.NullInt64
			dbPrice     sql.NullFloat64
			dbSale      sql.NullInt64
		)

		err := rows.Scan(&dbProductID, &dbCount, &dbPrice, &dbSale)
		if err != nil {
			continue
		}

		if !dbCount.Valid {
			return custErr.ErrNotEnoughProductCount
		}

		currentInv, ok := invMap[dbProductID]
		if !ok {
			continue
		}

		if int(dbCount.Int64) < currentInv.ProductCount {
			return custErr.ErrNotEnoughProductCount
		}

		if dbPrice.Valid {
			currentInv.ProductPrice = dbPrice.Float64
		}

		if dbSale.Valid {
			currentInv.ProductSale = int(dbSale.Int64)
		}

	}
	return nil
}

// updateProductCount обновляет количество продуктов в инвентаре.
//
// Если количество продуктов меньше нуля, то возвращает ErrNotEnoughProductCount.
func updateProductCount(ctx context.Context, tx pgx.Tx, invs []*domain.Inventory) error {
	warehouseID := invs[0].Warehouse.ID.String()
	for _, inv := range invs {
		productID := inv.Product.ID.String()
		want := inv.ProductCount

		stmt := `
		UPDATE inventory
		SET	product_count = product_count - $1
		WHERE warehouse_id = $2 AND product_id = $3 AND product_count >= $1
		`

		tag, err := tx.Exec(ctx, stmt, want, warehouseID, productID)
		if err != nil {
			return err
		}

		if tag.RowsAffected() < 1 {
			return custErr.ErrNotEnoughProductCount
		}
	}

	return nil
}
