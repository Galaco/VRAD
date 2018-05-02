package brush

import (
	"github.com/galaco/vrad/common"
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/primitives/leaf"
	bspNode "github.com/galaco/bsp/primitives/node"
)

func GetBrushRecursive(node int, list *[]int) {
	if node < 0 {
		leafs := *(*common.GLOBALGET_BSP().GetLump(bsp.LUMP_LEAFS).GetContents()).GetData().(*[]leaf.Leaf)
		leafBrushes := *(*common.GLOBALGET_BSP().GetLump(bsp.LUMP_LEAFBRUSHES).GetContents()).GetData().(*[]uint16)

		leafIndex := -1 - node

		// Add the solids in the leaf
		for i := 0; i < int(leafs[leafIndex].NumLeafBrushes); i++ {
			brushIndex := int(leafBrushes[int(leafs[leafIndex].FirstLeafBrush) + i])
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
		nodes := *(*common.GLOBALGET_BSP().GetLump(bsp.LUMP_NODES).GetContents()).GetData().(*[]bspNode.Node)
		n := &nodes[0 + node]

		GetBrushRecursive( int(n.Children[0]), list)
		GetBrushRecursive( int(n.Children[1]), list)
	}
}
