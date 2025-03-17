package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisClient() *RedisClient {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			viper.GetString("memory.host"),
			viper.GetInt("memory.port"),
		),
		Password: viper.GetString("memory.password"),
		DB:       viper.GetInt("memory.db"),
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
