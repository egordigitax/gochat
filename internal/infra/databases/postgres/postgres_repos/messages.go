package postgres_repos

import (
	"chat-service/internal/domain"
	"chat-service/internal/infra/databases/postgres"
	"log"
)

type PGMessagesRepository struct {
	postgresClient *postgres.PostgresClient
}

func NewPGMessagesRepository(
	postgresClient *postgres.PostgresClient,
) *PGMessagesRepository {
	return &PGMessagesRepository{
		postgresClient: postgresClient,
	}
}

func (m *PGMessagesRepository) SaveMessage(
	msg domain.Message,
) error {
	query := `
    INSERT INTO users_messages (message, user_uid, chat_uid)
    VALUES ($1, $2, $3);
    `

	_, err := m.postgresClient.C_RW.Exec(
		query,
		msg.Text,
		msg.UserID,
		msg.ChatID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *PGMessagesRepository) GetMessages(chatID string, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	query := `
    SELECT message, user_uid, chat_uid FROM users_messages WHERE chat_uid = $1
    `

	err := m.postgresClient.C_RO.Select(&messages, query, chatID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *PGMessagesRepository) GetUsersChats(
	user_uid string,
	limit,
	offset int,
) ([]domain.ChatInfo, error) {
	var chats []domain.ChatInfo
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
