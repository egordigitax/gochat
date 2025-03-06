package ws_api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: Remove CheckOrigin true
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}
