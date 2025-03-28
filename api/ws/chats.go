package ws_api

import (
	"chat-service/internal/chat_list"
	"chat-service/internal/types"
	"chat-service/internal/types/actions"
	"chat-service/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//TODO: Use worker pool instead goroutines directly

type ChatsHandlerFunc func(
	ctx context.Context,
	data interface{},
	client *chat_list.ChatsClient,
) error

type ChatsResponseFunc func(
	ctx context.Context,
	data types.Action,
	client *chat_list.ChatsClient,
) error

type ChatsWSController struct {
	hub       *chat_list.ChatsHub
	handlers  map[ActionType]ChatsHandlerFunc
	responses map[ActionType]ChatsResponseFunc
}

func NewChatsWSController(
	hub *chat_list.ChatsHub,
) *ChatsWSController {

	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) Handle() {

	c.responses = map[ActionType]ChatsResponseFunc{
		ActionType(types.REQUEST_CHATS): c.ResponseRequestChats,
	}

	http.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
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
		responseHandler, ok := c.responses[ActionType(msg.Action)]
		if !ok {
			log.Println("wrong action type recieved on ChatsWs")
			continue
		}

		err := responseHandler(ctx, msg, client)
		if err != nil {
			log.Println("failed to response user on ChatWs")
		}
	}
}

func (c *ChatsWSController) ResponseRequestChats(
	ctx context.Context,
	data types.Action,
	client *chat_list.ChatsClient,
) error {
	actionData, ok := data.Data.(actions.RequestUserChatsAction)
	if !ok {
		log.Println("Got wrong type of RequestUserChatsAction")
	}

	chats := make([]Chat, len(actionData.Items))
	for i, item := range actionData.Items {
		chats[i] = Chat{
			Title:       item.Title,
			UnreadCount: 0,
			LastMessage: item.LastMessage.Text,
			LastAuthor:  item.LastMessage.UserInfo.Nickname,
			MediaUrl:    item.MediaURL,
		}
	}

	payload := GetChatsResponse{
		Items: chats,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("failed to marshal payload of RequestUserChatsAction")
	}

	response, err := PackToRootMessage(jsonPayload, actionData)
	if err != nil {
		log.Println("failed to pack the root message of RequestUserChatsAction")
	}

	return client.Conn.WriteMessage(websocket.TextMessage, response)
}
