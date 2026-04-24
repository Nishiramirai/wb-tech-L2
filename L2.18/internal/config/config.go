package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env:"ADDRESS" env-default:"0.0.0.0:8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {

	// Пытаемся взять путь к конфигу из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")

	// Если переменная не задана, пробуем открыть файл для локальной среды
	if configPath == "" {
		const defaultPath = "config/local.yaml"

		if _, err := os.Stat(defaultPath); err != nil {
			log.Fatalf("CONFIG_PATH is not set. Open config by default: %s", err)
		}

		configPath = defaultPath
	}

	// Проверяем существование указанного в переменной окружения конфиг файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
