package memory

import (
	"chat-service/internal/types"
	"context"
	"fmt"
	"log"
)

type RedisMessagesCache struct {
	redisClient *RedisClient
}

func NewRedisMessagesCache(
	redisClient *RedisClient,
) *RedisMessagesCache {
	return &RedisMessagesCache{
		redisClient: redisClient,
	}
}

func (r RedisMessagesCache) GetMessagesByChatUid(chat_uid string) ([]types.Message, error) {
	ctx := context.Background()

	msgList, err := r.redisClient.Rdb.LRange(
		ctx,
		fmt.Sprintf("chat:%s:messages", chat_uid),
		0, -1,
	).Result()

	if err != nil {
		return nil, err
	}

	messages := make([]types.Message, len(msgList))

	for i, item := range msgList {
		msg, err := types.NewMessageFromJson(item)
		if err != nil {
			log.Println("cant decode msg json from redis")
			continue
		}
		messages[i] = msg
	}

	return messages, nil
}

func (r RedisMessagesCache) SaveMessage(msg types.Message) error {
	ctx := context.Background()

	err := r.redisClient.Rdb.LPush(
		ctx,
		fmt.Sprintf("chat:%s:messages", msg.ChatUid),
		msg.ToJSON(),
	).Err()

	return err
}

func (r RedisMessagesCache) DeleteMessage(msg types.Message) error {
	ctx := context.Background()

	// Redis lists don't support direct deletion of an element, so we use LREM
	err := r.redisClient.Rdb.LRem(
		ctx,
		fmt.Sprintf("chat:%s:messages", msg.ChatUid),
		1, // Remove only one matching occurrence
		msg.ToJSON(),
	).Err()

	if err != nil {
		log.Println(err)
	}

	return err
}
