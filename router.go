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
	gameInstance.InitChannel()
	gameInstance.InitGameInstanceMutexes()
	datastore.AddGameInstance(&gameInstance)

	p := game.Player{}
	p.UserName = "test1"
	p.Color = "red"
	gameInstance.AllPlayers = append(gameInstance.AllPlayers, p)
	p = game.Player{}
	p.UserName = "test2"
	p.Color = "green"
	gameInstance.AllPlayers = append(gameInstance.AllPlayers, p)
}

func main() {
	route := gin.Default()
	mock()
	route.GET("/new", api.CreateNewGame)
	route.GET("/join", api.JoinExistingGame)
	route.GET("/play", api.StartGamePlay)
	route.Run()
}
