package api

import (
	"log"
	"net/http"
	"time"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/game"
	"github.com/chainreaction/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
	gameInstance.InstanceID = uuid.NewV4().String()
	if gameInstance.PlayersCount < 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "At least two players needed"})
		return
	}
	gameInstance.InitBroadcast()
	datastore.AddGameInstance(gameInstance)
	c.JSON(http.StatusCreated, gin.H{"Game Instance": gameInstance})
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	ws, err := utils.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	instanceID := c.Query("instance_id")
	if instanceID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Please provide a game instance id"})
		return
	}
	gInstance, exists := datastore.GetGameInstance(instanceID)
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": "No such active game instance found"})
		return
	}
	if gInstance.CurrentActivePlayers == gInstance.PlayersCount {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Game is already full."})
		return
	}
	gInstance.AllPlayers[gInstance.CurrentActivePlayers] = game.Player{uuid.NewV4().String(), "green", ws}
	gInstance.CurrentActivePlayers++

	for {
		var move game.Move
		err := ws.ReadJSON(&move)
		if err != nil {
			log.Printf("error: %v", err)
			gInstance.AllPlayers[gInstance.CurrentActivePlayers-1].WsConnection = nil
			break
		}
		*gInstance.GetBroadcast() <- move
	}
}
