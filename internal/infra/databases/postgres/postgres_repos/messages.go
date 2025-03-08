package postgres_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/databases/postgres"
)

type PGMessagesStorage struct {
	postgresClient *postgres.PostgresClient
}

func NewPGMessagesStorage(
	postgresClient *postgres.PostgresClient,
) *PGMessagesStorage {
	return &PGMessagesStorage{
		postgresClient: postgresClient,
	}
}

func (m PGMessagesStorage) GetMessages(
	chat_uid string,
	limit int, offset int,
) ([]entities.Message, error) {
	var messages []entities.Message

	query := `
    SELECT 
        text, 
        user_uid, 
        chat_uid,
        created_at,
        uid,
    FROM users_chats_messages 
    WHERE chat_uid = $1 LIMIT $2;
    `

	err := m.postgresClient.C_RO.Select(
		&messages,
		query,
		chat_uid,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m PGMessagesStorage) SaveMessage(msg entities.Message) error {
	query := `
    INSERT INTO users_chats_messages (text, user_uid, chat_uid)
    SELECT $1, $2, $3
    `

	_, err := m.postgresClient.C_RW.Exec(
		query,
		msg.Text,
		msg.UserUid,
		msg.ChatUid,
	)

	if err != nil {
		return err
	}

	return nil
}
