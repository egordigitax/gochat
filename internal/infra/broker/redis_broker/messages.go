package redis_broker

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
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
	chats_uids ...string,
) (chan entities.Message, error) {

	ch, err := r.redisClient.Subscribe(chats_uids...)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan entities.Message, 100)

	go func() {
		defer close(msgChan)

		for {
			redisMsg := <-ch
			message, _ := entities.NewMessageFromJson(redisMsg)
			msgChan <- message
		}
	}()

	return msgChan, nil
}

func (r *RedisMessagesBroker) SendMessageToChat(
	msg entities.Message,
) error {
	return r.redisClient.Publish(msg.ChatUid, msg.ToJSON())
}
