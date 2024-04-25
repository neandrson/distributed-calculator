package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string         `yaml:"env"`
	StoragePath string         `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConfig     `yaml:"grpc"`
	GRPCWeb     string         `yaml:"grpcweb"`
	Token_time  time.Duration  `yaml:"token_time"`
	Postgres    PostgresConfig `yaml:"postgres"`
	CountWorker int            `yaml:"countWorker"`
}
type GRPCConfig struct {
	Port    string        `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}
type PostgresConfig struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	DBname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type ConfigTime struct {
	Plus     time.Duration `yaml:"plus"`
	Minus    time.Duration `yaml:"minus"`
	Division time.Duration `yaml:"division"`
	MultiP   time.Duration `yaml:"multp"`
	Exponent time.Duration `yaml:"exponent"`
}

func LoadConfig(path string) *Config {

	if path == "" {
		panic("config file empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config not exist")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("config not read")
	}
	return &cfg
}
func LoadConfigTime(path string) *ConfigTime {

	if path == "" {
		panic("config file empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config not exist")
	}
	var cfg ConfigTime
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("config not read")
	}
	return &cfg
}
