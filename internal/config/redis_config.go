package config

import "github.com/redis/go-redis/v9"

func GetRedisClient(address, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
