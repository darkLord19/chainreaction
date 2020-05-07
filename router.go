package main

import (
	"time"

	"github.com/chainreaction/api"
	"github.com/chainreaction/datastore"
	"github.com/chainreaction/game"
	"github.com/gin-gonic/gin"
)

func mock() {
	var gameInstance game.Instance
	gameInstance.CreatedOn = time.Now().UTC()
	gameInstance.ExpiresOn = gameInstance.CreatedOn.Add(time.Minute * time.Duration(25))
	gameInstance.CurrentTurn = 0
	gameInstance.RoomName = "test"
	gameInstance.Dimension = 4
	gameInstance.PlayersCount = 2
	gameInstance.Board = make([][]game.Pixel, gameInstance.Dimension)
	for i := 0; i < gameInstance.Dimension; i++ {
		gameInstance.Board[i] = make([]game.Pixel, gameInstance.Dimension)
	}
	gameInstance.InitBroadcasts()
	datastore.AddGameInstance(&gameInstance)

	gameInstance.AllPlayers = append(gameInstance.AllPlayers, game.Player{"test1", "red", nil})
	gameInstance.AllPlayers = append(gameInstance.AllPlayers, game.Player{"test2", "blue", nil})
}

func main() {
	route := gin.Default()
	mock()
	route.GET("/new", api.CreateNewGame)
	route.GET("/join", api.JoinExistingGame)
	route.GET("/play", api.StartGamePlay)
	route.Run()
}
