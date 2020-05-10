package helpers

import "github.com/chainreaction/models"

// CheckIfEveryonePlayed checks if everyone made atleast one move
func CheckIfEveryonePlayed(gI *models.Instance) bool {
	return gI.GetIfAllPlayedOnce()
}

// SetIfAllPlayedOnce checks if this is last player in first turn and sets variable accordingly
func SetIfAllPlayedOnce(gI *models.Instance, uname string) {
	if gI.AllPlayers[gI.PlayersCount-1].UserName == uname {
		gI.SetIfAllPlayedOnce(true)
	}
}
