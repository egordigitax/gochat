package managers

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/domain/repositories"
	"context"
	"log"
	"slices"
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

	msgChan, err := s.broker.GetMessagesFromQueue(ctx, constants.CHATS_QUEUE)
	if err != nil {
		cancel()
		return err
	}

    // И дальше надо перекидывать все сообщения из саба в очередь на сохранение с ack
    // Или от юзера сразу его складывать в очередь, а здесь забирать и сохранять, после чего отправлять паб в основной канал сообщений

	ticker := time.NewTicker(500 * time.Millisecond)

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

			toSaveArr = append(toSaveArr, msg)

		case <-ticker.C:

			if len(toSaveArr) == 0 {
				continue
			}

			slices.Reverse(toSaveArr)

			err := s.messages.MessagesStorage.SaveMessagesBulk(toSaveArr...)
			if err != nil {
				log.Println("Bulk save failed")
			}

			// Range over redis cache instead of chan OR push all redis cache on app start to msg chan <--

			for _, msg := range toSaveArr {
				err = s.broker.SendMessageToChannel(
					ctx,
					constants.SAVED_MESSAGES_CHANNEL,
					msg,
				)
				if err != nil {
					log.Println("Message dropped")
				}
			}

			s.savedCount += len(toSaveArr)
			log.Println("saved to db: ", s.savedCount)
			toSaveArr = nil
		}
	}
}
