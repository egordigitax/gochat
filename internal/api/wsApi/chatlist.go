package wsApi

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/domain"
	"chat-service/internal/infra"
	"log"
	"net/http"
)

// WebSocket хендлер для подписки на все чаты (главный экран)
func ServeMainWebSocket(hub *infra.ChatsHub, w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из JWT
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Апгрейдим соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}

	// Создаем клиента
	client := &infra.ChatsClient{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []domain.ChatInfo, 10),
	}

	client.Done = make(chan struct{}, 1) // ✅ Buffered channel
	hub.RegisterClient(client)

	// Запускаем горутины для получения/отправки сообщений
	go client.ReadPump()
	go client.WritePump()

	log.Printf("[INFO] Main WS connected: user=%s", userID)
}
