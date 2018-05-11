package cache

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/vrad/common/constants"
)

var faceEntity [constants.MAX_MAP_FACES]*types.Entity
var faceOffset [constants.MAX_MAP_FACES]mgl32.Vec3		// for rotating bmodels

func GetFaceEntities() *[constants.MAX_MAP_FACES]*types.Entity {
	return &faceEntity
}
func GetFaceOffsets() *[constants.MAX_MAP_FACES]mgl32.Vec3 {
	return &faceOffset
}