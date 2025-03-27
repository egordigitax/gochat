package ws_fb

import (
	"chat-service/gen/fbchat"
	messages2 "chat-service/internal/messages"
	resources2 "chat-service/internal/types"
	"chat-service/internal/types/dto"
	"chat-service/internal/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type MessageHandlerFunc func(
	ctx context.Context,
	data *fbchat.RootMessage,
	client *messages2.MessageClient,
) error

type MessageResponseFunc func(
	ctx context.Context,
	data resources2.Action,
	client *messages2.MessageClient,
) error

type MessagesWSController struct {
	hub       *messages2.MessageHub
	handlers  map[fbchat.ActionType]MessageHandlerFunc
	responses map[resources2.ActionType]MessageResponseFunc
}

func NewMessagesWSController(
	hub *messages2.MessageHub,
) *MessagesWSController {

	return &MessagesWSController{
		hub: hub,
	}
}

func (m *MessagesWSController) Handle() {

	http.HandleFunc("/fb/messages", func(w http.ResponseWriter, r *http.Request) {
		m.ServeMessagesWebSocket(w, r)
	})

	m.handlers = map[fbchat.ActionType]MessageHandlerFunc{
		fbchat.ActionTypeGET_MESSAGE: m.HandleSendMessageAction,
	}

	m.responses = map[resources2.ActionType]MessageResponseFunc{
		resources2.REQUEST_MESSAGE: m.ResponseRequestMessageAction,
	}

	// TODO: add responses instead of standalone serializer (?)
}

func (m *MessagesWSController) ServeMessagesWebSocket(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromHeader(
		r.Header.Get("Authorization"),
	)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatId := r.URL.Query().Get("chat_id")
	if chatId == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	upgrader := utils.GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := messages2.NewMessagesClient(
		m.hub,
		conn,
		userId,
		chatId,
	)

	m.hub.RegisterClient(client)

	go m.StartClientWrite(client)
	go m.StartClientRead(client)
}

func (m *MessagesWSController) StartClientWrite(
	client *messages2.MessageClient,
) {
	ctx := context.Background()

	defer func() {
		client.Hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		responseHandler, ok := m.responses[msg.Action]
		if !ok {
			log.Println("wrong type response recieved")
		}
		responseHandler(ctx, msg, client)

	}
}

func (m *MessagesWSController) StartClientRead(
	client *messages2.MessageClient,
) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		client.Hub.UnregisterClient(client)
		cancel()
	}()

	for {
		_, payload, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("error while reading message", err)
			break
		}

		if len(payload) < 4 {
			continue
		}

		root := fbchat.GetRootAsRootMessage(payload, 0)
		handler, ok := m.handlers[root.ActionType()]

		if !ok {
			// log this to client
			log.Println("unknown action type")
			continue
		}

		err = handler(ctx, root, client)
		if err != nil {
			log.Println("error while handling: ", root.ActionType())
		}
	}
}

func (m *MessagesWSController) HandleSendMessageAction(
	ctx context.Context,
	data *fbchat.RootMessage,
	client *messages2.MessageClient,
) error {

	payload := fbchat.GetRootAsGetMessageFromClientRequest(data.PayloadBytes(), 0)
	log.Println(string(payload.Text()))

	client.SendMessage(ctx, dto.SendMessagePayload{
		ChatUid:   client.ChatUid,
		AuthorUid: client.UserUid,
		CreatedAt: time.Now().String(),
		Text:      string(payload.Text()),
	})

	return nil
}

func (m *MessagesWSController) ResponseRequestMessageAction(
	ctx context.Context,
	data resources2.Action,
	client *messages2.MessageClient,
) error {

	actionData, ok := data.Data.(dto.RequestMessagePayload)
	if !ok {
		client.Conn.WriteJSON("Error while handling RequestMessage")
	}

	response := &fbchat.SendMessageToClientResponseT{
		Text:      actionData.Text,
		AuthorId:  actionData.AuthorUid,
		Nickname:  "unimplemented",
		CreatedAt: actionData.CreatedAt,
	}

	bytes := PackRootMessage(
		fbchat.ActionTypeGET_MESSAGE,
		response,
	)

	return client.Conn.WriteMessage(websocket.BinaryMessage, bytes)
}
