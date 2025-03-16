package main

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/api/ws"
	"chat-service/internal/application/services/chat"
	"chat-service/internal/application/services/history"
	"chat-service/internal/application/services/message"
	"chat-service/internal/infra/broker"
	"chat-service/internal/infra/broker/redis_broker"
	"chat-service/internal/infra/databases/postgres"
	"chat-service/internal/infra/databases/postgres/postgres_repos"
	"chat-service/internal/infra/memory"
	"chat-service/internal/infra/memory/redis_repos"
	"fmt"
	"log"
	"os"

	"net/http"
	// _ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	// uncomment for memory debug via pprof

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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()

	redisClient := memory.NewRedisClient()
	postgresClient := postgres.NewPostgresClient()

	messagesCache := redis_repos.NewRedisMessagesCache(redisClient)
	chatsCache := redis_repos.NewRedisChatsCache(redisClient)

	messagesStorage := postgres_repos.NewPGMessagesStorage(postgresClient)
	chatsStorage := postgres_repos.NewPGChatsStorage(postgresClient)

	messagesService := message.NewMessageService(
		messagesStorage,
		messagesCache,
	)
	chatsService := chat.NewChatsService(
		chatsStorage,
		chatsCache,
	)

	broker := broker.NewRedisBroker()
	messagesBroker := redis_broker.NewRedisMessagesBroker(broker)

	messagesHub := message.NewMessagesHub(messagesService, messagesBroker)
	chatsHub := chat.NewChatsHub(chatsService, messagesBroker)
	savingHub := history.NewSaveMessagesHub(messagesBroker, messagesCache, messagesStorage)

	messagesController := ws_api.NewMessagesWSController(messagesHub)
	chatsController := ws_api.NewChatsWSController(chatsHub)

	messagesController.Handle()
	chatsController.Handle()

	go chatsHub.StartPumpChats()
	go messagesHub.StartPumpMessages()
	go savingHub.StartSavingPump()

	log.Println("Server started on :8080")

	err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), nil)

	if err != nil {
		panic(err)
	}
}
