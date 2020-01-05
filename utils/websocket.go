package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WsUpgrader is used to upgrade HTTP connection to websocket connection
var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
