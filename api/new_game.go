package api

import (
	"log"
	"net/http"
	"time"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/game"
	"github.com/gin-gonic/gin"
)

// CreateNewGame provides endpoint to create a new game instance
func CreateNewGame(c *gin.Context) {
	var gameInstance game.Instance
	var ret gin.H
	if c.ShouldBindQuery(&gameInstance) != nil {
		ret = gin.H{"Error": "Bad request"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
	}
	gameInstance.CreatedOn = time.Now().UTC()
	gameInstance.ExpiresOn = gameInstance.CreatedOn.Add(time.Minute * time.Duration(25))
	gameInstance.CurrentTurn = 0
	gameInstance.RoomName = datastore.GetNewUniqueRoomName()
	if gameInstance.PlayersCount < 2 {
		ret = gin.H{"Error": "At least two players needed"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	if gameInstance.Dimension == 0 {
		ret = gin.H{"Error": "Provide valid dimension value"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	gameInstance.Board = make([][]game.Pixel, gameInstance.Dimension)
	for i := 0; i < gameInstance.Dimension; i++ {
		gameInstance.Board[i] = make([]game.Pixel, gameInstance.Dimension)
	}
	gameInstance.InitChannel()
	gameInstance.InitGameInstanceMutexes()
	datastore.AddGameInstance(&gameInstance)
	ret = gin.H{"GameRoomName": gameInstance.RoomName}
	c.JSON(http.StatusCreated, ret)
}
