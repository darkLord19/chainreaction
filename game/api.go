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
	if gameInstance.PlayersCount < 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "At least two players needed"})
	}
	addGameInstance(gameInstance)
	c.JSON(http.StatusCreated, gin.H{"Game Instance": gameInstance})
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	instanceID := c.Query("instance_id")
	if instanceID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Please provide a game instance id"})
	}
	gInstance, exists := getGameInstance(instanceID)
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "No such active game instance found"})
	}
	if gInstance.CurrentActivePlayers == gInstance.PlayersCount {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Game is already full."})
	}
	gInstance.CurrentActivePlayers++
	c.JSON(http.StatusOK, gin.H{"Game Instance": gInstance})
}
