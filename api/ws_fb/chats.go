package ws_fb

import (
	"chat-service/gen/fbchat"
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/application/use_cases/chat_list"
	"chat-service/internal/utils"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//TODO: Use worker pool instead goroutines directly

type ChatsHandlerFunc func(ctx context.Context, data interface{}, client *chat_list.ChatsClient) error
type ChatsResponseFunc func(ctx context.Context, data resources.Action, client *chat_list.ChatsClient) error

type ChatsWSController struct {
	hub              *chat_list.ChatsHub
	responseHandlers map[resources.ActionType]ChatsResponseFunc
}

func NewChatsWSController(
	hub *chat_list.ChatsHub,
) *ChatsWSController {
	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) Handle() {

	c.responseHandlers = map[resources.ActionType]ChatsResponseFunc{
		resources.REQUEST_CHATS: c.ResponseRequestChats,
	}

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
	ctx := context.Background()

	defer func() {
		c.hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		handler, ok := c.responseHandlers[msg.Action]
		if !ok {
			log.Println("wrong action type")
		}

        err := handler(ctx, msg, client)
        if err != nil {
            log.Println("error while handling")
        }
	}
}

func (c *ChatsWSController) ResponseRequestChats(
	ctx context.Context,
	data resources.Action,
	client *chat_list.ChatsClient,
) error {
	actionData, ok := data.Data.(dto.RequestUserChatsPayload)
	if !ok {
		log.Println("Got wrong type of RequestUserChatsPayload")
	}

	chats := make([]*fbchat.ChatT, len(actionData.Items))
	for i, item := range actionData.Items {
		chats[i] = &fbchat.ChatT{
			Title:       item.Title,
			UnreadCount: int32(item.UnreadCount),
			LastMessage: item.LastMessage.Text,
			LastAuthor:  item.LastMessage.Username,
			MediaUrl:    item.MediaUrl,
		}
	}

	payload := &fbchat.GetChatsResponseT{
		Items: chats,
	}

	bytes := PackRootMessage(
		fbchat.ActionTypeGET_CHATS,
		fbchat.RootMessagePayloadGetChatsResponse,
		payload,
	)

	return client.Conn.WriteMessage(websocket.BinaryMessage, bytes)
}
