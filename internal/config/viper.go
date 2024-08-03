package config

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	Server        ServerConfig
	Database      DatabaseConfig
	ElasticSearch ElasticSearchConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Schema      string
	Host        string
	Environment string
}

type ServerConfig struct {
	Port     string
	Debug    bool
	TimeZone string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	DBName   string
	UserName string
	Password string
	Debug    bool
	Pool     PoolConfig
}

type PoolConfig struct {
	Idle     int
	Max      int
	Lifetime int
}

type ElasticSearchConfig struct {
	Host       string
	UserName   string
	Password   string
	Index      string
	MaxTimeOut int
	Debug      bool
}

func LoadConfigPath(path string) (Config, error) {
	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}
