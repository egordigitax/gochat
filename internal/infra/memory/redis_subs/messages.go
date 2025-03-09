package redis_subs

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisMessageBroker struct {
	redisClient *redis.Client
}

func NewRedisMessageBroker() *RedisMessageBroker {

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

	log.Println("Connected to Redis Broker")

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	return &RedisMessageBroker{
		redisClient: rdb,
	}
}

func (r RedisMessageBroker) Publish(topic string, message string) error {
	ctx := context.Background()
	err := r.redisClient.Publish(ctx, topic, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisMessageBroker) Subscribe(topics ...string) (chan string, error) {
	ctx := context.Background()
	sub := r.redisClient.Subscribe(ctx, topics...)

	if _, err := sub.Receive(ctx); err != nil {
		return nil, err
	}

	msgCh := make(chan string)

	go func() {
		defer close(msgCh)
		for msg := range sub.Channel() {
			msgCh <- msg.Payload
		}
	}()

	return msgCh, nil
}
