package broker

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	redisClient *redis.Client
}

func NewRedisBroker() *RedisBroker {

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:            "localhost:6379", // Redis address
		Password:        "",               // No password
		DB:              0,
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

	return &RedisBroker{
		redisClient: rdb,
	}
}

func (r RedisBroker) Publish(ctx context.Context, topic string, message string) error {
	err := r.redisClient.Publish(ctx, topic, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisBroker) Subscribe(ctx context.Context, topics ...string) (chan string, error) {
	sub := r.redisClient.Subscribe(ctx, topics...)

    log.Println("CREATED SUBSCIBE")

	if sub == nil {
		return nil, errors.New("failed to subscribe")
	}

	msgCh := make(chan string, 1000)

	go func() {
		defer close(msgCh)
		defer sub.Close()

		for msg := range sub.Channel(redis.WithChannelSize(1000)) {
			msgCh <- msg.Payload
		}
	}()

	return msgCh, nil
}
