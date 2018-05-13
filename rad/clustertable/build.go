package clustertable

import (
	"github.com/galaco/vrad/cache"
	"github.com/galaco/source-tools-common/constants"
)


func BuildClusterTable() {
	var leafCount int
	var leafList [constants.MAX_MAP_LEAFS]int

	cache.SetClusterLeafsSize(int(cache.GetLumpCache().Visibility.NumClusters))

	for i := 0; i < int(cache.GetLumpCache().Visibility.NumClusters); i++ {
		leafCount = 0
		for j := 0; j < len(cache.GetLumpCache().Leafs); j++ {
			if int(cache.GetLumpCache().Leafs[j].Cluster) == i {
				leafList[leafCount] = j
				leafCount++
			}
		}

		(*cache.GetClusterLeafs())[i].LeafCount = leafCount
		if leafCount != 0 {
			(*cache.GetClusterLeafs())[i].Leafs = make([]int, leafCount)
			for j := 0; j < leafCount; j++ {
				(*cache.GetClusterLeafs())[i].Leafs[j] = leafList[j]
			}
		}
	}
}