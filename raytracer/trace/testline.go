package trace

import (
	"github.com/galaco/vrad/vmath/ssemath"
	"github.com/galaco/vrad/vmath/ssemath/simd"
	"log"
	"github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/raytracer"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/raytracer/types"
)

// @TODO Use -textureshadows flag...
var textureShadows = false
// @TODO Use -noskyboxrecurse flag...
var noSkyRecurse = false

func TestLineDoesHitSky(start *ssemath.FourVectors, stop ssemath.FourVectors,
	fractionVisible *simd.Flt4x, canRecurse bool, staticPropToSkip int, doDebug bool) {

	log.Panicln("TestLineDoesHitSky: NOT IMPLEMENTED")
	var myrays types.FourRays
	myrays.Origin = *start
	myrays.Direction = stop
	myrays.Direction = myrays.Direction.Sub(myrays.Origin)
	length := myrays.Direction.Length()
	myrays.Direction = *myrays.Direction.MultiplyFlt4x(simd.ReciprocalSIMD(&length))
	var rtResult types.RayTracingResult
	coverageCallback := types.NewCoverageCountTexture()

	var target types.ITransparentTriangleCallback
	if textureShadows == true {
		target = &coverageCallback
	}

	raytracer.GetEnvironment().Trace4Rays(&myrays, ssemath.FourZeros, length, &rtResult, raytracer.TRACE_ID_STATICPROP | staticPropToSkip, &target)

	if doDebug == true {
//		WriteTrace( "trace.txt", myrays, rt_result )
	}

	var aOcclusion [4]float32
	for i := 0; i < 4; i++ {
		aOcclusion[i] = 0.0
		if (rtResult.HitIds[i] != -1) && (rtResult.HitDistance[i] < length[i]) {
			id := raytracer.GetEnvironment().OptimizedTriangleList[rtResult.HitIds[i]].TriIntersectData.NTriangleID
			if 0 == (id & raytracer.TRACE_ID_SKY) {
				aOcclusion[i] = 1.0
			}
		}
	}
	occlusion := simd.LoadMultiUnalignedSIMD(aOcclusion)
	if textureShadows == true {
		occlusion = simd.MaxSIMD(occlusion, coverageCallback.GetCoverage())
	}

	fullyOccluded := simd.TestSignSIMD(simd.CmpGeSIMD(occlusion, ssemath.FourOnes)) == 0xF

	// if we hit sky, and we're not in a sky camera's area, try clipping into the 3D sky boxes
	if (!fullyOccluded) && canRecurse && !noSkyRecurse {
		dir := stop
		dir = dir.Sub(*start)
		dir.VectorNormalize()

		leafIndex := -1
		leafIndex = PointLeafnum(start.Vec(0))
		if leafIndex >= 0 {
			area := cache.GetLumpCache().Leafs[leafIndex].Area()
			if int(area) >= 0 && int(area) < len(cache.GetLumpCache().Areas) {
				if cache.GetAreaSkyCamera(int(area)) < 0 {
					var cam int
					for cam = 0; cam < cache.CountSkyCameras(); cam++ {
						var skystart, skystop ssemath.FourVectors
						skystart.DuplicateVector(&cache.GetSkyCamera(cam).Origin)
						skystop = *start
						skystop = skystop.Multiply(cache.GetSkyCamera(cam).WorldToSky)
						skystart = skystart.Add4Vectors(skystop)

						skystop = dir
						skystop = skystop.Multiply(constants.MAX_TRACE_LENGTH)
						skystop = skystop.Add4Vectors(skystart)
						TestLineDoesHitSky(&skystart, skystop, fractionVisible, false, staticPropToSkip, doDebug)
						occlusion = simd.AddSIMD(occlusion, ssemath.FourOnes)
						occlusion = simd.SubSIMD(occlusion, *fractionVisible)
					}
				}
			}
		}
	}

	occlusion = simd.MaxSIMD(occlusion, ssemath.FourZeros)
	occlusion = simd.MinSIMD(occlusion, ssemath.FourOnes)
	*fractionVisible = simd.SubSIMD(ssemath.FourOnes, occlusion)
}
