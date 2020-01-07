package game

import (
	"log"
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
	Board                [][]Pixel
	PlayersCount         int `json:"players_count" form:"players_count"`
	CurrentTurn          int `json:"current_turn"`
	AllPlayers           []Player
	InstanceID           string
	CurrentActivePlayers int
	Dimension            int `json:"dimension"`
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
	PlayerID string `json:"player_id"`
}

// InitBroadcast initializes brodcast channel
func (i *Instance) InitBroadcast() {
	i.broadcast = make(chan Move)
}

// GetBroadcast return brodcast channel
func (i *Instance) GetBroadcast() chan Move {
	return i.broadcast
}

// BroadcastMoves brodcasts move to all players
func (i *Instance) BroadcastMoves() {
	for {
		move := <-i.broadcast
		for _, p := range i.AllPlayers {
			err := p.WsConnection.WriteJSON(move)
			if err != nil {
				log.Printf("error: %v", err)
				p.WsConnection.Close()
				p.WsConnection = nil
			}
		}
	}
}

// CheckIfColorSelected checks if given color is already selected by another player
func (i *Instance) CheckIfColorSelected(color string) bool {
	for _, p := range i.AllPlayers {
		if p.Color == color {
			return true
		}
	}
	return false
}
