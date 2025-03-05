package main

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/api/ws_api"
	"chat-service/internal/application/services"
	"chat-service/internal/infra/databases/postgres"
	"chat-service/internal/infra/databases/postgres/postgres_repos"
	"chat-service/internal/infra/memory"
	"chat-service/internal/infra/memory/redis_repos"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(utils.GenerateJWT("6328040e-6a50-49b3-92fc-4d31e53c2dab"))

	if os.Getenv("SWAGGER_HOST") == "" {
		log.Println("Load ENV from file")
		err := godotenv.Load(".env")
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Println("Load ENV from OS")
	}

	redisClient := memory.NewRedisClient()
	MessagesCache := redis_repos.NewRedisMessagesCache(redisClient)
	ChatsCache := redis_repos.NewRedisChatsCache(redisClient)

	postgresClient := postgres.NewPostgresClient()
	MessagesStorage := postgres_repos.NewPGMessagesStorage(postgresClient)
	ChatsStorage := postgres_repos.NewPGChatsStorage(postgresClient)

	messagesService := services.NewMessageService(
		MessagesStorage,
		MessagesCache,
		ChatsStorage,
		ChatsCache,
	)

	hub := services.NewHub(MessagesStorage)

	chatsHub := services.NewChatsHub(messagesService)

	go chatsHub.Run()
	go hub.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ws_api.ServeWebSocket(hub, w, r)
	})

	http.HandleFunc("/chatlist", func(w http.ResponseWriter, r *http.Request) {
		ws_api.ServeMainWebSocket(chatsHub, w, r)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
