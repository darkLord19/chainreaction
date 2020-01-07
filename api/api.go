package api

import (
	"log"
	"net/http"
	"time"

	"github.com/chainreaction/simulate"

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
	c.JSON(http.StatusCreated, gin.H{"Game Instance": gameInstance})
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
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
	color := c.Query("color")
	if color == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Game is already full."})
		return
	}
	if gInstance.CheckIfColorSelected(color) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Color " + color + " is already selected by someone else"})
		return
	}

	ws, err := utils.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	gInstance.AllPlayers = append(gInstance.AllPlayers, game.Player{uuid.NewV4().String(), color, ws})
	gInstance.CurrentActivePlayers++

	go gInstance.BroadcastMoves()

	for {
		var move game.Move
		err := ws.ReadJSON(&move)
		if err != nil {
			log.Printf("error: %v", err)
			gInstance.AllPlayers[gInstance.CurrentActivePlayers-1].WsConnection = nil
			break
		}
		gInstance.GetBroadcast() <- move
		simulate.ChainReaction(gInstance, move)
	}
}
