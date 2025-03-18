package postgres_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/databases/postgres"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
)

type PGChatsStorage struct {
	postgresClient *postgres.PostgresClient
}

func NewPGChatsStorage(postgresClient *postgres.PostgresClient) *PGChatsStorage {
	return &PGChatsStorage{
		postgresClient: postgresClient,
	}
}

func (m *PGChatsStorage) GetAllUsersFromChatByUid(chatUid string) ([]string, error) {
	var userUids []string

	query := `
    SELECT
        users_uids
    FROM users_chats
    WHERE uid = $1
    `

	err := m.postgresClient.C_RO.QueryRow(query, chatUid).Scan(pq.Array(&userUids))
	if err != nil {
		return nil, err
	}

	return userUids, nil
}

func (m *PGChatsStorage) GetChatByUid(chatUid string) (entities.Chat, error) {
	var chat entities.Chat
	query := `
    SELECT 
        uc.id,
        uc.title,
        COALESCE(uc.media_url, '') AS media_url,
        uc.users_uids,
        uc.updated_at,
        uc.uid,
        uc.chat_type
    FROM users_chats uc
    WHERE uc.uid = $1;
    `

	err := m.postgresClient.C_RO.Get(&chat, query, chatUid)
	if err != nil {
		log.Println(err)
		return chat, err
	}

	return chat, nil
}

func (m *PGChatsStorage) GetChatsByUserUid(userUid string, limit, offset int) ([]entities.Chat, error) {
	var chats []entities.Chat

	query := `
    SELECT 
        uc.id,
        uc.title,
        COALESCE(uc.media_url, '') AS media_url,
        uc.updated_at,
        uc.uid,
        uc.chat_type
    FROM users_chats uc
    JOIN users u ON u.uid = ANY(uc.users_uids)
    WHERE u.uid = $1
    LIMIT $2 OFFSET $3;
    `

	err := m.postgresClient.C_RO.Select(&chats, query, userUid, limit, offset)
	if err != nil {
		log.Printf("Error fetching chats for user %s: %v\n", userUid, err)
		return nil, err
	}

	return chats, nil
}

func (m *PGChatsStorage) FetchChatsLastMessages(chats *[]entities.Chat) error {
	if len(*chats) == 0 {
		return errors.New("chats list is empty")
	}

	chatUIDs := make([]string, len(*chats))
	for i, chat := range *chats {
		chatUIDs[i] = chat.Uid
	}

	query := `
    SELECT chats.chat_uid, ucm.text, ucm.created_at, u.nickname, u.uid
    FROM (
        SELECT DISTINCT chat_uid FROM users_chats_messages WHERE chat_uid = ANY($1)
    ) AS chats
    CROSS JOIN LATERAL (
        SELECT chat_uid, text, created_at, user_uid
        FROM users_chats_messages 
        WHERE chat_uid = chats.chat_uid
        ORDER BY created_at DESC
        LIMIT 1
    ) AS ucm
    JOIN users u ON ucm.user_uid = u.uid;
    `

	rows, err := m.postgresClient.C_RO.Query(query, pq.Array(chatUIDs))
	if err != nil {
		log.Printf("Error fetching last messages: %v\n", err)
		return err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows")
		}
	}(rows)

	messages := make(map[string]entities.Message, len(*chats))

	for rows.Next() {
		var message entities.Message
		if err := rows.Scan(
			&message.ChatUid,
			&message.Text,
			&message.CreatedAt,
			&message.UserInfo.Nickname,
			&message.UserInfo.Uid,
		); err != nil {
			return err
		}
		messages[message.ChatUid] = message
	}

	for i := range *chats {
		if msg, exists := messages[(*chats)[i].Uid]; exists {
			(*chats)[i].LastMessage = msg
		}
	}

	return rows.Err()
}

func (m *PGChatsStorage) CheckIfUserHasAccess(
	userUid string,
	chatUid string,
) (bool, error) {
	query := `
    SELECT EXISTS (
        SELECT 1 FROM users_chats 
        WHERE uid = $1 
        AND $2 = ANY(users_uids)
    );`

	var isAccess bool
	err := m.postgresClient.C_RO.Get(
		&isAccess,
		query,
		chatUid,
		userUid,
	)

	return isAccess, err
}

//
// func (m *PGChatsStorage) CreateNewChat(chat domain.Chat) (string, error) {
//
// }
