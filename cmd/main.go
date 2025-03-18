package main

import (
	ws_api2 "chat-service/api/ws"
	"chat-service/internal/application/use_cases/chat_list"
	"chat-service/internal/application/use_cases/messages"
	"chat-service/internal/application/use_cases/save_history"
	"chat-service/internal/infra/broker"
	"chat-service/internal/infra/broker/redis_broker"
	"chat-service/internal/infra/databases/postgres"
	"chat-service/internal/infra/databases/postgres/postgres_repos"
	"chat-service/internal/infra/memory"
	"chat-service/internal/infra/memory/redis_repos"
	"chat-service/internal/utils"
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

	messagesService := messages.NewMessageService(
		messagesStorage,
		messagesCache,
	)
	chatsService := chat_list.NewChatsService(
		chatsStorage,
		chatsCache,
	)

	redisBroker := broker.NewRedisBroker()
	messagesBroker := redis_broker.NewRedisMessagesBroker(redisBroker)

	messagesHub := messages.NewMessagesHub(
		messagesService,
		messagesBroker,
	)
	chatsHub := chat_list.NewChatsHub(
		chatsService,
		messagesBroker,
	)
	savingHub := save_history.NewSaveMessagesHub(
		messagesBroker,
		messagesCache,
		messagesStorage,
	)

	messagesController := ws_api2.NewMessagesWSController(messagesHub)
	chatsController := ws_api2.NewChatsWSController(chatsHub)

	messagesController.Handle()
	chatsController.Handle()

	go chatsHub.StartPumpChats()
	go messagesHub.StartPumpMessages()
	go savingHub.StartSavingPump()

	log.Println(fmt.Sprintf("Server started on :%d", viper.GetInt("app.port")))

	err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), nil)

	if err != nil {
		panic(err)
	}
}
