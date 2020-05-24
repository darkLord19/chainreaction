package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chainreaction/datastore"
	"github.com/gin-gonic/gin"
)

// GetAvailableColors returns available colors for game to choose from
func GetAvailableColors(c *gin.Context) {
	var ret gin.H
	roomName := strings.ToLower(c.Param("name"))
	if roomName == "" {
		ret = gin.H{"Error": "Please provide a game instance id"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
	}
	gInstance, exists := datastore.GetGameInstance(roomName)
	if !exists {
		ret = gin.H{"Error": "No such active game instance found"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusNotFound, ret)
	}
	var clrs []string
	for k, v := range gInstance.AvailableColors {
		if v {
			clrs = append(clrs, k)
		}
	}
	fmt.Println(clrs, gInstance)
	c.JSON(http.StatusOK, gin.H{"colors": clrs})
}
