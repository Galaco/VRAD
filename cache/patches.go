package cache

import (
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/source-tools-common/constants"
)

var patches []types.Patch
var facePatches [constants.MAX_MAP_FACES]int
var faceParents [constants.MAX_MAP_FACES]int
var clusterChildren [constants.MAX_MAP_CLUSTERS]int

func AddPatchToCache(patch *types.Patch) int {
	patches = append(patches, *patch)

	return len(patches)-1
}

func GetPatches() *[]types.Patch{
	return &patches
}

func SetFacePatch(index int, facePatch int) {
	facePatches[index] = facePatch
}

func GetFacePatches() *[constants.MAX_MAP_FACES]int {
	return &facePatches
}


func SetFaceParent(index int, faceParent int) {
	facePatches[index] = faceParent
}

func GetFaceParents() *[constants.MAX_MAP_FACES]int {
	return &faceParents
}


func SetClusterChild(index int, clusterChild int) {
	clusterChildren[index] = clusterChild
}

func GetClusterChildren() *[constants.MAX_MAP_CLUSTERS]int {
	return &clusterChildren
}