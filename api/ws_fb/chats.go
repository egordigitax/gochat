package ws_fb

import (
	"chat-service/gen/fbchat"
	chat_list2 "chat-service/internal/chat_list"
	resources2 "chat-service/internal/types"
	"chat-service/internal/types/actions"
	"chat-service/internal/utils"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//TODO: Use worker pool instead goroutines directly

type ChatsHandlerFunc func(ctx context.Context, data interface{}, client *chat_list2.ChatsClient) error
type ChatsResponseFunc func(ctx context.Context, data resources2.Action, client *chat_list2.ChatsClient) error

type ChatsWSController struct {
	hub              *chat_list2.ChatsHub
	responseHandlers map[resources2.ActionType]ChatsResponseFunc
}

func NewChatsWSController(
	hub *chat_list2.ChatsHub,
) *ChatsWSController {
	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) Handle() {

	c.responseHandlers = map[resources2.ActionType]ChatsResponseFunc{
		resources2.REQUEST_CHATS: c.ResponseRequestChats,
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

	client := chat_list2.NewChatsClient(c.hub, conn, userID)

	c.hub.RegisterClient(client)

	go c.StartClientWrite(client)
}

func (c *ChatsWSController) StartClientWrite(client *chat_list2.ChatsClient) {
	ctx := context.Background()

	defer func() {
		c.hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		handler, ok := c.responseHandlers[msg.Action]
		if !ok {
			log.Println("wrong action type")
			continue
		}

		err := handler(ctx, msg, client)
		if err != nil {
			log.Println("error while handling")
		}
	}
}

func (c *ChatsWSController) ResponseRequestChats(
	ctx context.Context,
	data resources2.Action,
	client *chat_list2.ChatsClient,
) error {
	actionData, ok := data.Data.(actions.RequestUserChatsAction)
	if !ok {
		log.Println("Got wrong type of RequestUserChatsAction")
	}

	chats := make([]*fbchat.ChatT, len(actionData.Items))
	for i, item := range actionData.Items {
		chats[i] = &fbchat.ChatT{
			Title:       item.Title,
			UnreadCount: int32(0),
			LastMessage: item.LastMessage.Text,
			LastAuthor:  item.LastMessage.UserInfo.Nickname,
			MediaUrl:    item.MediaURL,
		}
	}

	payload := &fbchat.GetChatsResponseT{
		Items: chats,
	}

	bytes := PackRootMessage(
		fbchat.ActionTypeGET_CHATS,
		payload,
	)

	return client.Conn.WriteMessage(websocket.BinaryMessage, bytes)
}
