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
		for x := range i.AllPlayers {
			move.MsgType = constants.MoveBcastMsg
			err := i.AllPlayers[x].WriteToWebsocket(move)
			if err != nil {
				log.Printf("error: %v", err)
				i.AllPlayers[x].CleanupWs()
			}
		}
	}
}

// BroadcastBoardUpdates broadcasts board updates to users
func BroadcastBoardUpdates(i *models.Instance) {
	for {
		i.HandleMutexes("bbcast", "lock")
		val := i.ReadUnsafe("bbcastFlag")
		if val.(bool) {
			for x := range i.AllPlayers {
				msg := models.NewStateMsg{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.WriteUnsafe("bbcast", false)
		}
		i.HandleMutexes("bbcast", "unlock")
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
			for x := range i.AllPlayers {
				msg := models.WinnerMsg{constants.UserWonMsg, *i.Winner}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.IsOver = true
		}
		i.HandleMutexes("winnerBcast", "unlock")
	}
}
