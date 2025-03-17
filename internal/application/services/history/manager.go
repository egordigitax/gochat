package history

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/domain/repositories"
	"context"
	"log"
	"slices"
	"time"

	"github.com/spf13/viper"
)

type SaveMessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	memory     repositories.MessagesCache
	storage    repositories.MessagesStorage
	savedCount int
}

func NewSaveMessagesHub(
	broker events.BrokerMessagesAdaptor,
	memory repositories.MessagesCache,
	storage repositories.MessagesStorage,
) *SaveMessagesHub {
	return &SaveMessagesHub{
		broker:  broker,
		memory:  memory,
		storage: storage,
	}
}

func (s *SaveMessagesHub) StartSavingPump() error {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := s.broker.GetMessagesFromQueue(ctx, constants.CHATS_QUEUE)
	if err != nil {
		cancel()
		return err
	}

	ticker := time.NewTicker(
		viper.GetDuration("app.save_rate") * time.Millisecond,
	)

	defer func() {
		ticker.Stop()
		cancel()
	}()

	var toSaveArr []entities.Message

	log.Println("Saving pump started")

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			toSaveArr = append(toSaveArr, msg)
            s.broker.SendMessageToChannel(ctx, constants.CHATS_CHANNEL, msg)

		case <-ticker.C:

			if len(toSaveArr) == 0 {
				continue
			}

			slices.Reverse(toSaveArr)

			err := s.storage.SaveMessagesBulk(toSaveArr...)
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
			}

			s.savedCount += len(toSaveArr)
			log.Println("saved to db: ", s.savedCount)
			toSaveArr = nil
		}
	}
}
