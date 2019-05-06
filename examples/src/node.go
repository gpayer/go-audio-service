package examples

import (
	"github.com/gpayer/pixelext/nodes"
)

var sceneList map[string]nodes.Node

func AddScene(name string, n nodes.Node) {
	sceneList[name] = n
}

func GetScene(name string) (nodes.Node, bool) {
	s, ok := sceneList[name]
	return s, ok
}

func SwitchScene(name string) {
	s, ok := GetScene(name)
	if !ok {
		return
	}
	nodes.SceneManager().SetRoot(s)
}

func init() {
	sceneList = make(map[string]nodes.Node)
}
