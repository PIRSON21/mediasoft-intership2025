package postgresql

import (
	"context"
	"fmt"

	"github.com/PIRSON21/mediasoft-intership2025/pkg/config"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dbConfig config.DBConfig) (*Postgres, error) {
	const op = "repository.postgresql.NewPostgres"
	log := logger.GetLogger()
	log = log.With(zap.String("op", op))

	connOpts, err := parsePostgresOpts(dbConfig)
	if err != nil {
		log.Error("error while parsing config", zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, connOpts)
	if err != nil {
		log.Error("error while creating connection pool", zap.Error(err))
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Error("error while checking postgres connection", zap.Error(err))
		return nil, err
	}

	return &Postgres{
		pool: pool,
	}, nil
}

func parsePostgresOpts(cfg config.DBConfig) (*pgxpool.Config, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	pgxCfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	return pgxCfg, nil
}

func (db *Postgres) Close() {
	db.pool.Close()
}
