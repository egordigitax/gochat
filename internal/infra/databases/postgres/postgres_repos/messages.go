package postgres_repos

import (
	"chat-service/internal/domain"
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
) ([]domain.Message, error) {

	var messages []domain.Message
	query := `
    SELECT message, user_uid, chat_uid FROM users_messages WHERE chat_uid = $1
    `

	err := m.postgresClient.C_RO.Select(&messages, query, chat_uid)
	if err != nil {
		return nil, err
	}

	return messages, nil

}

func (m PGMessagesStorage) SaveMessage(msg domain.Message) error {
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
