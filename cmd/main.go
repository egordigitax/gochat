package main

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/api/ws_api"
	"chat-service/internal/application/managers"
	"chat-service/internal/application/services"
	"chat-service/internal/infra/databases/postgres"
	"chat-service/internal/infra/databases/postgres/postgres_repos"
	"chat-service/internal/infra/memory"
	"chat-service/internal/infra/memory/redis_repos"
	"chat-service/internal/infra/memory/redis_subs"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	MessagesCache := redis_repos.NewRedisMessagesCache(redisClient)
	ChatsCache := redis_repos.NewRedisChatsCache(redisClient)

	postgresClient := postgres.NewPostgresClient()
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

	updateChan := make(chan string, 100)

	broker := redis_subs.NewRedisMessageBroker()
	go func() {
		ch, err := broker.Subscribe("b5c6c767-cb9f-43a5-b8b5-7fe4ff083447")
		if err != nil {
			log.Println(err.Error())
		}

		for {
			select {
			case msg := <-ch:
                log.Println("got from sub: ", msg)
			}
		}
	}()

	time.Sleep(5 * time.Millisecond)

	broker.Publish("test-chat", "hello!")
	broker.Publish("test-chat", "hello!")
	broker.Publish("test-chat", "hello!")
	broker.Publish("test-chat", "im gay!")
	broker.Publish("test-chat", "hello!")
	broker.Publish("test-chat", "hello!")

	messagesHub := managers.NewMessagesHub(MessagesStorage, updateChan, broker)
	chatsHub := managers.NewChatsHub(messagesService, chatsService, updateChan)

	go chatsHub.Run()
	go messagesHub.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ws_api.ServeWebSocket(messagesHub, w, r)
	})

	http.HandleFunc("/chatlist", func(w http.ResponseWriter, r *http.Request) {
		ws_api.ServeMainWebSocket(chatsHub, w, r)
	})

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
