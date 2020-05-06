package main

import (
	"github.com/chainreaction/api"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	store := cookie.NewStore([]byte("iAmLordVoldemort"))
	route.Use(sessions.Sessions("chainreaction", store))
	route.GET("/new", api.CreateNewGame)
	route.GET("/join", api.JoinExistingGame)
	route.Run()
}
