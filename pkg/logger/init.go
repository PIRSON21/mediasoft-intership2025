package logger

import (
	"log"

	"github.com/PIRSON21/mediasoft-go/pkg/config"
	"go.uber.org/zap"
)

var instance *zap.Logger

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

func GetLogger() *zap.Logger {
	if instance == nil {
		log.Fatal("logger is not initialized")
	}

	return instance
}

func createLoggerConfig(config *config.LoggerConfig) zap.Config {
	var (
		level    = zap.NewAtomicLevel()
		err      error
		encoding = "console"
		cfg      zap.Config
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

func Sync() error {
	return instance.Sync()
}
