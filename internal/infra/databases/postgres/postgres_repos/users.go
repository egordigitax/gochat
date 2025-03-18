package postgres_repos

import (
	"chat-service/internal/domain/entities"
	// "chat-service/internal/domain/repositories"
	"chat-service/internal/infra/databases/postgres"
)

type PGUsersStorage struct {
	postgresClient *postgres.PostgresClient
}

func NewPGUsersStorage(
	postgresClient *postgres.PostgresClient,
) *PGUsersStorage {
	return &PGUsersStorage{
		postgresClient: postgresClient,
	}
}

// var _ repositories.UsersStorage = PGUsersStorage{}

func (p PGUsersStorage) GetUserByUid(userUid string) (entities.User, error) {
	user := entities.User{
		Uid:      userUid,
		Nickname: "",
		MediaUrl: "",
	}

	query := `
    SELECT 
        uid, 
        nickname,
        coalesce(media.url, '') as media_url
    FROM 
        users u 
    LEFT JOIN
        media ON media.id = u.photo_id
    WHERE uid = $1
    LIMIT 1;
    `

	err := p.postgresClient.C_RO.Get(&user, query, userUid)
	if err != nil {
		return user, err
	}

	return user, nil
}
