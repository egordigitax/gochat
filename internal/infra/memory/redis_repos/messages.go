package redis_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/memory"
	"context"
	"fmt"
	"log"
)

type RedisMessagesCache struct {
	redisClient *memory.RedisClient
}

func NewRedisMessagesCache(
	redisClient *memory.RedisClient,
) *RedisMessagesCache {
	return &RedisMessagesCache{
		redisClient: redisClient,
	}
}

func (r RedisMessagesCache) GetMessages(chat_uid string) ([]entities.Message, error) {
	ctx := context.Background()

	msgMap, err := r.redisClient.Rdb.HGetAll(
		ctx,
		fmt.Sprintf("chat:%s:messages", chat_uid),
	).Result()

	if err != nil {
		return nil, err
	}

	messages := make([]entities.Message, len(msgMap))

	i := 0
	for _, item := range msgMap {
		msg, err := entities.NewMessageFromJson(item)
		if err != nil {
			log.Println("cant decode msg json from redis")

		}
		messages[i] = msg
		i++
	}

	return messages, nil
}

func (r RedisMessagesCache) SaveMessage(msg entities.Message) error {
	ctx := context.Background()

	err := r.redisClient.Rdb.HSet(
		ctx,
		fmt.Sprintf("chat:%s:messages", msg.ChatUid),
		msg.Uid,
		msg.ToJSON(),
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r RedisMessagesCache) DeleteMessage(msg entities.Message) error {
	ctx := context.Background()

	err := r.redisClient.Rdb.HDel(
		ctx,
		fmt.Sprintf("chat:%s:messages", msg.ChatUid),
		msg.Uid,
	).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
