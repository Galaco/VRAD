package types

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/constants"
)

type BumpLights struct {
	Light [constants.NUM_BUMP_VECTS + 1] mgl32.Vec3
}
