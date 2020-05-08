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

// NotifyIndividual notifies individual player
func (p *Player) NotifyIndividual(val string) {
	tmp := Err{}
	tmp.MsgType = constants.InvalidMoveMsg
	tmp.ErrStr = val
	err := p.writeToWebsocket(tmp)
	if err != nil {
		log.Printf("error: %v", err)
		p.wsConnection.Close()
		p.wsConnection = nil
	}
}
