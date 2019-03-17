package examples

import (
	"github.com/faiface/pixel/pixelgl"
)

type INode interface {
	Init()
	Mounted()
	Unmounted()
	Update(win *pixelgl.Window, dt float32)
}

var sceneList map[string]INode
var currentRoot INode

func GetRoot() INode {
	return currentRoot
}

func SetRoot(n INode) {
	if currentRoot != nil {
		currentRoot.Unmounted()
	}
	currentRoot = n
	currentRoot.Mounted()
}

func AddScene(name string, n INode) {
	sceneList[name] = n
}

func GetScene(name string) (INode, bool) {
	s, ok := sceneList[name]
	return s, ok
}

func SwitchScene(name string) {
	s, ok := GetScene(name)
	if !ok {
		return
	}
	SetRoot(s)
}

func init() {
	sceneList = make(map[string]INode)
}
