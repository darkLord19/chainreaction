package game

import (
	"time"

	"github.com/gorilla/websocket"
)

// Pixel represents current state of one pixel on board
type Pixel struct {
	DotCount int    `json:"dot_count"`
	Color    string `json:"color"`
}

// Instance represents a single game instance
type Instance struct {
	Board                [32][32]Pixel
	PlayersCount         int `json:"players_count" form:"players_count"`
	CurrentTurn          int `json:"current_turn"`
	AllPlayers           [2]Player
	InstanceID           string
	CurrentActivePlayers int
	CreatedOn            time.Time
	ExpiresOn            time.Time
	broadcast            chan Move
}

// Player represents a single player
type Player struct {
	PlayerID     string
	Color        string
	WsConnection *websocket.Conn
}

// Move struct is used to get Move messages from websocket client
type Move struct {
	XPos     int    `json:"xpos"`
	YPos     int    `json:"ypos"`
	PlayerID string `player_id:"player_id"`
}

// InitBroadcast initializes brodcast channel
func (i *Instance) InitBroadcast() {
	i.broadcast = make(chan Move)
}

// GetBroadcast return brodcast channel
func (i *Instance) GetBroadcast() *chan Move {
	return &i.broadcast
}
