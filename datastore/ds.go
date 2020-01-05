package datastore

import "github.com/chainreaction/game"

var (
	allGameInstances    map[string]game.Instance
	activeGameInstances map[string]game.Instance
)

func init() {
	allGameInstances = make(map[string]game.Instance)
	activeGameInstances = make(map[string]game.Instance)
}

// GetGameInstance returns game instance from instance id
func GetGameInstance(iid string) (game.Instance, bool) {
	val, ok := activeGameInstances[iid]
	return val, ok
}

// AddGameInstance adds game instance in a data store indexed by instance id
func AddGameInstance(gameInstance game.Instance) {
	allGameInstances[gameInstance.InstanceID] = gameInstance
	activeGameInstances[gameInstance.InstanceID] = gameInstance
}
