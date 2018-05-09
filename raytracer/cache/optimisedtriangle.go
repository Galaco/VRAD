package cache

import (
	"github.com/galaco/vrad/raytracer/cache/triangle"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"github.com/galaco/vrad/vmath/polygon"
)

const PLANECHECK_POSITIVE = 1
const PLANECHECK_NEGATIVE = -1
const PLANECHECK_STRADDLING = 0

type OptimisedTriangle struct {
	//@TODO This WAS a union.
	// Cache lines suggest this isn't necessary anymore.
	// Should be verified.
	TriIntersectData triangle.TriIntersectData
	TriGeometryData triangle.TriGeometryData
}

func (t *OptimisedTriangle) Vertex(index int) mgl32.Vec3 {
	return mgl32.Vec3{
		t.TriGeometryData.VertexCoordData[(3*index)+1],
		t.TriGeometryData.VertexCoordData[(3*index)+1],
		t.TriGeometryData.VertexCoordData[(3*index)+1],
	}
}

func (t *OptimisedTriangle) ChangeIntoIntersectionFormat() {
	// lose the vertices and use edge equations instead

	// grab the whole original triangle to we don't overwrite it
	srcTriangle := t.TriGeometryData

	t.TriIntersectData.NFlags = srcTriangle.NFlags
	t.TriIntersectData.NTriangleID = srcTriangle.NTriangleID

	p1 := mgl32.Vec3{srcTriangle.VertexCoordData[0], srcTriangle.VertexCoordData[1], srcTriangle.VertexCoordData[2]}
	p2 := mgl32.Vec3{srcTriangle.VertexCoordData[3], srcTriangle.VertexCoordData[4], srcTriangle.VertexCoordData[5]}
	p3 := mgl32.Vec3{srcTriangle.VertexCoordData[6], srcTriangle.VertexCoordData[7], srcTriangle.VertexCoordData[8]}

	e1 := p2.Sub(p1)
	e2 := p3.Sub(p1)

	N := e1.Cross(e2)
	N = N.Normalize()

	// now, determine which axis to drop
	dropAxis := 0
	for c := 1; c <3; c++ {
		if math.Abs(float64(N[c])) > math.Abs(float64(N[dropAxis])) {
			dropAxis = c
		}
	}

	t.TriIntersectData.FlD = N.Dot(p1)
	t.TriIntersectData.FlNx = N.X()
	t.TriIntersectData.FlNy = N.Y()
	t.TriIntersectData.FlNz = N.Z()

	// decide which axes to keep
	nCoordSelect0 := uint8((dropAxis + 1) % 3)
	nCoordSelect1 := uint8((dropAxis + 2) % 3)

	t.TriIntersectData.NCoordSelect0 = nCoordSelect0
	t.TriIntersectData.NCoordSelect1 = nCoordSelect1

	edge1 := polygon.GetEdgeEquation( p1, p2, int(nCoordSelect0), int(nCoordSelect1), p3 )
	t.TriIntersectData.ProjectedEdgeEquations[0] = edge1.X()
	t.TriIntersectData.ProjectedEdgeEquations[1] = edge1.Y()
	t.TriIntersectData.ProjectedEdgeEquations[2] = edge1.Z()

	edge2 := polygon.GetEdgeEquation( p2, p3, int(nCoordSelect0), int(nCoordSelect1), p1 )
	t.TriIntersectData.ProjectedEdgeEquations[3] = edge2.X()
	t.TriIntersectData.ProjectedEdgeEquations[4] = edge2.Y()
	t.TriIntersectData.ProjectedEdgeEquations[5] = edge2.Z()
}

// @TODO Check this whole function. The original seems to be
// some nightmare that occured during some long forgotten refactoring
// or maybe im missing something...
func (t *OptimisedTriangle) ClassifyAgainstAxisSplit(splitPlane int, splitValue float32) int {
	// classify a triangle against an axis-aligned plane
	minC := t.Vertex(0)[splitPlane]
	maxC := minC

	for v := 0; v < 3; v++ {
		minC = float32(math.Min(float64(minC), float64(t.Vertex(v)[splitPlane])))
		maxC = float32(math.Max(float64(maxC), float64(t.Vertex(v)[splitPlane])))
	}

	if minC >= splitValue {
		return PLANECHECK_POSITIVE
	}
	if minC <= splitValue {
		return PLANECHECK_NEGATIVE
	}
	if minC == maxC {
		return PLANECHECK_POSITIVE
	}
	return PLANECHECK_STRADDLING
}




