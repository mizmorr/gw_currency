package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, host, port, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{
		Client: rdb,
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value float64, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (float64, error) {
	rateRaw, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("key %s not found", key)
	}
	if err != nil {
		return 0, err
	}
	rate, err := strconv.ParseFloat(rateRaw, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
