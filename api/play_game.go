package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/chainreaction/constants"

	"github.com/chainreaction/datastore"
	"github.com/chainreaction/helpers"
	"github.com/chainreaction/models"
	"github.com/chainreaction/simulate"
	"github.com/chainreaction/utils"
	"github.com/gin-gonic/gin"
)

// StartGamePlay start websocket connection with clients for game play
func StartGamePlay(c *gin.Context) {
	var ret gin.H
	roomname := strings.ToLower(c.Param("name"))
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

	if gInstance.GetCurrentActivePlayersCount() == gInstance.PlayersCount {
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

	player, _ := gInstance.GetPlayerByUsername(uname)

	if player == nil {
		ret = gin.H{"Error": "No such user exists in this game"}
		log.Println(ret)
		c.AbortWithStatusJSON(http.StatusBadRequest, ret)
		return
	}

	gInstance.IncCurrentActivePlayers()

	ws, err := utils.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	player.SetWsConnection(ws)

	go helpers.UpdatedBoardUpdates(gInstance)
	go helpers.BroadcastWinner(gInstance)

	for {
		if gInstance.GetCurrentActivePlayersCount() != gInstance.PlayersCount {
			continue
		}
		var move models.MoveMsg
		err := ws.ReadJSON(&move)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		if move.PlayerUserName == gInstance.AllPlayers[gInstance.CurrentTurn].UserName {
			gInstance.RecvMove <- move
			gInstance.CurrentTurn = (gInstance.CurrentTurn + 1) % gInstance.PlayersCount
			err = simulate.ChainReaction(gInstance, move)
			if err != nil {
				helpers.NotifyIndividual(player, models.ErrMsg{constants.InvalidMoveMsg, err.Error()})
			}
		}
	}

}
