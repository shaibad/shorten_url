package config

import (
	"log"
	"github.com/kelseyhightower/envconfig"
)

// Config is a general interfce for configuration types
type Config interface{}

type HandlerTypeConf struct {
	Path string `envconfig:"HANDLER_PATH"`
	Method string `envconfig:"HANDLER_METHOD"`
	Port string `envconfig:"HANDLER_PORT"`
}

type BaseUrlConf struct {
	Url string `envconfig:"BASE_URL"`
	Protocol string `envconfig:"BASE_URL_PROTOCOL"`
}

type RedisConf struct {
	Address string `envconfig:"REDIS_ADDRESS"`
	Password string `envconfig:"REDIS_PASSWORD"`
}

type PostgresConf struct {
	Host string `envconfig:"POSTGRES_HOST"`
	Port int `envconfig:"POSTGRES_PORT"`
	User string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
}

// GetEnv is a function to retrieve environemnt variables
func GetEnv(cfg Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Println("Couldn't process env config")
	}
}
