package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PIRSON21/mediasoft-intership2025/internal/handler"
	"github.com/PIRSON21/mediasoft-intership2025/internal/middleware"
	"github.com/PIRSON21/mediasoft-intership2025/internal/repository"
	"github.com/PIRSON21/mediasoft-intership2025/internal/service"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/config"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

func startApp() {
	cfg := config.MustParseConfig()
	log.Println("config successfully parsed")

	logger.MustCreateLogger(cfg.LoggerConfig)

	zlog := logger.GetLogger()
	zlog.Debug("logger successfully set up")
	zlog.Info("starting mediasoft-intership2025", zap.String("version", version), zap.String("environment", cfg.Environment))

	// подключение repositories
	zlog.Debug("trying to connect to repositories")
	repo := repository.MustInitRepository(context.Background(), cfg.DBConfig)
	defer repo.Close()

	zlog.Debug("repositories set up successfully")

	// инициализация services
	zlog.Debug("setting up the services")
	warehouseService := service.NewWarehouseService(repo)
	productService := service.NewProductService(repo, cfg.Address)
	analyticsService := service.NewAnalyticsService(repo)
	inventoryService := service.NewInventoryService(repo, analyticsService)

	// инициализация handlers
	zlog.Debug("setting up the handlers")
	warehouseHandlers := handler.NewWarehouseHandler(warehouseService)
	productHandlers := handler.NewProductHandler(productService)
	inventoryHandlers := handler.NewInventoryHandler(inventoryService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// задание роутингов
	zlog.Debug("creating router")
	mux := createRouter(warehouseHandlers, productHandlers, inventoryHandlers, analyticsHandler)

	// создание сервера
	zlog.Debug("creating server")
	srv := http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	stopCh := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		zlog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			zlog.Error("error while shutdown server", zap.Error(err))
		}
		close(stopCh)
	}()

	// запуск сервера
	zlog.Info("server ready to start", zap.String("addr", cfg.Address))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		zlog.Fatal("server error", zap.Error(err))
	}

	<-stopCh
}

func createRouter(warehouseHandlers *handler.WarehouseHandler, productHandlers *handler.ProductHandler, inventoryHandlers *handler.InventoryHandler, analyticsHandlers *handler.AnalyticsHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

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

	mux.Handle("/inventory/check_cart", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.CalculateCart),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/inventory/buy", chainMiddleware(
		http.HandlerFunc(inventoryHandlers.BuyProducts),
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

	// analytics
	mux.Handle("/analytics/", chainMiddleware(
		http.HandlerFunc(analyticsHandlers.GetWarehouseAnalytics),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.LoggingMiddleware,
	))

	mux.Handle("/analytics/top_warehouses", chainMiddleware(
		http.HandlerFunc(analyticsHandlers.GetTopWarehouses),
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
