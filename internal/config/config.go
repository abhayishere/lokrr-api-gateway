package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		HTTPPort string `mapstructure:"http_port"`
	} `mapstructure:"server"`
	Grpc struct {
		AuthServiceAddress string `mapstructure:"auth_service_address"`
	} `mapstructure:"grpc"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")   // Config file name without extension
	viper.SetConfigType("yaml")     // Config file type
	viper.AddConfigPath("./config") // Path to look for the config file
	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if os.Getenv("GRPC_AUTH_SERVER") != "" {
		cfg.Grpc.AuthServiceAddress = os.Getenv("GRPC_AUTH_SERVER")
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
