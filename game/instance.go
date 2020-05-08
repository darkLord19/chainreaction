package game

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	moveRcvMsg int = iota
	moveBcastMsg
	stateUpBcastMsg
	didUserWonMsg
)

type Mutex chan struct{}

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
	Winner               Player
	RoomName             string
	Dimension            int `json:"dimension" form:"dimension"`
	CreatedOn            time.Time
	ExpiresOn            time.Time
	CurrentActivePlayers int
	getMove              chan Move
	broadcastBoardFlag   bool
	didWin               bool
	bbcastMutex          Mutex //bbcastMutex protects read write to broadcastBoardFlag
}

// Player represents a single player
type Player struct {
	UserName     string
	Color        string
	mutex        Mutex
	wsConnection *websocket.Conn
}

// Move struct is used to get Move messages from websocket client
type Move struct {
	MsgType        int    `json:"msg_type"`
	XPos           int    `json:"xpos"`
	YPos           int    `json:"ypos"`
	PlayerUserName string `json:"player_username"`
}

// NewState struct is used to represent board update for websocket broadcast
type NewState struct {
	MsgType     int       `json:"msg_type"`
	NewCurrTurn string    `json:"new_currturn"`
	NewBoard    [][]Pixel `json:"new_board"`
}

// Winner struct is used to send winner notification to users
type Winner struct {
	MsgType int    `json:"msg_type"`
	Winner  Player `json:"winner"`
}

// Lock the kraken
func (m Mutex) Lock() {
	<-m
}

// Unlock the kraken
func (m Mutex) Unlock() {
	m <- struct{}{}
}

// InitMutex initializes mutex
func (p *Player) InitMutex() {
	p.mutex = make(Mutex, 1)
	p.mutex.Unlock()
}

// InitBbcastMutex initializes bbcastmutex
func (i *Instance) InitBbcastMutex() {
	i.bbcastMutex = make(Mutex, 1)
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

func (p *Player) writeToWebsocket(val interface{}) error {
	p.mutex.Lock()
	err := p.wsConnection.WriteJSON(val)
	p.mutex.Unlock()
	return err
}

// SetWsConnection sets ws connection field of Player
func (p *Player) SetWsConnection(ws *websocket.Conn) {
	p.wsConnection = ws
}

// BroadcastMoves brodcasts move to all players
func (i *Instance) BroadcastMoves() {
	for {
		move := <-i.getMove
		for _, p := range i.AllPlayers {
			move.MsgType = moveBcastMsg
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
				msg := NewState{stateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := p.writeToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.wsConnection.Close()
					p.wsConnection = nil
				}
			}
			i.SetBroadcastBoardFlag(false)
		}
		i.bbcastMutex.Unlock()
	}
}

// BroadcastWinner broadcasts winner to users
func (i *Instance) BroadcastWinner() {
	for {
		if i.didWin {
			for _, p := range i.AllPlayers {
				msg := Winner{didUserWonMsg, i.Winner}
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
