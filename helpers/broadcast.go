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
		if i.IsOver {
			continue
		}
		i.ReadMoveChan(&move)
		move.Color = i.GetPlayerByID(move.PlayerUserName).Color
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
		if i.GetBroadcastBoardFlag() && !i.IsOver {
			for x := range i.AllPlayers {
				msg := models.NewStateMsg{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName, i.Board}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.SetBroadcastBoardFlag(false)
		}
	}
}

// BroadcastWinner broadcasts winner to users
func BroadcastWinner(i *models.Instance) {
	for {
		if i.GetIfSomeoneWon() && !i.IsOver {
			for x := range i.AllPlayers {
				msg := models.WinnerMsg{constants.UserWonMsg, *i.Winner}
				err := (i.AllPlayers[x]).WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.IsOver = true
		}
	}
}
