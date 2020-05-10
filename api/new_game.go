package api

import (
	"log"
	"net/http"

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

	err := validators.ValidateInstance(&gameInstance)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	gameInstance.Init(datastore.GetNewUniqueRoomName())
	datastore.AddGameInstance(&gameInstance)
	ret = gin.H{"GameRoomName": gameInstance.RoomName}
	c.JSON(http.StatusCreated, ret)
}
