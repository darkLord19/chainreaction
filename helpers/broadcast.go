package helpers

import (
	"log"

	"github.com/chainreaction/constants"
	"github.com/chainreaction/models"
)

// BroadcastMoves brodcasts move to all players
func BroadcastMoves(i *models.Instance) {
	var move models.MoveMsg
	for {
		i.ReadMoveChan(&move)
		for _, p := range i.AllPlayers {
			move.MsgType = constants.MoveBcastMsg
			err := p.WriteToWebsocket(move)
			if err != nil {
				log.Printf("error: %v", err)
				p.CleanupWs()
			}
		}
	}
}

// BroadcastBoardUpdates broadcasts board updates to users
func BroadcastBoardUpdates(i *models.Instance) {
	for {
		i.HandleMutexes("bbcast", "lock")
		val := i.ReadUnsafe("bbcastFlag")
		if val == nil {
			continue
		}
		if val.(bool) {
			for _, p := range i.AllPlayers {
				msg := models.NewStateMsg{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := p.WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.CleanupWs()
				}
			}
			i.HandleMutexes("bbcast", "unlock")
			i.SetBroadcastBoardFlag(false)
		} else {
			i.HandleMutexes("bbcast", "unlock")
		}
	}
}

// BroadcastWinner broadcasts winner to users
func BroadcastWinner(i *models.Instance) {
	for {
		i.HandleMutexes("winnerBcast", "lock")
		val := i.ReadUnsafe("didWin")
		if val == nil {
			continue
		}
		if val.(bool) {
			for _, p := range i.AllPlayers {
				msg := models.WinnerMsg{constants.UserWonMsg, *i.Winner}
				err := p.WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					p.CleanupWs()
				}
			}
			i.IsOver = true
		}
		i.HandleMutexes("winnerBcast", "unlock")
	}
}
