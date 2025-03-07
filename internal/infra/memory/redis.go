package memory

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisClient() *RedisClient {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis")

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		Rdb: rdb,
	}
}
