package services

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisProvider() *redis.Client {
	config := ConfigProvider()
	var err error
	opts := redis.Options{
		Addr: config.RedisUri,
	}

	opts.DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		opts.DB = 0
	}

	if config.RedisPassword != "" {
		opts.Password = config.RedisPassword
	}

	client := redis.NewClient(&opts)

	return client
}
