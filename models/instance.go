package models

import (
	"time"

	"github.com/chainreaction/utils"
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
	IsOver                 bool
	currentActivePlayers   int
	getMove                chan MoveMsg
	broadcastBoardFlag     bool
	didWin                 bool
	allPlayedOnce          bool
	allPlayedMutex         utils.Mutex //allPlayedMutex protects read write to allPlayedOnce
	bbcastMutex            utils.Mutex //bbcastMutex protects read write to broadcastBoardFlag
	currActivePlayersMutex utils.Mutex //currActivePlayersMutex protects read write to CurrentActivePlayers
	winnerBcastMutex       utils.Mutex //winnerBcastMutex protects read write to didWin, Winner and IsOver
}

// Pixel represents current state of one pixel on board
type Pixel struct {
	DotCount int    `json:"dot_count"`
	Color    string `json:"color"`
}

// InitGameInstanceMutexes initializes bbcastmutex
func (i *Instance) InitGameInstanceMutexes() {
	i.bbcastMutex = make(utils.Mutex, 1)
	i.currActivePlayersMutex = make(utils.Mutex, 1)
	i.winnerBcastMutex = make(utils.Mutex, 1)
	i.allPlayedMutex = make(utils.Mutex, 1)
	i.currActivePlayersMutex.Unlock()
	i.bbcastMutex.Unlock()
	i.winnerBcastMutex.Unlock()
	i.allPlayedMutex.Unlock()
}

// HandleMutexes handles mutex lock unlock for outside the package
func (i *Instance) HandleMutexes(name string, op string) {
	switch op {
	case "lock":
		switch name {
		case "bbcast":
			i.bbcastMutex.Lock()
			break
		case "allPlayed":
			i.allPlayedMutex.Lock()
			break
		case "currActive":
			i.currActivePlayersMutex.Lock()
			break
		case "winnerBcast":
			i.winnerBcastMutex.Lock()
			break
		}
		break
	case "unlock":
		switch name {
		case "bbcast":
			i.bbcastMutex.Unlock()
			break
		case "allPlayed":
			i.allPlayedMutex.Unlock()
			break
		case "currActive":
			i.currActivePlayersMutex.Unlock()
			break
		case "winnerBcast":
			i.winnerBcastMutex.Unlock()
			break
		}
		break
	}
}

// ReadUnsafe reads values directly without locking
func (i *Instance) ReadUnsafe(name string) interface{} {
	switch name {
	case "didWin":
		return i.didWin
	case "bbcastFlag":
		return i.broadcastBoardFlag
	}
	return nil
}

// InitChannel initializes brodcast channel
func (i *Instance) InitChannel() {
	i.getMove = make(chan MoveMsg)
}

// WriteToMoveCh return brodcast channel
func (i *Instance) WriteToMoveCh(m MoveMsg) {
	i.getMove <- m
}

// ReadMoveChan reads value from move chan
func (i *Instance) ReadMoveChan(m *MoveMsg) {
	*m = <-i.getMove
}

// SetBroadcastBoardFlag sets broadcast board state flag safely
func (i *Instance) SetBroadcastBoardFlag(val bool) {
	i.bbcastMutex.Lock()
	i.broadcastBoardFlag = val
	i.bbcastMutex.Unlock()
}

// GetIfAllPlayedOnce returns if all player played once at least
func (i *Instance) GetIfAllPlayedOnce() bool {
	i.allPlayedMutex.Lock()
	val := i.allPlayedOnce
	i.allPlayedMutex.Unlock()
	return val
}

// SetIfAllPlayedOnce sets if everyone played once
func (i *Instance) SetIfAllPlayedOnce() {
	i.allPlayedMutex.Lock()
	i.allPlayedOnce = true
	i.allPlayedMutex.Unlock()
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

// GetPlayerByID returns Player struct from instance id
func (i *Instance) GetPlayerByID(username string) *Player {
	for a := range i.AllPlayers {
		if i.AllPlayers[a].UserName == username {
			return &i.AllPlayers[a]
		}
	}
	return nil
}

// CheckIfUserNameClaimed checks if given username is claimed by another user or not
func (i *Instance) CheckIfUserNameClaimed(username string) bool {
	for _, p := range i.AllPlayers {
		if p.UserName == username {
			return true
		}
	}
	return false
}

// SetWinner sets winner of game
func (i *Instance) SetWinner(p *Player) {
	i.winnerBcastMutex.Lock()
	i.didWin = true
	i.Winner = p
	i.winnerBcastMutex.Unlock()
}
