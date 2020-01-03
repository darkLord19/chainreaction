package game

var (
	allGameInstances    map[string]Instance
	activeGameInstances map[string]Instance
)

func init() {
	allGameInstances = make(map[string]Instance)
	activeGameInstances = make(map[string]Instance)
}

func getGameInstance(iid string) (Instance, bool) {
	val, ok := activeGameInstances[iid]
	return val, ok
}

func addGameInstance(gameInstance Instance) {
	allGameInstances[gameInstance.InstanceID] = gameInstance
	activeGameInstances[gameInstance.InstanceID] = gameInstance
}
