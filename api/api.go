package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/game"
	"github.com/chainreaction/simulate"
	"github.com/chainreaction/utils"
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
	gameInstance.InitBroadcasts()
	datastore.AddGameInstance(&gameInstance)
	ret = gin.H{"GameRoomName": gameInstance.RoomName}
	c.JSON(http.StatusCreated, ret)
}

// JoinExistingGame provides wndpoint to join already created game
func JoinExistingGame(c *gin.Context) {
	var ret gin.H
	roomName := strings.ToLower(c.Query("roomname"))
	if roomName == "" {
		ret = gin.H{"Error": "Please provide a game instance id"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	gInstance, exists := datastore.GetGameInstance(roomName)
	if !exists {
		ret = gin.H{"Error": "No such active game instance found"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusNotFound, ret)
		return
	}
	if gInstance.CurrentActivePlayers == gInstance.PlayersCount {
		ret = gin.H{"Error": "Game is already full."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}
	username := strings.ToLower(c.Query("username"))
	if username == "" {
		ret = gin.H{"Error": "username cannot be empty."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	if gInstance.CheckIfUserNameClaimed(username) {
		ret = gin.H{"Error": "Username `" + username + "` is already selected by someone else"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}
	color := c.Query("color")
	if color == "" {
		ret = gin.H{"Error": "Game is already full."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	if gInstance.CheckIfColorSelected(color) {
		ret = gin.H{"Error": "Color `" + color + "` is already selected by someone else"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusForbidden, ret)
		return
	}

	gInstance.AllPlayers = append(gInstance.AllPlayers, game.Player{username, color, nil})

	ret = gin.H{"Success": "You have joined the game mothafucka", "game instance": gInstance,
		"user": gin.H{"username": username, "color": color}}

	log.Println(ret)
	c.JSON(200, ret)
}

// StartGamePlay start websocket connection with clients for game play
func StartGamePlay(c *gin.Context) {
	var ret gin.H
	roomname := c.Query("roomname")
	if roomname == "" {
		ret = gin.H{"Error": "Please provide a game instance id"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}
	uname := strings.ToLower(c.Query("username"))
	if uname == "" {
		ret = gin.H{"Error": "username cannot be empty."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}

	gInstance, exists := datastore.GetGameInstance(roomname)

	if gInstance.CurrentActivePlayers == gInstance.PlayersCount {
		ret = gin.H{"Error": "Game is already full."}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}

	if !exists {
		ret = gin.H{"Error": "Wrong room name"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}

	player := gInstance.GetPlayerByID(uname)

	if player == nil {
		ret = gin.H{"Error": "No such user exists in this game"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}

	gInstance.CurrentActivePlayers++

	ws, err := utils.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	player.WsConnection = ws

	go gInstance.BroadcastMoves()
	go gInstance.BroadcastBoardUpdates()
	go gInstance.BroadcastWinner()

	for {
		if gInstance.CurrentActivePlayers != gInstance.PlayersCount {
			continue
		}
		var move game.Move
		err := ws.ReadJSON(&move)
		if err != nil {
			log.Printf("error: %v", err)
			gInstance.AllPlayers[gInstance.CurrentActivePlayers-1].WsConnection = nil
			gInstance.CurrentActivePlayers--
			break
		}
		if move.PlayerUserName == gInstance.AllPlayers[gInstance.CurrentTurn].UserName {
			gInstance.GetBroadcastMove() <- move
			gInstance.CurrentTurn = (gInstance.CurrentTurn + 1) % gInstance.PlayersCount
			simulate.ChainReaction(gInstance, move)
		}
	}

}
