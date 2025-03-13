package main

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/api/ws_api"
	"chat-service/internal/application/managers"
	"chat-service/internal/application/services"
	"chat-service/internal/infra/broker"
	"chat-service/internal/infra/broker/redis_broker"
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
	fmt.Println(utils.GenerateJWT("51929f93-fd17-4e9d-b38c-31f4c26fa51c"))

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
	postgresClient := postgres.NewPostgresClient()

	MessagesCache := redis_repos.NewRedisMessagesCache(redisClient)
	ChatsCache := redis_repos.NewRedisChatsCache(redisClient)

	MessagesStorage := postgres_repos.NewPGMessagesStorage(postgresClient)
	ChatsStorage := postgres_repos.NewPGChatsStorage(postgresClient)

	messagesService := services.NewMessageService(
		MessagesStorage,
		MessagesCache,
	)
	chatsService := services.NewChatsService(
		ChatsStorage,
		ChatsCache,
	)

	broker := broker.NewRedisBroker()
	messagesBroker := redis_broker.NewRedisMessagesBroker(broker)

	messagesHub := managers.NewMessagesHub(MessagesStorage, messagesBroker)
	chatsHub := managers.NewChatsHub(messagesService, chatsService, messagesBroker)
	savingHub := managers.NewSaveMessagesHub(messagesBroker, messagesService)

    go savingHub.StartSavingPump()

	MessagesController := ws_api.NewMessagesWSController(messagesHub)
	ChatsController := ws_api.NewChatsWSController(chatsHub)

	go chatsHub.Run()

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		MessagesController.ServeMessagesWebSocket(w, r)
	})

	http.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		ChatsController.ServeChatsWebSocket(w, r)
	})

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
