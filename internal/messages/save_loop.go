package messages

import (
	"chat-service/common/constants"
	"chat-service/internal/broker"
	"chat-service/internal/types"
	"context"
	"log"
	"slices"
	"time"

	"github.com/spf13/viper"
)

type SaveMessagesHub struct {
	broker     broker.MessagesAdaptor
	memory     MessagesCache
	storage    MessagesStorage
	savedCount int
}

func NewSaveMessagesHub(
	broker broker.MessagesAdaptor,
	memory MessagesCache,
	storage MessagesStorage,
) *SaveMessagesHub {
	return &SaveMessagesHub{
		broker:  broker,
		memory:  memory,
		storage: storage,
	}
}

func (s *SaveMessagesHub) StartSavingPump() {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := s.broker.GetMessagesFromQueue(ctx, constants.CHATS_QUEUE)
	if err != nil {
		cancel()
	}

	ticker := time.NewTicker(
		viper.GetDuration("app.save_rate") * time.Millisecond,
	)

	defer func() {
		ticker.Stop()
		cancel()
	}()

	var toSaveArr []types.Message

	log.Println("Saving pump started")

	for {
		select {
		case msg := <-msgChan:
			toSaveArr = append(toSaveArr, msg)
			err := s.broker.SendMessageToChannel(ctx, constants.CHATS_CHANNEL, msg)
			if err != nil {
				log.Println("Error saving pump:", err)
			}

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
