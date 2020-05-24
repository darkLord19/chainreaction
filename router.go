package main

import (
	"github.com/chainreaction/api"
	"github.com/chainreaction/datastore"
	"github.com/gin-gonic/gin"
)

func setupRouter() {
	router := gin.Default()
	router.GET("/new", api.CreateNewGame)
	// here it is /games/:name due to limitation of httprouter used by gin because
	// it doesn't support wildcard and static route at the same position
	// i.e no support for /login and /:userid as they conflict
	router.GET("/games/:name/join", api.JoinExistingGame)
	router.GET("/games/:name/play", api.StartGamePlay)
	router.GET("/games/:name/colors", api.GetAvailableColors)

	router.Run()
}

func main() {
	setupRouter()
	go datastore.Cleanup()
}
