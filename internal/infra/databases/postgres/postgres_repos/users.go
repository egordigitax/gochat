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

func (p PGUsersStorage) GetUserByUid(user_uid string) (entities.User, error) {
	panic("unimplemented")
}


