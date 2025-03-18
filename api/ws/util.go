package ws_api

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// TODO: Remove CheckOrigin true

func GetUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
}
