package logger

import (
	"log"

	"github.com/PIRSON21/mediasoft-intership2025/pkg/config"
	"go.uber.org/zap"
)

var instance *zap.Logger

// MustCreateLogger инициализирует логгер на основе конфигурации.
func MustCreateLogger(cfg config.LoggerConfig) {
	logCfg := createLoggerConfig(&cfg)
	// TODO: доделать норм инициализацию логера. сделать создание конфига исходя из настроек + создание логгера

	newLogger, err := logCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	instance = newLogger

	instance.Info("logger initialized")
}

// GetLogger возвращает инициализированный логгер для использования в приложении.
// Если логгер не инициализирован, вызывает панику.
func GetLogger() *zap.Logger {
	if instance == nil {
		log.Fatal("logger is not initialized")
	}

	return instance
}

// createLoggerConfig создает конфигурацию логгера на основе переданных настроек.
func createLoggerConfig(config *config.LoggerConfig) zap.Config {
	var (
		level = zap.NewAtomicLevel()
		err   error
		// TODO: по заданию нужно использовать JSON, нужно поменять на "json"
		encoding = "console"

		cfg zap.Config
	)

	if config.Debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	level, err = zap.ParseAtomicLevel(config.Level)
	if err != nil {
		log.Printf("error while parsing logger level: %v", err)
	}
	cfg.Level = level

	cfg.Encoding = encoding
	cfg.DisableStacktrace = true

	return cfg
}

// Sync синхронизирует логгер, чтобы все сообщения были записаны.
// Это нужно для корректного завершения работы приложения.
func Sync() error {
	return instance.Sync()
}
