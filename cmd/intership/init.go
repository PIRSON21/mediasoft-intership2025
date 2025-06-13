package main

import (
	"log"
	"net/http"

	"github.com/PIRSON21/mediasoft-go/internal/handler"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/internal/service"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
)

func initApp() {
	if err := logger.NewLogger(); err != nil {
		log.Fatal(err)
	}

	zlog := logger.GetLogger()
	zlog.Info("logger successfully set up")

	// подключение repositories
	zlog.Info("trying to connect to repositories")
	repo, err := repository.NewRepository()
	if err != nil {
		zlog.Error("error while connetcion to repository", "err", err)
		log.Fatal(err)
	}
	zlog.Info("repositories set up successfully")

	// инициализация services
	zlog.Debug("setting up the services")

	// инициализация handlers
	zlog.Debug("setting up the handlers")
	warehouseHandlers := &handler.WarehouseHandler{
		Service: service.NewWarehouseService(repo),
	}

	// задание роутингов
	mux := http.NewServeMux()

	mux.HandleFunc("/test", warehouseHandlers.GetWarehouses)

	// запуск сервера
	zlog.Info("server ready to start", "addr", ":8080")
	http.ListenAndServe(":8080", mux)
}
