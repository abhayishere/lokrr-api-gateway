package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		HTTPPort string `mapstructure:"http_port"`
	} `mapstructure:"server"`
	Grpc struct {
		AuthServiceAddress string `mapstructure:"auth_service_address"`
		DocServiceAddress  string `mapstructure:"doc_service_address"`
	} `mapstructure:"grpc"`
	CORS struct {
		AllowedOrigins string `mapstructure:"allowed_origins"`
		AllowedMethods string `mapstructure:"allowed_methods"`
		AllowedHeaders string `mapstructure:"allowed_headers"`
	} `mapstructure:"cors"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config") // Docker container path
	viper.AddConfigPath("./config")    // Local path
	viper.AddConfigPath("../config")   // Local development path
	viper.AddConfigPath("/config")     // Alternative container path

	// Read environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	fmt.Println("value for http port is ", os.Getenv("HTTP_PORT"))
	// Override with environment variables if set
	if os.Getenv("HTTP_PORT") != "" {
		cfg.Server.HTTPPort = os.Getenv("HTTP_PORT")
	}
	if os.Getenv("GRPC_AUTH_SERVER") != "" {
		cfg.Grpc.AuthServiceAddress = os.Getenv("GRPC_AUTH_SERVER")
	}
	if os.Getenv("GRPC_DOC_SERVER") != "" {
		cfg.Grpc.DocServiceAddress = os.Getenv("GRPC_DOC_SERVER")
	}

	return &cfg, nil
}
