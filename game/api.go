package game

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CreateNewGame provides endpoint to create a new game instance
func CreateNewGame(c *gin.Context) {
	var gameInstance Instance
	if c.ShouldBindQuery(&gameInstance) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Bad request")
	}
	gameInstance.CreatedOn = time.Now().UTC()
	gameInstance.ExpiresOn = gameInstance.CreatedOn.Add(time.Minute * time.Duration(25))
	gameInstance.CurrentTurn = 0
	gameInstance.InstanceID = uuid.NewV4().String()
	if gameInstance.PlayersCount == 0 {
		gameInstance.PlayersCount = 2
	}
	c.JSON(200, gin.H{"Game Instance": gameInstance})
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	instance, exists := c.Get("instance-id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "No such game instance found"})
	}
	gInstance := getGameInstance(instance)
	c.JSON(http.StatusOK, gin.H{"Game Instance": gInstance})
}
