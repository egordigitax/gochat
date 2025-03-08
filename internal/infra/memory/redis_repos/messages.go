package redis_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/memory"
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

func (r RedisMessagesCache) GetMessages(chat_uid string, limit int, offset int) ([]entities.Message, error) {
	panic("unimplemented")
}

func (r RedisMessagesCache) SaveMessage(msg entities.Message) error {
	panic("unimplemented")
}
