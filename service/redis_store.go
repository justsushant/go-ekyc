package service

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(dsn string) RedisStore {
	redisClient := db.NewRedisClient(dsn)
	return RedisStore{
		client: redisClient,
	}
}
