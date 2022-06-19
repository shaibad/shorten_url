package helpers

import (
	"fmt"
	"log"
	"github.com/go-redis/redis"
	"shorten_url/config"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

var redisClient *redis.Client
var db *sql.DB

func GetFromRedis(key string) (bool, string) {
	val, err := redisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return false, "Failed to get value from redis"
	}
	return true, val
}

func InsertToRedis(key string, value string) (bool) {
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

func InsertToPostgres(tableName string, valueKeys [2]string, values [2]string) bool {
	statement := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2)`, tableName, valueKeys[0], valueKeys[1])
	_, err := db.Exec(statement, values[0], values[1])
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetFromPostgres(valueToSelect string, tableName string, pkey string, value string) (bool, string) {
	var result string
	s := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, valueToSelect, tableName, pkey)
	err := db.QueryRow(s, value).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false, ""
	}
	return true, result
}

func init() {
	var RedisConf config.RedisConf
	config.GetEnv(&RedisConf)
    redisClient = redis.NewClient(&redis.Options{
		Addr: RedisConf.Address,
		Password: RedisConf.Password,
		DB: 0,
	})

	var PostgresConf config.PostgresConf
	var err error
	config.GetEnv(&PostgresConf)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", 
	PostgresConf.Host, PostgresConf.Port, PostgresConf.User, PostgresConf.Password)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}