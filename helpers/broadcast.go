package helpers

import (
	"log"

	"github.com/chainreaction/constants"
	"github.com/chainreaction/models"
)

// UpdatedBoardUpdates broadcasts board updates to users
func UpdatedBoardUpdates(i *models.Instance) {
	for {
		move := <-i.RecvMove
		val := <-i.UpdatedBoard
		if val != nil {
			p, _ := i.GetPlayerByUsername(move.PlayerUserName)
			for x := range i.AllPlayers {
				msg := models.NewStateMsg{constants.StateUpBcastMsg, i.AllPlayers[i.CurrentTurn].UserName,
					i.AllPlayers[i.CurrentTurn].Color, p.Color, move.PlayerUserName, val}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
					i.DecCurrentActivePlayers()
				}
			}
		}
	}
}

// BroadcastWinner broadcasts winner to users
func BroadcastWinner(i *models.Instance) {
	for {
		if w := i.GetWinner(); w != nil {
			for x := range i.AllPlayers {
				msg := models.WinnerMsg{constants.UserWonMsg, w.UserName, w.Color}
				err := i.AllPlayers[x].WriteToWebsocket(msg)
				if err != nil {
					log.Printf("error: %v", err)
					i.AllPlayers[x].CleanupWs()
					i.DecCurrentActivePlayers()
				}
			}
			i.IsOver = true
			return
		}
	}
}
