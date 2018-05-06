package light

import (
	"github.com/go-gl/mathgl/mgl32"
)

const MATERIAL_LIGHT_DISABLE = 0
const MATERIAL_LIGHT_POINT = 1
const MATERIAL_LIGHT_DIRECTIONAL = 2
const MATERIAL_LIGHT_SPOT =3

type Descriptor struct {
	Type uint8;
	Color mgl32.Vec3
	Position mgl32.Vec3
	Direction mgl32.Vec3
	Range float32
	Falloff float32
	Attenuation0 float32
	Attenuation1 float32
	Attenuation2 float32
	Theta float32
	Phi float32
	// These aren't used by DX8. . used for software lighting.
	ThetaDot float32
	PhiDot float32
	Flags uint32
}