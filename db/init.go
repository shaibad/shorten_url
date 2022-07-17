package db

import (
	"fmt"

	"url-shortener/config"

	"github.com/go-redis/redis"
)

var RedisClient Redis
var DbClient DB

// Init DB connections
func InitConnections() error{
	var RedisConf config.RedisConf
	var PostgresConf config.PostgresConf
	var err error

	config.GetEnv(&RedisConf)
	config.GetEnv(&PostgresConf)

	client := redis.NewClient(&redis.Options{
		Addr: RedisConf.Address,
		Password: RedisConf.Password,
		DB: 0,
	})
	RedisClient = NewRedis(client)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", 
	PostgresConf.Host, PostgresConf.Port, PostgresConf.User, PostgresConf.Password)	
	DbClient, err = NewDB("postgres", psqlInfo)
	if err != nil {
		return err
	}
	return nil
}
