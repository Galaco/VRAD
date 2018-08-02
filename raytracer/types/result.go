package types

import (
	"github.com/galaco/vrad/vmath/ssemath"
	"github.com/galaco/vrad/vmath/ssemath/simd"
)

type RayTracingResult struct {
	SurfaceNormal ssemath.FourVectors			// surface normal at intersection
	HitIds [4]int32								// -1=no hit. otherwise, triangle index
	HitDistance simd.Flt4x						// distance to intersection
}

