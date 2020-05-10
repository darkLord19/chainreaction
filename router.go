package main

import (
	"github.com/chainreaction/api"
	"github.com/chainreaction/datastore"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	go datastore.Cleanup()
	route.GET("/new", api.CreateNewGame)
	route.GET("/join", api.JoinExistingGame)
	route.GET("/play", api.StartGamePlay)
	route.GET("/colors", api.GetAvailableColors)
	route.Run()
}
