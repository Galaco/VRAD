package raytracer

import (
	"github.com/galaco/vrad/vmath/vector"
	"github.com/galaco/vrad/raytracer/cache"
)

type Environment struct {
	Flags uint32
	MinBound vector.Vec3
	MaxBound vector.Vec3
	BackgroundColour vector.Vec4					//< color where no intersection
	OptimizedKDTree []cache.OptimisedKDNode			//< the packed kdtree. root is 0
	OptimizedTriangleList []cache.OptimisedTriangle //< the packed triangles
	TriangleIndexList []int32						//< the list of triangle indices.
	LightList []LightDesc_t							//< the list of lights
	TriangleColors []vector.Vec3					//< color of tries
	TriangleMaterials []int32						//< material index of tries
}

func NewEnvironment() *Environment{
	env := Environment{}
	env.Flags = 0
	env.BackgroundColour = vector.Vec4{
		0,0,0,0,
	}

	return &env
}