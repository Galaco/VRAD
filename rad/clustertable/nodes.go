package clustertable

import (
	"github.com/galaco/source-tools-common/constants"
	"github.com/galaco/bsp/primitives/node"
	"github.com/galaco/vrad/cache"
)

var leafParents [constants.MAX_MAP_LEAFS]int
var nodeParents [constants.MAX_MAP_NODES]int

func GetLeafParents() [constants.MAX_MAP_LEAFS]int {
	return leafParents
}

func GetNodeParents() [constants.MAX_MAP_NODES]int {
	return nodeParents
}

func MakeParents (nodeNum int, parent int) {
	var j int
	var n *node.Node

	nodeParents[nodeNum] = parent
	n = &cache.GetLumpCache().Nodes[nodeNum]

	for i := 0; i < 2; i++ {
		j = int(n.Children[i])
		if j < 0 {
			leafParents[-j-1] = nodeNum
		} else {
			MakeParents(j, nodeNum)
		}
	}
}

