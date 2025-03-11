package repositories

import "chat-service/internal/domain/entities"

type UsersStorage interface {
    GetUserByUid(user_uid string) (entities.User, error)
}
