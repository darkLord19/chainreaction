package main

import (
	"github.com/chainreaction/api"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/new", api.CreateNewGame)
	route.GET("/join", api.JoinExistingGame)
	route.GET("/play", api.StartGamePlay)
	route.Run()
}
