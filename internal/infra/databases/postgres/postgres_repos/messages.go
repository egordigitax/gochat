package postgres_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/databases/postgres"
	"fmt"
	"strings"
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

func (m PGMessagesStorage) SaveMessagesBulk(msgs ...entities.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	query := `
    INSERT INTO users_chats_messages (text, user_uid, chat_uid)
    VALUES 
    `

	args := []interface{}{}
	values := []string{}

	for i, msg := range msgs {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		args = append(args, msg.Text, msg.UserUid, msg.ChatUid)
	}

	query += strings.Join(values, ", ")

	_, err := m.postgresClient.C_RW.Exec(query, args...)
	return err
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
        uid
    FROM users_chats_messages 
    WHERE chat_uid = $1 
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;
    `

	err := m.postgresClient.C_RO.Select(
		&messages,
		query,
		chat_uid,
		limit,
        offset,
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
