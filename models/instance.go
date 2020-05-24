package models

import (
	"sync"
	"time"

	"github.com/chainreaction/constants"
)

// Instance represents a single game instance
type Instance struct {
	Board                  [][]Pixel
	PlayersCount           int `json:"players_count" form:"players_count"`
	CurrentTurn            int `json:"current_turn"`
	AllPlayers             []Player
	Winner                 *Player
	RoomName               string
	Dimension              int `json:"dimension" form:"dimension"`
	CreatedOn              time.Time
	ExpiresOn              time.Time
	AvailableColors        map[string]bool
	IsOver                 bool
	currentActivePlayers   int
	allPlayedOnce          bool
	RecvMove               chan MoveMsg `json:"-"`
	UpdatedBoard           chan []int   `json:"-"`
	winnerMutex            sync.RWMutex //winnerMutex protects read write to Winner
	allPlayedMutex         sync.RWMutex //allPlayedMutex protects allPlayedOnce
	currActivePlayersMutex sync.RWMutex //currActivePlayersMutex protects read write to CurrentActivePlayers
}

// Pixel represents current state of one pixel on board
type Pixel struct {
	DotCount int    `json:"dot_count"`
	Color    string `json:"color"`
}

// Init initializes game instance
func (i *Instance) Init(name string) {
	i.CreatedOn = time.Now().UTC()
	i.ExpiresOn = i.CreatedOn.Add(time.Minute * time.Duration(25))
	i.CurrentTurn = 0
	i.Board = make([][]Pixel, i.Dimension)
	i.RoomName = name
	for a := 0; a < i.Dimension; a++ {
		i.Board[a] = make([]Pixel, i.Dimension)
	}
	i.AvailableColors = make(map[string]bool)
	for _, c := range constants.Colors {
		i.AvailableColors[c] = true
	}
	i.RecvMove = make(chan MoveMsg)
	i.UpdatedBoard = make(chan []int)
}

// CheckIfColorSelected checks if given color is already selected by another player
func (i *Instance) CheckIfColorSelected(color string) bool {
	_, v := i.AvailableColors[color]
	return !v
}

// CheckIfValidColor checks if recvd color is valid
func (i *Instance) CheckIfValidColor(color string) bool {
	for k := range i.AvailableColors {
		if k == color {
			return true
		}
	}
	return false
}

// IncCurrentActivePlayers increases current active players count safely
func (i *Instance) IncCurrentActivePlayers() {
	i.currActivePlayersMutex.Lock()
	i.currentActivePlayers++
	i.currActivePlayersMutex.Unlock()
}

// DecCurrentActivePlayers decreases active players count safely
func (i *Instance) DecCurrentActivePlayers() {
	i.currActivePlayersMutex.Lock()
	i.currentActivePlayers--
	i.currActivePlayersMutex.Unlock()
}

// GetCurrentActivePlayers gets the current active players count safely
func (i *Instance) GetCurrentActivePlayers() int {
	i.currActivePlayersMutex.Lock()
	val := i.currentActivePlayers
	i.currActivePlayersMutex.Unlock()
	return val
}

// GetPlayerByUsername returns Player struct from instance id
func (i *Instance) GetPlayerByUsername(username string) (*Player, bool) {
	for a := range i.AllPlayers {
		if i.AllPlayers[a].UserName == username {
			return &i.AllPlayers[a], true
		}
	}
	return nil, false
}

// SetWinner sets winner of game
func (i *Instance) SetWinner(p *Player) {
	i.winnerMutex.Lock()
	i.Winner = p
	i.winnerMutex.Unlock()
}

// GetWinner sets winner of game
func (i *Instance) GetWinner() *Player {
	i.winnerMutex.Lock()
	defer i.winnerMutex.Unlock()
	return i.Winner
}

// SetIfAllPlayedOnce sets if everyone played once
func (i *Instance) SetIfAllPlayedOnce(val bool) {
	i.allPlayedMutex.Lock()
	i.allPlayedOnce = val
	i.allPlayedMutex.Unlock()
}

// GetIfAllPlayedOnce sets if everyone played once
func (i *Instance) GetIfAllPlayedOnce() bool {
	i.allPlayedMutex.Lock()
	defer i.allPlayedMutex.Unlock()
	return i.allPlayedOnce
}
