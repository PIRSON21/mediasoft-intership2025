package main

import (
	"context"
	"log"
	"net/http"

	"github.com/PIRSON21/mediasoft-go/internal/handler"
	"github.com/PIRSON21/mediasoft-go/internal/middleware"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/internal/service"
	"github.com/PIRSON21/mediasoft-go/pkg/config"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

func initApp() {
	cfg := config.MustParseConfig()
	log.Println("config successfully parsed", cfg)

	logger.MustCreateLogger(cfg.LoggerConfig)

	zlog := logger.GetLogger()
	zlog.Info("logger successfully set up")

	// подключение repositories
	zlog.Info("trying to connect to repositories")
	repo := repository.MustInitRepository(context.TODO(), cfg.DBConfig)

	zlog.Info("repositories set up successfully")

	// инициализация services
	zlog.Debug("setting up the services")

	// инициализация handlers
	zlog.Debug("setting up the handlers")
	warehouseHandlers := &handler.WarehouseHandler{
		Service: service.NewWarehouseService(repo),
	}

	// задание роутингов
	mux := createRouter(warehouseHandlers)

	// запуск сервера TODO: убрать отсюда
	zlog.Info("server ready to start", zap.String("addr", cfg.Address))
	http.ListenAndServe(cfg.Address, mux)
}

func createRouter(warehouseHandlers *handler.WarehouseHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// warehouses
	mux.Handle("/warehouses", chainMiddleware(
		http.HandlerFunc(warehouseHandlers.WarehousesHandler),
		middleware.Recoverer,
		middleware.LoggingMiddleware,
	))

	return mux
}

func chainMiddleware(h http.Handler, mws ...func(http.Handler) http.HandlerFunc) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
