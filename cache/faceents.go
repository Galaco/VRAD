package cache

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/source-tools-common/constants"
)

var faceEntity [constants.MAX_MAP_FACES]*types.Entity
var faceOffset [constants.MAX_MAP_FACES]mgl32.Vec3		// for rotating bmodels
var faceCentroids [constants.MAX_MAP_EDGES]mgl32.Vec3

func GetFaceEntities() *[constants.MAX_MAP_FACES]*types.Entity {
	return &faceEntity
}
func GetFaceOffsets() *[constants.MAX_MAP_FACES]mgl32.Vec3 {
	return &faceOffset
}

func GetFaceCentroids() *[constants.MAX_MAP_EDGES]mgl32.Vec3 {
	return &faceCentroids
}
