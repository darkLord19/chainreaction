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
	broadcastMove        chan Move
	broadcastBoardFlag   bool
	broadcastBoard       chan<- NewState
	didWin               bool
	broadcastWinner      chan<- Winner
}

// Player represents a single player
type Player struct {
	UserName     string
	Color        string
	WsConnection *websocket.Conn
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

// InitBroadcasts initializes brodcast channel
func (i *Instance) InitBroadcasts() {
	i.broadcastMove = make(chan Move)
	i.broadcastBoard = make(chan NewState)
}

// GetBroadcastMove return brodcast channel
func (i *Instance) GetBroadcastMove() chan Move {
	return i.broadcastMove
}

// BroadcastMoves brodcasts move to all players
func (i *Instance) BroadcastMoves() {
	for {
		move := <-i.broadcastMove
		for _, p := range i.AllPlayers {
			move.MsgType = moveBcastMsg
			err := p.WsConnection.WriteJSON(move)
			if err != nil {
				log.Printf("error: %v", err)
				p.WsConnection.Close()
				p.WsConnection = nil
			}
		}
	}
}

// BroadcastBoardUpdates broadcasts board updates to users
func (i *Instance) BroadcastBoardUpdates() {
	for {
		if i.broadcastBoardFlag {
			for _, p := range i.AllPlayers {
				msg := NewState{stateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := p.WsConnection.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.WsConnection.Close()
					p.WsConnection = nil
				}
			}
			i.broadcastBoardFlag = false
		}
	}
}

// BroadcastWinner broadcasts winner to users
func (i *Instance) BroadcastWinner() {
	for {
		if i.didWin {
			for _, p := range i.AllPlayers {
				msg := Winner{didUserWonMsg, i.Winner}
				err := p.WsConnection.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.WsConnection.Close()
					p.WsConnection = nil
				}
			}
			i.broadcastBoardFlag = false
		}
	}
}

// SetBroadcastBoardFlag sets broadcast board state flag
func (i *Instance) SetBroadcastBoardFlag() {
	i.broadcastBoardFlag = true
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
