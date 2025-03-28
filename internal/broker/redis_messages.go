package broker

import (
	"chat-service/internal/types"
	"context"
	"log"

	"github.com/spf13/viper"
)

// implements events.MessagesAdaptor

type RedisMessagesBroker struct {
	redisClient types.BrokerBaseAdaptor
}

func NewRedisMessagesBroker(
	redisClient types.BrokerBaseAdaptor,
) *RedisMessagesBroker {
	return &RedisMessagesBroker{
		redisClient: redisClient,
	}
}

func (r *RedisMessagesBroker) GetMessagesFromQueue(ctx context.Context, topic string) (chan types.Message, error) {
	messages := make(chan types.Message, viper.GetInt("app.global_buff"))

	go func() {
		for {
			msg, err := r.redisClient.FromQueue(ctx, topic)
			if err != nil {
				log.Println("error during fetching data from queue: ", err)
			}
			message, err := types.NewMessageFromJson(msg)
			if err != nil {
				log.Println("error while unmarshal message from queue")
			}
			messages <- message
		}
	}()

	return messages, nil
}

func (r *RedisMessagesBroker) SendMessageToQueue(ctx context.Context, topic string, msg types.Message) error {
	return r.redisClient.ToQueue(ctx, topic, msg.ToJSON())
}

func (r *RedisMessagesBroker) GetMessagesFromChannel(
	ctx context.Context,
	chatsUids ...string,
) (chan types.Message, error) {

	ch, err := r.redisClient.Subscribe(ctx, chatsUids...)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan types.Message, viper.GetInt("app.global_buff"))

	go func() {
		defer close(msgChan)

		for {
			select {
			case redisMsg := <-ch:
				message, _ := types.NewMessageFromJson(redisMsg)
				msgChan <- message
			case <-ctx.Done():
				return
			}
		}
	}()

	return msgChan, nil
}

func (r *RedisMessagesBroker) SendMessageToChannel(
	ctx context.Context,
	topic string,
	msg types.Message,
) error {
	return r.redisClient.Publish(ctx, topic, msg.ToJSON())
}
