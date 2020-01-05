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

func GetGameInstance(iid string) (game.Instance, bool) {
	val, ok := activeGameInstances[iid]
	return val, ok
}

func AddGameInstance(gameInstance game.Instance) {
	allGameInstances[gameInstance.InstanceID] = gameInstance
	activeGameInstances[gameInstance.InstanceID] = gameInstance
}
