package config

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `yaml:"env" env:"ENV" required:"true"`
	Storage struct {
		Database string `yaml:"database" env:"DATABASE"  required:"true"`
	} `yaml:"storage" required:"true"`

	HttpServer struct {
		Address string `yaml:"address" env:"ADDRESS"  required:"true"`
	} `yaml:"http_server" env:"HTTPSERVER"  required:"true"`
}

var (
	instance *Config
	once     sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		configPath := flag.String("config", "", "path of config file ")
		flag.Parse()

		if *configPath == "" {
			*configPath = os.Getenv("CONFIG_PATH")
		}

		if *configPath == "" {
			log.Fatal("config path does not provided!")
		}

		if _, err := os.Stat(*configPath); err != nil {
			log.Fatalf("%s file does not exits!", *configPath)
		}

		instance = &Config{}
		cleanenv.ReadConfig(*configPath, instance)
	})
	slog.Info("configration read successfully")
	// println(instance.Storage.Database)
	return instance
}
