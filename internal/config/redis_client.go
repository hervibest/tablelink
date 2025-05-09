package config

import "github.com/redis/go-redis/v9"

func NewRedis(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
}
