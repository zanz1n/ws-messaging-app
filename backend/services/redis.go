package services

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisProvider() *redis.Client {
	var err error
	opts := redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
	}

	opts.DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		opts.DB = 0
	}

	if (os.Getenv("REDIS_PASSWORD") != "") {
		opts.Password = os.Getenv("REDIS_PASSWORD")
	}

	client := redis.NewClient(&opts)

	return client
}
