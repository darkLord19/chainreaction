package api

import (
	"log"
	"net/http"
	"time"

	"github.com/chainreaction/api/validators"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/models"
	"github.com/gin-gonic/gin"
)

// CreateNewGame provides endpoint to create a new game instance
func CreateNewGame(c *gin.Context) {
	var gameInstance models.Instance
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

	err := validators.ValidateInstance(&gameInstance)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	gameInstance.Board = make([][]models.Pixel, gameInstance.Dimension)
	for i := 0; i < gameInstance.Dimension; i++ {
		gameInstance.Board[i] = make([]models.Pixel, gameInstance.Dimension)
	}

	gameInstance.InitChannel()
	datastore.AddGameInstance(&gameInstance)
	ret = gin.H{"GameRoomName": gameInstance.RoomName}
	c.JSON(http.StatusCreated, ret)
}
