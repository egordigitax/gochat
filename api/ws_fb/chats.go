package ws_fb

import (
	"chat-service/gen/fbchat"
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/application/use_cases/chat_list"
	"chat-service/internal/utils"
	"log"
	"net/http"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

//TODO: Use worker pool instead goroutines directly

type ChatsWSController struct {
	hub *chat_list.ChatsHub
}

func NewChatsWSController(
	hub *chat_list.ChatsHub,
) *ChatsWSController {
	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) Handle() {
	http.HandleFunc("/fb/chats", func(w http.ResponseWriter, r *http.Request) {
		c.ServeChatsWebSocket(w, r)
	})
}

func (c *ChatsWSController) ServeChatsWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	upgrader := utils.GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}

	client := chat_list.NewChatsClient(c.hub, conn, userID)

	c.hub.RegisterClient(client)

	go c.StartClientWrite(client)
}

func (c *ChatsWSController) StartClientWrite(client *chat_list.ChatsClient) {
	defer func() {
		c.hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		// handle different actions and parse to schema

		if msg.Action == resources.REQUEST_CHATS {
			message, ok := msg.Data.(dto.RequestUserChatsPayload)
			if !ok {
				log.Println("[ERROR] wrong type of data")
				continue
			}

			items := make([]*fbchat.ChatT, len(message.Items))
			for i, item := range message.Items {
				items[i] = &fbchat.ChatT{
					Title:       item.Title,
					UnreadCount: int32(item.UnreadCount),
					LastMessage: item.LastMessage.Text,
					LastAuthor:  item.LastMessage.AuthorUid,
					MediaUrl:    item.MediaUrl,
				}
			}

			data := fbchat.GetChatsResponseT{
				Items: items,
			}

			response := fbchat.RootMessageT{
				ActionType: fbchat.ActionTypeGET_CHATS,
				Payload: &fbchat.RootMessagePayloadT{
					Type:  fbchat.RootMessagePayloadGetChatsResponse,
					Value: &data,
				},
			}

			builder := flatbuffers.NewBuilder(1024)
			payload := response.Pack(builder)
			builder.Finish(payload)
			fBytes := builder.FinishedBytes()

			if err := client.Conn.WriteMessage(websocket.BinaryMessage, fBytes); err != nil {
				break
			}
		}
	}
}
