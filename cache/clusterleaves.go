package cache

import (
	"github.com/galaco/vrad/rad/clustertable/list"
)

var clusterLeafs []list.ClusterList

func SetClusterLeafsSize(size int) {
	clusterLeafs = make([]list.ClusterList, size)
}

func AddClusterLeaf(clusterList *list.ClusterList) {
	clusterLeafs = append(clusterLeafs, *clusterList)
}

func GetClusterLeafs() *[]list.ClusterList {
	return &clusterLeafs
}