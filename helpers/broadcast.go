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
		p, _ := i.GetPlayerByUsername(move.PlayerUserName)
		move.Color = p.Color
		for x := range i.AllPlayers {
			move.MsgType = constants.MoveBcastMsg
			err := i.AllPlayers[x].WriteToWebsocket(move)
			if err != nil {
				log.Printf("error: %v", err)
				i.AllPlayers[x].CleanupWs()
			}
		}
		if i.IsOver {
			return
		}
	}
}

// BroadcastBoardUpdates broadcasts board updates to users
func BroadcastBoardUpdates(i *models.Instance) {
	for {
		if i.GetBroadcastBoardFlag() {
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
		if i.IsOver {
			return
		}
	}
}

// BroadcastWinner broadcasts winner to users
func BroadcastWinner(i *models.Instance) {
	for {
		if i.IsOver {
			return
		}
		if i.GetIfSomeoneWon() {
			for x := range i.AllPlayers {
				msg := models.WinnerMsg{constants.UserWonMsg, i.Winner.UserName, i.Winner.Color}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.IsOver = true
		}
	}
}
