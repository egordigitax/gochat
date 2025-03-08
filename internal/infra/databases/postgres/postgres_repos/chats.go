package postgres_repos

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/infra/databases/postgres"
	"github.com/lib/pq"
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
) ([]entities.Chat, error) {
	var chats []entities.Chat
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
    JOIN users u 
        ON u.uid = ANY(uc.users_uids)
    WHERE u.uid = $1
    LIMIT $2 OFFSET $3;
    `

	rows, err := m.postgresClient.C_RO.Query(query, user_uid, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		chat := entities.Chat{}
		err := rows.Scan(
			&chat.Id,
			&chat.Title,
			&chat.MediaURL,
			pq.Array(&chat.UsersUids),
			&chat.UpdatedAt,
			&chat.Uid,
			&chat.ChatType,
		)
		chats = append(chats, chat)

		if err != nil {
			return nil, err
		}
	}

	err = m.FetchChatsLastMessages(&chats)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (m *PGChatsStorage) FetchChatsLastMessages(chats *[]entities.Chat) error {
	if len(*chats) == 0 {
		return nil
	}

	chatIDs := make([]string, len(*chats))
	for i, chat := range *chats {
		chatIDs[i] = chat.Uid
	}

	query := `
    SELECT DISTINCT ON (ucm.chat_uid) 
        ucm.chat_uid, 
        ucm.text,
        ucm.created_at,
        u.nickname
    FROM users_chats_messages ucm
    JOIN users u ON ucm.user_uid = u.uid
    WHERE ucm.chat_uid = ANY($1)
    ORDER BY ucm.chat_uid, ucm.created_at DESC;
    `

	rows, err := m.postgresClient.C_RO.Query(query, pq.Array(chatIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	lastMessages := make(map[string]entities.Message, len(*chats))

	for rows.Next() {
		var message entities.Message
		if err := rows.Scan(
			&message.ChatUid,
			&message.Text,
			&message.CreatedAt,
			&message.UserInfo.Nickname,
		); err != nil {
			return err
		}
		lastMessages[message.ChatUid] = message
	}

	for i := range *chats {
		if msg, exists := lastMessages[(*chats)[i].Uid]; exists {
			(*chats)[i].LastMessage = msg
		}
	}

	return rows.Err()
}

func (m *PGChatsStorage) CheckIfUserHasAccess(
	user_uid string,
	chat_uid string,
) (bool, error) {
	query := `
    SELECT EXISTS (
        SELECT 1 FROM users_chats 
        WHERE chat_uid = $1 
        AND $2 = ANY(users_uids)
    );`

	var isAccess bool
	err := m.postgresClient.C_RO.Get(
		&isAccess,
		query,
		chat_uid,
		user_uid,
	)

	return isAccess, err
}

//
// func (m *PGChatsStorage) CreateNewChat(chat domain.Chat) (string, error) {
//
// }
