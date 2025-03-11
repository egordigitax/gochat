package redis_broker

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"context"
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

func (r *RedisMessagesBroker) GetMessagesFromChats(
	ctx context.Context,
	chats_uids ...string,
) (chan entities.Message, error) {

	ch, err := r.redisClient.Subscribe(ctx, chats_uids...)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan entities.Message, 100)

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

func (r *RedisMessagesBroker) SendMessageToChat(
	ctx context.Context,
	topic string,
	msg entities.Message,
) error {
	return r.redisClient.Publish(ctx, topic, msg.ToJSON())
}
