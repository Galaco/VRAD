package math

import (
	"math"
	"github.com/galaco/vrad/raytracer/cache"
	"github.com/go-gl/mathgl/mgl32"
)

func CalculateTriangleListBounds(optimisedTriangleList *[]cache.OptimisedTriangle, triangles *[]int, minOut *mgl32.Vec3, maxOut *mgl32.Vec3) {
	minOut = &mgl32.Vec3{1.0e23, 1.0e23, 1.0e23}
	maxOut = &mgl32.Vec3{-1.0e23, -1.0e23, -1.0e23}
	for i := 0; i < len(*optimisedTriangleList); i++ {
		tri := (*optimisedTriangleList)[(*triangles)[i]]
		for v := 0; v <3; v++ {
			for c := 0; c <3; c++ {
				minOut[c] = float32(math.Min(float64(minOut[c]), float64(tri.Vertex(v)[c])))
				maxOut[c] = float32(math.Max(float64(maxOut[c]), float64(tri.Vertex(v)[c])))
			}
		}
	}
}
