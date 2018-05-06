package brush

import (
	"github.com/galaco/vrad/cache"
)

func GetBrushRecursive(node int, list *[]int) {
	if node < 0 {
		leafs := &cache.GetLumpCache().Leafs
		leafBrushes := &cache.GetLumpCache().LeafBrushes

		leafIndex := -1 - node

		// Add the solids in the leaf
		for i := 0; i < int((*leafs)[leafIndex].NumLeafBrushes); i++ {
			brushIndex := int((*leafBrushes)[int((*leafs)[leafIndex].FirstLeafBrush) + i])
			found := false
			for _,v := range *list {
				if v == brushIndex {
					found = true
					break
				}
			}
			if found == false {
				*list = append(*list, brushIndex)
			}
		}
	} else {
		// recurse
		nodes := &cache.GetLumpCache().Nodes
		n := &(*nodes)[0 + node]

		GetBrushRecursive( int(n.Children[0]), list)
		GetBrushRecursive( int(n.Children[1]), list)
	}
}
