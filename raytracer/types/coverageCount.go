package types

import (
	"github.com/galaco/vrad/raytracer/cache/triangle"
	"github.com/galaco/vrad/vmath/ssemath/simd"
	"github.com/galaco/vrad/vmath/ssemath"
	"github.com/galaco/vrad/rad/staticprops"
	"github.com/go-gl/mathgl/mgl32"
)

type ITransparentTriangleCallback interface {
	VisitTriangleShouldContinue(data *triangle.TriIntersectData, rays *FourRays, hitMask *simd.Flt4x, b0 *simd.Flt4x, b1 *simd.Flt4x, b2 *simd.Flt4x, hitID int32, color *mgl32.Vec3) bool
}


type CoverageCount struct {
	Coverage simd.Flt4x
}

func NewCoverageCount() CoverageCount {
	return CoverageCount{
		Coverage: ssemath.FourZeros,
	}
}

//NOTE: Avoid globals. Its more important that we pass the color of the triangle in.
func (coverage *CoverageCount) VisitTriangleShouldContinue(triangle *triangle.TriIntersectData, rays *FourRays, hitMask *simd.Flt4x, b0 *simd.Flt4x, b1 *simd.Flt4x, b2 *simd.Flt4x, hitID int32, color *mgl32.Vec3) bool {
	//color := raytracer.GetEnvironment().GetTriangleColor(int(hitID)).X()
	coverage.Coverage = simd.AddSIMD(coverage.Coverage, simd.AndSIMD( *hitMask, simd.ReplicateX4(color.X())))
	coverage.Coverage = simd.MinSIMD(coverage.Coverage, ssemath.FourOnes)

	onesMask := simd.CmpEqSIMD(coverage.Coverage, ssemath.FourOnes)

	// we should continue if the ones that hit the triangle have onesMask set to zero
	// so hitMask & onesMask != hitMask
	// so hitMask & onesMask == hitMask means we're done
	// so ts(hitMask & onesMask == hitMask) != 0xF says go on
	tmp := simd.CmpEqSIMD(simd.AndSIMD(*hitMask, onesMask), *hitMask)
	return 0xF != simd.TestSignSIMD(&tmp)
}

func (coverage *CoverageCount) GetCoverage() simd.Flt4x {
	return coverage.Coverage
}

func (coverage *CoverageCount)  GetFractionVisible() simd.Flt4x {
	return simd.SubSIMD(ssemath.FourOnes, coverage.Coverage)
}


// this will sample the texture to get a coverage at the ray intersection point
type CoverageCountTexture struct {
	*CoverageCount
}

func NewCoverageCountTexture() CoverageCountTexture {
	cc := NewCoverageCount()
	return CoverageCountTexture{
		&cc,
	}
}

func (coverage *CoverageCountTexture) VisitTriangleShouldContinue(triangle *triangle.TriIntersectData, rays *FourRays, hitMask *simd.Flt4x, b0 *simd.Flt4x, b1 *simd.Flt4x, b2 *simd.Flt4x, hitID int32, color *mgl32.Vec3) bool {
	sign := simd.TestSignSIMD(hitMask)
	var addedCoverage [4]float32
	for s := 0; s < 4; s++ {
		addedCoverage[s] = 0.0
		if 0 != (sign >> uint(s)) & 0x1 {
			addedCoverage[s] = float32(staticprops.ComputeCoverageFromTexture(b0[s], b1[s], b2[s], int(hitID)))
		}
	}
	coverage.Coverage = simd.AddSIMD(coverage.Coverage, simd.LoadMultiUnalignedSIMD(addedCoverage))
	coverage.Coverage = simd.MinSIMD(coverage.Coverage, ssemath.FourOnes)
	onesMask := simd.CmpEqSIMD(coverage.Coverage, ssemath.FourOnes)

	// we should continue if the ones that hit the triangle have onesMask set to zero
	// so hitMask & onesMask != hitMask
	// so hitMask & onesMask == hitMask means we're done
	// so ts(hitMask & onesMask == hitMask) != 0xF says go on
	tmp := simd.CmpEqSIMD(simd.AndSIMD( *hitMask, onesMask ), *hitMask)
	return 0xF != simd.TestSignSIMD(&tmp)
}
