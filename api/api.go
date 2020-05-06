package api

import (
	"net/http"
	"time"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/game"
	"github.com/gin-gonic/gin"
)

// CreateNewGame provides endpoint to create a new game instance
func CreateNewGame(c *gin.Context) {
	var gameInstance game.Instance
	if c.ShouldBindQuery(&gameInstance) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Bad request")
	}
	gameInstance.CreatedOn = time.Now().UTC()
	gameInstance.ExpiresOn = gameInstance.CreatedOn.Add(time.Minute * time.Duration(25))
	gameInstance.CurrentTurn = 0
	gameInstance.RoomName = datastore.GetNewUniqueRoomName()
	if gameInstance.PlayersCount < 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "At least two players needed"})
		return
	}
	if gameInstance.Dimension == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Provide valid dimension value"})
		return
	}
	gameInstance.Board = make([][]game.Pixel, gameInstance.Dimension)
	for i := 0; i < gameInstance.Dimension; i++ {
		gameInstance.Board[i] = make([]game.Pixel, gameInstance.Dimension)
	}
	gameInstance.InitBroadcast()
	datastore.AddGameInstance(&gameInstance)
	c.JSON(http.StatusCreated, gin.H{"GameRoomName": gameInstance.RoomName})
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	roomName := c.Query("room_name")
	if roomName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Please provide a game instance id"})
		return
	}
	gInstance, exists := datastore.GetGameInstance(roomName)
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "No such active game instance found"})
		return
	}
	if gInstance.CurrentActivePlayers == gInstance.PlayersCount {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Game is already full."})
		return
	}
	username := c.Query("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "username cannot be empty."})
		return
	}
	if gInstance.CheckIfUserNameClaimed(username) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Username `" + username + "` is already selected by someone else"})
		return
	}
	color := c.Query("color")
	if color == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Game is already full."})
		return
	}
	if gInstance.CheckIfColorSelected(color) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Color `" + color + "` is already selected by someone else"})
		return
	}

	gInstance.AllPlayers = append(gInstance.AllPlayers, game.Player{username, color, ws})
	gInstance.CurrentActivePlayers++
}
