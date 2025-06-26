package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (db *Postgres) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	log := logger.GetLogger().With(zap.String("op", "repository.Postgres.GetProduct"))

	stmt := `
	SELECT product_id, product_name, product_description, product_weight, product_params, product_barcode
	FROM product
	`

	rows, err := db.pool.Query(ctx, stmt)
	if err != nil {
		log.Error("error while getting products from DB", zap.String("err", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Weight, &product.Params, &product.Barcode)
		if err != nil {
			log.Error("error while parsing product", zap.String("err", err.Error()))
			continue
		}
		products = append(products, &product)
	}
	if rows.Err() != nil {
		log.Error("error while getting rows", zap.String("err", rows.Err().Error()))
		return nil, rows.Err()
	}

	return products, nil
}

func (db *Postgres) AddProduct(ctx context.Context, p *domain.Product) error {
	stmt := `
	INSERT INTO product(product_name, product_description, product_weight, product_params, product_barcode)
	VALUES ($1, $2, $3, $4, $5)
	`

	tag, err := db.pool.Exec(ctx, stmt, p.Name, p.Description, p.Weight, p.Params, p.Barcode)
	if err != nil {
		pgError := new(pgconn.PgError)
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return custErr.ErrProductAlreadyExists
			}
		}
		return err
	}

	if tag.RowsAffected() < 1 {
		return fmt.Errorf("table wasn't affected")
	}

	return nil
}

func (db *Postgres) UpdateProduct(ctx context.Context, product *domain.Product) error {
	var (
		query         []string
		currentCursor int = 1
		args          []any
	)

	if product.Name != "" {
		query = append(query, fmt.Sprintf("product_name = $%d", currentCursor))
		args = append(args, product.Name)
		currentCursor++
	}

	if product.Description != "" {
		query = append(query, fmt.Sprintf("product_description = $%d", currentCursor))
		args = append(args, product.Description)
		currentCursor++
	}

	if product.Weight != 0 {
		query = append(query, fmt.Sprintf("product_weight = $%d", currentCursor))
		args = append(args, product.Weight)
		currentCursor++
	}

	if product.Params != nil {
		query = append(query, fmt.Sprintf("product_params = $%d", currentCursor))
		args = append(args, product.Params)
		currentCursor++
	}

	if product.Barcode != "" {
		query = append(query, fmt.Sprintf("product_barcode = $%d", currentCursor))
		args = append(args, product.Barcode)
		currentCursor++
	}

	stmt := "UPDATE product SET " + strings.Join(query, ", ") + fmt.Sprintf(" WHERE product_id = $%d", currentCursor)
	args = append(args, product.ID)

	tag, err := db.pool.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return custErr.ErrProductNotFound
	}

	return nil
}
