package memory

import (
	entities "chat-service/internal/types"
	"context"
	"encoding/json"
	"fmt"
)

type RedisChatsCache struct {
	redisClient *RedisClient
}

func NewRedisChatsCache(
	redisClient *RedisClient,
) *RedisChatsCache {
	return &RedisChatsCache{
		redisClient: redisClient,
	}
}

func (r RedisChatsCache) GetUsersChats(user_uid string, limit int, offset int) ([]entities.Chat, error) {
	ctx := context.Background()

	chats, err := r.redisClient.Rdb.HGet(ctx, fmt.Sprintf("users_chats:%s", user_uid), "chats").Result()
	if err != nil {
		return nil, err
	}

	var userChats []entities.Chat
	err = json.Unmarshal([]byte(chats), &userChats)
	if err != nil {
		return nil, err
	}

	// if len(userChats) == 0 {
	// 	return nil, errors.New("nocache")
	// }

	return userChats, nil
}

func (r RedisChatsCache) SetUsersChats(user_uid string, chats []entities.Chat) error {
	ctx := context.Background()

	data, err := json.Marshal(chats)
	if err != nil {
		return err
	}

	err = r.redisClient.Rdb.HSet(ctx, fmt.Sprintf("users_chats:%s", user_uid), "chats", data).Err()
	if err != nil {
		return err
	}

	return nil
}
