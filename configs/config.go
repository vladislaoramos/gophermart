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
	Name                 string `env-required:"true" env-default:"gophemart"`
	Version              string `default:"0.0.1"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" env-required:"true"`
}

type Client struct {
	Address string `env-required:"true" env:"RUN_ADDRESS"`
}

type Logger struct {
	Level string `yaml:"level" yaml-default:"debug"`
}

type Database struct {
	URI string `env:"DATABASE_URI" env-required:"true"`
}

func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.Client.Address, "a", cfg.Client.Address, "address to listen on")
	flag.StringVar(&cfg.Database.URI, "d", cfg.Database.URI, "database URI")
	flag.StringVar(&cfg.App.AccrualSystemAddress, "r", cfg.App.AccrualSystemAddress, "accrual address")

	flag.Parse()
}

func (cfg *Config) updateConfig(v *Config) {
	if v.AccrualSystemAddress != "" && v.AccrualSystemAddress != cfg.AccrualSystemAddress {
		cfg.AccrualSystemAddress = v.AccrualSystemAddress
	}

	if v.Database.URI != "" && v.Database.URI != cfg.Database.URI {
		cfg.Database.URI = v.Database.URI
	}

	if v.Address != "" && v.Address != cfg.Address {
		cfg.Address = v.Address
	}
}

const (
	appName            = "gophemart"
	appVersion         = "0.0.1"
	loggerDefaultLevel = "debug"
)

func defaultConfig() *Config {
	return &Config{
		App: App{
			Name:    appName,
			Version: appVersion,
		},
		Logger: Logger{Level: loggerDefaultLevel},
	}
}

func NewConfig() *Config {
	var (
		cfg   *Config
		envs  *Config
		flags *Config
	)

	cfg = defaultConfig()

	flags = new(Config)
	flags.parseFlags()
	cfg.updateConfig(flags)

	envs = new(Config)
	_ = cleanenv.ReadEnv(envs)
	cfg.updateConfig(envs)

	return cfg
}
