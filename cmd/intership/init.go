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
	productHandlers := handler.NewProductHandler(*service.NewProductService(repo, cfg.Address))
	inventoryHandlers := handler.NewInventoryHandler(service.NewInventoryService(repo))

	// задание роутингов
	mux := createRouter(warehouseHandlers, productHandlers, inventoryHandlers)

	// запуск сервера TODO: убрать отсюда
	zlog.Info("server ready to start", zap.String("addr", cfg.Address))
	http.ListenAndServe(cfg.Address, mux)
}

func createRouter(warehouseHandlers *handler.WarehouseHandler, productHandlers *handler.ProductHandler, inventoryHandlers *handler.InventoryHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// warehouses
	mux.Handle("/warehouses", chainMiddleware(
		http.HandlerFunc(warehouseHandlers.WarehousesHandler),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	// products
	mux.Handle("/products", chainMiddleware(
		http.HandlerFunc(productHandlers.ProductsHandler),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/product/", chainMiddleware(
		http.HandlerFunc(productHandlers.UpdateProduct),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	// inventory
	mux.Handle("/inventory/change_count", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.ChangeProductCount),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/inventory/add_discount", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.AddDiscountToProduct),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/warehouse/", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.GetProductFromWarehouse),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/inventory", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.CreateInventory),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	// static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}

func chainMiddleware(h http.Handler, mws ...func(http.Handler) http.HandlerFunc) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
