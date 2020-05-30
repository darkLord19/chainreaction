package models

import (
	"sync"
	"time"

	"github.com/chainreaction/constants"
)

// Instance represents a single game instance
type Instance struct {
	Board                     []Pixel   `json:"-"`
	PlayersCount              int       `json:"players_count" form:"players_count"`
	CurrentTurn               int       `json:"current_turn"`
	AllPlayers                []Player  `json:"all_players"`
	Winner                    *Player   `json:"-"`
	RoomName                  string    `json:"room_name"`
	Dimension                 int       `json:"dimension" form:"dimension"`
	CreatedOn                 time.Time `json:"-"`
	ExpiresOn                 time.Time `json:"-"`
	AvailableColors           [8]string `json:"-"`
	IsOver                    bool      `json:"-"`
	joinedPlayers             int
	currActivePlayersCnt      int
	allPlayedOnce             bool
	RecvMove                  chan MoveMsg `json:"-"`
	UpdatedBoard              chan []int   `json:"-"`
	winnerMutex               sync.RWMutex //winnerMutex protects read write to Winner
	allPlayedMutex            sync.RWMutex //allPlayedMutex protects allPlayedOnce
	currActivePlayersCntMutex sync.RWMutex //currActivePlayersMutex protects read write to CurrentActivePlayers
	joinedPlayersMutex        sync.RWMutex
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
	i.Board = make([]Pixel, i.Dimension*i.Dimension)
	i.RoomName = name
	i.AvailableColors = constants.Colors
	i.RecvMove = make(chan MoveMsg)
	i.UpdatedBoard = make(chan []int)
}

// GetAndIncJoinPlayers gets count of currently joined players count and increase it by one
func (i *Instance) GetAndIncJoinPlayers() int {
	i.joinedPlayersMutex.Lock()
	defer i.joinedPlayersMutex.Unlock()
	v := i.joinedPlayers
	i.joinedPlayers++
	return v
}

// IncCurrentActivePlayers increases current active players count safely
func (i *Instance) IncCurrentActivePlayers() {
	i.currActivePlayersCntMutex.Lock()
	i.currActivePlayersCnt++
	i.currActivePlayersCntMutex.Unlock()
}

// DecCurrentActivePlayersCount decreases active players count safely
func (i *Instance) DecCurrentActivePlayersCount() {
	i.currActivePlayersCntMutex.Lock()
	i.currActivePlayersCnt--
	i.currActivePlayersCntMutex.Unlock()
}

// GetCurrentActivePlayersCount gets the current active players count safely
func (i *Instance) GetCurrentActivePlayersCount() int {
	i.currActivePlayersCntMutex.Lock()
	val := i.currActivePlayersCnt
	i.currActivePlayersCntMutex.Unlock()
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

// IncCellCountOfPlayer accepts color and increases cell count
// for the given player who has that color
func (i *Instance) IncCellCountOfPlayer(color string) {
	for a := range i.AllPlayers {
		if i.AllPlayers[a].Color == color {
			i.AllPlayers[a].CellCount++
		}
	}
}

// DecCellCountOfPlayer accepts color and increases cell count
// for the given player who has that color
func (i *Instance) DecCellCountOfPlayer(color string) {
	for a := range i.AllPlayers {
		if i.AllPlayers[a].Color == color {
			i.AllPlayers[a].CellCount--
		}
	}
}
