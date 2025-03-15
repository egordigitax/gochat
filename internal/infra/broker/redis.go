package broker

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisBroker struct {
	redisClient *redis.Client
}

func NewRedisBroker() *RedisBroker {

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			viper.GetString("broker.host"),
			viper.GetInt("broker.port"),
		),
		Password: viper.GetString("broker.password"),
		DB:       viper.GetInt("broker.db"),
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

	msgCh := make(
		chan string,
		viper.GetInt("app.global_buff"),
	)

	go func() {
		defer close(msgCh)
		defer sub.Close()

		for msg := range sub.Channel(
			redis.WithChannelSize(
				viper.GetInt("app.global_buff"),
			),
		) {
			msgCh <- msg.Payload
		}
	}()

	return msgCh, nil
}

func (r RedisBroker) FromQueue(ctx context.Context, topic string) (string, error) {
	task, err := r.redisClient.BLPop(ctx, 0, topic).Result()
	if err != nil {
		return "", err
	}

	if len(task) < 2 {
		return "", errors.New("no values found")
	}

	return task[1], err
}

// ToQueue implements events.BrokerBaseAdaptor.
func (r RedisBroker) ToQueue(ctx context.Context, topic string, message string) error {
	err := r.redisClient.RPush(ctx, topic, message).Err()
	if err != nil {
		return err
	}
	return nil
}
