package helpers

import (
	"log"

	"github.com/chainreaction/constants"
	"github.com/chainreaction/models"
)

// BroadcastBoardUpdates broadcasts board updates to users
func BroadcastBoardUpdates(i *models.Instance) {
	var move models.MoveMsg
	var val []int
	for {
		i.ReadMoveChan(&move)
		i.ReadBcastBoardChan(&val)
		if val != nil {
			p, _ := i.GetPlayerByUsername(move.PlayerUserName)
			for x := range i.AllPlayers {
				msg := models.NewStateMsg{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName,
					p.Color, move.PlayerUserName, val}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
		}
	}
}

// BroadcastWinner broadcasts winner to users
func BroadcastWinner(i *models.Instance) {
	var val bool
	for {
		i.ReadDidWinChan(&val)
		if val {
			for x := range i.AllPlayers {
				winner := i.GetWinner()
				msg := models.WinnerMsg{constants.UserWonMsg, winner.UserName, winner.Color}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
				}
			}
			i.IsOver = true
			return
		}
	}
}
