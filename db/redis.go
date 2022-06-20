package db

import (
	"log"
	"github.com/go-redis/redis"
	"url-shortener/config"
	"time"
)

var redisClient *redis.Client

func GetFromRedis(key string) (bool, string) {
	val, err := redisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return false, "Failed to get value from redis"
	}
	return true, val
}

func InsertToRedis(key, value string) (bool) {
	err := redisClient.Set(key, value, 0).Err()
    if err != nil {
		log.Println(err)
		return false
	}
	_, err = redisClient.Expire(key, 1 * time.Hour).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func init() {
	var RedisConf config.RedisConf
	config.GetEnv(&RedisConf)
    redisClient = redis.NewClient(&redis.Options{
		Addr: RedisConf.Address,
		Password: RedisConf.Password,
		DB: 0,
	})
}
