package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - общая конфигурация приложения.
type Config struct {
	Environment string `env:"ENV" env-default:"prod"`
	Address     string `env:"ADDRESS" env-default:":8080"`
	DBConfig
	LoggerConfig
}

// DBConfig - конфигурация базы данных.
type DBConfig struct {
	DBName     string `env:"DBNAME" env-required:"true"`
	DBUser     string `env:"DBUSER" env-required:"true"`
	DBPassword string `env:"DBPASSWORD" env-required:"true"`
	DBHost     string `env:"DBHOST" env-default:"localhost"`
	DBPort     uint16 `env:"DBPORT" env-default:"5432"`
}

// LoggerConfig - конфигурация логгера.
type LoggerConfig struct {
	Debug bool
	Level string `env:"LEVEL" env-default:"INFO"`
}

// MustParseConfig читает данные конфига из переменных окружения.
//
// При ошибке возвращает панику.
func MustParseConfig() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Environment == "dev" {
		cfg.Debug = true
	}

	return &cfg
}
