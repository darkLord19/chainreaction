package main

import (
	"github.com/chainreaction/game"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/new", game.CreateNewGame)
	route.Run()
}
