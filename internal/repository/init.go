package repository

import (
	"context"
	"os"

	"github.com/PIRSON21/mediasoft-intership2025/internal/repository/postgresql"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/config"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

// Repository вмещает в себя все репозитории приложения.
// Так как в рамках этого приложения используется одна база данных,
// то репозитории могут быть объединены в один интерфейс.
//
// Это освобождает меня от необходимости создавать для каждого интерфейса новое соединение.
type Repository interface {
	CloserRepository

	WarehouseRepository
	ProductRepository
	InventoryRepository

	AnalyticsRepository
}

// CloserRepository - интерфейс для репозиториев, которые нужно закрывать.
type CloserRepository interface {
	Close()
}

// MustInitRepository инициализирует репозитории приложения.
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
