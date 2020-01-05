package game

import "time"

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
}

// Player represents a single player
type Player struct {
	PlayerID string
	Color    string
}
