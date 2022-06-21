package db

import (
	"time"
	"github.com/go-redis/redis"
)

type Redis interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
}

type _redis struct {
	Client redis.Cmdable
}

func NewRedis(Client redis.Cmdable) Redis {
	return &_redis{Client}
}

func (r *_redis) Get(key string) (string, error) {
	res, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func (r *_redis) Set(key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(key, value, exp).Err()
}
