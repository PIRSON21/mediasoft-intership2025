package repository

import (
	"context"
	"os"

	"github.com/PIRSON21/mediasoft-go/internal/repository/postgresql"
	"github.com/PIRSON21/mediasoft-go/pkg/config"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

type Repository interface {
	WarehouseRepository
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
