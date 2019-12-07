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
	PlayersCount         int `json:"players_count"`
	CurrentTurn          int `json:"current_turn"`
	InstanceID           string
	CurrentActivePlayers int
	CreatedOn            time.Time
	ExpiresOn            time.Time
}
