package services

import (
	"chat-service/internal/domain"
	"chat-service/internal/infra/databases/postgres/postgres_repos"
	"chat-service/internal/infra/memory/redis_repos"
)

type MessageService struct {
	MsgPGRepo    *postgres_repos.PGMessagesRepository
	MsgRedisRepo *redis_repos.RedisMessagesRepo
}

func NewMessageService(
	pgMessagesRepo *postgres_repos.PGMessagesRepository,
	redisMessagesRepo *redis_repos.RedisMessagesRepo,
) *MessageService {
	return &MessageService{
		MsgPGRepo:    pgMessagesRepo,
		MsgRedisRepo: redisMessagesRepo,
	}
}

func (m *MessageService) GetChatsList(user_uid string) ([]domain.ChatInfo, error) {
	msgs, err := m.MsgPGRepo.GetUsersChats(user_uid, 10, 0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
