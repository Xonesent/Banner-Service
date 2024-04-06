package config

import (
	"avito/assignment/pkg/constant"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
)

const configPath = "./config/config.json"

type Config struct {
	Server struct {
		Host string `validate:"required"`
	}
	OpenTelemetry struct {
		URL         string `validate:"required"`
		ServiceName string `validate:"required"`
	}
	Postgres struct {
		Host     string `validate:"required"`
		Port     string `validate:"required"`
		User     string `validate:"required"`
		DbName   string `validate:"required"`
		Password string `validate:"required"`
		SSLMode  string `validate:"required"`
	}
	Redis struct {
		Host         string `validate:"required"`
		Port         string `validate:"required"`
		MinIdleConns int    `validate:"required"`
		PoolSize     int    `validate:"required"`
		PoolTimeout  int    `validate:"required"`
		Password     string `validate:"required"`
	}
	BannerSettings struct {
		BannerTTLSeconds int `validate:"required"`
	}
}

func LoadConfig() (c *Config, err error) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(jsonFile).Decode(&c)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}
	constant.Host = c.Server.Host
	return
}
