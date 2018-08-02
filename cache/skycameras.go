package cache

import (
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/source-tools-common/constants"
)

var skyCameras [constants.MAX_MAP_AREAS]types.SkyCamera
var numSkyCameras = 0
var areaSkyCameras [constants.MAX_MAP_AREAS]int

func GetSkyCameras() *[constants.MAX_MAP_AREAS]types.SkyCamera {
	return &skyCameras
}

func GetSkyCamera(index int) *types.SkyCamera {
	if index >= len(skyCameras) {
		return nil
	}
	return &(skyCameras[index])
}

func AddSkyCamera(camera *types.SkyCamera) int {
	skyCameras[numSkyCameras] = *camera
	numSkyCameras++

	return numSkyCameras-1
}

func CountSkyCameras() int {
	return len(skyCameras)
}

func GetAreaSkyCamera(index int) int {
	return areaSkyCameras[index]
}

func SetAreaSkyCamera(index int, value int) {
	areaSkyCameras[index] = value
}