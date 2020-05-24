package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Player represents a single player
type Player struct {
	UserName     string `json:"username"`
	Color        string `json:"color"`
	mutex        sync.RWMutex
	wsConnection *websocket.Conn
}

// WriteToWebsocket writes to player's websocket
func (p *Player) WriteToWebsocket(val interface{}) error {
	p.mutex.Lock()
	err := p.wsConnection.WriteJSON(val)
	p.mutex.Unlock()
	return err
}

// CleanupWs closes websocket connection and frees memory for wsConnection
func (p *Player) CleanupWs() {
	p.mutex.Lock()
	p.wsConnection.Close()
	p.wsConnection = nil
	p.mutex.Unlock()
}

// SetWsConnection sets ws connection field of Player
func (p *Player) SetWsConnection(ws *websocket.Conn) {
	p.mutex.Lock()
	p.wsConnection = ws
	p.mutex.Unlock()
}
