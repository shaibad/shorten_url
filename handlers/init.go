package handlers

import (
	"fmt"

	"url-shortener/config"
	"url-shortener/db"

	"github.com/go-redis/redis"
)

var redisClient db.Redis
var dbClient db.DB

func init () {
	// Init DB connections
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
	redisClient = db.NewRedis(client)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", 
	PostgresConf.Host, PostgresConf.Port, PostgresConf.User, PostgresConf.Password)	
	dbClient, err = db.NewDB("postgres", psqlInfo)
	if err != nil {
		panic("Error, could not connect to DB")
	}
}
