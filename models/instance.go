package models

import (
	"log"
	"time"

	"github.com/chainreaction/constants"
	"github.com/chainreaction/utils"
)

// Instance represents a single game instance
type Instance struct {
	Board                  [][]Pixel
	PlayersCount           int `json:"players_count" form:"players_count"`
	CurrentTurn            int `json:"current_turn"`
	AllPlayers             []Player
	Winner                 Player
	RoomName               string
	Dimension              int `json:"dimension" form:"dimension"`
	CreatedOn              time.Time
	ExpiresOn              time.Time
	IsOver                 bool
	currentActivePlayers   int
	getMove                chan Move
	broadcastBoardFlag     bool
	didWin                 bool
	bbcastMutex            utils.Mutex //bbcastMutex protects read write to broadcastBoardFlag
	currActivePlayersMutex utils.Mutex //currActivePlayersMutex protects read write to CurrentActivePlayers
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
	i.currActivePlayersMutex.Unlock()
	i.bbcastMutex.Unlock()
}

// InitChannel initializes brodcast channel
func (i *Instance) InitChannel() {
	i.getMove = make(chan Move)
}

// WriteToMoveCh return brodcast channel
func (i *Instance) WriteToMoveCh(m Move) {
	i.getMove <- m
}

// SetBroadcastBoardFlag sets broadcast board state flag safely
func (i *Instance) SetBroadcastBoardFlag(val bool) {
	i.bbcastMutex.Lock()
	i.broadcastBoardFlag = val
	i.bbcastMutex.Unlock()
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

// BroadcastMoves brodcasts move to all players
func (i *Instance) BroadcastMoves() {
	for {
		move := <-i.getMove
		for _, p := range i.AllPlayers {
			move.MsgType = constants.MoveBcastMsg
			err := p.writeToWebsocket(move)
			if err != nil {
				log.Printf("error: %v", err)
				p.wsConnection.Close()
				p.wsConnection = nil
			}
		}
	}
}

// BroadcastBoardUpdates broadcasts board updates to users
func (i *Instance) BroadcastBoardUpdates() {
	for {
		i.bbcastMutex.Lock()
		if i.broadcastBoardFlag {
			for _, p := range i.AllPlayers {
				msg := NewState{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := p.writeToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.wsConnection.Close()
					p.wsConnection = nil
				}
			}
			i.broadcastBoardFlag = false
		}
		i.bbcastMutex.Unlock()
	}
}

// BroadcastWinner broadcasts winner to users
func (i *Instance) BroadcastWinner() {
	for {
		if i.didWin {
			for _, p := range i.AllPlayers {
				msg := Winner{constants.UserWonMsg, i.Winner}
				err := p.writeToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.wsConnection.Close()
					p.wsConnection = nil
				}
			}
			i.didWin = false
		}
	}
}
