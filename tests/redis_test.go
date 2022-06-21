package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"url-shortener/db"
)

var (
	client *redis.Client
)

var (
	key = Url
	val = Shorter
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("Error '%s' occured when opening a stub database connection", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	code := m.Run()
	os.Exit(code)
}

func TestSet(t *testing.T) {
	exp := time.Duration(0)
	mock := redismock.NewNiceMock(client)
	mock.On("Set", key, val, exp).Return(redis.NewStatusResult("", nil))
	r := db.NewRedis(mock)
	err := r.Set(key, val, exp)

	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	mock.On("Get", key).Return(redis.NewStringResult(val, nil))
	r := db.NewRedis(mock)
	res, err := r.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, val, res)
}
