package redis_repos

import "chat-service/internal/infra/memory"

type RedisMessagesRepo struct {
	redisClient *memory.RedisClient
}

func NewRedisMessagesRepo(
	redisClient *memory.RedisClient,
) *RedisMessagesRepo {
	return &RedisMessagesRepo{
		redisClient: redisClient,
	}
}
