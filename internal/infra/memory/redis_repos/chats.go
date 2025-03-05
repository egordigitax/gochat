package redis_repos

import (
	"chat-service/internal/domain"
	"chat-service/internal/infra/memory"
)

type RedisChatsCache struct {
	redisClient *memory.RedisClient
}

func NewRedisChatsCache(
	redisClient *memory.RedisClient,
) *RedisChatsCache {
	return &RedisChatsCache{
		redisClient: redisClient,
	}
}

// GetUsersChats implements interfaces.ChatsCache.
func (r RedisChatsCache) GetUsersChats(user_uid string, limit int, offset int) ([]domain.Chat, error) {
	panic("unimplemented")
}
