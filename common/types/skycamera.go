package types

import "github.com/go-gl/mathgl/mgl32"

type SkyCamera struct {
	Origin mgl32.Vec3
	WorldToSky float32
	SkyToWorld float32
	Area int
}
