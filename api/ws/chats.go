package ws_api

import (
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/application/use_cases/chat_list"
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
	data resources.Action,
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
		ActionType(resources.REQUEST_CHATS): c.ResponseRequestChats,
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
	data resources.Action,
	client *chat_list.ChatsClient,
) error {
	actionData, ok := data.Data.(dto.RequestUserChatsPayload)
	if !ok {
		log.Println("Got wrong type of RequestUserChatsPayload")
	}

	chats := make([]Chat, len(actionData.Items))
	for i, item := range actionData.Items {
		chats[i] = Chat{
			Title:       item.Title,
			UnreadCount: item.UnreadCount,
			LastMessage: item.LastMessage.Text,
			LastAuthor:  item.LastMessage.Username,
			MediaUrl:    item.MediaUrl,
		}
	}

	payload := GetChatsResponse{
		Items: chats,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("failed to marshal payload of RequestUserChatsPayload")
	}

	response, err := PackToRootMessage(jsonPayload, actionData)
	if err != nil {
		log.Println("failed to pack the root message of RequestUserChatsPayload")
	}

	return client.Conn.WriteMessage(websocket.TextMessage, response)
}
