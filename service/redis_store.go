package service

import (
	"context"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(dsn string) *RedisStore {
	redisClient := db.NewRedisClient(dsn)
	return &RedisStore{
		client: redisClient,
	}
}

func (r *RedisStore) GetObject(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisStore) SetObject(key, val string) error {
	err := r.client.Set(context.Background(), key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
