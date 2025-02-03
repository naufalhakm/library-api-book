package database

import "github.com/redis/go-redis/v9"

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "34.142.158.122:6379",
		DB:   0,
	})

	return client
}
