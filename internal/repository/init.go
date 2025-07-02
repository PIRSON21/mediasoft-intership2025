package repository

import (
	"context"
	"os"

	"github.com/PIRSON21/mediasoft-intership2025/internal/repository/postgresql"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/config"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

type Repository interface {
	WarehouseRepository
	ProductRepository
	InventoryRepository

	AnalyticsRepository
	Close()
}

func MustInitRepository(ctx context.Context, dbCfg config.DBConfig) Repository {
	const op = "repository.NewRepository"
	log := logger.GetLogger().With(zap.String("op", op))

	repo, err := postgresql.NewPostgres(ctx, dbCfg)
	if err != nil {
		log.Error("error while creating postgres repo", zap.String("err", err.Error()))
		os.Exit(1)
	}

	return repo
}
