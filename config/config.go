package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Mode     string
	Server   Server
	Postgres Postgres
}

type Server struct {
	Addr    string
	Port    string
	Timeout time.Duration
}

type Postgres struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDbname   string
}

func loadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
