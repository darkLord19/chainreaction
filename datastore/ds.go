package datastore

import (
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/chainreaction/models"
)

var (
	allGameInstances    map[string]*models.Instance
	activeGameInstances map[string]*models.Instance
)

func init() {
	allGameInstances = make(map[string]*models.Instance)
	activeGameInstances = make(map[string]*models.Instance)
}

// GetGameInstance returns game instance from instance id
func GetGameInstance(iid string) (*models.Instance, bool) {
	val, ok := activeGameInstances[iid]
	return val, ok
}

// AddGameInstance adds game instance in a data store indexed by instance id
func AddGameInstance(gameInstance *models.Instance) {
	allGameInstances[gameInstance.RoomName] = gameInstance
	activeGameInstances[gameInstance.RoomName] = gameInstance
}

// GetNewUniqueRoomName returns new random unique name for game room
func GetNewUniqueRoomName() string {
	name := randomdata.SillyName()
	for activeGameInstances[name] != nil {
		name = randomdata.SillyName()
	}
	return strings.ToLower(name)
}

// Cleanup removes expired and over games from data store
func Cleanup() {
	for {
		time.Sleep(1 * time.Minute)
		for k, v := range activeGameInstances {
			if v.ExpiresOn.Sub(time.Now().UTC()) < 0 && v.IsOver {
				delete(activeGameInstances, k)
			}
		}
	}
}
