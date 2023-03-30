package services

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisProvider() *redis.Client {
	config := ConfigProvider()
	opts := redis.Options{
		Addr: config.RedisUri,
	}

	opts.DB = config.RedisDb

	if config.RedisPassword != "" {
		opts.Password = config.RedisPassword
	}

	client := redis.NewClient(&opts)

	return client
}
