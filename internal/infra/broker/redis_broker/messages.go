package redis_broker

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"context"
	"log"
)

type RedisMessagesBroker struct {
	redisClient events.BrokerBaseAdaptor
}

func NewRedisMessagesBroker(
	redisClient events.BrokerBaseAdaptor,
) *RedisMessagesBroker {
	return &RedisMessagesBroker{
		redisClient: redisClient,
	}
}

func (r RedisMessagesBroker) GetMessagesFromQueue(ctx context.Context, topic string) (chan entities.Message, error) {
	messages := make(chan entities.Message, 1000)

	go func() {
		for {
			msg, err := r.redisClient.FromQueue(ctx, topic)
			if err != nil {
				log.Println("error during fetching data from queue: ", err)
			}
			message, err := entities.NewMessageFromJson(msg)
			if err != nil {
				log.Println("error while unmarshal message from queue")
			}
			messages <- message
		}
	}()

	return messages, nil
}

func (r RedisMessagesBroker) SendMessageToQueue(ctx context.Context, topic string, msg entities.Message) error {
	return r.redisClient.ToQueue(ctx, topic, msg.ToJSON())
}

func (r *RedisMessagesBroker) GetMessagesFromChannel(
	ctx context.Context,
	chats_uids ...string,
) (chan entities.Message, error) {

	ch, err := r.redisClient.Subscribe(ctx, chats_uids...)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan entities.Message, 1000)

	go func() {
		defer close(msgChan)

		for {
			select {
			case redisMsg := <-ch:
				message, _ := entities.NewMessageFromJson(redisMsg)
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
	msg entities.Message,
) error {
	return r.redisClient.Publish(ctx, topic, msg.ToJSON())
}
