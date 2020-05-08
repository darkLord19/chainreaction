package models

import (
	"log"

	"github.com/chainreaction/constants"

	"github.com/chainreaction/utils"
	"github.com/gorilla/websocket"
)

// Player represents a single player
type Player struct {
	UserName     string
	Color        string
	mutex        utils.Mutex
	wsConnection *websocket.Conn
}

// InitMutex initializes mutex
func (p *Player) InitMutex() {
	p.mutex = make(utils.Mutex, 1)
	p.mutex.Unlock()
}

// WriteToWebsocket writes to player's websocket
func (p *Player) WriteToWebsocket(val interface{}) error {
	p.mutex.Lock()
	err := p.wsConnection.WriteJSON(val)
	p.mutex.Unlock()
	return err
}

// CleanupWs closes websocket connection and frees memory for wsConnection
func (p *Player) CleanupWs() {
	p.wsConnection.Close()
	p.wsConnection = nil
}

// SetWsConnection sets ws connection field of Player
func (p *Player) SetWsConnection(ws *websocket.Conn) {
	p.wsConnection = ws
}

// NotifyIndividual notifies individual player
func (p *Player) NotifyIndividual(val string) {
	tmp := Err{}
	tmp.MsgType = constants.InvalidMoveMsg
	tmp.ErrStr = val
	err := p.WriteToWebsocket(tmp)
	if err != nil {
		log.Printf("error: %v", err)
		p.wsConnection.Close()
		p.wsConnection = nil
	}
}
