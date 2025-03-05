package postgres_repos

import (
	"chat-service/internal/domain"
	"chat-service/internal/infra/databases/postgres"
	"log"
)

type PGChatsStorage struct {
	postgresClient *postgres.PostgresClient
}

func NewPGChatsStorage(
	postgresClient *postgres.PostgresClient,
) *PGChatsStorage {
	return &PGChatsStorage{
		postgresClient: postgresClient,
	}
}

func (m *PGChatsStorage) GetUsersChats(
	user_uid string,
	limit,
	offset int,
) ([]domain.Chat, error) {
	var chats []domain.Chat
	log.Println(user_uid)
	query := `
        SELECT 
            uc.user_uid,
            uc.chat_uid,
            uc.message,
            uc.created_at,
            u.nickname as title,
            coalesce('', '') as image_url
        FROM 
            users_messages uc
        JOIN users u 
            ON u.uid = uc.user_uid
        WHERE 
            uc.user_uid = $1
    `

	err := m.postgresClient.C_RO.Select(&chats, query, user_uid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(chats)

	return chats, nil
}
