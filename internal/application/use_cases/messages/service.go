package messages

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/repositories"
	"log"
)

type MessageService struct {
	MessagesStorage repositories.MessagesStorage
	MessagesCache   repositories.MessagesCache
}

func NewMessageService(
	messagesStorage repositories.MessagesStorage,
	messagesCache repositories.MessagesCache,
) *MessageService {
	return &MessageService{
		MessagesStorage: messagesStorage,
		MessagesCache:   messagesCache,
	}
}

func (m *MessageService) GetMessagesHistory(chatUID string, limit, offset int) ([]entities.Message, error) {
	cacheMsgs, err := m.MessagesCache.GetMessagesByChatUid(chatUID)
	if err != nil {
		return nil, err
	}

	lenCache := len(cacheMsgs)
	log.Println(lenCache)

	// TODO: also messages in cache dont have created_at and cant be sorted

	var result []entities.Message

	if offset < lenCache {

		cacheSlice := cacheMsgs[offset:]

		if len(cacheSlice) >= limit {
			log.Println("got all from cache")
			return cacheSlice[:limit], nil
		}

		result = append(result, cacheSlice...)
		dbMsgs, err := m.MessagesStorage.GetMessages(chatUID, limit-len(cacheSlice), 0)
		if err != nil {
			return nil, err
		}

		result = append(result, dbMsgs...)
		log.Printf("got %d from cache and %d from db", lenCache, len(dbMsgs))
		return result, nil

	} else {
		newOffset := offset - lenCache
		dbMsgs, err := m.MessagesStorage.GetMessages(chatUID, limit, newOffset)
		if err != nil {
			return nil, err
		}
		log.Println("got all from db")
		return dbMsgs, nil
	}
}

func (m *MessageService) SaveMessagesBulk(msgs ...entities.Message) error {
	return m.MessagesStorage.SaveMessagesBulk(msgs...)
}
