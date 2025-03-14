package managers

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/domain/repositories"
	"context"
	"log"
	"time"
)

type SaveMessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	memory     repositories.MessagesCache
	messages   *services.MessageService
	savedCount int
}

func NewSaveMessagesHub(
	broker events.BrokerMessagesAdaptor,
	messagesService *services.MessageService,
	memory repositories.MessagesCache,
) *SaveMessagesHub {
	return &SaveMessagesHub{
		broker:   broker,
		messages: messagesService,
		memory:   memory,
	}
}

func (s *SaveMessagesHub) StartSavingPump() error {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := s.broker.GetMessagesFromChannel(ctx, constants.CHATS_CHANNEL)
	if err != nil {
		cancel()
		return err
	}

	ticker := time.NewTicker(2 * time.Second)

	defer func() {
		ticker.Stop()
		cancel()
	}()

	var toSaveArr []entities.Message

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			err := s.memory.SaveMessage(msg)
			if err != nil {
				log.Println("Message dropped")
				continue
			}

			toSaveArr = append(toSaveArr, msg)

		case <-ticker.C:

			if len(toSaveArr) == 0 {
				continue
			}

			err := s.messages.MessagesStorage.SaveMessagesBulk(toSaveArr...)
			if err != nil {
				log.Println("Bulk save failed")
			}

			for _, msg := range toSaveArr {
				err = s.broker.SendMessageToChannel(
					ctx,
					constants.SAVED_MESSAGES_CHANNEL,
					msg,
				)
				if err != nil {
					log.Println("Message dropped")
				}

				err := s.memory.DeleteMessage(msg)
				if err != nil {
					log.Println("Cant delete message from redis")
				}
			}

			s.savedCount += len(toSaveArr)
			log.Println("saved to db: ", s.savedCount)
			toSaveArr = nil
		}
	}
}
