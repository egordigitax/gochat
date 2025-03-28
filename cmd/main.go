package main

import (
	"chat-service/api/ws"
	"chat-service/api/ws_fb"
	"chat-service/internal/broker"
	"chat-service/internal/chat_list"
	"chat-service/internal/memory"
	"chat-service/internal/messages"
	"chat-service/internal/storage/postgres"
	"chat-service/internal/storage/postgres/postgres_repos"
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

	viper.SetConfigName("config.prod")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()

	redisClient := memory.NewRedisClient()
	postgresClient := postgres.NewPostgresClient()

	messagesCache := memory.NewRedisMessagesCache(redisClient)
	chatsCache := memory.NewRedisChatsCache(redisClient)

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
	messagesBroker := broker.NewRedisMessagesBroker(redisBroker)

	messagesHub := messages.NewMessagesHub(
		messagesService,
		messagesBroker,
	)
	chatsHub := chat_list.NewChatsHub(
		chatsService,
		messagesBroker,
	)
	savingHub := messages.NewSaveMessagesHub(
		messagesBroker,
		messagesCache,
		messagesStorage,
	)

	messagesController := ws_api.NewMessagesWSController(messagesHub)
	chatsController := ws_api.NewChatsWSController(chatsHub)
	fbChatsController := ws_fb.NewChatsWSController(chatsHub)
	fbMessagesController := ws_fb.NewMessagesWSController(messagesHub)

	messagesController.Handle()
	chatsController.Handle()
	fbChatsController.Handle()
	fbMessagesController.Handle()

	go chatsHub.StartPumpChats()
	go messagesHub.StartPumpMessages()
	go savingHub.StartSavingPump()

	log.Println(fmt.Sprintf("Server started on :%d", viper.GetInt("app.port")))

	err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), nil)

	if err != nil {
		panic(err)
	}
}
