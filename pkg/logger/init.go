package logger

import (
	"log"

	"go.uber.org/zap"
)

var instance *zap.SugaredLogger

func NewLogger() error {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer newLogger.Sync()

	sugarred := newLogger.Sugar()
	instance = sugarred

	instance.Info("logger initialized")
	return nil
}

func GetLogger() *zap.SugaredLogger {
	if instance == nil {
		log.Fatal("logger is not initialized")
	}

	return instance
}
