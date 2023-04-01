package configs

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      `yaml:"app"`
	Client   `yaml:"client"`
	Logger   `yaml:"logger"`
	Database `yaml:"database"`
}

type App struct {
	Name                 string `env-required:"true" yaml:"name"`
	Version              string `env-required:"true" yaml:"version"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

type Client struct {
	Address string `env-required:"true" yaml:"address" env:"RUN_ADDRESS"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Database struct {
	URI string `env:"DATABASE_URI"`
}

func NewConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Client.Address, "a", cfg.Client.Address, "address to listen on")
	flag.StringVar(&cfg.Database.URI, "d", cfg.Database.URI, "database URI")
	flag.StringVar(&cfg.App.AccrualSystemAddress, "r", cfg.App.AccrualSystemAddress, "accrual address")

	_ = cleanenv.ReadConfig("./config/config.yaml", cfg)

	_ = cleanenv.ReadEnv(cfg)

	return cfg
}
