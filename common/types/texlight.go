package types

import "github.com/go-gl/mathgl/mgl32"

type TexLight struct {
	Name string
	Value mgl32.Vec3
	Filename *string
}