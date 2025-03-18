package repositories

import "chat-service/internal/domain/entities"

type UsersStorage interface {
	GetUserByUid(userUid string) (entities.User, error)
}
